package mcp

import (
	"reflect"
	"strings"
)

// buildInputSchema 由 req 类型反射生成 MCP inputSchema（JSON Schema）。
// 读 struct 字段的 json（参数名）、desc（描述）、validate（required）tag。
// 非 struct 或无字段时返回 {type:object} 空对象 schema。
func buildInputSchema(reqType reflect.Type) map[string]interface{} {
	return buildSchema(derefType(reqType))
}

// buildOutputSchema 由 resp 类型反射生成 outputSchema。
func buildOutputSchema(respType reflect.Type) map[string]interface{} {
	return buildSchema(derefType(respType))
}

func buildSchema(t reflect.Type) map[string]interface{} {
	schema := map[string]interface{}{"type": "object"}
	if t == nil || t.Kind() != reflect.Struct {
		return schema
	}
	props := map[string]interface{}{}
	required := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous || !f.IsExported() {
			continue
		}
		name := fieldName(f)
		if name == "-" || name == "" {
			continue
		}
		prop := map[string]interface{}{"type": jsonType(f.Type)}
		if desc := f.Tag.Get("desc"); desc != "" {
			prop["description"] = desc
		}
		props[name] = prop
		if strings.Contains(f.Tag.Get("validate"), "required") {
			required = append(required, name)
		}
	}
	schema["properties"] = props
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

// fieldName 取 json tag 名（去掉 omitempty 等选项），无 json tag 用字段名。
func fieldName(f reflect.StructField) string {
	if tag := f.Tag.Get("json"); tag != "" {
		return strings.Split(tag, ",")[0]
	}
	return f.Name
}

// jsonType Go 类型 → JSON Schema 类型。
func jsonType(t reflect.Type) string {
	switch derefType(t).Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	default:
		return "string"
	}
}

// derefType 解引用指针到元素类型。
func derefType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
