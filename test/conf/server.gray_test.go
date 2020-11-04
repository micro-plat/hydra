package conf

import (
	"net/url"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/cron"
	"github.com/micro-plat/hydra/conf/server/gray"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestGrayNew(t *testing.T) {
	tests := []struct {
		name string
		opts []gray.Option
		want *gray.Gray
	}{
		{name: "初始化默认对象", opts: []gray.Option{}, want: &gray.Gray{Disable: false}},
		{name: "设置disable对象", opts: []gray.Option{gray.WithDisable()}, want: &gray.Gray{Disable: true}},
		{name: "设置enable对象", opts: []gray.Option{gray.WithEnable()}, want: &gray.Gray{Disable: false}},
		{name: "设置全量对象", opts: []gray.Option{gray.WithDisable(), gray.WithFilter("tao"), gray.WithUPCluster("dd")},
			want: &gray.Gray{Disable: true, Filter: "tao", UPCluster: "dd"}},
	}
	for _, tt := range tests {
		got := gray.New(tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestGray_Check(t *testing.T) {
	var f1 func(name string) string = func(name string) string {
		return name
	}
	var f2 func() string = func() string {
		return time.Now().Format("20060102150405")
	}

	type args struct {
		funcs map[string]interface{}
		i     interface{}
	}
	tests := []struct {
		name    string
		opts    []gray.Option
		args    args
		want    bool
		wantErr bool
	}{
		{name: "需要灰度", opts: []gray.Option{gray.WithFilter("true"), gray.WithUPCluster("UPCluster")},
			args: args{funcs: map[string]interface{}{"getString": f1, "getTime": f2}, i: nil}, want: true, wantErr: false},
		{name: "不需要灰度", opts: []gray.Option{gray.WithFilter("flase"), gray.WithUPCluster("UPCluster")},
			args: args{funcs: map[string]interface{}{"getString": f1, "getTime": f2}, i: nil}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		g := gray.New(tt.opts...)
		got, err := g.Check(tt.args.funcs, tt.args.i)
		assert.Equal(t, tt.wantErr, (err != nil), tt.name+",err")
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestGrayGetConf(t *testing.T) {
	type test struct {
		name    string
		cnf     conf.IServerConf
		want    *gray.Gray
		wantErr bool
	}

	conf := mocks.NewConfBy("hydra", "graytest")
	confB := conf.API(":8090")
	test1 := test{name: "灰度节点不存在", cnf: conf.GetAPIConf().GetServerConf(), want: &gray.Gray{Disable: true}, wantErr: false}
	grayObj, err := gray.GetConf(test1.cnf)
	assert.Equal(t, test1.wantErr, (err != nil), test1.name)
	assert.Equal(t, test1.want, grayObj, test1.name)

	confB.Gray(gray.WithDisable(), gray.WithUPCluster("graytest"))
	test2 := test{name: "灰度节点存在,filter不存在", cnf: conf.GetAPIConf().GetServerConf(), want: nil, wantErr: true}
	grayObj, err = gray.GetConf(test2.cnf)
	assert.Equal(t, test2.wantErr, (err != nil), test2.name+",err")
	assert.Equal(t, test2.want, grayObj, test2.name)

	confB.Gray(gray.WithDisable(), gray.WithFilter("tao"))
	test3 := test{name: "灰度节点存在,UPCluster不存在", cnf: conf.GetAPIConf().GetServerConf(), want: nil, wantErr: true}
	grayObj, err = gray.GetConf(test3.cnf)
	assert.Equal(t, test3.wantErr, (err != nil), test3.name+",err")
	assert.Equal(t, test3.want, grayObj, test3.name)

	confB.Gray(gray.WithDisable(), gray.WithFilter("tao"), gray.WithUPCluster("graytest"))
	test4 := test{name: "灰度节点存在", cnf: conf.GetAPIConf().GetServerConf(), want: &gray.Gray{Disable: true, Filter: "tao", UPCluster: "graytest"}, wantErr: false}
	grayObj, err = gray.GetConf(test4.cnf)
	assert.Equal(t, test4.wantErr, (err != nil), test4.name+",err")
	assert.Equal(t, test4.want.Disable, grayObj.Disable, test4.name+",Disable")
	assert.Equal(t, test4.want.Filter, grayObj.Filter, test4.name+",Filter")
	assert.Equal(t, test4.want.UPCluster, grayObj.UPCluster, test4.name+",UPCluster")
}

func TestGray_Allow(t *testing.T) {
	type test struct {
		name   string
		fields *gray.Gray
		want   bool
	}

	test1 := test{name: "无集群对象获取", fields: gray.New(gray.WithFilter("true")), want: false}
	got := test1.fields.Allow()
	assert.Equal(t, test1.want, got, test1.name)

	conf := mocks.NewConfBy("test", "gray1")
	conf.API(":8090")
	grayObj, _ := gray.GetConf(conf.GetAPIConf().GetServerConf())
	test2 := test{name: "api服务集群获取", fields: grayObj, want: false}
	got = test2.fields.Allow()
	assert.Equal(t, test2.want, got, test2.name)

	conf = mocks.NewConfBy("test", "gray2")
	conf.RPC(":8090")
	grayObj, _ = gray.GetConf(conf.GetRPCConf().GetServerConf())
	test3 := test{name: "rpc服务集群获取", fields: grayObj, want: false}
	got = test3.fields.Allow()
	assert.Equal(t, test3.want, got, test3.name)

	conf = mocks.NewConfBy("test", "gray3")
	conf.MQC("redis://192.168.0.1")
	grayObj, _ = gray.GetConf(conf.GetMQCConf().GetServerConf())
	test4 := test{name: "mqc服务集群获取", fields: grayObj, want: false}
	got = test4.fields.Allow()
	assert.Equal(t, test4.want, got, test4.name)

	conf = mocks.NewConfBy("test", "gray4")
	conf.CRON(cron.WithTrace())
	grayObj, _ = gray.GetConf(conf.GetCronConf().GetServerConf())
	test5 := test{name: "cron服务集群获取", fields: grayObj, want: false}
	got = test5.fields.Allow()
	assert.Equal(t, test5.want, got, test5.name)
}

func TestGray_Next(t *testing.T) {
	type test struct {
		name    string
		fields  *gray.Gray
		wantU   *url.URL
		wantErr bool
	}

	test1 := test{name: "nil集群对象", fields: gray.New(gray.WithFilter("true")), wantU: nil, wantErr: true}
	gotU, err := test1.fields.Next()
	assert.Equal(t, test1.wantErr, (err != nil), test1.name+",err")
	assert.Equal(t, test1.wantU, gotU, test1.name+",url")

	conf := mocks.NewConfBy("hydra", "graytest")
	conf.API(":8090").Gray(gray.WithFilter("tao"), gray.WithUPCluster("graytest"))
	nomalObj, _ := gray.GetConf(conf.GetAPIConf().GetServerConf())
	test2 := test{name: "无服务器集群对象", fields: nomalObj, wantU: nil, wantErr: true}
	gotU, err = test2.fields.Next()
	assert.Equal(t, test2.wantErr, (err != nil), test2.name+",err")
	assert.Equal(t, test2.wantU, gotU, test2.name+",url")

	path := conf.GetAPIConf().GetServerConf().GetServerPubPath("graytest")
	conf.Registry.CreateTempNode(path+":123456", "错误的服务器地址")
	time.Sleep(2 * time.Second)
	test3 := test{name: "错误配置服务器集群对象", fields: nomalObj, wantU: nil, wantErr: true}
	gotU, err = test3.fields.Next()
	assert.Equal(t, test3.wantErr, (err != nil), test3.name+",err")
	assert.Equal(t, test3.wantU, gotU, test3.name+",url")

	url, _ := url.Parse("http://192.168.5.94:8090")
	conf.Registry.CreateTempNode(path+":123456", "http://192.168.5.94:8090")
	time.Sleep(2 * time.Second)
	test4 := test{name: "正确配置服务器集群对象", fields: nomalObj, wantU: url, wantErr: true}
	gotU, err = test4.fields.Next()
	assert.Equal(t, test4.wantErr, (err != nil), test4.name+",err")
	assert.Equal(t, test4.wantU, gotU, test4.name+",url")
}
