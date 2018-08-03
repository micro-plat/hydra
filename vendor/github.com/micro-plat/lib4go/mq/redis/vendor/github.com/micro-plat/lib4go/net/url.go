package net

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

var LocalIP string

func init() {
	LocalIP = GetLocalIPAddress()
}

//QueryStringToMap 将URL查询字符串中的参数转换成map
func QueryStringToMap(urlQuery string) (result map[string]interface{}, err error) {
	index := strings.IndexAny(urlQuery, "?")
	if index == -1 || index >= len(urlQuery)-1 {
		return
	}
	values, err := url.ParseQuery(urlQuery[index+1:])
	if err != nil {
		err = fmt.Errorf("url ParseQuery fail: %v", err)
		return
	}
	result = make(map[string]interface{})
	for k, v := range values {
		if len(v) == 1 {
			result[k] = v[0]
		} else {
			result[k] = v
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
