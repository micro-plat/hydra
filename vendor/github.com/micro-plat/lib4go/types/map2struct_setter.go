package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type setOptions struct {
	isDefaultExists bool
	defaultValue    string
}

//SetWithProperType 设置字段的值
func SetWithProperType(val interface{}, value reflect.Value, field reflect.StructField) error {

	switch value.Kind() {
	case reflect.Int:
		return setIntField(fmt.Sprint(val), 0, value)
	case reflect.Int8:
		return setIntField(fmt.Sprint(val), 8, value)
	case reflect.Int16:
		return setIntField(fmt.Sprint(val), 16, value)
	case reflect.Int32:
		return setIntField(fmt.Sprint(val), 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(fmt.Sprint(val), value, field)
		}
		return setIntField(fmt.Sprint(val), 64, value)
	case reflect.Uint:
		return setUintField(fmt.Sprint(val), 0, value)
	case reflect.Uint8:
		return setUintField(fmt.Sprint(val), 8, value)
	case reflect.Uint16:
		return setUintField(fmt.Sprint(val), 16, value)
	case reflect.Uint32:
		return setUintField(fmt.Sprint(val), 32, value)
	case reflect.Uint64:
		return setUintField(fmt.Sprint(val), 64, value)
	case reflect.Bool:
		return setBoolField(fmt.Sprint(val), value)
	case reflect.Float32:
		return setFloatField(fmt.Sprint(val), 32, value)
	case reflect.Float64:
		return setFloatField(fmt.Sprint(val), 64, value)
	case reflect.String:
		value.SetString(fmt.Sprint(val))
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Ptr:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(fmt.Sprint(val), field, value)
		}
		buff, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return json.Unmarshal(buff, value.Addr().Interface())
	case reflect.Map:
		buff, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return json.Unmarshal(buff, value.Addr().Interface())
	default:
		return fmt.Errorf("%s %s %+v", errUnknownType.Error(), value.Kind().String(), value)
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	nval := val
	if dx, err := NewDecimalFromString(val); err == nil {
		nval = dx.String()
	}
	intVal, err := strconv.ParseInt(nval, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	nval := val
	if dx, err := NewDecimalFromString(val); err == nil {
		nval = dx.String()
	}
	uintVal, err := strconv.ParseUint(nval, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	nval := val
	if dx, err := NewDecimalFromString(val); err == nil {
		nval = dx.String()
	}
	floatVal, err := strconv.ParseFloat(nval, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		nval := val
		if dx, err := NewDecimalFromString(val); err == nil {
			nval = dx.String()
		}
		tv, err := strconv.ParseInt(nval, 10, 0)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil

	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

//SetArray 设置数组的值
func SetArray(vals []interface{}, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := SetWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}

//SetSlice 设置slice的值
func SetSlice(vals []interface{}, value reflect.Value, field reflect.StructField) error {
	slice := reflect.MakeSlice(value.Type(), len(vals), len(vals))
	err := SetArray(vals, slice, field)
	if err != nil {
		return err
	}
	value.Set(slice)
	return nil
}

func setTimeDuration(val string, value reflect.Value, field reflect.StructField) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}
func head(str, sep string) (head string, tail string) {
	idx := strings.Index(str, sep)
	if idx < 0 {
		return str, ""
	}
	return str[:idx], str[idx+len(sep):]
}

var errUnknownType = errors.New("unknown type")
