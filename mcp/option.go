package mcp

// Option MCP 工具配置项，修改 Config 的某个字段。
type Option func(*Config)

// Name 设置 tool 名（覆盖路径生成的默认名）。优先级见 service-design §8。
func Name(name string) Option {
	return func(c *Config) { c.Name = name }
}

// Desc 设置 tool 描述。MCP client/LLM 依赖描述判断何时调用，建议填写。
func Desc(description string) Option {
	return func(c *Config) { c.Desc = description }
}

// InputSchema 覆盖默认从 Req 反射得到的输入 schema。
func InputSchema(schema map[string]interface{}) Option {
	return func(c *Config) { c.InputSchema = schema }
}

// OutputSchema 覆盖默认从 Resp 反射得到的输出 schema。
func OutputSchema(schema map[string]interface{}) Option {
	return func(c *Config) { c.OutputSchema = schema }
}

// Annotations 设置工具标注（MCP 2024-serviceservice-05 规范的 hints）。
// 逐 key 并入 Config.Annotations（非整体替换），可与 ReadOnly/Destructive 等 builder 叠加。
func Annotations(a map[string]bool) Option {
	return func(c *Config) {
		if a == nil {
			return
		}
		if c.Annotations == nil {
			c.Annotations = map[string]bool{}
		}
		for k, v := range a {
			c.Annotations[k] = v
		}
	}
}

// ReadOnly 标注为只读查询工具：readOnlyHint=true、destructiveHint=false、idempotentHint=true。
// 供调用方（如 ai-agent 安全护栏）判定可安全放行。
func ReadOnly() Option {
	return func(c *Config) {
		Annotations(map[string]bool{
			"readOnlyHint":    true,
			"destructiveHint": false,
			"idempotentHint":  true,
		})(c)
	}
}

// Destructive 标注为破坏性（写/删）工具：destructiveHint=true、idempotentHint=false。
func Destructive() Option {
	return func(c *Config) {
		Annotations(map[string]bool{
			"destructiveHint": true,
			"idempotentHint":  false,
		})(c)
	}
}

// Idempotent 标注为幂等写入工具：idempotentHint=true、destructiveHint=false。
func Idempotent() Option {
	return func(c *Config) {
		Annotations(map[string]bool{
			"idempotentHint":  true,
			"destructiveHint": false,
		})(c)
	}
}

// WithAll 一次性声明 name+desc+标注，返回单个 Option。
// extra 接收 ReadOnly()/Destructive()/Idempotent()/Annotations(...) 等，逐个并入。
func WithAll(name, desc string, extra ...Option) Option {
	return func(c *Config) {
		if name != "" {
			c.Name = name
		}
		if desc != "" {
			c.Desc = desc
		}
		for _, o := range extra {
			if o != nil {
				o(c)
			}
		}
	}
}

// WithName 仅声明 tool 名（覆盖路径生成的默认名）。
func WithName(name string) Option { return func(c *Config) { c.Name = name } }

// WithDesc 声明 desc + 标注（无显式 name，对象方法级常用）。
func WithDesc(desc string, extra ...Option) Option {
	return func(c *Config) {
		if desc != "" {
			c.Desc = desc
		}
		for _, o := range extra {
			if o != nil {
				o(c)
			}
		}
	}
}
