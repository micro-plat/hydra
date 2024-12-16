package middleware

import (
	"fmt"

	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

type dispWriter struct {
	dispatcher.ResponseWriter
	gzip *gzipWriter
}

func newDispWriter(w dispatcher.ResponseWriter, ctx IMiddleContext, level int) *dispWriter {
	return &dispWriter{
		ResponseWriter: w,
		gzip:           newGzipWriter(w, ctx, level),
	}
}
func (g *dispWriter) WriteString(s string) (int, error) {
	g.ResponseWriter.Header().Del("Content-Length")
	return g.gzip.WriteString(s)
}

func (g *dispWriter) Write(data []byte) (int, error) {
	g.ResponseWriter.Header().Del("Content-Length")
	return g.gzip.Write(data)
}
func (g *dispWriter) WriteHeader(code int) {
	g.ResponseWriter.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

func (g *dispWriter) Close() {
	g.ResponseWriter.Header().Del("Content-Length")

	g.gzip.Close()
	g.ResponseWriter.Header().Add("Content-Length", fmt.Sprint(g.ResponseWriter.Size()))
}
