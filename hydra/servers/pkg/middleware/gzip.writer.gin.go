package middleware

import (
	"github.com/gin-gonic/gin"
)

type ginWriter struct {
	gin.ResponseWriter
	gzip *gzipWriter
}

func newGinWriter(w gin.ResponseWriter, ctx IMiddleContext, level int) *ginWriter {
	return &ginWriter{
		ResponseWriter: w,
		gzip:           newGzipWriter(w, ctx, level),
	}
}
func (g *ginWriter) WriteString(s string) (int, error) {
	// 仅真走 gzip 压缩时才需删 Content-Length（压缩后长度不可预知）；
	// 未压缩透传时保留业务 handler 设置的 Content-Length（文件下载走 identity 定长帧，
	// 避免部分中间层无法转发 chunked 大响应导致浏览器 net::ERR_FAILED）
	if g.gzip.needCompress {
		g.ResponseWriter.Header().Del("Content-Length")
	}
	// 返回压缩前的长度，保持与原始行为一致
	return g.gzip.WriteString(s)
}

func (g *ginWriter) Write(data []byte) (int, error) {
	if g.gzip.needCompress {
		g.ResponseWriter.Header().Del("Content-Length")
	}
	// 返回压缩前的长度，保持与原始行为一致
	return g.gzip.Write(data)
}

func (g *ginWriter) WriteHeader(code int) {
	if g.gzip.needCompress {
		g.ResponseWriter.Header().Del("Content-Length")
	}
	g.ResponseWriter.WriteHeader(code)
}

// Size 覆盖原始的 Size() 方法
// 如果启用了 gzip，返回实际写入的字节数（压缩后的大小）
// 否则返回原始大小
func (g *ginWriter) Size() int {
	if g.gzip.WrittenSize() > 0 {
		return g.gzip.WrittenSize()
	}
	return g.ResponseWriter.Size()
}

// Close 关闭 gzip writer 并完成响应
// 修复说明：gzip 压缩后实际大小与原始大小不同，无法在压缩前预知 Content-Length。
// 修复方案：删除 Content-Length，让 HTTP 使用 Transfer-Encoding: chunked 自动处理。
func (g *ginWriter) Close() {
	if g.gzip.needCompress {
		g.ResponseWriter.Header().Del("Content-Length")
	}
	g.gzip.Close()
}
