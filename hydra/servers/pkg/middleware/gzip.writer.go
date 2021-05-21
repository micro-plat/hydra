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

type gzipWriter struct {
	respWriter
	gzPool  sync.Pool
	ctx     IMiddleContext
	cwriter interface{}
	isgzip  bool
}

func newGzipWriter(w respWriter, ctx IMiddleContext, level int) *gzipWriter {
	writer := &gzipWriter{
		respWriter: w,
		ctx:        ctx,
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
	if !shouldCompress(g.ctx) || l == 0 {
		g.cwriter = g.respWriter
		return g.respWriter
	}
	g.isgzip = true
	g.ctx.Response().Header("Content-Encoding", "gzip")
	g.ctx.Response().Header("Vary", "Accept-Encoding")

	g.ctx.Response().AddSpecial("gzip")

	gw := g.gzPool.Get().(*gzip.Writer)
	gw.Reset(g.respWriter.(io.Writer))
	g.cwriter = gw
	return gw

}
func (g *gzipWriter) Write(data []byte) (int, error) {
	writer := g.getWriter(len(data))
	s, err := writer.Write(data)
	if err != nil {
		return s, err
	}
	return s, nil
}
func (g *gzipWriter) Close() {
	if !g.isgzip {
		return
	}
	writer := g.cwriter.(*gzip.Writer)
	writer.Close()
	g.ctx.Response().Header("Content-Length", "")
	g.ctx.Response().Header("Content-Length", fmt.Sprint(g.respWriter.Size()))
	g.ctx.Response().AddSpecial(fmt.Sprint(g.respWriter.Size()))
	writer.Reset(ioutil.Discard)
	g.gzPool.Put(writer)
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
