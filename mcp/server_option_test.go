package mcp

import (
	"testing"

	"github.com/micro-plat/hydra/services"
)

// TestServerOption_Default 验证默认配置：启用、路径 /mcp。
func TestServerOption_Default(t *testing.T) {
	c := services.DefaultMCPConf()
	if !c.Enable {
		t.Fatal("默认应启用 MCP")
	}
	if c.Path != "/mcp" {
		t.Fatalf("默认路径应为 /mcp, 实际 %q", c.Path)
	}
}

// TestServerOption_WithDisable 验证 WithDisable 关闭启用开关。
func TestServerOption_WithDisable(t *testing.T) {
	c := services.DefaultMCPConf()
	WithDisable()(c)
	if c.Enable {
		t.Fatal("WithDisable 后应禁用 MCP")
	}
}

// TestServerOption_WithPath 验证 WithPath 自定义端点路径。
func TestServerOption_WithPath(t *testing.T) {
	c := services.DefaultMCPConf()
	WithPath("/v1/mcp")(c)
	if c.Path != "/v1/mcp" {
		t.Fatalf("路径应为 /v1/mcp, 实际 %q", c.Path)
	}
}

// TestServerOption_Combine 验证多选项组合：禁用 + 自定义路径（启用开关应保留禁用态）。
func TestServerOption_Combine(t *testing.T) {
	c := services.DefaultMCPConf()
	WithDisable()(c)
	WithPath("/v1/mcp")(c)
	if c.Enable {
		t.Fatal("组合后应禁用")
	}
	if c.Path != "/v1/mcp" {
		t.Fatalf("路径应为 /v1/mcp, 实际 %q", c.Path)
	}
}
