package main

import (
	"fmt"
	"strings"
)

func main() {
	ss := "/taosytest/test-rgtredis/api/taosy/conf/taosytest:test-rgtredis:api:taosy:conf:router"
	fmt.Println(JoinR(ss))
}

func JoinR(elem ...string) string {
	var builder strings.Builder
	for _, v := range elem {
		if v == "/" || v == "\\" || strings.TrimSpace(v) == "" {
			continue
		}
		builder.WriteString(strings.Trim(v, "/"))
		builder.WriteString(":")
	}

	str := strings.ReplaceAll(builder.String(), "/", ":")
	return strings.TrimSuffix(str, ":")
}
