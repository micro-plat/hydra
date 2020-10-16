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
)

func TestNewJWT(t *testing.T) {
	tests := []struct {
		name string
		opts []jwt.Option
		want *jwt.JWTAuth
	}{
		{name: "设置secert",
			opts: []jwt.Option{jwt.WithSecret("12345678")},
			want: &jwt.JWTAuth{Name: "Authorization-Jwt", Mode: "HS512", Secret: "12345678", ExpireAt: 86400, Source: "COOKIE", PathMatch: conf.NewPathMatch()},
		},
		{name: "设置disable",
			opts: []jwt.Option{jwt.WithSecret("12345678"), jwt.WithDisable()},
			want: &jwt.JWTAuth{Name: "Authorization-Jwt", Mode: "HS512", Secret: "12345678", Disable: true, ExpireAt: 86400, Source: "COOKIE", PathMatch: conf.NewPathMatch()},
		},
		{name: "设置Enable",
			opts: []jwt.Option{jwt.WithSecret("12345678"), jwt.WithEnable()},
			want: &jwt.JWTAuth{Name: "Authorization-Jwt", Mode: "HS512", Secret: "12345678", Disable: false, ExpireAt: 86400, Source: "COOKIE", PathMatch: conf.NewPathMatch()},
		},
		{name: "设置自定义对象",
			opts: []jwt.Option{jwt.WithSecret("12345678"), jwt.WithHeader(), jwt.WithExcludes("/t1/**"), jwt.WithExpireAt(1000), jwt.WithMode("ES256"), jwt.WithName("test"), jwt.WithRedirect("1111")},
			want: &jwt.JWTAuth{Name: "test", Redirect: "1111", Mode: "ES256", Secret: "12345678", ExpireAt: 1000, Source: "HEADER", Excludes: []string{"/t1/**"}, PathMatch: conf.NewPathMatch("/t1/**")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jwt.NewJWT(tt.opts...)
			if !reflect.DeepEqual(got, tt.want) {
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
