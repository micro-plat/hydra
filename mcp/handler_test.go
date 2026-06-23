package mcp

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/micro-plat/hydra/mock"
)

// rpcReq 构造 JSON-RPC 2.0 请求字符串（测试 fixture）。
func rpcReq(id interface{}, method string, params interface{}) string {
	m := map[string]interface{}{"jsonrpc": "2.0", "id": id, "method": method}
	if params != nil {
		m["params"] = params
	}
	b, _ := json.Marshal(m)
	return string(b)
}

// ---- 纯函数 ----

func TestHandleInitialize(t *testing.T) {
	r := handleInitialize()
	if r["protocolVersion"] != protocolVersion {
		t.Fatalf("protocolVersion=%v want %s", r["protocolVersion"], protocolVersion)
	}
	caps, ok := r["capabilities"].(map[string]interface{})
	if !ok {
		t.Fatalf("capabilities not a map: %T", r["capabilities"])
	}
	if _, ok := caps["tools"]; !ok {
		t.Fatal("missing tools capability")
	}
	si, ok := r["serverInfo"].(map[string]interface{})
	if !ok || si["name"] == "" {
		t.Fatalf("serverInfo invalid: %+v", r["serverInfo"])
	}
}

func TestHandleToolsListPure(t *testing.T) {
	resetTools()
	registerTool("/v1/add", addFn, "加法")

	r := handleToolsList()
	list := r["tools"].([]map[string]interface{})
	if len(list) != 1 {
		t.Fatalf("tools len=%d want 1", len(list))
	}
	item := list[0]
	if item["name"] != "v1_add" {
		t.Errorf("name=%v", item["name"])
	}
	if item["description"] != "加法" {
		t.Errorf("desc=%v", item["description"])
	}
	if item["inputSchema"] == nil {
		t.Error("inputSchema missing")
	}
}

func TestWrapToolResult(t *testing.T) {
	id := json.RawMessage("1")

	// 成功：struct 序列化为 text
	r := wrapToolResult(id, addResp{Sum: 7})
	if r.Error != nil {
		t.Fatalf("unexpected error: %+v", r.Error)
	}
	res := r.Result.(map[string]interface{})
	if res["isError"] != false {
		t.Errorf("isError=%v want false", res["isError"])
	}
	text := res["content"].([]map[string]interface{})[0]["text"].(string)
	if text != `{"sum":7}` {
		t.Errorf("text=%q want {\"sum\":7}", text)
	}

	// 错误：isError=true
	r = wrapToolResult(id, errors.New("boom"))
	res = r.Result.(map[string]interface{})
	if res["isError"] != true {
		t.Errorf("isError=%v want true", res["isError"])
	}

	// nil：空 text
	r = wrapToolResult(id, nil)
	res = r.Result.(map[string]interface{})
	if got := res["content"].([]map[string]interface{})[0]["text"]; got != "" {
		t.Errorf("nil text=%v want empty", got)
	}
}

func TestNormID(t *testing.T) {
	if got := normID(nil); string(got) != "null" {
		t.Errorf("nil id=%q want null", got)
	}
	if got := normID(json.RawMessage("0")); string(got) != "0" {
		t.Errorf("id 0=%q", got)
	}
}

// ---- 经 mock 端到端 ----

func TestJSONRPCHandleInitialize(t *testing.T) {
	ctx := mock.NewContext(rpcReq(1, "initialize", nil))
	resp := JSONRPCHandle(ctx).(*jsonrpcResponse)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
	if resp.Result.(map[string]interface{})["protocolVersion"] != protocolVersion {
		t.Fatal("protocolVersion mismatch")
	}
}

func TestJSONRPCHandlePing(t *testing.T) {
	ctx := mock.NewContext(rpcReq(7, "ping", nil))
	resp := JSONRPCHandle(ctx).(*jsonrpcResponse)
	if resp.Error != nil {
		t.Fatalf("ping error: %+v", resp.Error)
	}
	if string(resp.ID) != "7" {
		t.Errorf("id=%q want 7", resp.ID)
	}
}

func TestJSONRPCHandleToolsList(t *testing.T) {
	resetTools()
	registerTool("/v1/add", addFn, "加法")

	ctx := mock.NewContext(rpcReq(1, "tools/list", nil))
	resp := JSONRPCHandle(ctx).(*jsonrpcResponse)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %+v", resp.Error)
	}
	r := resp.Result.(map[string]interface{})
	list := r["tools"].([]map[string]interface{})
	if len(list) != 1 || list[0]["name"] != "v1_add" {
		t.Fatalf("tools/list result: %+v", list)
	}
}

func TestJSONRPCHandleToolsCall(t *testing.T) {
	resetTools()
	registerTool("/v1/add", addFn, "加法")

	body := rpcReq(1, "tools/call", map[string]interface{}{
		"name":      "v1_add",
		"arguments": map[string]interface{}{"a": 5},
	})
	ctx := mock.NewContext(body)
	resp := JSONRPCHandle(ctx).(*jsonrpcResponse)
	if resp.Error != nil {
		t.Fatalf("unexpected rpc error: %+v", resp.Error)
	}

	result := resp.Result.(map[string]interface{})
	if result["isError"] != false {
		t.Fatalf("isError=%v want false (result=%+v)", result["isError"], result)
	}
	content := result["content"].([]map[string]interface{})
	if len(content) != 1 || content[0]["type"] != "text" {
		t.Fatalf("content invalid: %+v", content)
	}
	var got addResp
	if err := json.Unmarshal([]byte(content[0]["text"].(string)), &got); err != nil {
		t.Fatalf("unmarshal text %q: %v", content[0]["text"], err)
	}
	if got.Sum != 10 {
		t.Fatalf("sum=%d want 10 (a=5, addFn doubles)", got.Sum)
	}
}

func TestJSONRPCHandleUnknownTool(t *testing.T) {
	resetTools()
	body := rpcReq(1, "tools/call", map[string]interface{}{
		"name":      "no_such_tool",
		"arguments": map[string]interface{}{},
	})
	ctx := mock.NewContext(body)
	resp := JSONRPCHandle(ctx).(*jsonrpcResponse)
	if resp.Error == nil || resp.Error.Code != errInvalidParams {
		t.Fatalf("want invalid-params error, got: %+v", resp)
	}
}

func TestJSONRPCHandleUnknownMethod(t *testing.T) {
	ctx := mock.NewContext(rpcReq(1, "foo/bar", nil))
	resp := JSONRPCHandle(ctx).(*jsonrpcResponse)
	if resp.Error == nil || resp.Error.Code != errMethodNotFound {
		t.Fatalf("want method-not-found error, got: %+v", resp)
	}
}

func TestJSONRPCHandleParseError(t *testing.T) {
	ctx := mock.NewContext("{not json")
	resp := JSONRPCHandle(ctx).(*jsonrpcResponse)
	if resp.Error == nil || resp.Error.Code != errParseError {
		t.Fatalf("want parse error, got: %+v", resp)
	}
}
