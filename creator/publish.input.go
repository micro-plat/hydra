package creator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	logs "github.com/lib4dev/cli/logger"
	"github.com/manifoldco/promptui"
	vc "github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/types"
)

func getCustomString(path, label, value, tname string, input types.XMap) (string, error) {
	if strings.HasPrefix(value, vc.ByInstall) || strings.EqualFold(value, fmt.Sprint(vc.ByInstallI)) {
		fname := getFullName(path, label, tname)
		valueKey := types.GetString(strings.TrimPrefix(value, vc.ByInstall), fname)
		if v := input.GetString(valueKey); v != "" {
			return v, nil
		}
		newValue, err := readFromCli(fname, "required", tname, "", false)
		if err != nil {
			return "", err
		}
		input[valueKey] = newValue
		return newValue, nil
	}
	return value, nil
}

//检查输入参数，并处理用户输入
func checkAndInput(path, label string, value reflect.Value, tnames []string, input types.XMap) error {
	//处理参数类型
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
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
		return checkStruct(path, value, tnames, input)
	case reflect.Map:
		return checkMap(path, value, input)
	default:
		return checkString(path, label, value, tnames, input)
	}
}

func checkString(path, label string, m reflect.Value, tnames []string, input types.XMap) error {
	if !strings.Contains(path, "/var/") || !m.CanSet() {
		return nil
	}
	//处理var自定义配置
	tname := ""
	if len(tnames) > 0 {
		tname = tnames[0]
	}
	v, err := getCustomString(path, label, m.String(), tname, input)
	if err != nil {
		return err
	}
	m.SetString(v)
	return nil
}

func checkMap(path string, m reflect.Value, input types.XMap) (err error) {
	return nil
}

func checkStruct(path string, value reflect.Value, tnames []string, input types.XMap) (err error) {
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
			if err := setSliceValue(path, vfield, tfield, tnames, input); err != nil {
				return err
			}
		case reflect.Map:
			if isRequire(tfield) && (!vfield.IsValid() || vfield.IsZero()) {
				return fmt.Errorf("未配置%s", path)
			}
			if err := checkMap(path, vfield, input); err != nil { //属性类型为map，没有进行验证处理
				return err
			}
		case reflect.Struct:
			if !vfield.IsValid() || vfield.IsZero() {
				if !isRequire(tfield) {
					continue
				}
				setZeroField(vfield, tfield)
			}
			tnames = append(tnames, tfield.Name)
			if err := checkStruct(path, vfield, tnames, input); err != nil {
				return err
			}
		case reflect.Ptr:
			if !vfield.IsValid() || vfield.IsZero() {
				if !isRequire(tfield) {
					continue
				}
				setZeroField(vfield, tfield)
			}
			tnames = append(tnames, tfield.Name)
			if err := checkStruct(path, vfield.Elem(), tnames, input); err != nil {
				return err
			}
		default:
			if err := setFieldValue(path, vfield, tfield, tnames, input); err != nil {
				return err
			}
		}

	}
	return nil
}

func getValues(path string, vfield reflect.Value, tfield reflect.StructField, tnames []string, input types.XMap) (value interface{}, err error) {
	validTagName := tfield.Tag.Get("valid")
	label, msg := getLable(tfield)
	tnames = append(tnames, tfield.Name)
	fname := getFullName(path, label, strings.Join(tnames, "."))

	isArray := vfield.Kind() == reflect.Array || vfield.Kind() == reflect.Slice
	check := func(key string) (interface{}, error) {
		valueKey := types.GetString(key, fname)
		if v, ok := input.Get(valueKey); ok {
			return v, nil
		}
		v, err := readFromCli(fname, validTagName, label, msg, isArray)
		if err != nil {
			return nil, err
		}
		input[valueKey] = v
		return v, nil
	}
	svalue := getSValue(vfield, isArray)
	switch {
	case isArray:
		if err := validateArray(svalue, validTagName, label, msg); err != nil {
			logs.Log.Errorf("%s验证不通过:%+v", path, err)
			return check("")
		}
	case strings.Contains(svalue, vc.ByInstall) || strings.EqualFold(svalue, fmt.Sprint(vc.ByInstallI)):
		if strings.Contains(path, "/var/") { //处理var的自定义配置
			return check("")
		}
		return check(strings.TrimPrefix(svalue, vc.ByInstall))
	case isRequire(tfield) && (!vfield.IsValid() || vfield.IsZero()):
		return check("")
	case !isArray && vfield.IsValid() && !vfield.IsZero():
		if err := validate(svalue, validTagName, label, msg); err != nil {
			logs.Log.Errorf("%s验证不通过:%+v", path, err)
			return check("")
		}
	}
	return nil, nil

}

//为 zero value 设置一个它对应类型的空值，并引导用户填入
func setZeroField(vfield reflect.Value, tfield reflect.StructField) {
	isPtr := false
	t := tfield.Type

	fieldType := t.Kind() //filed的类型
	if fieldType == reflect.Ptr {
		fieldType = t.Elem().Kind()
	}

	if t.Kind() == reflect.Slice { //取数组的元素类型
		t = t.Elem()
	}

	if t.Kind() == reflect.Ptr { //取具体的类型
		isPtr = true
		t = t.Elem()
	}

	v := reflect.New(t)

	if t.Kind() == reflect.Struct {
		initializeStruct(t, v.Elem())
	}

	if !isPtr {
		v = v.Elem()
	}

	switch fieldType {
	case reflect.Slice:
		vfield.Set(reflect.Append(vfield, v))
		return
	case reflect.Struct:
		vfield.Set(v)
		return
	default:
		panic("设置空值错误")
	}
}

func initializeStruct(t reflect.Type, v reflect.Value) {
	if v.Type().Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		if !isRequire(ft) {
			continue
		}
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}

func setSliceValue(path string, vfield reflect.Value, tfield reflect.StructField, tnames []string, input types.XMap) (err error) {

	//处理多个数据值问题
	if vfield.Len() == 0 { //数组为空,元素为结构体时,添加的一个新元素
		if !isRequire(tfield) {
			return nil
		}
		setZeroField(vfield, tfield)
	}
	label, _ := getLable(tfield)
	var v interface{}
	listValue := make([]string, 0, 1)
	for i := 0; i < vfield.Len(); i++ {
		t := vfield.Index(i)
		rootNames := append(tnames, fmt.Sprintf("%s[%d]", tfield.Name, i))
		err := checkAndInput(path, label, t, rootNames, input)
		if err != nil {
			return err
		}
		rootNames = tnames
		listValue = append(listValue, fmt.Sprint(t.Interface()))
	}
	v, err = getValues(path, reflect.ValueOf(listValue), tfield, tnames, input)

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

func setFieldValue(path string, vfield reflect.Value, tfield reflect.StructField, tnames []string, input types.XMap) (err error) {
	v, err := getValues(path, vfield, tfield, tnames, input)
	if err != nil || v == nil {
		return err
	}
	if err := types.SetWithProperType(v, vfield, tfield); err != nil {
		return err
	}
	return nil
}

func readFromCli(name string, tagName string, label string, msg string, isArray bool) (string, error) {
	if tagName == "-" {
		return "", nil
	}

	if isArray {
		//input数组输入项
		prompt := promptui.Prompt{
			Label:    name,
			Validate: func(input string) error { return validateArray(input, tagName, label, msg) },
		}
		result, err := prompt.Run()
		return result, err
	}

	//检查in参数，包括in则使用select,否则为input
	ps := regexp.MustCompile(`^in\((.*)\)`).FindStringSubmatch(tagName)
	if len(ps) == 0 {
		//input输入项
		prompt := promptui.Prompt{
			Label:    name,
			Validate: func(input string) error { return validate(input, tagName, label, msg) },
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

func isRequire(tfield reflect.StructField) bool {
	tagName := tfield.Tag.Get("valid")
	return strings.Contains(tagName, "required")
}

func validateArray(input string, tagName string, label string, msg string) error {
	if input == "" {
		return nil
	}
	//数据元素重复
	items := map[string]bool{}
	for _, data := range strings.Split(input, "|") {
		if _, ok := items[data]; ok {
			return fmt.Errorf("%s", types.GetString(msg, fmt.Sprintf("请输入不重复的%s", label)))
		}
		items[data] = true
		err := validate(data, tagName, label, msg)
		if err != nil {
			return err
		}
	}

	return nil
}

//validate 验证值是否合法
func validate(input string, tagName string, label string, msg string) error {
	if tagName == "" || tagName == "-" {
		return nil
	}
	if len(input) < 1 {
		return errors.New("至少包含1个字符")
	}
	in := map[string]interface{}{"name": input}
	ck := map[string]interface{}{"name": tagName}
	if ok, err := govalidator.ValidateMap(in, ck); !ok {
		return fmt.Errorf("%s (%w)", types.GetString(msg, fmt.Sprintf("请输入正确的%s", label)), err)
	}
	return nil
}

func getSValue(vfield reflect.Value, isArray bool) string {
	svalue := fmt.Sprint(vfield.Interface())
	if isArray {
		listValue := make([]string, 0, 1)
		for i := 0; i < vfield.Len(); i++ {
			t := vfield.Index(i)
			listValue = append(listValue, fmt.Sprint(t.Interface()))
		}
		svalue = strings.Join(listValue, "|")
	}
	return svalue
}

func getFullName(path string, label string, tname string) string {
	return fmt.Sprintf("%s(%s,%s)", label, tname, path)
}

func getLable(tfield reflect.StructField) (string, string) {
	tagName := tfield.Tag.Get("label")
	if tagName == "" || tagName == "-" {
		return tfield.Name, ""
	}
	list := strings.Split(tagName, "|")
	if len(list) > 1 {
		return list[0], list[1]
	}
	return list[0], ""
}
