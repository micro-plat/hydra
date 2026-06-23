package mcp

import (
	"reflect"
	"runtime"
	"strings"
)

// Service 对象服务嵌入基类，用方法引用声明各 Handle 方法的 MCP 描述（编译期校验方法存在）。
//
//	type BizSvc struct{ mcp.Service }
//	func (s *BizSvc) PostHandle(ctx hydra.IContext, req Req) (*Resp, error) { ... }
//	s.MCPDesc(s.PostHandle, "批量保存向导表单数据")
type Service struct {
	configs map[string]*Config
}

// IMCPService 由嵌入 mcp.Service 的对象服务实现，供注册时读取方法级 MCP 配置。
type IMCPService interface {
	MCPConfigs() map[string]*Config
}

// MCPConfigs 返回方法级 MCP 配置（方法名 → Config）。实现 IMCPService。
func (s *Service) MCPConfigs() map[string]*Config {
	if s.configs == nil {
		return map[string]*Config{}
	}
	return s.configs
}

// MCPDesc 声明某 Handle 方法的描述。method 为方法引用（如 s.PostHandle），编译期校验方法存在。
func (s *Service) MCPDesc(method interface{}, desc string) {
	s.store(method, func(c *Config) { c.Desc = desc })
}

// MCP 声明某 Handle 方法的完整配置，opts 规则同 buildConfig（string/Option/Config）。
func (s *Service) MCP(method interface{}, opts ...interface{}) {
	cfg, err := buildConfig(opts...)
	if err != nil || cfg == nil {
		return
	}
	s.store(method, func(c *Config) { mergeConfig(c, cfg) })
}

func (s *Service) store(method interface{}, apply func(*Config)) {
	name := methodRefName(method)
	if name == "" {
		return
	}
	if s.configs == nil {
		s.configs = map[string]*Config{}
	}
	c, ok := s.configs[name]
	if !ok {
		c = &Config{}
		s.configs[name] = c
	}
	apply(c)
}

// methodRefName 从方法引用提取方法名（如 "PostHandle"），
// 利用 runtime.FuncForPC 取方法值对应的函数名。
func methodRefName(method interface{}) string {
	v := reflect.ValueOf(method)
	if !v.IsValid() || v.Kind() != reflect.Func {
		return ""
	}
	pc := runtime.FuncForPC(v.Pointer())
	if pc == nil {
		return ""
	}
	// 形如 pkg.(*BizSvc).PostHandle-fm 或 pkg.BizSvc.PostHandle-fm
	name := pc.Name()
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[i+1:]
	}
	return strings.TrimSuffix(name, "-fm")
}
