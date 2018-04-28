package utility

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strings"
)

//GetMapWithQuery 将URL参数转换为map
func GetMapWithQuery(query string) (r map[string]string, err error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return
	}
	r = make(map[string]string)
	for k, v := range values {
		count := len(v)
		if count >= 0 {
			r[k] = v[count-1]
		}
	}
	return
}
func GetQueryWithMap(data map[string]string) (res string, err error) {
	buf := bytes.Buffer{}
	for k, v := range data {
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(v))
		buf.WriteString("&")
	}
	res = buf.String()
	res = strings.TrimRight(res, "&")
	return
}

//GetJSONWithQuery 将URL参数转换为JSON
func GetJSONWithQuery(query string) (res string, err error) {
	result, err := GetMapWithQuery(query)
	buffer, err := json.Marshal(&result)
	if err != nil {
		return
	}
	return string(buffer), nil
}
