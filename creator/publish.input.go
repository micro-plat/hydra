package creator

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/manifoldco/promptui"
	"github.com/micro-plat/lib4go/types"
	"github.com/zkfy/log"
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
			v, err = readFromCli(path, "", fmt.Sprint(key.Interface()), "", value)
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
		if tfield.PkgPath != "" && !tfield.Anonymous { // unexported
			continue
		}
		if !vfield.IsValid() || vfield.IsZero() {
			continue
		}
		if !strings.HasPrefix(fmt.Sprint(vfield.Interface()), "#") {
			continue
		}
		v, ok := input[tfield.Name]
		if !ok {
			v, err = readFromCli(path, vt.Name(), tfield.Name, tfield.Tag.Get("valid"), vfield)
			if err != nil {
				return err
			}
		}
		if err := types.SetWithProperType(v, vfield, tfield); err != nil {
			return err
		}

	}
	return nil
}

var once sync.Once
var logger = log.New(os.Stdout, "", log.Lall|log.Llongcolor)

func readFromCli2(path string, tname, fname string) (string, error) {
	once.Do(func() {
		logger.Println()
		logger.Printf(`----------------------------------------------------------------`)
		logger.Println()
		logger.Printf(` 配置中使用了"#"占位符，需要输入对应的值才能完成安装`)
		logger.Println()
		logger.Printf(` 按回车进行下一个配置的设置 Q:退出`)
		logger.Printf(`----------------------------------------------------------------`)
	})
	rname := strings.Join([]string{tname, fname}, ".")
	logger.Println()
	logger.Warnf("? 请输入%s的值(%s):", rname, path)
	var value string
	fmt.Scan(&value)
	if strings.ToUpper(value) == "Q" {
		return "", fmt.Errorf("未设置%s(%s)值", rname, path)
	}
	if value == "" {
		return readFromCli2(path, tname, fname)
	}
	return value, nil
}
func readFromCli(path string, tname, fname string, tagName string, v reflect.Value) (string, error) {
	rname := strings.Join([]string{tname, fname}, ".")
	validate := func(input string) error {
		if tagName == "" || tagName == "-" {
			return nil
		}
		in := map[string]interface{}{"name": input}
		ck := map[string]interface{}{"name": tagName}

		if ok, err := govalidator.ValidateMap(in, ck); !ok {
			return fmt.Errorf("err %w", err)
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("%s(%s)", rname, path),
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err

}
