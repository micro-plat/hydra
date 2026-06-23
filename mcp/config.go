package mcp

import "fmt"

// Config MCP 工具配置。Name/Desc 为空时走默认值；InputSchema/OutputSchema 为空时走 Req/Resp 反射。
type Config struct {
	Name         string
	Desc         string
	InputSchema  map[string]interface{}
	OutputSchema map[string]interface{}
	// Annotations MCP 工具标注（2024-serviceservice-05 规范）：readOnlyHint / destructiveHint /
	// idempotentHint / openWorldHint。为空表示未声明；tools/list 仅在非空时输出该字段。
	Annotations map[string]bool
}

// buildConfig 解析注册参数：string 当 Desc；Option 应用；*Config/Config 合并。
// 优先级：参数列表中后出现的覆盖先出现的。返回的 Config 空字段由调用方补默认值。
func buildConfig(opts ...interface{}) (*Config, error) {
	c := &Config{}
	for _, o := range opts {
		switch v := o.(type) {
		case nil:
			continue
		case string:
			c.Desc = v
		case Option:
			v(c)
		case *Config:
			if v != nil {
				mergeConfig(c, v)
			}
		case Config:
			mergeConfig(c, &v)
		default:
			return nil, fmt.Errorf("mcp: 不支持的配置参数类型 %T", o)
		}
	}
	return c, nil
}

// mergeConfig 把 src 的非零字段合并到 dst（空字符串/nil 不覆盖）。
func mergeConfig(dst, src *Config) {
	if src.Name != "" {
		dst.Name = src.Name
	}
	if src.Desc != "" {
		dst.Desc = src.Desc
	}
	if src.InputSchema != nil {
		dst.InputSchema = src.InputSchema
	}
	if src.OutputSchema != nil {
		dst.OutputSchema = src.OutputSchema
	}
	if src.Annotations != nil {
		if dst.Annotations == nil {
			dst.Annotations = map[string]bool{}
		}
		for k, v := range src.Annotations {
			dst.Annotations[k] = v
		}
	}
}
