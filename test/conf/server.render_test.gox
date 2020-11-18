package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"

	"github.com/micro-plat/hydra/conf/server/render"
)

func TestNewRender(t *testing.T) {
	tests := []struct {
		name string
		opts []render.Option
		want *render.Render
	}{
		{name: "初始化空对象-nil", opts: nil, want: &render.Render{Tmplts: make(map[string]*render.Tmplt)}},
		{name: "初始化空对象", opts: []render.Option{}, want: &render.Render{Tmplts: make(map[string]*render.Tmplt)}},
		{name: "设置disable对象", opts: []render.Option{render.WithDisable()}, want: &render.Render{Disable: true, Tmplts: make(map[string]*render.Tmplt)}},
		{name: "设置enable对象", opts: []render.Option{render.WithEnable()}, want: &render.Render{Disable: false, Tmplts: make(map[string]*render.Tmplt)}},
		{name: "设置单个模板对象对象", opts: []render.Option{render.WithTmplt("/path1", "success", render.WithContentType("tpltm1"))},
			want: &render.Render{Disable: false, Tmplts: map[string]*render.Tmplt{"/path1": &render.Tmplt{ContentType: "tpltm1", Content: "success", Status: ""}}}},
		{name: "设置多个模板对象对象", opts: []render.Option{render.WithTmplt("/path1", "success", render.WithContentType("tpltm1")), render.WithTmplt("/path2", "fail", render.WithContentType("tpltm2"))},
			want: &render.Render{Disable: false, Tmplts: map[string]*render.Tmplt{"/path1": &render.Tmplt{ContentType: "tpltm1", Content: "success", Status: ""}, "/path2": &render.Tmplt{ContentType: "tpltm2", Content: "fail", Status: ""}}}},
	}
	for _, tt := range tests {
		got := render.NewRender(tt.opts...)
		assert.Equal(t, tt.want.Disable, got.Disable, tt.name)
		assert.Equal(t, tt.want.Tmplts, got.Tmplts, tt.name)
	}
}

func TestRenderGetConf(t *testing.T) {
	type test struct {
		name       string
		cnf        conf.IServerConf
		want       *render.Render
		wantErr    bool
		wantErrStr string
	}

	conf := mocks.NewConfBy("hydra", "graytest")
	confB := conf.API(":8090")
	test1 := test{name: "render节点不存在", cnf: conf.GetAPIConf().GetServerConf(), want: &render.Render{Disable: true}, wantErr: false}
	renderObj, err := render.GetConf(test1.cnf)
	assert.Equal(t, test1.wantErr, (err != nil), test1.name)
	assert.Equal(t, test1.want, renderObj, test1.name)

	confB.Render(render.WithTmplt("/p1", "succ", render.WithStatus("")))
	test2 := test{name: "render节点存在,数据错误", cnf: conf.GetAPIConf().GetServerConf(), want: nil, wantErr: true, wantErrStr: "render Tmplt配置数据有误"}
	renderObj, err = render.GetConf(test2.cnf)
	assert.Equal(t, test2.wantErr, (err != nil), test2.name+",err")
	assert.Equal(t, test2.wantErrStr, err.Error()[:len(test2.wantErrStr)], test2.name+",err1")
	assert.Equal(t, test2.want, renderObj, test2.name+",obj")

	confB.Render(render.WithTmplt("/p1", "succ", render.WithStatus("200")))
	test3 := test{name: "render节点存在,正确节点", cnf: conf.GetAPIConf().GetServerConf(),
		want: render.NewRender(render.WithTmplt("/p1", "succ", render.WithStatus("200"))), wantErr: false}
	renderObj, err = render.GetConf(test3.cnf)
	assert.Equal(t, test3.wantErr, (err != nil), test3.name+",err")
	assert.Equal(t, test3.want, renderObj, test3.name+",obj")
}

type TFuncs map[string]interface{}

func TestRender_Get(t *testing.T) {
	tmplFuncs := TFuncs{}
	tmplFuncs["getStatus"] = func(status string) string {
		return status
	}
	tmplFuncs["getRequestID"] = func() string {
		return "20060102150405"
	}

	tmplFuncs["getRequestID1"] = func(ss string) string {
		return ss
	}
	type args struct {
		path  string
		funcs map[string]interface{}
		i     interface{}
	}
	tests := []struct {
		name       string
		fields     *render.Render
		args       args
		want       bool
		want1      int
		want2      string
		want3      string
		wantErr    bool
		wantErrStr string
	}{
		{name: "路径匹配失败", fields: render.NewRender(render.WithTmplt("/p1/p2", "succ", render.WithStatus("{{getStatus .Status}}"))),
			args: args{path: "/t1/t2", funcs: nil, i: nil}, want: false, want1: 0, want2: "", want3: "", wantErr: false},
		{name: "路径匹配通过,匹配status错误", fields: render.NewRender(render.WithTmplt("/p1/p2", "succ", render.WithStatus("{{getStatus}}"))),
			args: args{path: "/p1/p2", funcs: tmplFuncs, i: nil}, want: false, want1: 0, want2: "", want3: "", wantErr: true, wantErrStr: "status模板"},
		{name: "路径匹配通过,匹配status不是数字", fields: render.NewRender(render.WithTmplt("/p1/p2", "succ", render.WithStatus("{{getStatus .Status}}"))),
			args: args{path: "/p1/p2", funcs: tmplFuncs, i: render.Tmplt{Status: "ttt", ContentType: "text/xml", Content: "success"}},
			want: false, want1: 0, want2: "", want3: "", wantErr: true, wantErrStr: "status模板"},
		{name: "路径匹配通过,匹配status正确", fields: render.NewRender(render.WithTmplt("/p1/p2", "", render.WithStatus("{{getStatus .Status}}"))),
			args: args{path: "/p1/p2", funcs: tmplFuncs, i: render.Tmplt{Status: "222", ContentType: "text/xml", Content: "success"}},
			want: true, want1: 222, want2: "", want3: "", wantErr: false},
		{name: "路径匹配通过,匹配content错误", fields: render.NewRender(render.WithTmplt("/p1/p2", "{{getRequestID1}}")),
			args: args{path: "/p1/p2", funcs: tmplFuncs, i: nil}, want: false, want1: 0, want2: "", want3: "", wantErr: true, wantErrStr: "响应内容模板"},
		{name: "路径匹配通过,匹配content错误", fields: render.NewRender(render.WithTmplt("/p1/p2", "{{getRequestID}}")),
			args: args{path: "/p1/p2", funcs: tmplFuncs, i: nil}, want: true, want1: 0, want2: "", want3: "20060102150405", wantErr: false},
		{name: "路径匹配通过,匹配content-type错误", fields: render.NewRender(render.WithTmplt("/p1/p2", "", render.WithContentType("{{getRequestID1}}"))),
			args: args{path: "/p1/p2", funcs: tmplFuncs, i: render.Tmplt{Status: "ttt", ContentType: "text/xml", Content: "success"}},
			want: false, want1: 0, want2: "", want3: "", wantErr: true, wantErrStr: "content_type模板"},
		{name: "路径匹配通过,匹配content-type错误", fields: render.NewRender(render.WithTmplt("/p1/p2", "", render.WithContentType("{{getRequestID1 .ContentType}}"))),
			args: args{path: "/p1/p2", funcs: tmplFuncs, i: render.Tmplt{Status: "ttt", ContentType: "text/xml", Content: "success"}},
			want: true, want1: 0, want2: "text/xml", want3: "", wantErr: false},
		{name: "路径匹配通过,全匹配", fields: render.NewRender(render.WithTmplt("/p1/p2", "{{getRequestID}}", render.WithStatus("{{getStatus .Status}}"), render.WithContentType("{{getRequestID1 .ContentType}}"))),
			args: args{path: "/p1/p2", funcs: tmplFuncs, i: render.Tmplt{Status: "211", ContentType: "text/xml", Content: ""}},
			want: true, want1: 211, want2: "text/xml", want3: "20060102150405", wantErr: false},
	}
	for _, tt := range tests {
		got, got1, got2, got3, err := tt.fields.Get(tt.args.path, tt.args.funcs, tt.args.i)
		assert.Equal(t, tt.wantErr, (err != nil), tt.name+",err")
		if !tt.wantErr {
			assert.Equal(t, tt.want, got, tt.name+",res")
			assert.Equal(t, tt.want1, got1, tt.name+",status")
			assert.Equal(t, tt.want2, got2, tt.name+",ContentType")
			assert.Equal(t, tt.want3, got3, tt.name+",content")
		} else {
			assert.Equal(t, tt.wantErrStr, err.Error()[:len(tt.wantErrStr)], tt.name+",err1")
		}
	}
}
