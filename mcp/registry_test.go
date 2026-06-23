package mcp

import (
	"reflect"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/context"
)

// resetTools 清空全局 tool 表（测试隔离）。
func resetTools() {
	toolsMu.Lock()
	tools = map[string]*Tool{}
	toolsMu.Unlock()
}

// ---- toolName ----

func TestToolName(t *testing.T) {
	cases := map[string]string{
		"/v1/my/meta/data/query":       "v1_my_meta_data_query",
		"/v1/my/meta/:operationId/biz": "v1_my_meta_operationId_biz",
		"v1/x/y":                       "v1_x_y",
		"//v1//x//":                    "v1_x",
	}
	for in, want := range cases {
		if got := toolName(in); got != want {
			t.Errorf("toolName(%q)=%q want %q", in, got, want)
		}
	}
}

// ---- buildConfig：opts 解析 ----

func TestBuildConfig(t *testing.T) {
	// string -> Desc
	c, err := buildConfig("描述")
	if err != nil || c.Desc != "描述" {
		t.Fatalf("string opts: %+v err=%v", c, err)
	}
	// Option 应用
	c, _ = buildConfig(Name("n"), Desc("d"))
	if c.Name != "n" || c.Desc != "d" {
		t.Fatalf("option opts: %+v", c)
	}
	// Config 合并
	c, _ = buildConfig(&Config{Name: "cn", Desc: "cd"})
	if c.Name != "cn" || c.Desc != "cd" {
		t.Fatalf("config opts: %+v", c)
	}
	// 后出现的覆盖先出现的
	c, _ = buildConfig(Desc("first"), Desc("second"))
	if c.Desc != "second" {
		t.Fatalf("override: %+v", c)
	}
	// 非法类型报错
	if _, err := buildConfig(123); err == nil {
		t.Fatal("want error for int opts")
	}
}

// ---- buildInputSchema：struct 反射 ----

type schemaReq struct {
	Name string `json:"name" validate:"required" desc:"姓名"`
	Age  int    `json:"age" desc:"年龄"`
}

func TestBuildInputSchema(t *testing.T) {
	s := buildInputSchema(reflect.TypeOf(schemaReq{}))
	if s["type"] != "object" {
		t.Fatalf("type=%v", s["type"])
	}
	props := s["properties"].(map[string]interface{})
	name := props["name"].(map[string]interface{})
	if name["type"] != "string" || name["description"] != "姓名" {
		t.Errorf("name prop=%+v", name)
	}
	req := s["required"].([]string)
	if len(req) != 1 || req[0] != "name" {
		t.Errorf("required=%+v", req)
	}
}

// ---- registerTool：函数形式端到端 ----

type addReq struct {
	A int `json:"a" validate:"required" desc:"参数A"`
}
type addResp struct {
	Sum int `json:"sum"`
}

func addFn(ctx context.IContext, req addReq) (addResp, error) {
	return addResp{Sum: req.A + req.A}, nil
}

func TestRegisterToolFunction(t *testing.T) {
	resetTools()
	registerTool("/v1/add", addFn, "加法")

	tool, ok := Get("v1_add")
	if !ok {
		t.Fatal("tool v1_add not registered")
	}
	if tool.Config.Desc != "加法" {
		t.Errorf("desc=%q", tool.Config.Desc)
	}
	if tool.Path != "/v1/add" {
		t.Errorf("path=%q", tool.Path)
	}
	props := tool.Config.InputSchema["properties"].(map[string]interface{})
	if props["a"].(map[string]interface{})["type"] != "number" {
		t.Errorf("a schema=%+v", props["a"])
	}
}

// ---- registerTool：Name 覆盖 ----

func TestRegisterToolNameOverride(t *testing.T) {
	resetTools()
	registerTool("/v1/add", addFn, Name("custom_add"))
	if _, ok := Get("custom_add"); !ok {
		t.Fatal("custom_add not registered")
	}
	if _, ok := Get("v1_add"); ok {
		t.Fatal("default name should be overridden")
	}
}

// ---- registerTool：对象服务 + MCPDesc ----

type greetReq struct {
	Name string `json:"name" validate:"required" desc:"姓名"`
}
type greetResp struct {
	Msg string `json:"msg"`
}

type greetSvc struct {
	Service
}

func (s *greetSvc) PostHandle(ctx context.IContext, req greetReq) (*greetResp, error) {
	return &greetResp{Msg: "hi " + req.Name}, nil
}

func newGreetSvc() *greetSvc {
	s := &greetSvc{}
	s.MCPDesc(s.PostHandle, "问候")
	return s
}

func TestRegisterToolObject(t *testing.T) {
	resetTools()
	registerTool("/v1/greet", newGreetSvc()) // 无 opts，描述取自 MCPDesc

	tool, ok := Get("v1_greet")
	if !ok {
		t.Fatal("tool v1_greet not registered")
	}
	if tool.Config.Desc != "问候" {
		t.Errorf("desc=%q want 问候", tool.Config.Desc)
	}
	// 对象动词方法源路径 = basePath（不拼接 /post）
	if tool.Path != "/v1/greet" {
		t.Errorf("path=%q want /v1/greet", tool.Path)
	}
}

// ---- List 排序 ----

func TestListSorted(t *testing.T) {
	resetTools()
	registerTool("/v1/zzz", addFn)
	registerTool("/v1/aaa", addFn)
	list := List()
	if len(list) != 2 || list[0].Name != "v1_aaa" || list[1].Name != "v1_zzz" {
		t.Fatalf("list order: %+v", list)
	}
}

// ---- registerTool：非 typed 显式登记应启动期 panic（Logic 1）----

// legacyNonTyped 旧签名 handler（func(ctx) interface{}），非 typed，仅用于验证报错。
func legacyNonTyped(ctx context.IContext) interface{} { return nil }

func TestRegisterToolNonTypedPanics(t *testing.T) {
	resetTools()
	// 显式 .WithMCP()/.MCP() 登记非 typed handler 属编程错误，应 fail-fast panic
	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				msg, _ := r.(string)
				if !strings.Contains(msg, "非 typed") {
					t.Errorf("panic 信息应含「非 typed」，实际: %v", r)
				}
			}
		}()
		registerTool("/v1/legacy", legacyNonTyped)
	}()
	if !panicked {
		t.Fatal("registerTool 对非 typed handler 应 panic，实际未 panic")
	}
	// panic 发生在登记前，确认未被登记
	if _, ok := Get("v1_legacy"); ok {
		t.Fatal("非 typed handler 不应被登记为 tool")
	}
}
