package creator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/manifoldco/promptui"
	"github.com/micro-plat/lib4go/types"
)

//检查输入参数，并处理用户输入
func checkAndInput(path string, v interface{}, input map[string]interface{}) error {

	//处理参数类型
	value := reflect.ValueOf(v)
	vt := reflect.TypeOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		vt = vt.Elem()
	}
	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value.Interface())
	}

	//2. 根据类型分别处理
	switch value.Kind() {
	case reflect.Struct:
		if !reflect.ValueOf(value).IsValid() || reflect.ValueOf(value).IsZero() {
			return nil
		}
		return checkStruct(path, value, input)
	case reflect.Map:
		return checkMap(path, value, input)
	default:
		return checkString(path, value, input)
	}
}

func checkString(path string, m reflect.Value, input map[string]interface{}) error {
	return nil
}

func checkMap(path string, m reflect.Value, input map[string]interface{}) (err error) {
	keys := m.MapKeys()
	for _, key := range keys {
		value := m.MapIndex(key)
		if !reflect.ValueOf(value).IsValid() || reflect.ValueOf(value).IsZero() {
			continue
		}
		if !strings.HasPrefix(fmt.Sprint(value.Interface()), "#") {
			continue
		}
		v, ok := input[fmt.Sprint(key.Interface())]
		if !ok {
			v, err = readFromCli(path, "", fmt.Sprint(key.Interface()), "")
			if err != nil {
				return err
			}
		}
		m.SetMapIndex(key, reflect.ValueOf(v))
		continue

	}
	return nil
}
func checkStruct(path string, value reflect.Value, input map[string]interface{}) (err error) {
	vt := reflect.TypeOf(value.Interface())
	for i := 0; i < value.NumField(); i++ {
		tfield := vt.Field(i)
		vfield := value.Field(i)

		//私有字段
		if tfield.PkgPath != "" && !tfield.Anonymous { // unexported
			continue
		}
		tagName := tfield.Tag.Get("valid")
		svalue := fmt.Sprint(vfield.Interface())
		check := func() error {
			v, ok := input[tfield.Name]
			if !ok {
				v, err = readFromCli(path, vt.Name(), tfield.Name, tagName)
				if err != nil {
					return err
				}
			}
			if err := types.SetWithProperType(v, vfield, tfield); err != nil {
				return err
			}
			return nil
		}
		switch {
		case strings.HasPrefix(svalue, "#"):
			err = check()
		case isRequire(vfield, tagName) && (!vfield.IsValid() || vfield.IsZero()):
			err = check()
		case vfield.IsValid() && !vfield.IsZero() && validate(svalue, tagName) != nil:
			err = check()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func readFromCli(path string, tname, fname string, tagName string) (string, error) {
	rname := strings.Join([]string{tname, fname}, ".")

	//检查in参数，包括in则使用select,否则为input
	ps := regexp.MustCompile(`^in\((.*)\)`).FindStringSubmatch(tagName)
	if len(ps) == 0 {

		//input输入项
		prompt := promptui.Prompt{
			Label:    fmt.Sprintf("%s(%s)", rname, path),
			Validate: func(input string) error { return validate(input, tagName) },
		}
		result, err := prompt.Run()
		return result, err
	}

	//select选择项
	prompt := promptui.Select{
		Label: fmt.Sprintf("%s(%s)", rname, path),
		Items: strings.Split(ps[1:][0], "|"),
	}

	_, result, err := prompt.Run()
	return result, err

}

func isRequire(input reflect.Value, tagName string) bool {
	return strings.Contains(tagName, "required")
}

//validate 验证值是否合法
func validate(input string, tagName string) error {
	if tagName == "" || tagName == "-" {
		return nil
	}
	if len(input) < 1 {
		return errors.New("至少包含1个字符")
	}
	in := map[string]interface{}{"name": input}
	ck := map[string]interface{}{"name": tagName}
	if ok, err := govalidator.ValidateMap(in, ck); !ok {
		return fmt.Errorf("err %w", err)
	}
	return nil
}
