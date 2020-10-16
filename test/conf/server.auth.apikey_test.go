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

	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/micro-plat/lib4go/security/sha1"
	"github.com/micro-plat/lib4go/security/sha256"
)

var md5secret = "12345678"
var sha1secret = "1234567812345678"
var sha256secret = "9876543210222222222"
var rawData = "taosy hydra test"
var md5Sign = md5.Encrypt(rawData + md5secret)
var sha1Sign = sha1.Encrypt(rawData + sha1secret)
var sha256Sign = sha256.Encrypt(rawData + sha256secret)

func TestAPIKeyNew(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		opts   []apikey.Option
		want   *apikey.APIKeyAuth
	}{
		{name: "初始化默认对象", secret: "", opts: []apikey.Option{}, want: &apikey.APIKeyAuth{Mode: "MD5", PathMatch: conf.NewPathMatch()}},
		{name: "设置密钥和路径", secret: "1111", opts: []apikey.Option{apikey.WithSecret("123456"), apikey.WithExcludes("/t/tw", "/t1/t2")}, want: &apikey.APIKeyAuth{Secret: "123456", Excludes: []string{"/t/tw", "/t1/t2"}, Mode: "MD5", PathMatch: conf.NewPathMatch("/t/tw", "/t1/t2")}},
		{name: "设置sha1", secret: "1111", opts: []apikey.Option{apikey.WithSHA1Mode(), apikey.WithDisable()}, want: &apikey.APIKeyAuth{Secret: "1111", Mode: "SHA1", Disable: true, PathMatch: conf.NewPathMatch()}},
		{name: "设置sha256", secret: "", opts: []apikey.Option{apikey.WithSHA256Mode(), apikey.WithEnable()}, want: &apikey.APIKeyAuth{Mode: "SHA256", Disable: false, PathMatch: conf.NewPathMatch()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apikey.New(tt.secret, tt.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIKeyAuth_Verify(t *testing.T) {

	type args struct {
		raw  string
		sign string
	}
	tests := []struct {
		name    string
		mode    string
		secret  string
		args    args
		wantErr bool
	}{
		{name: "不支持的签名方式", mode: "md4", secret: md5secret, args: args{raw: rawData, sign: md5Sign}, wantErr: true},
		{name: "签名方式不正确", mode: "md5", secret: sha1secret, args: args{raw: rawData, sign: sha1Sign}, wantErr: true},
		{name: "签名数据错误", mode: "md5", secret: md5secret, args: args{raw: "rawData", sign: md5Sign}, wantErr: true},
		{name: "密钥错误", mode: "md5", secret: "md5secret", args: args{raw: rawData, sign: md5Sign}, wantErr: true},
		{name: "md5签名成功", mode: "md5", secret: md5secret, args: args{raw: rawData, sign: md5Sign}, wantErr: false},
		{name: "sha1签名成功", mode: "sha1", secret: sha1secret, args: args{raw: rawData, sign: sha1Sign}, wantErr: false},
		{name: "sha256签名成功", mode: "sha256", secret: sha256secret, args: args{raw: rawData, sign: sha256Sign}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apikey.APIKeyAuth{
				Mode:   tt.mode,
				Secret: tt.secret,
			}
			if err := a.Verify(tt.args.raw, tt.args.sign); (err != nil) != tt.wantErr {
				t.Errorf("APIKeyAuth.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApikeyGetConf(t *testing.T) {

	tests := []struct {
		name string
		args func() conf.IMainConf
		want *apikey.APIKeyAuth
	}{
		{name: "未设置apikey节点", args: func() conf.IMainConf {
			conf := mocks.NewConf()
			conf.API(":8081")
			return conf.GetAPIConf().GetMainConf()
		}, want: &apikey.APIKeyAuth{Disable: true, PathMatch: conf.NewPathMatch()}},
		{name: "配置参数有误", args: func() conf.IMainConf {
			conf := mocks.NewConf()
			conf.API(":8081").APIKEY("")
			return conf.GetAPIConf().GetMainConf()
		}, want: nil},
		{name: "配置参数正确", args: func() conf.IMainConf {
			conf := mocks.NewConf()
			conf.API(":8081").APIKEY("123456", apikey.WithMD5Mode())
			return conf.GetAPIConf().GetMainConf()
		}, want: apikey.New("123456", apikey.WithMD5Mode())},
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
			got := apikey.GetConf(tt.args())
			assert.Equal(t, got, tt.want, tt.name)
		})
	}
}
