package aimiddleware

import (
	"io"
	"net/http"
)

// RawJSON 表示已经按上游协议编码好的JSON响应。
type RawJSON struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// SSEProxy 表示OpenAI-compatible SSE透传响应。
type SSEProxy struct {
	StatusCode int
	Header     http.Header
	Body       io.Reader
}

// NewRawJSON 构建RawJSON响应。
func NewRawJSON(statusCode int, body []byte, headers ...http.Header) RawJSON {
	header := http.Header{}
	if len(headers) > 0 && headers[0] != nil {
		header = headers[0]
	}
	return RawJSON{StatusCode: statusCode, Header: header, Body: body}
}

// NewSSEProxy 构建SSE透传响应。
func NewSSEProxy(body io.Reader, statusCode ...int) SSEProxy {
	code := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		code = statusCode[0]
	}
	return SSEProxy{StatusCode: code, Body: body}
}
