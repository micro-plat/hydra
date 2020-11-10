package context

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func Test_request_Bind(t *testing.T) {
	type result struct {
		Key   string `json:"key" valid:"required"`
		Value string `json:"value" valid:"required"`
	}
	tests := []struct {
		name       string
		body       string
		out        interface{}
		wantErrStr string
		want       interface{}
	}{
		{name: "参数非指针,无法进行数据绑定", out: map[string]string{}, wantErrStr: "输入参数非指针 map"},
		{name: "参数类型非struct,无法进行数据绑定", out: &map[string]string{}, wantErrStr: "输入参数非struct map"},
		{name: "绑定数据为空", body: "", out: &result{}, wantErrStr: "unexpected end of JSON input"},
		{name: "绑定数据验证错误", body: `{"key":"","value":"2"}`, out: &result{}, wantErrStr: "输入参数有误 key: non zero value required"},
		{name: "正确绑定", body: `{"key":"1","value":"2"}`, out: &result{}, want: &result{Key: "1", Value: "2"}},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	for _, tt := range tests {
		r := ctx.NewRequest(&mocks.TestContxt{Form: url.Values{"__body_": []string{tt.body}}}, serverConf, conf.NewMeta())

		err := r.Bind(tt.out)
		if tt.wantErrStr != "" {
			assert.Equal(t, tt.wantErrStr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, tt.want, tt.out, tt.name)
	}
}

func Test_request_Bind_WithHttp(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		want        string
	}{
		{name: "绑定xml数据", contentType: "application/xml;charset=UTF-8", //gin 不支持gbk编码的xml绑定
			body: `<?xml version="1.0" encoding="utf-8" ?><data><key>数据绑定bind!@#$%^&amp;*()_+</key><value>12</value></data>`,
			want: `{"key":"数据绑定bind!@#$%^\u0026*()_+","value":"12"}`},
		{name: "绑定json数据", contentType: "application/json;charset=utf-8",
			body: `{"key":"数据绑定bind!@#$%^&*()_+","value":"12"}`,
			want: `{"key":"数据绑定bind!@#$%^\u0026*()_+","value":"12"}`},
		{name: "绑定yaml数据", contentType: "application/x-yaml; charset=utf-8",
			body: "key: 数据绑定bind!@#$%^&*()_+ \nvalue: 12",
			want: `{"key":"数据绑定bind!@#$%^\u0026*()_+","value":"12"}`},
		{name: "绑定form数据", contentType: "application/x-www-form-urlencoded; charset=utf-8",
			body: `key=%E6%95%B0%E6%8D%AE%E7%BB%91%E5%AE%9Abind!%40%23%24%25%5E%26*()_%2B&value=12`,
			want: `{"key":"数据绑定bind!@#$%^\u0026*()_+","value":"12"}`},
	}

	startServer()
	for _, tt := range tests {
		resp, err := http.Post("http://localhost:9091/request/bind", tt.contentType, strings.NewReader(tt.body))
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, "application/json; charset=UTF-8", resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, "200 OK", resp.Status, tt.name)
		assert.Equal(t, 200, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}

func Test_request_Check(t *testing.T) {

	tests := []struct {
		name       string
		ctx        context.IInnerContext
		args       []string
		wantErr    bool
		wantErrStr string
	}{
		{name: "验证非空数据判断", ctx: &mocks.TestContxt{
			Body:       `{"key1":"value1","key2":"value2"}`,
			HttpHeader: http.Header{"Content-Type": []string{context.JSONF}},
		}, args: []string{"key1", "key2"}, wantErr: false},
		{name: "验证空数据判断", ctx: &mocks.TestContxt{
			Body:       `{"key1":"","key2":"value2"}`,
			HttpHeader: http.Header{"Content-Type": []string{context.JSONF}},
		}, args: []string{"key1", "key2"}, wantErrStr: "输入参数:key1值不能为空", wantErr: true},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	for _, tt := range tests {
		r := ctx.NewRequest(tt.ctx, serverConf, conf.NewMeta())
		err := r.Check(tt.args...)
		if err != nil {
			assert.Equal(t, tt.wantErr, err != nil, tt.name)
			assert.Equal(t, tt.wantErrStr, err.Error(), tt.name)
		}
	}
}

func Test_request_Check_WithHttp(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		want        string
	}{
		// {name: "检查xml数据", contentType: "application/xml;charset=UTF-8", //只能验证根节点
		// 	body: `<?xml version="1.0" encoding="utf-8" ?><data><key>12</key><value>12</value></data>`,
		// 	want: `{"data":"success"}`},
		{name: "检查json数据", contentType: "application/json;charset=utf-8",
			body: `{"key":"12","value":"12"}`,
			want: `{"data":"success"}`},
		{name: "检查form数据", contentType: "application/x-www-form-urlencoded; charset=utf-8",
			body: `key=12&value=12`,
			want: `{"data":"success"}`},
		{name: "检查yaml数据", contentType: "application/x-yaml;charset=utf-8",
			body: "key: key \nvalue: value",
			want: `{"data":"success"}`},
	}

	startServer()
	for _, tt := range tests {
		resp, err := http.Post("http://localhost:9091/request/check", tt.contentType, strings.NewReader(tt.body))
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, "application/json; charset=UTF-8", resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, "200 OK", resp.Status, tt.name)
		assert.Equal(t, 200, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(body))
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}

func Test_request_GetKeys(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	r := ctx.NewRequest(&mocks.TestContxt{
		Form:       url.Values{"key3": []string{}},
		Body:       `{"key1":"value1","key2":"value2"}`,
		HttpHeader: http.Header{"Content-Type": []string{context.JSONF}},
	}, serverConf, conf.NewMeta())

	//获取所有key
	got := r.GetKeys()
	sort.Strings(got)
	assert.Equal(t, []string{"key1", "key2", "key3"}, got, "获取所有key")

}

func Test_request_GetMap(t *testing.T) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置

	r := ctx.NewRequest(&mocks.TestContxt{
		Form:       url.Values{"key3": []string{"value3"}},
		Body:       `{"key1":"value1","key2":"value2"}`,
		HttpHeader: http.Header{"Content-Type": []string{context.JSONF}},
	}, serverConf, conf.NewMeta())

	//获取所有key
	got, err := r.GetMap()

	assert.Equal(t, false, (err != nil), "获取所有map")
	assert.Equal(t, map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"}, got, "获取所有map")

}

func Test_request_Get(t *testing.T) {

	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantOk     bool
	}{
		{name: "通过BodyMap获取key对应的value", args: args{name: "key1"}, wantResult: "value1", wantOk: true},
		{name: "通过FormValue获取key对应的value", args: args{name: "key2"}, wantResult: "  value2", wantOk: true},
		{name: "获取不存在key的值", args: args{name: "key3"}, wantResult: "", wantOk: false},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	r := ctx.NewRequest(&mocks.TestContxt{
		Body:       `{"key1":"value1"}`,
		Form:       url.Values{"key2": []string{"%20+value2"}},
		HttpHeader: http.Header{"Content-Type": []string{context.JSONF}},
	}, serverConf, conf.NewMeta())

	for _, tt := range tests {
		gotResult, gotOk := r.Get(tt.args.name)
		assert.Equal(t, tt.wantResult, gotResult, tt.name)
		assert.Equal(t, tt.wantOk, gotOk, tt.name)
	}

}
