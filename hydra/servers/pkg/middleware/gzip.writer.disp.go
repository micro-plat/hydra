package middleware

import (
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
	return g.gzip.WriteString(s)
}

func (g *dispWriter) Write(data []byte) (int, error) {
	return g.gzip.Write(data)
}
func (g *dispWriter) Close() {
	g.gzip.Close()
}
