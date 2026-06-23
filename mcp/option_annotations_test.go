package mcp

import "testing"

// ---- 语义化 builder：Annotations 映射正确 ----

func TestReadOnlyOption(t *testing.T) {
	c := &Config{}
	ReadOnly()(c)
	if c.Annotations["readOnlyHint"] != true ||
		c.Annotations["destructiveHint"] != false ||
		c.Annotations["idempotentHint"] != true {
		t.Fatalf("ReadOnly annotations=%+v", c.Annotations)
	}
	if len(c.Annotations) != 3 {
		t.Errorf("ReadOnly 应恰好设 3 个 hint，实际 %d 个: %+v", len(c.Annotations), c.Annotations)
	}
}

func TestDestructiveOption(t *testing.T) {
	c := &Config{}
	Destructive()(c)
	if c.Annotations["destructiveHint"] != true ||
		c.Annotations["idempotentHint"] != false {
		t.Fatalf("Destructive annotations=%+v", c.Annotations)
	}
}

func TestIdempotentOption(t *testing.T) {
	c := &Config{}
	Idempotent()(c)
	if c.Annotations["idempotentHint"] != true ||
		c.Annotations["destructiveHint"] != false {
		t.Fatalf("Idempotent annotations=%+v", c.Annotations)
	}
}

// Annotations(a) 逐 key 并入，可多次叠加（不整体替换）。
func TestAnnotationsOptionMerge(t *testing.T) {
	c := &Config{}
	Annotations(map[string]bool{"readOnlyHint": true})(c)
	Annotations(map[string]bool{"openWorldHint": true})(c)
	if c.Annotations["readOnlyHint"] != true || c.Annotations["openWorldHint"] != true {
		t.Fatalf("叠加后 annotations=%+v", c.Annotations)
	}
}

// Annotations(nil) 安全无副作用。
func TestAnnotationsOptionNil(t *testing.T) {
	c := &Config{Annotations: map[string]bool{"readOnlyHint": true}}
	Annotations(nil)(c)
	if c.Annotations["readOnlyHint"] != true {
		t.Fatalf("nil annotations 应为无操作: %+v", c.Annotations)
	}
}

// ---- mergeConfig：Annotations 合并 ----

func TestMergeConfigAnnotations(t *testing.T) {
	// dst 已有 readOnlyHint；src 补另两项；已有 key 不丢。
	dst := &Config{Annotations: map[string]bool{"readOnlyHint": true}}
	mergeConfig(dst, &Config{Annotations: map[string]bool{
		"destructiveHint": false,
		"idempotentHint":  true,
	}})
	if dst.Annotations["readOnlyHint"] != true ||
		dst.Annotations["destructiveHint"] != false ||
		dst.Annotations["idempotentHint"] != true {
		t.Fatalf("合并后 dst annotations=%+v", dst.Annotations)
	}
}

// mergeConfig 对 nil src.Annotations 安全：dst 为 nil 时不应初始化出空 map。
func TestMergeConfigAnnotationsNilSrc(t *testing.T) {
	dst := &Config{}            // Annotations nil
	mergeConfig(dst, &Config{}) // src.Annotations nil
	if dst.Annotations != nil {
		t.Fatalf("src 为 nil 时 dst.Annotations 应保持 nil，实际 %+v", dst.Annotations)
	}
}

// ---- handleToolsList：annotations 输出分支 ----

func TestHandleToolsListAnnotationsPresent(t *testing.T) {
	resetTools()
	registerTool("/vservice/q", addFn, Desc("查询"), ReadOnly())

	list := handleToolsList()["tools"].([]map[string]interface{})
	if len(list) == 0 {
		t.Fatalf("tools 为空，登记失败")
	}
	ann, ok := list[0]["annotations"].(map[string]bool)
	if !ok {
		t.Fatalf("annotations 缺失或类型错误: %+v", list[0]["annotations"])
	}
	if ann["readOnlyHint"] != true {
		t.Errorf("readOnlyHint=%v want true", ann["readOnlyHint"])
	}
}

func TestHandleToolsListAnnotationsAbsent(t *testing.T) {
	resetTools()
	registerTool("/vservice/q", addFn, "查询") // 未声明 annotations

	list := handleToolsList()["tools"].([]map[string]interface{})
	if len(list) == 0 {
		t.Fatalf("tools 为空，登记失败")
	}
	if _, ok := list[0]["annotations"]; ok {
		t.Fatalf("未声明时不应输出 annotations key: %+v", list[0])
	}
}

// ---- WithAll/WithName/WithDesc builder ----

func TestWithAll(t *testing.T) {
	c := &Config{}
	WithAll("test_name", "test_desc", ReadOnly())(c)
	if c.Name != "test_name" {
		t.Errorf("Name want test_name got %s", c.Name)
	}
	if c.Desc != "test_desc" {
		t.Errorf("Desc want test_desc got %s", c.Desc)
	}
	if c.Annotations["readOnlyHint"] != true {
		t.Errorf("readOnlyHint want true got %v", c.Annotations["readOnlyHint"])
	}
}

func TestWithAllEmptyNameAndDesc(t *testing.T) {
	c := &Config{Name: "old_name", Desc: "old_desc"}
	WithAll("", "", ReadOnly())(c)
	if c.Name != "old_name" {
		t.Errorf("空 name 不应覆盖旧值: got %s", c.Name)
	}
	if c.Desc != "old_desc" {
		t.Errorf("空 desc 不应覆盖旧值: got %s", c.Desc)
	}
}

func TestWithAllNilExtra(t *testing.T) {
	c := &Config{}
	WithAll("n", "d", nil, ReadOnly(), nil)(c) // 忽略 nil extra
	if c.Name != "n" || c.Desc != "d" {
		t.Fatalf("WithAll(n,d,nil) Name/Desc 应正常")
	}
	if c.Annotations["readOnlyHint"] != true {
		t.Fatalf("WithAll extra 里的 ReadOnly() 应生效")
	}
}

func TestWithName(t *testing.T) {
	c := &Config{Name: "old_name"}
	WithName("new_name")(c)
	if c.Name != "new_name" {
		t.Errorf("WithName want new_name got %s", c.Name)
	}
	if c.Desc != "" {
		t.Errorf("WithName 不应碰 Desc")
	}
}

func TestWithDesc(t *testing.T) {
	c := &Config{Desc: "old_desc"}
	WithDesc("new_desc", Destructive())(c)
	if c.Desc != "new_desc" {
		t.Errorf("WithDesc want new_desc got %s", c.Desc)
	}
	if c.Annotations["destructiveHint"] != true {
		t.Errorf("WithDesc extra Destructive() 应生效")
	}
}

func TestWithAllIntegration(t *testing.T) {
	resetTools()
	// 通过 registerTool 集成，验证 buildConfig 支持 WithAll
	registerTool("/test/echo", addFn, WithAll("test_echo", "测试 echo", ReadOnly()))
	list := handleToolsList()["tools"].([]map[string]interface{})
	if len(list) == 0 {
		t.Fatalf("tools 为空，登记失败")
	}
	tool := list[0]
	if tool["name"] != "test_echo" {
		t.Errorf("name want test_echo got %v", tool["name"])
	}
	if tool["description"] != "测试 echo" {
		t.Errorf("description want 测试 echo got %v", tool["description"])
	}
}
