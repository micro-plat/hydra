package creator

//检查输入参数，并处理用户输入
func checkAndInput(v interface{}) ([]string, error) {

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
			return
		}
	case reflect.Map:

	default:

	}
}

func checkStruct(value reflect.Value) (string, error) {
	var builder = &strings.Builder{}
	var newValue interface{}
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
	}

	return builder.String(), nil
}
