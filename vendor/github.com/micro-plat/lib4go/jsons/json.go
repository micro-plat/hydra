package jsons

import (
	"encoding/json"
	"strings"
)

// Escape 把编码 \\u0026，\\u003c，\\u003e 替换为 &,<,>
func Escape(input string) string {
	r := strings.Replace(input, "\\u0026", "&", -1)
	r = strings.Replace(r, "\\u003c", "<", -1)
	r = strings.Replace(r, "\\u003e", ">", -1)
	r = strings.Replace(r, "\n", "", -1)
	return r
}

//Unmarshal 反序列化JSON
func Unmarshal(buf []byte) (c map[string]interface{}, err error) {
	c = make(map[string]interface{})
	err = json.Unmarshal(buf, &c)
	return
}

//Marshal 序列化JSON
func Marshal(data interface{}) (b []byte, err error) {
	return json.Marshal(data)
}
