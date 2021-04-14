package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/hydra/servers/pkg/dispatcher"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

func Gzip(level int, options ...Option) Handler {
	handler := newGzipHandler(level, options...).Handle
	return func(ctx IMiddleContext) {
		handler(ctx)
	}
}

type gzipGinWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipGinWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write([]byte(s))
}

func (g *gzipGinWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *gzipGinWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

type gzipDispatcherWriter struct {
	dispatcher.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipDispatcherWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write([]byte(s))
}

func (g *gzipDispatcherWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *gzipDispatcherWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

var (
	DefaultExcludedExtentions = NewExcludedExtensions([]string{
		".png", ".gif", ".jpeg", ".jpg",
	})
	DefaultOptions = &GzipOptions{
		ExcludedExtensions: DefaultExcludedExtentions,
	}
)

type GzipOptions struct {
	ExcludedExtensions   ExcludedExtensions
	ExcludedPaths        ExcludedPaths
	ExcludedPathesRegexs ExcludedPathesRegexs
	DecompressFn         func(c IMiddleContext)
}

type Option func(*GzipOptions)

func WithExcludedExtensions(args []string) Option {
	return func(o *GzipOptions) {
		o.ExcludedExtensions = NewExcludedExtensions(args)
	}
}

func WithExcludedPaths(args []string) Option {
	return func(o *GzipOptions) {
		o.ExcludedPaths = NewExcludedPaths(args)
	}
}

func WithExcludedPathsRegexs(args []string) Option {
	return func(o *GzipOptions) {
		o.ExcludedPathesRegexs = NewExcludedPathesRegexs(args)
	}
}

func WithDecompressFn(decompressFn func(c IMiddleContext)) Option {
	return func(o *GzipOptions) {
		o.DecompressFn = decompressFn
	}
}

// Using map for better lookup performance
type ExcludedExtensions map[string]bool

func NewExcludedExtensions(extensions []string) ExcludedExtensions {
	res := make(ExcludedExtensions)
	for _, e := range extensions {
		res[e] = true
	}
	return res
}

func (e ExcludedExtensions) Contains(target string) bool {
	_, ok := e[target]
	return ok
}

type ExcludedPaths []string

func NewExcludedPaths(paths []string) ExcludedPaths {
	return ExcludedPaths(paths)
}

func (e ExcludedPaths) Contains(requestURI string) bool {
	for _, path := range e {
		if strings.HasPrefix(requestURI, path) {
			return true
		}
	}
	return false
}

type ExcludedPathesRegexs []*regexp.Regexp

func NewExcludedPathesRegexs(regexs []string) ExcludedPathesRegexs {
	result := make([]*regexp.Regexp, len(regexs), len(regexs))
	for i, reg := range regexs {
		result[i] = regexp.MustCompile(reg)
	}
	return result
}

func (e ExcludedPathesRegexs) Contains(requestURI string) bool {
	for _, reg := range e {
		if reg.MatchString(requestURI) {
			return true
		}
	}
	return false
}

func DefaultDecompressHandle(c *gin.Context) {
	if c.Request.Body == nil {
		return
	}
	r, err := gzip.NewReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Request.Header.Del("Content-Encoding")
	c.Request.Header.Del("Content-Length")
	c.Request.Body = r
}

type gzipHandler struct {
	*GzipOptions
	gzPool sync.Pool
}

func newGzipHandler(level int, options ...Option) *gzipHandler {
	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	handler := &gzipHandler{
		GzipOptions: DefaultOptions,
		gzPool:      gzPool,
	}
	for _, setter := range options {
		setter(handler.GzipOptions)
	}
	return handler
}

func (g *gzipHandler) Handle(ctx IMiddleContext) {
	if fn := g.DecompressFn; fn != nil && ctx.Request().Headers().GetString("Content-Encoding") == "gzip" {
		fn(ctx)
	}

	if !g.shouldCompress(ctx) {
		return
	}

	gz := g.gzPool.Get().(*gzip.Writer)
	defer g.gzPool.Put(gz)
	defer gz.Reset(ioutil.Discard)
	gz.Reset(ctx.GetWriter().(io.Writer))

	ctx.Response().Header("Content-Encoding", "gzip")
	ctx.Response().Header("Vary", "Accept-Encoding")

	ctx.Response().AddSpecial("gzip")
	switch strings.ToLower(ctx.GetType()) {
	case "gin":
		writer := ctx.GetWriter().(gin.ResponseWriter)
		ctx.SetWriter(&gzipGinWriter{writer, gz})
		defer func() {
			gz.Close()
			ctx.Response().Header("Content-Length", fmt.Sprint(writer.Size()))
		}()
	default:
		writer := ctx.GetWriter().(dispatcher.ResponseWriter)
		ctx.SetWriter(&gzipDispatcherWriter{writer, gz})
		defer func() {
			gz.Close()
			ctx.Response().Header("Content-Length", fmt.Sprint(writer.Size()))
		}()
	}

	ctx.Next()
}

func (g *gzipHandler) shouldCompress(ctx IMiddleContext) bool {
	if !strings.Contains(ctx.Request().Headers().GetString("Accept-Encoding"), "gzip") ||
		strings.Contains(ctx.Request().Headers().GetString("Connection"), "Upgrade") ||
		strings.Contains(ctx.Request().Headers().GetString("Content-Type"), "text/event-stream") {

		return false
	}

	extension := filepath.Ext(ctx.Request().Path().GetURL().Path)
	if g.ExcludedExtensions.Contains(extension) {
		return false
	}

	if g.ExcludedPaths.Contains(ctx.Request().Path().GetURL().Path) {
		return false
	}
	if g.ExcludedPathesRegexs.Contains(ctx.Request().Path().GetURL().Path) {
		return false
	}

	return true
}
