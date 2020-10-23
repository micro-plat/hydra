/*
author:taoshouyin
time:2020-10-16
*/

package conf

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/auth/basic"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/encoding/base64"
)

func createAuth(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.Encode(base)
}

func TestNewBasic(t *testing.T) {
	tests := []struct {
		name string
		opts []basic.Option
		want *basic.BasicAuth
	}{
		{name: "初始化空对象",
			opts: []basic.Option{},
			want: &basic.BasicAuth{Excludes: []string{}, Members: map[string]string{}, PathMatch: conf.NewPathMatch([]string{}...)},
		},
		{name: "初始化对象disable设置",
			opts: []basic.Option{basic.WithDisable()},
			want: &basic.BasicAuth{Disable: true, Excludes: []string{}, Members: map[string]string{}, PathMatch: conf.NewPathMatch([]string{}...)},
		},
		{name: "初始化对象enable设置",
			opts: []basic.Option{basic.WithEnable()},
			want: &basic.BasicAuth{Disable: false, Excludes: []string{}, Members: map[string]string{}, PathMatch: conf.NewPathMatch([]string{}...)},
		},
		{name: "添加用户名和密码的对象",
			opts: []basic.Option{basic.WithUP("t1", "123"), basic.WithUP("t2", "321")},
			want: &basic.BasicAuth{Excludes: []string{}, Members: map[string]string{"t1": "123", "t2": "321"}, PathMatch: conf.NewPathMatch([]string{}...)},
		},
		{name: "添加验证路径的对象",
			opts: []basic.Option{basic.WithExcludes("/t1/t2", "/t1/t2/*")},
			want: &basic.BasicAuth{Excludes: []string{"/t1/t2", "/t1/t2/*"}, Members: map[string]string{}, PathMatch: conf.NewPathMatch([]string{"/t1/t2", "/t1/t2/*"}...)},
		},
	}
	for _, tt := range tests {
		got := basic.NewBasic(tt.opts...)
		assert.Equal(t, tt.want.Disable, got.Disable, tt.name+",Disable")
		assert.Equal(t, tt.want.Excludes, got.Excludes, tt.name+",Excludes")
		assert.Equal(t, tt.want.Members, got.Members, tt.name+",Members")
		assert.Equal(t, tt.want.PathMatch, got.PathMatch, tt.name+",PathMatch")
	}
}

func TestBasicAuth_Verify(t *testing.T) {
	tests := []struct {
		name   string
		fields *basic.BasicAuth
		args   string
		want   string
		want1  bool
	}{
		{name: "空数据认证", fields: basic.NewBasic(basic.WithUP("", "")), args: createAuth("", ""), want: "", want1: false},
		{name: "空数据认证1", fields: basic.NewBasic(basic.WithUP("t1", "123")), args: createAuth("", ""), want: "", want1: false},
		{name: "数据认证通过", fields: basic.NewBasic(basic.WithUP("t1", "123")), args: createAuth("t1", "123"), want: "t1", want1: true},
	}
	for _, tt := range tests {
		got, got1 := tt.fields.Verify(tt.args)
		assert.Equal(t, tt.want, got, tt.name+",username")
		assert.Equal(t, tt.want1, got1, tt.name+",bool")
	}
}

func TestBasicGetConf(t *testing.T) {

	tests := []struct {
		name string
		opts []basic.Option
		want *basic.BasicAuth
	}{
		{name: "basic节点不存在", opts: []basic.Option{}, want: &basic.BasicAuth{Disable: true}},
		{name: "Members==0节点", opts: []basic.Option{}, want: &basic.BasicAuth{Disable: true}},
		// {name: "Members!=0的空节点", opts: []basic.Option{basic.WithUP("t1", "123")}, want: basic.NewBasic(basic.WithUP("t1", "123"))},
		{name: "配置参数正确", opts: []basic.Option{basic.WithUP("t1", "123"), basic.WithExcludes("/t1/t12"), basic.WithDisable()},
			want: basic.NewBasic(basic.WithUP("t1", "123"), basic.WithExcludes("/t1/t12"), basic.WithDisable())},
	}

	conf := mocks.NewConf()
	confB := conf.API(":8081")
	for _, tt := range tests {
		if !strings.EqualFold(tt.name, "basic节点不存在") {
			confB.Basic(tt.opts...)
		}
		got, err := basic.GetConf(conf.GetAPIConf().GetMainConf())
		assert.NotEqual(t, nil, err, tt.name+",err")
		assert.Equal(t, got, tt.want, tt.name)
	}
}
