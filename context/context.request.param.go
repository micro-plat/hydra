package context

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type inputParams struct {
	data IData
}

//Check 检查是否包含指定的参数
func (i *inputParams) Check(names ...string) error {
	for _, v := range names {
		if r, err := i.Get(v); err != nil || r == "" {
			return fmt.Errorf("%s值不能为空", v)
		}
	}
	return nil
}

func (i *inputParams) Get(name string) (string, error) {
	return i.data.Get(name)
}
func (i *inputParams) GetString(name string, p ...string) string {
	v, err := i.Get(name)
	if err == nil {
		return v
	}
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

//GetInt 获取int数字
func (i *inputParams) GetInt(name string, p ...int) int {
	value, err := i.Get(name)
	var v int
	if err == nil {
		v, err = strconv.Atoi(value)
	}
	if err == nil {
		return v
	}
	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//GetInt64 获取int64数字
func (i *inputParams) GetInt64(name string, p ...int64) int64 {
	value, err := i.Get(name)
	var v int64
	if err == nil {
		v, err = strconv.ParseInt(value, 10, 64)
	}
	if err == nil {
		return v
	}
	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//GetFloat64 获取float64数字
func (i *inputParams) GetFloat64(name string, p ...float64) float64 {
	value, err := i.Get(name)
	var v float64
	if err == nil {
		v, err = strconv.ParseFloat(value, 64)
	}
	if err == nil {
		return v
	}
	if len(p) > 0 {
		return p[0]
	}
	return 0
}
func (i *inputParams) GetDataTime(name string, p ...time.Time) (time.Time, error) {
	return i.GetDataTimeByFormat(name, "20060102150405", p...)
}

//GetFloat64 获取float64数字
func (i *inputParams) GetDataTimeByFormat(name string, format string, p ...time.Time) (time.Time, error) {
	value, err := i.Get(name)
	var v time.Time
	if err == nil {
		v, err = time.Parse(format, value)
	}
	if err == nil {
		return v, nil
	}
	if len(p) > 0 {
		return p[0], nil
	}
	return v, err
}

//Translate 翻译带有@变量的字符串
func (i *inputParams) Translate(format string, a bool) (string, int) {
	brackets, _ := regexp.Compile(`\{@\w+\}`)
	v := 0
	result := brackets.ReplaceAllStringFunc(format, func(s string) string {
		if v, err := i.Get(s[2 : len(s)-1]); err == nil {
			return v
		}
		v++
		if a {
			return ""
		}
		return s
	})
	word, _ := regexp.Compile(`@\w+`)
	result = word.ReplaceAllStringFunc(result, func(s string) string {
		if v, err := i.Get(s[1:]); err == nil {
			return v
		}
		v++
		if a {
			return ""
		}
		return s
	})
	return result, v
}
