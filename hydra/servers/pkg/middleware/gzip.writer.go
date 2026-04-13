package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

type respWriter interface {
	io.Writer
	Size() int
	Header() http.Header
	WriteHeader(statusCode int)
}

type countWriter struct {
	writer io.Writer
	count  int
}

func (c *countWriter) Write(p []byte) (int, error) {
	n, err := c.writer.Write(p)
	c.count += n
	return n, err
}

func (c *countWriter) Count() int {
	return c.count
}

type gzipWriter struct {
	respWriter
	gzPool       sync.Pool
	ctx          IMiddleContext
	cwriter      interface{}
	isgzip       bool
	needCompress bool
	countWriter  *countWriter  // 用于追踪实际写入的字节数
	closed       bool          // 追踪是否已关闭
	mu           sync.Mutex    // 保护并发访问
}

func newGzipWriter(w respWriter, ctx IMiddleContext, level int) *gzipWriter {
	writer := &gzipWriter{
		respWriter:   w,
		ctx:          ctx,
		needCompress: shouldCompress(ctx),
	}
	writer.gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	return writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.Write([]byte(s))
}
func (g *gzipWriter) getWriter(l int) io.Writer {
	if g.cwriter != nil {
		return g.cwriter.(io.Writer)
	}
	if !g.needCompress {
		g.cwriter = g.respWriter
		return g.respWriter
	}
	g.isgzip = true
	g.ctx.Response().Header("Content-Encoding", "gzip")
	g.ctx.Response().Header("Vary", "Accept-Encoding")

	g.ctx.Response().AddSpecial("gzip")

	// 使用 countWriter 包装底层 writer，用于追踪实际写入的字节数
	g.countWriter = &countWriter{writer: g.respWriter.(io.Writer)}
	gw := g.gzPool.Get().(*gzip.Writer)
	gw.Reset(g.countWriter)
	g.cwriter = gw
	return gw

}
func (g *gzipWriter) Write(data []byte) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// 如果 gzip writer 已关闭，直接写入底层 writer
	if g.closed {
		return g.respWriter.Write(data)
	}

	writer := g.getWriter(len(data))
	n, err := writer.Write(data)

	// 立即刷新以确保数据被发送
	if g.isgzip && g.cwriter != nil {
		if gz, ok := g.cwriter.(*gzip.Writer); ok {
			gz.Flush()
		}
	}

	return n, err
}

func (g *gzipWriter) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.closed || !g.isgzip {
		g.closed = true
		return
	}

	writer := g.cwriter.(*gzip.Writer)
	writer.Close()
	g.ctx.Response().AddSpecial(fmt.Sprint(g.respWriter.Size()))
	writer.Reset(ioutil.Discard)
	g.gzPool.Put(writer)

	g.closed = true
}

// WrittenSize 返回实际写入的字节数（压缩后的大小）
func (g *gzipWriter) WrittenSize() int {
	if g.countWriter != nil {
		return g.countWriter.Count()
	}
	// 如果没有使用压缩，返回原始大小
	return g.respWriter.Size()
}

func shouldCompress(ctx IMiddleContext) bool {
	if !strings.Contains(ctx.Request().Headers().GetString("Accept-Encoding"), "gzip") ||
		strings.Contains(ctx.Request().Headers().GetString("Connection"), "Upgrade") ||
		strings.Contains(ctx.Request().Headers().GetString("Content-Type"), "text/event-stream") ||
		ctx.Response().HasSpecial("gz") {
		return false
	}

	extension := filepath.Ext(ctx.Request().Path().GetURL().Path)
	for _, ext := range DefaultExcludedExtentions {
		if strings.EqualFold(ext, extension) {
			return false
		}
	}
	return true
}
