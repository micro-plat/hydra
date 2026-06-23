package mcp

import (
	"encoding/json"

	"github.com/micro-plat/hydra/context"
	ctx "github.com/micro-plat/hydra/context/ctx"
)

// JSON-RPC 2.0 错误码（规范定义）。
const (
	errParseError     = -32700 // 解析错误
	errInvalidReq     = -32600 // 无效请求
	errMethodNotFound = -32601 // 方法不存在
	errInvalidParams  = -32602 // 无效参数
	errInternal       = -32603 // 内部错误
)

// 支持的 MCP 协议版本。
const protocolVersion = "2024-11-05"

// jsonrpcRequest JSON-RPC 2.0 请求信封。Params 保留原始 JSON，由各方法按需解析。
type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// jsonrpcResponse JSON-RPC 2.0 响应信封。Result 与 Error 互斥。
type jsonrpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

// jsonrpcError JSON-RPC 错误对象。
type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// toolsCallParams tools/call 的参数。
type toolsCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// NewJSONRPCHandler 返回 /mcp 端点的业务处理函数，处理 MCP JSON-RPC 2.0。
//
// 注册方式（应用层）：
//
//	hydra.S.Micro("/mcp", mcp.NewJSONRPCHandler())
//
// 挂在 web/api server 上，自动复用完整中间件链（含 JwtAuth）；handler 直接返回
// JSON-RPC 响应对象，由框架 WriteAny 序列化为 application/json（无信封包装）。
func NewJSONRPCHandler() func(context.IContext) interface{} {
	return JSONRPCHandle
}

// JSONRPCHandle /mcp JSON-RPC 分派：initialize / tools/list / tools/call / ping。
func JSONRPCHandle(ctx context.IContext) interface{} {
	raw := ctx.Request().RawBody()

	var req jsonrpcRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return rpcErr(nil, errParseError, "parse error: "+err.Error())
	}
	if req.JSONRPC != "2.0" {
		return rpcErr(req.ID, errInvalidReq, "invalid request: jsonrpc must be 2.0")
	}

	switch req.Method {
	case "initialize":
		return rpcResult(req.ID, handleInitialize())
	case "tools/list":
		return rpcResult(req.ID, handleToolsList())
	case "tools/call":
		return handleToolsCall(ctx, req.ID, req.Params)
	case "ping":
		return rpcResult(req.ID, map[string]interface{}{})
	default:
		return rpcErr(req.ID, errMethodNotFound, "method not found: "+req.Method)
	}
}

// handleInitialize 返回 server 能力声明（仅声明 tools 能力）。
func handleInitialize() map[string]interface{} {
	return map[string]interface{}{
		"protocolVersion": protocolVersion,
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "hydra-mcp",
			"version": "1.0.0",
		},
	}
}

// handleToolsList 返回全部已注册 tool 的名+描述+输入 schema（不执行业务）。
func handleToolsList() map[string]interface{} {
	list := List()
	tools := make([]map[string]interface{}, 0, len(list))
	for _, t := range list {
		item := map[string]interface{}{
			"name":        t.Name,
			"description": t.Config.Desc,
			"inputSchema": t.Config.InputSchema,
		}
		if len(t.Config.Annotations) > 0 {
			item["annotations"] = t.Config.Annotations
		}
		tools = append(tools, item)
	}
	return map[string]interface{}{"tools": tools}
}

// handleToolsCall 执行指定 tool：构造 mcpCtx（注入 arguments）→ NewCtx → 迁移登录态 → 调用。
//
// 兼容性要点：
//   - 鉴权迁移 toolCtx.User().Auth().Request(srcCtx...) 复用 JwtAuth 在握手阶段写入的登录态，
//     使需登录业务经 MCP 调用与 HTTP 调用获得一致的用户态。
//   - tool.Handler 为 typed 包装，恒返回 resp/error，无需经响应流回退。
//   - 不调用 toolCtx.Close()：其 context.Del() 会按 X-Request-Id 移除源握手 ctx 的缓存槽位；
//     toolCtx 因 Cache 的 SetIfAbsent 未实际入缓存，cancelFunc 定时器在 timeout 后自释放。
func handleToolsCall(srcCtx context.IContext, id json.RawMessage, params json.RawMessage) *jsonrpcResponse {
	var p toolsCallParams
	if len(params) == 0 {
		return rpcErr(id, errInvalidParams, "missing params")
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return rpcErr(id, errInvalidParams, "invalid params: "+err.Error())
	}

	tool, ok := Get(p.Name)
	if !ok {
		return rpcErr(id, errInvalidParams, "unknown tool: "+p.Name)
	}

	// 鉴权由本服务器 JwtAuth 中间件在 /mcp 入口统一强制（/mcp 不豁免 JWT，MCP 客户端须登录）；
	// 此处不再做 per-tool 匿名判定，直接执行。

	// 构造 tool 执行上下文：body 注入 arguments，path 取源 HTTP 路径。
	argsBytes, _ := json.Marshal(p.Arguments)
	mc := newMCPCtx(srcCtx, argsBytes, tool.Path)
	toolCtx := ctx.NewCtx(mc, srcCtx.APPConf().GetServerConf().GetServerType())

	// 鉴权迁移：复用握手阶段的登录态，业务侧 Auth().Bind 零改动即可拿到当前用户。
	toolCtx.User().Auth().Request(srcCtx.User().Auth().Request())

	result := tool.Handler(toolCtx)
	return wrapToolResult(id, result)
}

// wrapToolResult 把 tool 执行返回值包成 tools/call 的 result（content/isError）。
// tool 业务错误用 isError=true 表达（仍是 JSON-RPC 成功响应）；仅协议级错误用 jsonrpcError。
func wrapToolResult(id json.RawMessage, result interface{}) *jsonrpcResponse {
	text := ""
	if result != nil {
		if err, ok := result.(error); ok && err != nil {
			return rpcResult(id, map[string]interface{}{
				"content": []map[string]interface{}{{"type": "text", "text": err.Error()}},
				"isError": true,
			})
		}
		b, err := json.Marshal(result)
		if err != nil {
			return rpcErr(id, errInternal, "marshal result: "+err.Error())
		}
		text = string(b)
	}
	return rpcResult(id, map[string]interface{}{
		"content": []map[string]interface{}{{"type": "text", "text": text}},
		"isError": false,
	})
}

// rpcResult 构造成功响应。
func rpcResult(id json.RawMessage, result interface{}) *jsonrpcResponse {
	return &jsonrpcResponse{JSONRPC: "2.0", ID: normID(id), Result: result}
}

// rpcErr 构造错误响应。
func rpcErr(id json.RawMessage, code int, msg string) *jsonrpcResponse {
	return &jsonrpcResponse{JSONRPC: "2.0", ID: normID(id), Error: &jsonrpcError{Code: code, Message: msg}}
}

// normID 规范化 id：空则 null，保证响应可序列化。
func normID(id json.RawMessage) json.RawMessage {
	if len(id) == 0 {
		return json.RawMessage("null")
	}
	return id
}
