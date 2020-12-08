package rpc

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func TestResolvePath(t *testing.T) {
	type args struct {
		address     string
		defPlatName string
	}
	tests := []struct {
		name         string
		args         args
		wantIsip     bool
		wantService  string
		wantPlatName string
		wantErr      bool
	}{
		{name: "1. 正常-通过ip访问", args: args{address: "/path1/path2@tcp://192.168.0.1:8888", defPlatName: ""}, wantIsip: true, wantService: "/path1/path2", wantPlatName: "tcp://192.168.0.1:8888", wantErr: false},
		{name: "2. 正常-新版-通过服务路径访问", args: args{address: "/path1/path2@platname", defPlatName: ""}, wantIsip: false, wantService: "/path1/path2", wantPlatName: "platname", wantErr: false},
		{name: "3. 正常-旧版-通过服务路径访问", args: args{address: "/path1/path2@servername.platname", defPlatName: ""}, wantIsip: false, wantService: "/path1/path2", wantPlatName: "servername.platname", wantErr: false},
		{name: "4. 正常-包含默认platname", args: args{address: "/path1", defPlatName: "platname"}, wantIsip: false, wantService: "/path1", wantPlatName: "platname", wantErr: false},
		{name: "4. 异常-服务地址为空-1", args: args{address: "", defPlatName: ""}, wantIsip: false, wantService: "", wantPlatName: "", wantErr: true},
		{name: "5. 异常-服务地址为空-2", args: args{address: "@platname", defPlatName: ""}, wantIsip: false, wantService: "", wantPlatName: "", wantErr: true},
	}
	for _, tt := range tests {
		gotIsip, gotService, gotPlatName, err := ResolvePath(tt.args.address, tt.args.defPlatName)
		assert.Equalf(t, tt.wantErr, err != nil, tt.name+";err:%w", err)
		assert.Equalf(t, tt.wantIsip, gotIsip, tt.name+";Isip,expect:%v;got:%v", tt.wantIsip, gotIsip)
		assert.Equalf(t, tt.wantService, gotService, tt.name+";gotService,expect:%v;got:%v", tt.wantService, gotService)
		assert.Equalf(t, tt.wantPlatName, gotPlatName, tt.name+";gotPlatName,expect:%v;got:%v", tt.wantPlatName, gotPlatName)
	}
}
