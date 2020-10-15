package context

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_body_GetBody(t *testing.T) {

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	rpath := ctx.NewRpath(&mocks.TestContxt{}, serverConf, conf.NewMeta())
	w1 := ctx.NewBody(&mocks.TestContxt{Body: "%20+body"}, rpath)
	w2 := ctx.NewBody(&mocks.TestContxt{Body: "TEST_BODY_READ_ERR"}, rpath)
	w3 := ctx.NewBody(&mocks.TestContxt{Body: "%-+body"}, rpath)

	type args struct {
		e []string
	}
	tests := []struct {
		name    string
		time    string
		args    args
		wantS   string
		wantErr bool
		err     error
	}{
		{name: "1-首次读取body且读取和解码无错误", time: "1", args: args{}, wantS: "  body", wantErr: false},
		{name: "1-再次读取body且返回body", time: "1", args: args{}, wantS: "  body", wantErr: false},
		{name: "2-首次读取body且读取错误", time: "2", args: args{}, wantS: "", wantErr: true, err: fmt.Errorf("获取body发生错误:读取出错")},
		{name: "2-再次读取body且返回的读取错误", time: "2", args: args{}, wantS: "", wantErr: true, err: fmt.Errorf("读取出错")},
		{name: "3-首次读取body,读取正确,解码错误", time: "3", args: args{}, wantS: "", wantErr: true, err: fmt.Errorf(`url.unescape出错:invalid URL escape "%%-+"`)},
		{name: "3-再次读取body且返回的解码错误", time: "3", args: args{}, wantS: "", wantErr: true, err: fmt.Errorf(`invalid URL escape "%%-+"`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotS string
			var err error
			if tt.time == "1" {
				gotS, err = w1.GetBody(tt.args.e...)
			}
			if tt.time == "2" {
				gotS, err = w2.GetBody(tt.args.e...)
			}
			if tt.time == "3" {
				gotS, err = w3.GetBody(tt.args.e...)
			}
			if (err != nil) == tt.wantErr && tt.err != nil {
				if tt.err.Error() != err.Error() {
					t.Errorf("body.GetBody() error = %v, wantErr %v", err, tt.err)
					return
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("body.GetBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotS != tt.wantS {
				t.Errorf("body.GetBody() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func Test_body_GetBodyMap(t *testing.T) {

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	rpath := ctx.NewRpath(&mocks.TestContxt{}, serverConf, conf.NewMeta())

	type fields struct {
		ctx context.IInnerContext
	}
	type args struct {
		encoding []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{name: "读取正确xml格式数据", fields: fields{ctx: &mocks.TestContxt{
			Body:            `<xml><key1>1&amp;$</key1><key2>value2</key2></xml>`,
			HttpContentType: context.XMLF,
		}}, args: args{encoding: []string{"gbk"}}, want: map[string]interface{}{"key1": "1&$", "key2": "value2"}},
		{name: "读取正确json格式数据", fields: fields{ctx: &mocks.TestContxt{
			Body:            `{"key1":"value1","key2":"value2"}`,
			HttpContentType: context.JSONF,
		}}, args: args{}, want: map[string]interface{}{"key1": "value1", "key2": "value2"}},
		{name: "读取正确yaml格式数据", fields: fields{ctx: &mocks.TestContxt{
			Body:            `key1: value1`,
			HttpContentType: context.YAMLF,
		}}, args: args{}, want: map[string]interface{}{"key1": "value1"}},
		{name: "读取错误的不匹配的格式数据", fields: fields{ctx: &mocks.TestContxt{
			Body:            `{"key1:"value1"}`,
			HttpContentType: context.JSONF,
		}}, args: args{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := ctx.NewBody(tt.fields.ctx, rpath)
			defer func() {
				if r := recover(); r != nil {
					if tt.wantErr {
						return
					}
					t.Errorf("body.GetBodyMap() 错误%+v", r)
				}
			}()
			got, err := w.GetBodyMap(tt.args.encoding...)
			if (err != nil) != tt.wantErr {
				t.Errorf("body.GetBodyMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("body.GetBodyMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
