package conf

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/static"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestStaticNew(t *testing.T) {
	defaultObj := &static.Static{
		FileMap:   map[string]static.FileInfo{},
		Dir:       static.DefaultSataticDir,
		FirstPage: static.DefaultFirstPage,
		Rewriters: static.DefaultRewriters,
		Exclude:   static.DefaultExclude,
		Exts:      []string{},
	}
	enObj := &static.Static{
		FileMap:   map[string]static.FileInfo{},
		Dir:       "./test",
		FirstPage: "index1.html",
		Rewriters: []string{"/", "indextest.htm", "defaulttest.html"},
		Exclude:   []string{"/views/", ".exe", ".so", ".zip"},
		Exts:      []string{".htm", ".js"},
		Archive:   "testsss.zip",
		Prefix:    "ssss",
		Disable:   true,
	}
	tests := []struct {
		name string
		opts []static.Option
		want *static.Static
	}{
		{name: "1. Conf-StaticNew-初始化nil对象", opts: nil, want: defaultObj},
		{name: "2. Conf-StaticNew-初始化空对象", opts: []static.Option{}, want: defaultObj},
		{name: "3. Conf-StaticNew-初始化image对象", opts: []static.Option{static.WithImages()}, want: defaultObj},
		{name: "4. Conf-StaticNew-初始化设置全量对象", opts: []static.Option{static.WithRoot("./test"), static.WithFirstPage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
			static.WithExts(".htm"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("ssss"), static.WithDisable(), static.WithExclude("/views/", ".exe", ".so", ".zip")},
			want: enObj},
	}
	for _, tt := range tests {
		got := static.New(tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestStatic_AllowRequest(t *testing.T) {
	tests := []struct {
		name   string
		fields *static.Static
		args   string
		want   bool
	}{
		{name: "1. Conf-StaticAllowRequest-Get支持的方法", fields: static.New(), args: http.MethodGet, want: true},
		{name: "2. Conf-StaticAllowRequest-Head支持的方法", fields: static.New(), args: http.MethodHead, want: true},
		{name: "3. Conf-StaticAllowRequest-Post不支持的方法", fields: static.New(), args: http.MethodPost, want: false},
		{name: "4. Conf-StaticAllowRequest-PUT不支持的方法", fields: static.New(), args: http.MethodPut, want: false},
		{name: "5. Conf-StaticAllowRequest-PATCH不支持的方法", fields: static.New(), args: http.MethodPatch, want: false},
		{name: "6. Conf-StaticAllowRequest-DELETE不支持的方法", fields: static.New(), args: http.MethodDelete, want: false},
		{name: "7. Conf-StaticAllowRequest-CONNECT不支持的方法", fields: static.New(), args: http.MethodConnect, want: false},
		{name: "8. Conf-StaticAllowRequest-OPTIONS不支持的方法", fields: static.New(), args: http.MethodOptions, want: false},
		{name: "9. Conf-StaticAllowRequest-TRACE不支持的方法", fields: static.New(), args: http.MethodTrace, want: false},
		{name: "10. Conf-StaticAllowRequest-other不支持的方法", fields: static.New(), args: "other", want: false},
	}
	for _, tt := range tests {
		got := tt.fields.AllowRequest(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestStaticGetConf(t *testing.T) {
	type test struct {
		name    string
		cnf     conf.IServerConf
		want    *static.Static
		wantErr bool
	}

	conf := mocks.NewConfBy("hydra", "graytest")
	confB := conf.API(":8090")
	test1 := test{name: "static节点不存在", cnf: conf.GetAPIConf().GetServerConf(), want: &static.Static{Disable: true, FileMap: map[string]static.FileInfo{}}, wantErr: false}
	staticObj, err := static.GetConf(test1.cnf)
	assert.Equal(t, test1.wantErr, (err != nil), test1.name+",err")
	assert.Equal(t, test1.want, staticObj, test1.name+",obj")

	confB.Static(static.WithArchive("车uowu"))
	test2 := test{name: "static节点存在,数据错误", cnf: conf.GetAPIConf().GetServerConf(), want: nil, wantErr: true}
	staticObj, err = static.GetConf(test2.cnf)
	assert.Equal(t, test2.wantErr, (err != nil), test2.name+",err")
	assert.Equal(t, test2.want, staticObj, test2.name+",obj")

	confB.Static(static.WithArchive("dddd"))
	test3 := test{name: "static节点存在,数据正确", cnf: conf.GetAPIConf().GetServerConf(), want: static.New(static.WithArchive("dddd")), wantErr: false}
	staticObj, err = static.GetConf(test3.cnf)
	assert.Equal(t, test3.wantErr, (err != nil), test3.name+",err")
	assert.Equal(t, test3.want, staticObj, test3.name+",obj")

	//处理归档文件
}

func TestStatic_IsFavRobot(t *testing.T) {
	tests := []struct {
		name   string
		fields *static.Static
		rPath  string
		wantB  bool
	}{
		{name: "1. Conf-StaticIsFavRobot-/favicon.ico文件", fields: static.New(), rPath: "/favicon.ico", wantB: true},
		{name: "2. Conf-StaticIsFavRobot-/robots.txt文件", fields: static.New(), rPath: "/robots.txt", wantB: true},
		{name: "3. Conf-StaticIsFavRobot-favicon.ico文件", fields: static.New(), rPath: "favicon.ico", wantB: false},
		{name: "4. Conf-StaticIsFavRobot-robots.txt文件", fields: static.New(), rPath: "robots.txt", wantB: false},
		{name: "5. Conf-StaticIsFavRobot-other文件", fields: static.New(), rPath: "test", wantB: false},
	}
	for _, tt := range tests {
		gotB := tt.fields.IsFavRobot(tt.rPath)
		assert.Equal(t, tt.wantB, gotB, tt.name)
	}
}

func TestStatic_IsExclude(t *testing.T) {
	tests := []struct {
		name   string
		fields *static.Static
		rPath  string
		want   bool
	}{
		{name: "1. Conf-StaticIsExclude-空Exclude对象", fields: static.New(), rPath: "/test", want: false},
		{name: "2. Conf-StaticIsExclude-Exclude对象,路径匹配成功", fields: static.New(static.WithExclude("/test")), rPath: "/test", want: true},
		{name: "3. Conf-StaticIsExclude-Exclude对象，扩展名匹配成功", fields: static.New(static.WithExclude(".so")), rPath: "/test1.so", want: true},
		{name: "4. Conf-StaticIsExclude-Exclude对象，路径匹配失败", fields: static.New(static.WithExclude("/test11")), rPath: "/test1", want: false},
		{name: "5. Conf-StaticIsExclude-Exclude对象，扩展名匹配失败", fields: static.New(static.WithExclude(".so")), rPath: "/test11.txt", want: false},
	}
	for _, tt := range tests {
		got := tt.fields.IsExclude(tt.rPath)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestStatic_HasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		fields *static.Static
		rPath  string
		want   bool
	}{
		{name: "1. Conf-StaticHasPrefix-空Prefix对象", fields: static.New(), rPath: "/test", want: false},
		{name: "2. Conf-StaticHasPrefix-Prefix对象，是前缀", fields: static.New(static.WithPrefix("/t")), rPath: "/test", want: true},
		{name: "3. Conf-StaticHasPrefix-Prefix对象，不是前缀", fields: static.New(static.WithPrefix("xxx")), rPath: "tatest", want: false},
		{name: "4. Conf-StaticHasPrefix-Prefix对象，与前缀相同", fields: static.New(static.WithPrefix("xxx")), rPath: "xxx", want: true},
	}
	for _, tt := range tests {
		got := tt.fields.HasPrefix(tt.rPath)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestStatic_IsContainExt(t *testing.T) {
	tests := []struct {
		name   string
		fields *static.Static
		rPath  string
		want   bool
	}{
		{name: "1. Conf-StaticIsContainExt-空exts对象，不限制，错误入参", fields: static.New(), rPath: "/test", want: false},
		{name: "2. Conf-StaticIsContainExt-空exts对象，不限制，正确入参", fields: static.New(), rPath: "/test.xx", want: true},
		{name: "3. Conf-StaticIsContainExt-错误的如参路径,没有扩展数据", fields: static.New(static.WithExts("*")), rPath: "test", want: false},
		{name: "4. Conf-StaticIsContainExt-错误的如参路径,*通配", fields: static.New(static.WithExts("*")), rPath: "test.xxx", want: true},
		{name: "5. Conf-StaticIsContainExt-错误的如参路径,指定扩展对象,匹配失败", fields: static.New(static.WithExts(".xx", ".tt")), rPath: "test.xxw", want: false},
		{name: "6. Conf-StaticIsContainExt-错误的如参路径,指定扩展对象，匹配成功", fields: static.New(static.WithExts(".xx", ".tt")), rPath: "test.xx", want: true},
	}
	for _, tt := range tests {
		got := tt.fields.IsContainExt(tt.rPath)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestStatic_NeedRewrite(t *testing.T) {
	tests := []struct {
		name   string
		fields *static.Static
		args   string
		want   bool
	}{
		{name: "1. Conf-StaticNeedRewrite-空Rewriters对象", fields: static.New(), args: "/test", want: false},
		{name: "2. Conf-StaticNeedRewrite-单个Rewriters对象", fields: static.New(static.WithRewriters("sasa")), args: "/test", want: false},
		{name: "3. Conf-StaticNeedRewrite-多个Rewriters对象", fields: static.New(static.WithRewriters("sasa", "sasaq")), args: "sasaq", want: true},
		{name: "4. Conf-StaticNeedRewrite-单个Rewriters对象,包含路径", fields: static.New(static.WithRewriters("sasa")), args: "sasaa", want: false},
		{name: "5. Conf-StaticNeedRewrite-单个Rewriters对象,被包含路径", fields: static.New(static.WithRewriters("sasa")), args: "sa", want: false},
	}
	for _, tt := range tests {
		got := tt.fields.NeedRewrite(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestStatic_GetGzFile(t *testing.T) {
	defaultObj := static.New()
	enObj := static.New()
	enObj.FileMap["t1"] = static.FileInfo{GzFile: "t1.gz", HasGz: true}
	enObj.FileMap["t2"] = static.FileInfo{GzFile: "t2.gz", HasGz: true}
	enObj.FileMap["t3"] = static.FileInfo{GzFile: "t3.txt", HasGz: false}

	tests := []struct {
		name   string
		fields *static.Static
		rPath  string
		want   string
	}{
		{name: "1. Conf-StaticGetGzFile-空对象获取", fields: defaultObj, rPath: "/test", want: ""},
		{name: "2. Conf-StaticGetGzFile-对象获取，失败", fields: enObj, rPath: "/test", want: ""},
		{name: "3. Conf-StaticGetGzFile-对象获取，成功", fields: enObj, rPath: "t2", want: "t2.gz"},
		{name: "4. Conf-StaticGetGzFile-对象获取，不是压缩文件", fields: enObj, rPath: "t3", want: ""},
	}
	for _, tt := range tests {
		got := tt.fields.GetGzFile(tt.rPath)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestStatic_IsStatic(t *testing.T) {
	enobj := static.New(static.WithRoot("./test"), static.WithFirstPage("index1.html"), static.WithRewriters("/", "indextest.htm", "defaulttest.html"),
		static.WithExts(".html"), static.WithArchive("testsss"), static.AppendExts(".js"), static.WithPrefix("/ssss"), static.WithExclude("/views/", ".exe", ".so", ".zip"))
	disableObj := static.New(static.WithDisable())
	type args struct {
		rPath  string
		method string
	}
	tests := []struct {
		name      string
		fields    *static.Static
		args      args
		want      bool
		wantXname string
	}{
		{name: "1. Conf-StaticIsStatic-是IsFavRobot", fields: disableObj, args: args{rPath: "/favicon.ico", method: ""}, want: true, wantXname: filepath.Join(disableObj.Dir, "/favicon.ico")},
		{name: "2. Conf-StaticIsStatic-是Disable", fields: disableObj, args: args{rPath: "/sdsdsd.ico", method: ""}, want: false, wantXname: ""},
		{name: "3. Conf-StaticIsStatic-不允许的请求方式", fields: enobj, args: args{rPath: "/sdsdsd.ico", method: "POST"}, want: false, wantXname: ""},
		{name: "4. Conf-StaticIsStatic-是排除路径", fields: enobj, args: args{rPath: "/views/dd", method: "GET"}, want: false, wantXname: ""},
		{name: "5. Conf-StaticIsStatic-是允许的扩展文件", fields: enobj, args: args{rPath: "/ttt/dd.html", method: "GET"}, want: true, wantXname: filepath.Join(enobj.Dir, "/ttt/dd.html")},
		{name: "6. Conf-StaticIsStatic-不允许的扩展文件，但是是指定前缀", fields: enobj, args: args{rPath: "/ssss/dd.xx", method: "GET"}, want: true, wantXname: filepath.Join(enobj.Dir, strings.TrimPrefix("/ssss/dd.xx", enobj.Prefix))},
		{name: "7. Conf-StaticIsStatic-不允许的扩展文件且不是指定前缀，需要转发", fields: enobj, args: args{rPath: "indextest.htm", method: "GET"}, want: true, wantXname: filepath.Join(enobj.Dir, enobj.FirstPage)},
		{name: "8. Conf-StaticIsStatic-所有条件都不满足", fields: enobj, args: args{rPath: "test.htm", method: "GET"}, want: false, wantXname: ""},
	}
	for _, tt := range tests {
		gotB, gotXname := tt.fields.IsStatic(tt.args.rPath, tt.args.method)
		assert.Equal(t, tt.want, gotB, tt.name+",bool")
		assert.Equal(t, tt.wantXname, gotXname, tt.name+",xname")
	}
}

func TestStatic_RereshData(t *testing.T) {
	tests := []struct {
		name   string
		fields *static.Static
	}{}
	for _, tt := range tests {
		tt.fields.RereshData()
	}
}
