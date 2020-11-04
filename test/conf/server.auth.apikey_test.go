/*
author:taoshouyin
time:2020-10-16
*/

package conf

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/auth/apikey"
	"github.com/micro-plat/hydra/registry"
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
		{name: "设置md5", secret: "1111", opts: []apikey.Option{apikey.WithMD5Mode()}, want: &apikey.APIKeyAuth{Secret: "1111", Mode: "MD5", Disable: false, PathMatch: conf.NewPathMatch()}},
		{name: "设置sha1", secret: "1111", opts: []apikey.Option{apikey.WithSHA1Mode(), apikey.WithDisable()}, want: &apikey.APIKeyAuth{Secret: "1111", Mode: "SHA1", Disable: true, PathMatch: conf.NewPathMatch()}},
		{name: "设置sha256", secret: "", opts: []apikey.Option{apikey.WithSHA256Mode(), apikey.WithEnable()}, want: &apikey.APIKeyAuth{Mode: "SHA256", Disable: false, PathMatch: conf.NewPathMatch()}},
	}
	for _, tt := range tests {
		got := apikey.New(tt.secret, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
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
		a := &apikey.APIKeyAuth{Mode: tt.mode, Secret: tt.secret}
		err := a.Verify(tt.args.raw, tt.args.sign)
		assert.Equal(t, tt.wantErr, (err != nil), tt.name)
	}
}

func TestApikeyGetConf(t *testing.T) {
	type test struct {
		name string
		opts []apikey.Option
		want *apikey.APIKeyAuth
	}

	apiConf := mocks.NewConf()
	confB := apiConf.API(":8081")
	test1 := test{name: "未设置apikey节点", want: &apikey.APIKeyAuth{Disable: true, PathMatch: conf.NewPathMatch()}}
	got, err := apikey.GetConf(apiConf.GetAPIConf().GetServerConf())
	assert.Equal(t, nil, err, test1.name+",err")
	assert.Equal(t, got, test1.want, test1.name)

	test2 := test{name: "配置参数正确", opts: []apikey.Option{apikey.WithMD5Mode(), apikey.WithDisable(), apikey.WithExcludes("/t1/t2"), apikey.WithSecret("123456")},
		want: apikey.New("123456", apikey.WithMD5Mode(), apikey.WithDisable(), apikey.WithExcludes("/t1/t2"))}
	confB.APIKEY("", test2.opts...)
	got, err = apikey.GetConf(apiConf.GetAPIConf().GetServerConf())
	assert.Equal(t, nil, err, test2.name+",err")
	assert.Equal(t, got, test2.want, test2.name)
}

//@todo
func xTestApikeyGetConf1(t *testing.T) {
	type test struct {
		name string
		opts []apikey.Option
		want *apikey.APIKeyAuth
	}

	defer func() {
		e := recover()
		if e != nil && strings.Contains(e.(error).Error(), "apikey配置有误1") {
			return
		}
		t.Errorf("节点密钥不存在,验证异常,%v", e)
	}()
	apiConf := mocks.NewConf()
	confB := apiConf.API(":8081")
	test1 := test{name: "节点密钥不存在,验证异常", opts: []apikey.Option{apikey.WithMD5Mode()}, want: apikey.New("", apikey.WithMD5Mode())}
	confB.APIKEY("", test1.opts...)
	apikey.GetConf(apiConf.GetAPIConf().GetServerConf())
	t.Errorf("%s,没有验证参数合法性错误", test1.name)
}

func TestApikeyGetConf2(t *testing.T) {
	type test struct {
		name string
		opts []apikey.Option
		want *apikey.APIKeyAuth
	}

	defer func() {
		e := recover()
		if e != nil && strings.Contains(e.(error).Error(), "apikey配置有误") {
			return
		}
		t.Errorf("apikey修改为错误json串,%v", e)
	}()
	apiConf := mocks.NewConf()
	confB := apiConf.API(":8081")
	test1 := test{name: "apikey修改为错误json串", opts: []apikey.Option{apikey.WithMD5Mode(), apikey.WithDisable(), apikey.WithExcludes("/t1/t2"), apikey.WithSecret("123456")},
		want: apikey.New("123456", apikey.WithMD5Mode(), apikey.WithDisable(), apikey.WithExcludes("/t1/t2"))}
	confB.APIKEY("", test1.opts...)
	// 修改json数据不合法
	path := apiConf.GetAPIConf().GetServerConf().GetSubConfPath("auth", "apikey")
	// ch, _ := apiConf.Registry.WatchValue(path)
	apiConf.Registry.Update(path, "错误的json字符串")

	apiConf = mocks.NewConf()
	ttx, err := apiConf.GetAPIConf().GetServerConf().GetSubConf(registry.Join("auth", "apikey"))
	t.Errorf("111111111111:%v,err:%v", string(ttx.GetRaw()), err)
	// select {
	// case <-time.After(3 * time.Second):
	// 	return
	// case <-ch:
	// 	ttx, err := apiConf.GetAPIConf().GetServerConf().GetSubConf(registry.Join("auth", "apikey"))
	// 	t.Errorf("111111111111:%v,err:%v", string(ttx.GetRaw()), err)
	// 	apikey.GetConf(apiConf.GetAPIConf().GetServerConf())
	// 	t.Errorf("%s,没有验证参数合法性错误", test1.name)
	// }
}
