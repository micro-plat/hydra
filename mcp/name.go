package mcp

import "strings"

// handleSuffix 对象服务主处理方法的后缀，镜像 services.defHandler。
const handleSuffix = "Handle"

// toolName 把注册路径转为 MCP tool 名：去首尾斜杠，斜杠转下划线，路径参数去冒号。
//	/v1/my/meta/data/query          -> v1_my_meta_data_query
//	/v1/my/meta/:operationId/biz    -> v1_my_meta_operationId_biz
func toolName(path string) string {
	path = strings.Trim(path, "/")
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	path = strings.ReplaceAll(path, ":", "")
	path = strings.ReplaceAll(path, "/", "_")
	return strings.Trim(path, "_")
}

// camelToPath 驼峰转路径段：大写字母前插'/'并转小写。镜像 services.camelToPath。
//	"Post"        -> "post"
//	"OperationId" -> "operation/id"
func camelToPath(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('/')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// joinPath 拼接路径段并规范化斜杠（用于非动词对象方法的子路径，镜像 services registry.Join 语义）。
func joinPath(base, sub string) string {
	base = strings.TrimRight(base, "/")
	sub = strings.Trim(sub, "/")
	if sub == "" {
		return base
	}
	return base + "/" + sub
}
