package dispatcher

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/servers/rpc/pb"
	"github.com/micro-plat/lib4go/ut"
)

func TestEngine(t *testing.T) {
	engine := New()
	engine.POST("/:a/:b", func(c *Context) {
		c.Header("Content-Type", "text/plain")
		c.JSON(200, "abc")
	})
	request := &pb.RequestContext{Service: "/order/request", Method: "POST"}
	p, err := engine.HandleRequest(request)
	ut.Expect(t, err, nil)
	ut.Expect(t, p.status, 200)
	ut.Expect(t, string(p.data), `"abc"`)
	ut.Expect(t, len(p.header), 1)
}
func TestEngineError(t *testing.T) {
	engine := New()
	engine.POST("/:a/:b", func(c *Context) {
		c.AbortWithError(500, fmt.Errorf("errr:", ""))
	})
	request := &pb.RequestContext{Service: "/order/request", Method: "POST"}
	p, err := engine.HandleRequest(request)
	ut.Refute(t, err, nil)
	ut.Expect(t, p.status, 500)
}

func TestEngine4041(t *testing.T) {
	engine := New()
	engine.POST("/:a/:b", func(c *Context) {
	})
	request := &pb.RequestContext{Service: "/order/request/abc", Method: "POST"}
	p, err := engine.HandleRequest(request)
	ut.Expect(t, err, nil)
	ut.Expect(t, p.status, 404)
}
func TestEngine4042(t *testing.T) {
	engine := New()
	engine.POST("/:a/:b", func(c *Context) {
	})
	request := &pb.RequestContext{Service: "/order/request/abc", Method: "GET"}
	p, err := engine.HandleRequest(request)
	ut.Expect(t, err, nil)
	ut.Expect(t, p.status, 404)
}
