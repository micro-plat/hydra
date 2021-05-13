package middleware

import "github.com/gin-gonic/gin"

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
	return g.gzip.WriteString(s)
}

func (g *ginWriter) Write(data []byte) (int, error) {
	return g.gzip.Write(data)
}
func (g *ginWriter) Close() {
	g.gzip.Close()
}
