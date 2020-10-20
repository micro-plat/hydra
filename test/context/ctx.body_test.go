package context

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

type testBody struct {
	name    string
	wantS   string
	wantErr bool
	err     error
}

func Test_body_GetBody(t *testing.T) {
	//测试读取body正确
	tests := []testBody{
		{name: "1-首次读取body且读取和解码无错误", wantS: "  body", wantErr: false},
		{name: "1-再次读取body且返回body", wantS: "  body", wantErr: false},
	}
	testGetBody(t, "%20+body", tests)

	//测试读取错误
	testReadErr := []testBody{
		{name: "2-首次读取body且读取错误", wantS: "", wantErr: true, err: fmt.Errorf("获取body发生错误:读取出错")},
		{name: "2-再次读取body且返回的读取错误", wantS: "", wantErr: true, err: fmt.Errorf("读取出错")},
	}
	testGetBody(t, "TEST_BODY_READ_ERR", testReadErr)

	//测试解码错误
	testUnescapeErr := []testBody{
		{name: "3-首次读取body,读取正确,解码错误", wantS: "", wantErr: true, err: fmt.Errorf(`url.unescape出错:invalid URL escape "%%-+"`)},
		{name: "3-再次读取body且返回的解码错误", wantS: "", wantErr: true, err: fmt.Errorf(`invalid URL escape "%%-+"`)},
	}
	testGetBody(t, "%-+body", testUnescapeErr)

}

func testGetBody(t *testing.T, body string, tests []testBody) {
	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	rpath := ctx.NewRpath(&mocks.TestContxt{}, serverConf, conf.NewMeta())
	w := ctx.NewBody(&mocks.TestContxt{Body: body}, rpath)

	for _, tt := range tests {
		gotS, err := w.GetBody()
		if (err != nil) == tt.wantErr && tt.err != nil {
			assert.Equal(t, tt.err.Error(), err.Error(), tt.name)
		}
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.Equal(t, tt.wantS, gotS, tt.name)
	}
}

func Test_body_GetBodyMap(t *testing.T) {

	tests := []struct {
		name     string
		ctx      context.IInnerContext
		encoding []string
		want     map[string]interface{}
		wantErr  bool
	}{
		{name: "读取正确xml格式数据", ctx: &mocks.TestContxt{
			Body:            `<xml><key1>1&amp;$</key1><key2>value2</key2></xml>`,
			HttpContentType: context.XMLF,
		}, encoding: []string{"gbk"}, want: map[string]interface{}{"key1": "1&$", "key2": "value2"}},
		{name: "读取正确json格式数据", ctx: &mocks.TestContxt{
			Body:            `{"key1":"value1","key2":"value2"}`,
			HttpContentType: context.JSONF,
		}, want: map[string]interface{}{"key1": "value1", "key2": "value2"}},
		{name: "读取正确yaml格式数据", ctx: &mocks.TestContxt{
			Body:            `key1: value1`,
			HttpContentType: context.YAMLF,
		}, want: map[string]interface{}{"key1": "value1"}},
		{name: "读取错误的不匹配的格式数据", ctx: &mocks.TestContxt{
			Body:            `{"key1:"value1"}`,
			HttpContentType: context.JSONF,
		}, wantErr: true},
	}

	confObj := mocks.NewConf()         //构建对象
	confObj.API(":8080")               //初始化参数
	serverConf := confObj.GetAPIConf() //获取配置
	rpath := ctx.NewRpath(&mocks.TestContxt{}, serverConf, conf.NewMeta())

	for _, tt := range tests {
		w := ctx.NewBody(tt.ctx, rpath)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, tt.wantErr, true, tt.name)
			}
		}()
		got, err := w.GetBodyMap(tt.encoding...)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
