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
		if r, b := i.Get(v); !b || r == "" {
			return fmt.Errorf("%s值不能为空", v)
		}
	}
	return nil
}

func (i *inputParams) Get(name string) (string, bool) {
	if c, ok := i.data.Get(name); ok {
		return fmt.Sprint(c), ok
	}
	return "", false
}
func (i *inputParams) GetString(name string, p ...string) string {
	v, b := i.Get(name)
	if b {
		return v
	}
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

//GetInt 获取int数字
func (i *inputParams) GetInt(name string, p ...int) int {
	value, b := i.Get(name)
	if b {
		v, err := strconv.Atoi(value)
		if err == nil {
			return v
		}
	}

	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//GetInt64 获取int64数字
func (i *inputParams) GetInt64(name string, p ...int64) int64 {
	value, b := i.Get(name)
	if b {
		v, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return v
		}
	}
	if len(p) > 0 {
		return p[0]
	}
	return 0
}

//GetFloat64 获取float64数字
func (i *inputParams) GetFloat64(name string, p ...float64) float64 {
	value, b := i.Get(name)
	if b {
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return v
		}
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
	value, b := i.Get(name)
	var v time.Time
	var err error
	if b {
		v, err = time.Parse(format, value)
		if err == nil {
			return v, nil
		}
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
		if v, b := i.Get(s[2 : len(s)-1]); b {
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
		if v, b := i.Get(s[1:]); b {
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
func (i *inputParams) Clear() {
	i.data = nil
}
