package tests

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
)

func Test_response_getString(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	meta := conf.NewMeta()

	type data struct {
		Key   string `json:"key" xml:"key"`
		Value string `json:"value" yaml:"value" xml:"value"`
	}
	param := &data{Key: "key", Value: "value"}

	type fields struct {
		ctx context.IInnerContext
	}
	type args struct {
		ctp string
		v   interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "获取nil序列化的值", fields: fields{&mocks.TestContxt{}}, args: args{ctp: context.XMLF, v: nil}, want: ""},
		{name: "获取xml序列化值", fields: fields{&mocks.TestContxt{}}, args: args{ctp: context.XMLF, v: param}, want: "<data><key>key</key><value>value</value></data>"},
		{name: "获取yaml序列化的值", fields: fields{&mocks.TestContxt{}}, args: args{ctp: context.YAMLF, v: param}, want: "key: key\nvalue: value\n"},
		{name: "获取json序列化的值", fields: fields{&mocks.TestContxt{}}, args: args{ctp: context.JSONF, v: param}, want: `{"key":"key","value":"value"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.GetSession(serverConf.GetMainConf().GetServerName(), ctx.NewUser(tt.fields.ctx, meta).GetRequestID())
			c := ctx.NewResponse(tt.fields.ctx, serverConf, log, meta)
			got := c.getString(tt.args.ctp, tt.args.v)
			if got != tt.want {
				t.Errorf("response.getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_response_writeNow(t *testing.T) {
// 	type fields struct {
// 		ctx         context.IInnerContext
// 		conf        server.IServerConf
// 		noneedWrite bool
// 		log         logger.ILogger
// 		asyncWrite  func() error
// 		specials    []string
// 	}
// 	type args struct {
// 		status  int
// 		ctyp    string
// 		content string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &response{
// 				ctx:         tt.fields.ctx,
// 				conf:        tt.fields.conf,
// 				path:        tt.fields.path,
// 				raw:         tt.fields.raw,
// 				final:       tt.fields.final,
// 				noneedWrite: tt.fields.noneedWrite,
// 				log:         tt.fields.log,
// 				asyncWrite:  tt.fields.asyncWrite,
// 				specials:    tt.fields.specials,
// 			}
// 			if err := c.writeNow(tt.args.status, tt.args.ctyp, tt.args.content); (err != nil) != tt.wantErr {
// 				t.Errorf("response.writeNow() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func Test_response_getContentType(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	meta := conf.NewMeta()

	type fields struct {
		ctx context.IInnerContext
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "获取context中的content-type", fields: fields{ctx: &mocks.TestContxt{
			HttpHeader: http.Header{
				"X-Request-Id": []string{"123456"},
				"Content-Type": []string{"ctp"},
			},
		}}, want: "ctp"},
		{name: "获取serverConf中的content-type", fields: fields{ctx: &mocks.TestContxt{}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.GetSession(serverConf.GetMainConf().GetServerName(), ctx.NewUser(tt.fields.ctx, meta).GetRequestID())
			c := ctx.NewResponse(tt.fields.ctx, serverConf, log, meta)
			if got := c.getContentType(); got != tt.want {
				t.Errorf("response.getContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_response_swapBytp(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	meta := conf.NewMeta()
	log := logger.GetSession(serverConf.GetMainConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, meta).GetRequestID())
	c := ctx.NewResponse(&mocks.TestContxt{}, serverConf, log, meta)
	type args struct {
		status  int
		content interface{}
	}
	tests := []struct {
		name   string
		args   args
		wantRs int
		wantRc interface{}
	}{
		{name: "成功返回", args: args{status: 0, content: nil}, wantRs: 200, wantRc: ""},
		{name: "返回包含错误码的错误", args: args{status: 0, content: errs.NewError(999, "错误")}, wantRs: 999, wantRc: "错误"},
		{name: "返回状态码在200到400之间的错误", args: args{status: 300, content: errors.New("err")}, wantRs: 400, wantRc: "err"},
		{name: "默认返回", args: args{status: 500, content: "content"}, wantRs: 500, wantRc: "content"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRs, gotRc := c.swapBytp(tt.args.status, tt.args.content)
			if gotRs != tt.wantRs {
				t.Errorf("response.swapBytp() gotRs = %v, want %v", gotRs, tt.wantRs)
			}
			if !reflect.DeepEqual(gotRc, tt.wantRc) {
				t.Errorf("response.swapBytp() gotRc = %v, want %v", gotRc, tt.wantRc)
			}
		})
	}
}
