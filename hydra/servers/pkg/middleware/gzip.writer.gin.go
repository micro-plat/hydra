package middleware

import (
	"fmt"

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
	g.ResponseWriter.Header().Del("Content-Length")
	return g.gzip.WriteString(s)
}

func (g *ginWriter) Write(data []byte) (int, error) {
	g.ResponseWriter.Header().Del("Content-Length")
	return g.gzip.Write(data)
}

func (g *ginWriter) WriteHeader(code int) {
	g.ResponseWriter.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

func (g *ginWriter) Close() {
	g.ResponseWriter.Header().Del("Content-Length")
	g.gzip.Close()
	g.ResponseWriter.Header().Add("Content-Length", fmt.Sprint(g.ResponseWriter.Size()))
}
