package context

import (
	"fmt"
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
		fmt.Println(err)
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

func Test_rpath_GetCookies(t *testing.T) {

	tests := []struct {
		name string
		ctx  context.IInnerContext
		want map[string]string
	}{
		{name: "获取全部cookies", ctx: &mocks.TestContxt{
			Cookie: []*http.Cookie{&http.Cookie{Name: "cookie1", Value: "value1"}, &http.Cookie{Name: "cookie2", Value: "value2"}},
		}, want: map[string]string{"cookie1": "value1", "cookie2": "value2"}},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	for _, tt := range tests {
		c := ctx.NewRequest(tt.ctx, serverConf, conf.NewMeta())
		got := c.GetCookies()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_rpath_GetCookies_WithHttp(t *testing.T) {

	tests := []struct {
		name            string
		contentType     string
		cookie          http.Cookie
		want            string
		wantStatus      string
		wantContentType string
		wantStatusCode  int
	}{
		//net/http: invalid byte 'Ö' in Cookie.Value; dropping invalid bytes
		{name: "cookie编码为中文GBK,无法提交的cookie", contentType: "application/json;charset=gbk", wantContentType: "application/json; charset=gbk",
			cookie:     http.Cookie{Name: "cname", Value: Utf8ToGbk("中文")},
			wantStatus: "200 OK", wantStatusCode: 200, want: `{"cname":""}`},
		//net/http: invalid byte 'ä' in Cookie.Value; dropping invalid bytes
		{name: "cookie编码为中文UTF-8,无法提交的cookie", contentType: "application/json;charset=utf-8", wantContentType: "application/json; charset=utf-8",
			cookie:     http.Cookie{Name: "cname", Value: "中文"},
			wantStatus: "200 OK", wantStatusCode: 200, want: `{"cname":""}`},
		{name: "cookie不带中文", contentType: "application/json;charset=utf-8", wantContentType: "application/json; charset=utf-8",
			cookie:     http.Cookie{Name: "cname", Value: "value!@#$%^&*()_+"},
			wantStatus: "200 OK", wantStatusCode: 200, want: `{"cname":"value!@#$%^\u0026*()_+"}`},
	}

	startServer()
	for _, tt := range tests {
		r, err := http.NewRequest("POST", "http://localhost:9091/getcookies/encoding", strings.NewReader(""))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//添加cookie
		r.AddCookie(&tt.cookie)

		client := &http.Client{}
		resp, err := client.Do(r)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, tt.wantContentType, resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, "200 OK", resp.Status, tt.name)
		assert.Equal(t, 200, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}

func Test_rpath_GetHeader_Encoding(t *testing.T) {

	tests := []struct {
		name            string
		contentType     string
		hName           string
		hValue          string
		want            string
		wantStatus      string
		wantContentType string
		wantStatusCode  int
	}{
		{name: "头部编码为GBK,内容为GBK", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "510 Not Extended", wantStatusCode: 510, want: "Server Error"},
		{name: "头部编码为GBK,内容为UTF8", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部编码为UTF-8,内容为UTF8", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
		{name: "头部编码为UTF-8,内容为GBK", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部编码未设置,内容为utf-8", contentType: "application/json", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
		{name: "头部编码未设置,内容为gbk", contentType: "application/json", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
	}

	startServer()
	for _, tt := range tests {
		r, err := http.NewRequest("POST", "http://localhost:9091/getheaders/encoding", strings.NewReader(""))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//添加cookie
		r.Header.Set(tt.hName, tt.hValue)

		client := &http.Client{}
		resp, err := client.Do(r)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, tt.wantContentType, resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, tt.wantStatus, resp.Status, tt.name)
		assert.Equal(t, tt.wantStatusCode, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}

func Test_rpath_GetHeader_Encoding_GBK(t *testing.T) {
	tests := []struct {
		name            string
		contentType     string
		hName           string
		hValue          string
		want            string
		wantStatus      string
		wantContentType string
		wantStatusCode  int
	}{
		{name: "头部为gbk,内容为gbk", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "510 Not Extended", wantStatusCode: 510, want: "Server Error"},
		{name: "头部为gbk,内容为utf-8", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部为utf-8,内容为gbk", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "510 Not Extended", wantStatusCode: 510, want: "Server Error"},
		{name: "头部为utf-8,内容为utf-8", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部编码未设置,内容为gbk", contentType: "application/json", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "510 Not Extended", wantStatusCode: 510, want: "Server Error"},
		{name: "头部编码未设置,内容为utf-8", contentType: "application/json", wantContentType: "text/plain; charset=gbk",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
	}

	startServer()
	for _, tt := range tests {
		r, err := http.NewRequest("POST", "http://localhost:9091/getheaders/encoding/gbk", strings.NewReader(""))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//添加cookie
		r.Header.Set(tt.hName, tt.hValue)

		client := &http.Client{}
		resp, err := client.Do(r)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, tt.wantContentType, resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, tt.wantStatus, resp.Status, tt.name)
		assert.Equal(t, tt.wantStatusCode, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}
func Test_rpath_GetHeader_Encoding_UTF8(t *testing.T) {
	tests := []struct {
		name            string
		contentType     string
		hName           string
		hValue          string
		want            string
		wantStatus      string
		wantContentType string
		wantStatusCode  int
	}{
		{name: "头部为gbk,内容为gbk", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部为gbk,内容为utf-8", contentType: "application/json;charset=gbk", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
		{name: "头部为utf-8,内容为gbk", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部为utf-8,内容为utf-8", contentType: "application/json;charset=utf-8", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
		{name: "头部编码未设置,内容为gbk", contentType: "application/json", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: Utf8ToGbk("中文!@#$%^&*("),
			wantStatus: "200 OK", wantStatusCode: 200, want: Utf8ToGbk("中文!@#$%^&*(")},
		{name: "头部编码未设置,内容为utf-8", contentType: "application/json", wantContentType: "text/plain; charset=utf-8",
			hName: "hname", hValue: "中文!@#$%^&*(",
			wantStatus: "200 OK", wantStatusCode: 200, want: "中文!@#$%^&*("},
	}

	startServer()
	for _, tt := range tests {
		r, err := http.NewRequest("POST", "http://localhost:9091/getheaders/encoding/utf8", strings.NewReader(""))
		assert.Equal(t, nil, err, "构建请求")

		//设置content-type
		r.Header.Set("Content-Type", tt.contentType)

		//添加cookie
		r.Header.Set(tt.hName, tt.hValue)

		client := &http.Client{}
		resp, err := client.Do(r)
		assert.Equal(t, false, err != nil, tt.name)
		defer resp.Body.Close()
		assert.Equal(t, tt.wantContentType, resp.Header["Content-Type"][0], tt.name)
		assert.Equal(t, tt.wantStatus, resp.Status, tt.name)
		assert.Equal(t, tt.wantStatusCode, resp.StatusCode, tt.name)
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.want, string(body), tt.name)
	}
}
func Test_rpath_GetCookie(t *testing.T) {

	tests := []struct {
		name       string
		cookieName string
		want       string
		want1      bool
	}{
		{name: "获取存在cookies", cookieName: "cookie2", want: "value2", want1: true},
		{name: "获取不存在cookies", cookieName: "cookie3", want: "", want1: false},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	rpath := ctx.NewRequest(&mocks.TestContxt{
		Cookie: []*http.Cookie{&http.Cookie{Name: "cookie1", Value: "value1"}, &http.Cookie{Name: "cookie2", Value: "value2"}},
	}, serverConf, conf.NewMeta())

	for _, tt := range tests {
		got, got1 := rpath.GetCookie(tt.cookieName)
		assert.Equal(t, tt.want, got, tt.name)
		assert.Equal(t, tt.want1, got1, tt.name)
	}
}
