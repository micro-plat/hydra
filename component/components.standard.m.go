package component

import "strings"

type list []string

var requestMethods list = []string{"get", "post", "put", "delete", "options"}

func (m list) Contains(name string) bool {
	for _, i := range m {
		if i == name {
			return true
		}
	}
	return false
}
func xContains(l list, name string) bool {
	return l.Contains(name)
}
func getMethod(text string) (string, []string) {
	index := strings.LastIndex(text, "/$")
	if index == -1 {
		return text, []string{}
	}
	rname := text[0:index]
	method := text[index+2:]
	index = strings.LastIndex(rname, "/$")
	if index == -1 {
		return rname, []string{method}
	}
	url, lst := getMethod(rname)
	return url, append(lst, method)
}
