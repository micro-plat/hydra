/*
author:taoshouyin
time:2020-10-16
*/

package basic

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf"
)

func TestNewBasic(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want *BasicAuth
	}{
		{name: "添加用户名和密码的对象",
			opts: []Option{WithUP("t1", "123"), WithUP("t2", "321")},
			want: &BasicAuth{
				Excludes:      []string{},
				Members:       map[string]string{"t1": "123", "t2": "321"},
				PathMatch:     conf.NewPathMatch([]string{}...),
				authorization: newAuthorization(map[string]string{"t1": "123", "t2": "321"}),
			},
		},
		{name: "添加验证路径的对象",
			opts: []Option{WithExcludes("/t1/t2", "/t1/t2/*")},
			want: &BasicAuth{
				Excludes:      []string{"/t1/t2", "/t1/t2/*"},
				Members:       map[string]string{},
				PathMatch:     conf.NewPathMatch([]string{"/t1/t2", "/t1/t2/*"}...),
				authorization: newAuthorization(map[string]string{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBasic(tt.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicAuth_Verify(t *testing.T) {
	type args struct {
		authValue string
	}
	tests := []struct {
		name   string
		fields *BasicAuth
		args   args
		want   string
		want1  bool
	}{
		{name: "空数据认证", fields: NewBasic(WithUP("", "")), args: args{authValue: createAuth("", "")}, want: "", want1: false},
		{name: "空数据认证1", fields: NewBasic(WithUP("t1", "123")), args: args{authValue: createAuth("", "")}, want: "", want1: false},
		{name: "数据认证通过", fields: NewBasic(WithUP("t1", "123")), args: args{authValue: createAuth("t1", "123")}, want: "t1", want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BasicAuth{
				Excludes:      tt.fields.Excludes,
				Members:       tt.fields.Members,
				Disable:       tt.fields.Disable,
				PathMatch:     tt.fields.PathMatch,
				authorization: tt.fields.authorization,
			}
			got, got1 := b.Verify(tt.args.authValue)
			if got != tt.want {
				t.Errorf("BasicAuth.Verify() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BasicAuth.Verify() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
