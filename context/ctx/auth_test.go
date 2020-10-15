package ctx

import (
	"reflect"
	"testing"

	"github.com/micro-plat/lib4go/errs"
)

func Test_auth_Response(t *testing.T) {
	type fields struct {
		request  interface{}
		response interface{}
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{name: "参数为空", fields: fields{}, args: args{}, want: nil},
		{name: "request为空,参数不为空", fields: fields{}, args: args{v: []interface{}{1}}, want: 1},
		{name: "request不为空,参数不为空", fields: fields{response: 1}, args: args{v: []interface{}{2, 3, 4}}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &auth{
				request:  tt.fields.request,
				response: tt.fields.response,
			}
			if got := c.Response(tt.args.v...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auth.Response() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_Request(t *testing.T) {
	type fields struct {
		request  interface{}
		response interface{}
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{name: "参数为空", fields: fields{}, args: args{}, want: nil},
		{name: "response为空,参数不为空", fields: fields{}, args: args{v: []interface{}{1}}, want: 1},
		{name: "response不为空,参数不为空", fields: fields{request: 1}, args: args{v: []interface{}{2, 3, 4}}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &auth{
				request:  tt.fields.request,
				response: tt.fields.response,
			}
			if got := c.Request(tt.args.v...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auth.Request() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_Bind(t *testing.T) {
	type result struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	type fields struct {
		request  interface{}
		response interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		want           interface{}
		wantPanicError bool
		def            bool
	}{
		{name: "request为空", fields: fields{}, wantPanicError: true},
		{name: "request为func返回空值", fields: fields{request: func() interface{} {
			return nil
		}}, wantPanicError: true},
		{name: "request为func返回非空值", fields: fields{request: func() interface{} {
			return result{Key: "1", Value: "1"}
		}}, want: result{Key: "1", Value: "1"}, wantPanicError: false},
		{name: "request为错误的json字符串", fields: fields{request: `{"key":"1",v}`}, wantPanicError: true},
		{name: "request为json字符串", fields: fields{request: `{"key":"1","value":"1"}`}, want: result{Key: "1", Value: "1"}, wantPanicError: false},
		{name: "默认情况", fields: fields{request: map[string]string{"key": "value"}}, def: true, want: map[string]string{"key": "value"}, wantPanicError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &auth{
				request:  tt.fields.request,
				response: tt.fields.response,
			}
			defer func() {
				if r := recover(); r != nil {
					if e, ok := r.(*errs.Error); ok {
						if tt.wantPanicError && e.GetCode() == 401 && e.Error() == "请求中未包含用户信息,用户未登录" {
							return
						}
					}
					if tt.wantPanicError {
						return
					}
					t.Errorf("recover:Bind() = %v", r)
				}
			}()

			if !tt.def {
				out := result{}
				c.Bind(&out)
				if !reflect.DeepEqual(out, tt.want) {
					t.Errorf("auth.Bind() out= %v, want %v", out, tt.want)
				}
				return
			}

			out := map[string]string{}
			c.Bind(&out)
			if !reflect.DeepEqual(out, tt.want) {
				t.Errorf("auth.Bind() out= %v, want %v", out, tt.want)
			}
		})
	}
}
