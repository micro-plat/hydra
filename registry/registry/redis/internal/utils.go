package internal

import (
	"strings"

	"github.com/micro-plat/hydra/registry"
)

//SwapKey 将”/"转换回 “：”拆分
func SwapKey(elem ...string) string {
	var builder strings.Builder
	for _, v := range elem {
		if v == "/" || v == "\\" || strings.TrimSpace(v) == "" {
			continue
		}
		v =  strings.ReplaceAll(v, ":", "###")
		builder.WriteString(strings.Trim(v, "/"))
		builder.WriteString(":")
	}

	str := strings.ReplaceAll(builder.String(), "/", ":")
	return strings.TrimSuffix(str, ":")
}

//SplitKey 拆分“：”key
func SplitKey(key string) []string {
	return strings.Split(key, ":")
}

//SwapPath 将“：”转换回 ”/"拆分
func SwapPath(elem ...string) string {
	var builder strings.Builder
	for _, v := range elem {
		if v == "/" || v == "\\" || strings.TrimSpace(v) == "" {
			continue
		}
		builder.WriteString(strings.Trim(v, "/"))
		builder.WriteString("/")
	}

	str := strings.ReplaceAll(builder.String(), ":", "/")
	str =  strings.ReplaceAll(str, "###", ":")
	return registry.Format(str)
}
 