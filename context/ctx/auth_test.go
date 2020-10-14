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
		{name: "1", fields: fields{}, args: args{}, want: nil},
		{name: "2", fields: fields{}, args: args{v: []interface{}{1}}, want: 1},
		{name: "3", fields: fields{response: 1}, args: args{v: []interface{}{2, 3, 4}}, want: 2},
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
		{name: "1", fields: fields{}, args: args{}, want: nil},
		{name: "2", fields: fields{}, args: args{v: []interface{}{1}}, want: 1},
		{name: "3", fields: fields{request: 1}, args: args{v: []interface{}{2, 3, 4}}, want: 2},
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
	type fields struct {
		request  interface{}
		response interface{}
	}
	type args struct {
		out interface{}
	}
	tests := []struct {
		name           string
		fields         fields
		want           interface{}
		args           args
		wantPanicError bool
	}{
		{name: "request is nil", fields: fields{}, args: args{}, wantPanicError: true},
		{name: "request return nil", fields: fields{request: func() interface{} {
			return nil
		}}, args: args{}, wantPanicError: true},
		{name: "request return non-nil value", fields: fields{request: func() interface{} {
			a := 1
			return a
		}}, args: args{out: 2}, want: 1, wantPanicError: false},
		{name: "request is err string", fields: fields{request: `{"key":"value"}`}, args: args{out: map[string]string{"1": "1"}}, wantPanicError: false},
		{name: "request is string", fields: fields{request: `{"key":"value"}`}, args: args{out: 1}, want: map[string]string{"key": "value"}, wantPanicError: false},
		{name: "default", fields: fields{request: map[string]string{"a": "b"}}, args: args{out: map[string]string{}}, want: map[string]string{"a": "b"}, wantPanicError: false},
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
			c.Bind(&tt.args.out)
			if !reflect.DeepEqual(tt.args.out, tt.want) {
				t.Errorf("auth.Bind() out= %v, want %v", tt.args.out, tt.want)
			}
		})
	}
}
