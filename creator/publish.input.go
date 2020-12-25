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
		skey := fmt.Sprint(key.Interface())
		if !reflect.ValueOf(value).IsValid() || reflect.ValueOf(value).IsZero() {
			continue
		}
		if !strings.HasPrefix(fmt.Sprint(value.Interface()), "#") {
			continue
		}
		v, ok := input[skey]
		if !ok {
			fname := getFullName(path, skey, "")
			v, err = readFromCli(fname, skey, fname, "")
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
		switch vfield.Kind() {
		case reflect.Array, reflect.Slice: //目标字段为数组，将数据源为非数据转化为数组
			if err := setSliceValue(path, vfield, tfield, input); err != nil {
				return err
			}
		case reflect.Map:
			if err := checkMap(path, vfield, input); err != nil {
				return err
			}
		case reflect.Struct:
			if err := checkStruct(path, vfield, input); err != nil {
				return err
			}
		default:
			if err := setFieldValue(path, vfield, tfield, input); err != nil {
				return err
			}
		}

	}
	return nil
}

func getValues(path string, vfield reflect.Value, tfield reflect.StructField, input map[string]interface{}) (value interface{}, err error) {

	validTagName := tfield.Tag.Get("valid")
	lable, msg := getLable(tfield)
	fname := getFullName(path, lable, tfield.Name)
	svalue := fmt.Sprint(vfield.Interface())
	check := func() (interface{}, error) {
		v, ok := input[tfield.Name]
		if !ok {
			v, err = readFromCli(fname, validTagName, lable, msg)
			if err != nil {
				return nil, err
			}
		}
		return v, nil
	}
	switch {
	case strings.HasPrefix(svalue, "#"):
		return check()
	case isRequire(vfield, validTagName) && (!vfield.IsValid() || vfield.IsZero()):
		return check()
	case vfield.IsValid() && !vfield.IsZero() && validate(svalue, validTagName, lable, msg) != nil:
		return check()
	}
	return nil, nil

}

func setSliceValue(path string, vfield reflect.Value, tfield reflect.StructField, input map[string]interface{}) (err error) {
	v, err := getValues(path, vfield, tfield, input)
	if err != nil || v == nil {
		return err
	}
	if r, ok := v.([]interface{}); ok {
		return types.SetSlice(r, vfield, tfield)
	}
	if r, ok := v.(string); ok {
		slist := strings.Split(r, "|")
		vi := make([]interface{}, 0, len(slist))
		for _, l := range slist {
			vi = append(vi, l)
		}
		return types.SetSlice(vi, vfield, tfield)
	}
	return types.SetSlice([]interface{}{v}, vfield, tfield)
}

func setFieldValue(path string, vfield reflect.Value, tfield reflect.StructField, input map[string]interface{}) (err error) {
	v, err := getValues(path, vfield, tfield, input)
	if err != nil || v == nil {
		return err
	}
	if err := types.SetWithProperType(v, vfield, tfield); err != nil {
		return err
	}
	return nil
}

func readFromCli(name string, tagName string, lable string, msg string) (string, error) {

	//检查in参数，包括in则使用select,否则为input
	ps := regexp.MustCompile(`^in\((.*)\)`).FindStringSubmatch(tagName)
	if len(ps) == 0 {

		//input输入项
		prompt := promptui.Prompt{
			Label:    name,
			Validate: func(input string) error { return validate(input, tagName, lable, msg) },
		}
		result, err := prompt.Run()
		return result, err
	}

	//select选择项
	prompt := promptui.Select{
		Label: name,
		Items: strings.Split(ps[1:][0], "|"),
	}

	_, result, err := prompt.Run()
	return result, err

}

func isRequire(input reflect.Value, tagName string) bool {
	return strings.Contains(tagName, "required")
}

//validate 验证值是否合法
func validate(input string, tagName string, lable string, msg string) error {
	if tagName == "" || tagName == "-" {
		return nil
	}
	if len(input) < 1 {
		return errors.New("至少包含1个字符")
	}
	in := map[string]interface{}{"name": input}
	ck := map[string]interface{}{"name": tagName}
	if ok, err := govalidator.ValidateMap(in, ck); !ok {
		return fmt.Errorf("%s (%w)", types.GetString(msg, fmt.Sprintf("请输入正确的%s", lable)), err)
	}
	return nil
}

func getFullName(path string, name string, tname string) string {
	return fmt.Sprintf("%s(%s)", name, strings.Join([]string{tname, path}, ","))

}

func getLable(tfield reflect.StructField) (string, string) {
	tagName := tfield.Tag.Get("lable")
	if tagName == "" || tagName == "-" {
		return tfield.Name, ""
	}
	list := strings.Split(tagName, "|")
	if len(list) > 1 {
		return list[0], list[1]
	}
	return list[0], ""
}
