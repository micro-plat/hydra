package registry

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context/ctx"
	"github.com/micro-plat/hydra/registry"

	//_ "github.com/micro-plat/hydra/registry/registry/etcd"
	_ "github.com/micro-plat/hydra/registry/registry/filesystem"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	_ "github.com/micro-plat/hydra/registry/registry/redis"
	_ "github.com/micro-plat/hydra/registry/registry/zookeeper"
	"github.com/micro-plat/lib4go/logger"

	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

// 私有变量的测试
// func Test_getAddrByUserPass(t *testing.T) {
// 	type args struct {
// 		addr string
// 	}
// 	tests := []struct {
// 		name        string
// 		args        args
// 		wantU       string
// 		wantP       string
// 		wantAddress string
// 		wantErr     bool
// 	}{
// 		{name: "正确格式的地址", args: args{addr: "root:123456@192.168.5.115:9091"}, wantU: "root", wantP: "123456", wantAddress: "192.168.5.115:9091", wantErr: false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotU, gotP, gotAddress, err := getAddrByUserPass(tt.args.addr)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("getAddrByUserPass() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotU != tt.wantU {
// 				t.Errorf("getAddrByUserPass() gotU = %v, want %v", gotU, tt.wantU)
// 			}
// 			if gotP != tt.wantP {
// 				t.Errorf("getAddrByUserPass() gotP = %v, want %v", gotP, tt.wantP)
// 			}
// 			if gotAddress != tt.wantAddress {
// 				t.Errorf("getAddrByUserPass() gotAddress = %v, want %v", gotAddress, tt.wantAddress)
// 			}
// 		})
// 	}
// }

func TestParse(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name      string
		args      args
		wantProto string
		wantRaddr []string
		wantU     string
		wantP     string
		wantErr   bool
	}{
		{name: "1. Parse-地址带有多个://", args: args{address: "zk://192.://168.0.101"}, wantU: "", wantP: "", wantErr: true},
		{name: "2. Parse-解析zk地址", args: args{address: "zk://192.168.0.101"}, wantProto: "zk", wantRaddr: []string{"192.168.0.101"}, wantU: "", wantP: "", wantErr: false},
		{name: "3. Parse-解析zk地址", args: args{address: "zk://192.168.0.101"}, wantProto: "zk", wantRaddr: []string{"192.168.0.101"}, wantU: "", wantP: "", wantErr: false},
		{name: "4. Parse-解析多个zk地址", args: args{address: "zk://192.168.0.101,192.168.0.102"}, wantProto: "zk", wantRaddr: []string{"192.168.0.101", "192.168.0.102"}, wantU: "", wantP: "", wantErr: false},
		{name: "5. Parse-解析lm地址", args: args{address: "lm://."}, wantProto: "lm", wantRaddr: []string{"."}, wantU: "", wantP: "", wantErr: false},
		{name: "6. Parse-解析fs地址", args: args{address: "fs://../a/b/c"}, wantProto: "fs", wantRaddr: []string{"../a/b/c"}, wantU: "", wantP: "", wantErr: false},
		{name: "7. Parse-解析etcd地址", args: args{address: "etcd://192.168.0.101:9099"}, wantProto: "etcd", wantRaddr: []string{"192.168.0.101:9099"}, wantU: "", wantP: "", wantErr: false},
		{name: "8. Parse-解析redis地址", args: args{address: "redis://192.168.0.101:6379"}, wantProto: "redis", wantRaddr: []string{"192.168.0.101:6379"}, wantU: "", wantP: "", wantErr: false},
		{name: "9. Parse-解析带有用户名和密码的地址", args: args{address: "redis://root:123456@192.168.0.101:6379"}, wantProto: "redis", wantRaddr: []string{"192.168.0.101:6379"}, wantU: "root", wantP: "123456", wantErr: false},
	}
	for _, tt := range tests {
		gotProto, gotRaddr, gotU, gotP, _, err := registry.Parse(tt.args.address)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			continue
		}
		assert.Equal(t, tt.wantProto, gotProto, tt.name)
		assert.Equal(t, tt.wantRaddr, gotRaddr, tt.name)
		assert.Equal(t, tt.wantU, gotU, tt.name)
		assert.Equal(t, tt.wantP, gotP, tt.name)
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		elem []string
		want string
	}{
		{name: "1. Join-参数的以/开头", elem: []string{"/", "path", "\\"}, want: "/path"},
		{name: "2. Join-参数的以/结尾", elem: []string{"/", "path", "/"}, want: "/path"},
		{name: "3. Join-参数的以/结尾,最后参数为空", elem: []string{"/", "path/", ""}, want: "/path"},
		{name: "4. Join-参数的以转义符开头", elem: []string{"\\", "/!@#$%^&*()_+{}:><?, }}", "/dsd"}, want: "/!@#$%^&*()_+{}:><?, }}/dsd"},
		{name: "5. Join-参数均为空,地址拼接", elem: []string{"", ""}, want: ""},
		{name: "6. Join-参数第一个为空,地址拼接", elem: []string{"", "a/"}, want: "/a"},
		{name: "7. Join-参数带有特殊符号,地址拼接", elem: []string{"a", "b", "!@#$%^&*c"}, want: "/a/b/!@#$%^&*c"},
		{name: "8. Join-参数带有相对地址,地址拼接", elem: []string{"..", "a/b", "c/"}, want: "/../a/b/c"},
		{name: "9. Join-参数带有转义符,地址拼接", elem: []string{"..", "", "\\", "c/"}, want: "/../c"},
	}
	for _, tt := range tests {
		got := registry.Join(tt.elem...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestNewRegistry(t *testing.T) {

	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     string
	}{
		{name: "1. NewRegistry-获取不支持的注册中心", args: args{address: "cloud://../"}, wantErr: true, err: "不支持的协议类型[cloud]"},
		{name: "2. NewRegistry-获取zk的注册中心", args: args{address: "zk://192.168.0.101"}, wantErr: false},
		{name: "3. NewRegistry-获取lm的注册中心", args: args{address: "lm://."}, wantErr: false},
		{name: "4. NewRegistry-获取fs的注册中心", args: args{address: "fs://../"}, wantErr: false},
	}

	confObj := mocks.NewConfBy("hydra_rgst_test", "rgsttest") //构建对象
	confObj.API(":8080")                                      //初始化参数
	serverConf := confObj.GetAPIConf()                        //获取配置
	meta := conf.NewMeta()
	log := logger.GetSession(serverConf.GetServerConf().GetServerName(), ctx.NewUser(&mocks.TestContxt{}, "", meta).GetRequestID())

	for _, tt := range tests {
		gotR, err := registry.GetRegistry(tt.args.address, log)
		if tt.wantErr {
			assert.Equal(t, tt.err, err.Error(), tt.name)
		}
		if !tt.wantErr {
			assert.IsNil(t, false, gotR, tt.name)
		}
	}
}
