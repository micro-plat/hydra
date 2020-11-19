package servers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/micro-plat/hydra/global"
	shttp "github.com/micro-plat/hydra/hydra/servers/http"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestHandler_GinFunc(t *testing.T) {
	var h middleware.Handler
	h = func(middleware.IMiddleContext) { return }
	tps := []string{global.API}

	confObj := mocks.NewConfBy("middleware_main_test", "middle") //构建对象
	serverConf := confObj.GetAPIConf()                           //获取配置
	_, _ = shttp.NewResponsive(serverConf)

	r, err := http.NewRequest("POST", "http://localhost:9091/getbody", strings.NewReader(""))
	assert.Equal(t, nil, err, "构建请求")

	//设置content-type
	r.Header.Set("X-Request-Id", "123456")
	//替换gin上下文的请求
	c := &gin.Context{}
	c.Request = r
	router := gin.New()
	router.HandleContext(c)
	//ginFunc执行
	f := h.GinFunc(tps...)
	f(c)

	got, ok := c.Get("__middle_context__")

	assert.Equal(t, true, got != nil, "获取中间件上下文")
	assert.Equal(t, true, ok, "获取中间件上下文")
}
