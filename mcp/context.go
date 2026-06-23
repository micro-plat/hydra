package mcp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/micro-plat/hydra/context"
)

// mcpCtx 实现 context.IInnerContext，作为 tools/call 执行业务时的上下文适配层。
//
// 设计要点（对照 hydra-mcp-implementation.md §6，参照 middle.gin.go 的 ginCtx 模板）：
//   - 请求侧：body 注入 tool 的 arguments；header/cookie/IP/原生 HTTP 对象从源握手 ctx 捕获，
//     保证业务侧 ctx.Request().Headers()、ctx.User().GetClientIP() 等与 HTTP 调用同源。
//   - 响应侧：Data/WStatus 只捕获不落 socket。typed Handle 恒走返回值（不经响应写入），
//     此处仅为满足 IInnerContext 契约；handler.go 以 tool.Handler 返回值为准。
//   - 鉴权：不经此处理。由 handler.go 在 NewCtx 后做 Auth().Request() 迁移（复用 JwtAuth 登录态）。
//
// 全部请求侧数据经 context.IContext 公共 API 捕获，无需访问 *ctx.Ctx 的未导出 inner 字段，
// 因此不改动核心 context 包（开闭原则）。
type mcpCtx struct {
	// ---- 请求侧 ----
	body     []byte             // tools/call 的 arguments 序列化，GetBody 来源
	path     string             // tool 源 HTTP 路径，GetURL/GetRouterPath 来源
	clientIP string             // 源握手客户端 IP
	headers  http.Header        // 源握手 header
	cookies  []*http.Cookie     // 源握手 cookie
	httpReq  *http.Request      // 源 HTTP 请求（GetHTTPReqResp）
	httpResp http.ResponseWriter // 源 HTTP 响应（GetHTTPReqResp）

	// ---- 响应侧（捕获）----
	status      int
	written     bool
	wheaders    http.Header
	respContent []byte
	respCType   string
}

// 编译期断言：mcpCtx 必须实现 IInnerContext 的全部 24 个方法。
var _ context.IInnerContext = (*mcpCtx)(nil)

// newMCPCtx 由源握手 ctx 捕获请求侧属性，构造 tools/call 的执行上下文。
// body 为 tool 的 arguments 序列化字节；path 为 tool 源 HTTP 路径。
func newMCPCtx(src context.IContext, body []byte, path string) *mcpCtx {
	c := &mcpCtx{
		body:     body,
		path:     path,
		clientIP: src.User().GetClientIP(),
		headers:  http.Header{},
		httpReq:  src.Request().GetHTTPRequest(),
		httpResp: src.Response().GetHTTPReponse(),
	}
	for k, v := range src.Request().Headers() {
		c.headers.Set(k, fmt.Sprint(v))
	}
	for k, v := range src.Request().Cookies() {
		c.cookies = append(c.cookies, &http.Cookie{Name: k, Value: fmt.Sprint(v)})
	}
	return c
}

// ---------------- 请求侧 ----------------

// GetBody 返回 arguments 序列化（每次新 reader，支持重复读）。
func (c *mcpCtx) GetBody() io.ReadCloser {
	return io.NopCloser(bytes.NewReader(c.body))
}

// GetMethod 固定 POST（JSON-RPC over HTTP 约定）。
func (c *mcpCtx) GetMethod() string { return http.MethodPost }

// GetURL 返回 tool 源 HTTP 路径。
func (c *mcpCtx) GetURL() *url.URL { return &url.URL{Path: c.path} }

// ContentType 固定 application/json。
func (c *mcpCtx) ContentType() string { return "application/json" }

// GetParams MCP 不使用路径参数。
func (c *mcpCtx) GetParams() map[string]interface{} { return map[string]interface{}{} }

// GetPostForm MCP 不使用表单。
func (c *mcpCtx) GetPostForm() url.Values { return url.Values{} }

// GetRawForm MCP 不使用表单。
func (c *mcpCtx) GetRawForm() map[string]interface{} { return map[string]interface{}{} }

// GetHeaders 返回源握手 header。
func (c *mcpCtx) GetHeaders() http.Header { return c.headers }

// GetCookies 返回源握手 cookie。
func (c *mcpCtx) GetCookies() []*http.Cookie { return c.cookies }

// ClientIP 返回源握手客户端 IP。
func (c *mcpCtx) ClientIP() string { return c.clientIP }

// GetHTTPReqResp 返回源 HTTP 请求/响应对象。
func (c *mcpCtx) GetHTTPReqResp() (*http.Request, http.ResponseWriter) {
	return c.httpReq, c.httpResp
}

// ClearAuth 空实现：tool 执行不再经 JwtAuth，鉴权由 handler.go 迁移登录态完成。
func (c *mcpCtx) ClearAuth(...bool) bool { return false }

// ---------------- 响应侧（捕获，不落 socket）----------------

// Header 捕获响应头设置。
func (c *mcpCtx) Header(k, v string) {
	if c.wheaders == nil {
		c.wheaders = http.Header{}
	}
	c.wheaders.Set(k, v)
}

// Abort 空实现（typed Handle 不通过响应流控制流程）。
func (c *mcpCtx) Abort() {}

// WStatus 捕获状态码。
func (c *mcpCtx) WStatus(s int) { c.status = s }

// Status 返回捕获的状态码。
func (c *mcpCtx) Status() int { return c.status }

// Written 是否已写入响应。
func (c *mcpCtx) Written() bool { return c.written }

// WHeaders 返回捕获的响应头。
func (c *mcpCtx) WHeaders() http.Header {
	if c.wheaders == nil {
		c.wheaders = http.Header{}
	}
	return c.wheaders
}

// WHeader 返回捕获的某个响应头。
func (c *mcpCtx) WHeader(k string) string { return c.WHeaders().Get(k) }

// Data 捕获响应内容（不落 socket）。
func (c *mcpCtx) Data(status int, contentType string, b []byte) {
	c.status = status
	c.respCType = contentType
	c.respContent = b
	c.written = true
}

// ServeContent MCP 不支持文件输出，返回 200 占位。
func (c *mcpCtx) ServeContent(filepath string, fs http.FileSystem) int { return http.StatusOK }

// Redirect 空实现。
func (c *mcpCtx) Redirect(int, string) {}

// GetRouterPath 返回 tool 源 HTTP 路径。
func (c *mcpCtx) GetRouterPath() string { return c.path }

// GetFile MCP 不支持文件上传。
func (c *mcpCtx) GetFile(fileKey string) (string, io.ReadCloser, int64, error) {
	return "", nil, 0, fmt.Errorf("mcp: file upload not supported")
}
