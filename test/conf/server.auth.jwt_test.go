/*
author:taoshouyin
time:2020-10-16
*/

package conf

import (
	"reflect"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/auth/jwt"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/utility"
)

func TestNewJWT(t *testing.T) {
	type args struct {
		opts []jwt.Option
	}
	tests := []struct {
		name string
		args args
		want *jwt.JWTAuth
	}{
		{name: "设置默认对象", args: args{opts: []jwt.Option{}}, want: &jwt.JWTAuth{Name: "Authorization-Jwt",
			Mode:      "HS512",
			Secret:    utility.GetGUID(),
			ExpireAt:  86400,
			Source:    "COOKIE",
			PathMatch: conf.NewPathMatch()}},
		{name: "设置自定义对象", args: args{opts: []jwt.Option{
			jwt.WithHeader(), jwt.WithDisable(), jwt.WithExcludes("/t1/**"), jwt.WithExpireAt(1000), jwt.WithMode("ES256"), jwt.WithName("test"), jwt.WithRedirect("1111"), jwt.WithSecret("132459678"),
		}}, want: &jwt.JWTAuth{Name: "test",
			Redirect:  "1111",
			Mode:      "ES256",
			Secret:    "132459678",
			ExpireAt:  1000,
			Source:    "HEADER",
			Disable:   true,
			Excludes:  []string{"/t1/**"},
			PathMatch: conf.NewPathMatch("/t1/**")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jwt.NewJWT(tt.args.opts...)
			if !reflect.DeepEqual(got.Mode, tt.want.Mode) ||
				!reflect.DeepEqual(got.Disable, tt.want.Disable) ||
				!reflect.DeepEqual(got.Excludes, tt.want.Excludes) ||
				!reflect.DeepEqual(got.ExpireAt, tt.want.ExpireAt) ||
				!reflect.DeepEqual(got.Name, tt.want.Name) ||
				!reflect.DeepEqual(*got.PathMatch, *tt.want.PathMatch) ||
				!reflect.DeepEqual(got.Redirect, tt.want.Redirect) ||
				!reflect.DeepEqual(got.Source, tt.want.Source) {
				t.Errorf("NewJWT() = %v, want %v", got, tt.want)
			}

			if tt.name == "设置自定义对象" && !reflect.DeepEqual(got.Secret, tt.want.Secret) {
				t.Errorf("NewJWT()1 = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConf(t *testing.T) {
	tests := []struct {
		name string
		args func() conf.IMainConf
		want *jwt.JWTAuth
	}{
		{name: "未设置jwt节点", args: func() conf.IMainConf {
			conf := mocks.NewConf()
			conf.API(":8081")
			return conf.GetAPIConf().GetMainConf()
		}, want: &jwt.JWTAuth{Disable: true}},
		{name: "配置参数有误", args: func() conf.IMainConf {
			conf := mocks.NewConf()
			conf.API(":8081").Jwt(jwt.WithMode("xxxx"))
			return conf.GetAPIConf().GetMainConf()
		}, want: nil},
		{name: "配置参数正确", args: func() conf.IMainConf {
			conf := mocks.NewConf()
			conf.API(":8081").Jwt(jwt.WithExpireAt(123), jwt.WithSecret("11111"))
			return conf.GetAPIConf().GetMainConf()
		}, want: jwt.NewJWT(jwt.WithExpireAt(123), jwt.WithSecret("11111"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					err1 := err.(error)
					if !(tt.name == "配置参数有误" && strings.Contains(err1.Error(), "配置有误")) {
						t.Errorf("apiKeyGetConf 获取配置对对象失败,err: %v", err)
					}
				}
			}()

			if got := jwt.GetConf(tt.args()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
