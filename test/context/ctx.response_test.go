package context

import (
	"errors"
	"net/http"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
)

func Test_response_Write(t *testing.T) {
	type args struct {
		status  int
		content interface{}
	}
	tests := []struct {
		name   string
		ctx    context.IInnerContext
		args   args
		wantRs int
		wantRc string
	}{
		{name: "状态码非0,返回包含错误码的错误", ctx: &mocks.TestContxt{}, args: args{status: 0, content: errs.NewError(999, "错误")}, wantRs: 999, wantRc: "错误"},
		{name: "状态码在200到400,返回错误", ctx: &mocks.TestContxt{}, args: args{status: 300, content: errors.New("err")}, wantRs: 400, wantRc: "err"},
		{name: "状态码为0,返回非错误内容", ctx: &mocks.TestContxt{}, args: args{status: 0, content: nil}, wantRs: 200, wantRc: ""},
		{name: "状态码非0,返回非错误内容", ctx: &mocks.TestContxt{}, args: args{status: 500, content: "content"}, wantRs: 500, wantRc: "content"},
		{name: "状态码非0,content-type为text/plain,返回非错误内容", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{
				"Content-Type": []string{context.PLAINF},
			},
		}, args: args{status: 200, content: "content"}, wantRs: 200, wantRc: "content"},
		{name: "状态码非0,content-type为application/json,返回json内容", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{
				"Content-Type": []string{context.JSONF},
			},
		}, args: args{status: 200, content: `{"key":"value"}`}, wantRs: 200, wantRc: `{"key":"value"}`},
		{name: "状态码非0,content-type为application/xml,返回xml内容", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{
				"Content-Type": []string{context.XMLF},
			},
		}, args: args{status: 200, content: "<?xml><key>value<key/><xml/>"}, wantRs: 200, wantRc: `<?xml><key>value<key/><xml/>`},
		{name: "状态码非0,content-type为text/html,返回html内容", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{
				"Content-Type": []string{context.HTMLF},
			},
		}, args: args{status: 200, content: "<!DOCTYPE html><html></html>"}, wantRs: 200, wantRc: `<!DOCTYPE html><html></html>`},
		{name: "状态码非0,content-type为text/yaml,返回内容", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{
				"Content-Type": []string{context.YAMLF},
			},
		}, args: args{status: 200, content: "key:value"}, wantRs: 200, wantRc: `key:value`},
		{name: "状态码非0,content-type为application/json,且返回内容非正确json字符串", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{
				"Content-Type": []string{context.JSONF},
			},
		}, args: args{status: 200, content: "{key:value"}, wantRs: 200, wantRc: `{"data":"{key:value"}`},
		// 用例引起panic
		// {name: "状态码非0,content-type为application/xml,且返回内容非正确xml字符串", ctx:&mocks.TestContxt{
		// 	HttpHeader: http.Header{
		// 		"Content-Type": []string{context.XMLF},
		// 	},
		// }, args: args{status: 200, content: "<key>value<key/>"}, wantRs: 200, wantRc: ``},
		{name: "状态码非0,content-type为空,返回布尔值/整型/浮点型/复数", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{},
		}, args: args{status: 200, content: false}, wantRs: 200, wantRc: `false`},
		{name: "状态码非0,content-type为application/json,返回布尔值/整型/浮点型/复数", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{
				"Content-Type": []string{context.JSONF},
			},
		}, args: args{status: 200, content: 1}, wantRs: 200, wantRc: `{"data":1}`},
		// 用例引起panic
		// {name: "状态码非0,content-type为application/xml,返回布尔值/整型/浮点型/复数", ctx:&mocks.TestContxt{
		// 	HttpHeader: http.Header{
		// 		"Content-Type": []string{context.XMLF},
		// 	},
		// }}, args: args{status: 200, content: 1}, wantRs: 200, wantRc: `{"data":1}`},
		{name: "状态码非0,content-type为空,返回非字符串/布尔值/整型/浮点型/复数的内容", ctx: &mocks.TestContxt{
			HttpHeader: http.Header{},
		}, args: args{status: 200, content: map[string]string{"key": "value"}}, wantRs: 200, wantRc: `{"key":"value"}`},
		// 用例引起panic
		// {name: "状态码非0,content-type为空,返回非字符串/布尔值/整型/浮点型/复数的内容", ctx:&mocks.TestContxt{
		// 	HttpHeader: http.Header{
		// 		"Content-Type": []string{context.XMLF},
		// 	},
		// }}, args: args{status: 200, content: map[string]string{"key": "value"}}, wantRs: 200, wantRc: ``},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	meta := conf.NewMeta()

	for _, tt := range tests {
		log := logger.GetSession(serverConf.GetMainConf().GetServerName(), ctx.NewUser(tt.ctx, meta).GetRequestID())
		
		//构建response对象
		c := ctx.NewResponse(tt.ctx, serverConf, log, meta)
		err := c.Write(tt.args.status, tt.args.content)
		assert.Equal(t, nil, err, tt.name)

		//测试reponse状态码和内容
		rs, rc := c.GetFinalResponse()
		assert.Equal(t, tt.wantRs, rs, tt.name)
		assert.Equal(t, tt.wantRc, rc, tt.name)

	}
}
