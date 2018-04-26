package component

import (
	"testing"

	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/lib4go/ut"
)

//CPUUPCollect CPU使用率收集
func testHandler() StandardServiceFunc {
	return func(name string, mode string, service string, ctx *context.Context) (response *context.StandardResponse, err error) {
		response = context.GetStandardResponse()
		return
	}
}

func Test_AddMicroService(t *testing.T) {
	c := NewStandardComponent("t", nil)
	c.AddMicroService("/order/request", testHandler)
	ut.Expect(t, len(c.funcs), 1)
	err := c.LoadServices()
	ut.Expect(t, err, nil)
	ut.Expect(t, len(c.funcs), 0)
	ut.Expect(t, len(c.Handlers), 1)
	ut.Expect(t, len(c.FallbackHandlers), 0)
	ut.Expect(t, len(c.FallbackHandlers), 0)
	ut.Expect(t, len(c.GroupServices), 1)
	ut.Expect(t, len(c.ServiceGroup), 1)
	ut.Expect(t, len(c.Services), 1)
	ut.Expect(t, len(c.ServicePages), 0)
	ut.Expect(t, len(c.ServiceTagPages), 0)
	ut.Expect(t, len(c.TagServices), 0)

}
func Test_AddPageService(t *testing.T) {
	c := NewStandardComponent("t", nil)
	c.AddPageService("/order/request", testHandler)
	ut.Expect(t, len(c.funcs), 1)
	err := c.LoadServices()
	ut.Expect(t, err, nil)
	ut.Expect(t, len(c.funcs), 0)
	ut.Expect(t, len(c.Handlers), 1)
	ut.Expect(t, len(c.FallbackHandlers), 0)
	ut.Expect(t, len(c.FallbackHandlers), 0)
	ut.Expect(t, len(c.GroupServices), 1)
	ut.Expect(t, len(c.ServiceGroup), 1)
	ut.Expect(t, len(c.Services), 1)
	ut.Expect(t, len(c.ServicePages), 1)
	ut.Expect(t, len(c.ServiceTagPages), 0)
	ut.Expect(t, len(c.TagServices), 0)

}

func Test_MixService(t *testing.T) {
	c := NewStandardComponent("t", nil)
	c.AddPageService("/order/request", testHandler)
	c.AddMicroService("/order/request", testHandler)
	ut.Expect(t, len(c.funcs), 2)
	err := c.LoadServices()
	ut.Expect(t, err, nil)
	ut.Expect(t, len(c.funcs), 0)
	ut.Expect(t, len(c.Handlers), 1)
	ut.Expect(t, len(c.FallbackHandlers), 0)
	ut.Expect(t, len(c.FallbackHandlers), 0)
	ut.Expect(t, len(c.GroupServices), 2)
	ut.Expect(t, len(c.ServiceGroup), 1)
	ut.Expect(t, len(c.ServiceGroup["/order/request"]), 2)
	ut.Expect(t, len(c.Services), 1)
	ut.Expect(t, len(c.ServicePages), 1)
	ut.Expect(t, len(c.ServiceTagPages), 0)
	ut.Expect(t, len(c.TagServices), 0)

}
