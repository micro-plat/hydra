package global

import (
	"strings"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
)

func Test_global_GetLongAppName(t *testing.T) {

	type args struct {
		n []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "1. 测试未传入名称-nil", args: args{n: nil}, want: "global.test"},
		{name: "2. 测试未传入名称-空数组", args: args{n: []string{}}, want: "global.test"},
		{name: "3. 测试有传入名称-短", args: args{n: []string{"name"}}, want: "name"},
		{name: "4. 测试有传入名称-超过32", args: args{n: []string{"name123456789012345678901234567890"}}, want: "name1234567890123456789"},
	}
	for _, tt := range tests {
		m := &global{}
		longName := m.GetLongAppName(tt.args.n...)
		assert.Equal(t, true, strings.Contains(longName, tt.want), tt.name)
	}
}

func Test_global_HasServerType(t *testing.T) {
	type fields struct {
		ServerTypes []string
	}
	type args struct {
		tp string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{

		{name: "1. 单个ServerTypes匹配", fields: fields{ServerTypes: []string{"api"}}, args: args{tp: "api"}, want: true},
		{name: "2. 单个ServerTypes不匹配", fields: fields{ServerTypes: []string{"api"}}, args: args{tp: "xapi"}, want: false},
		{name: "3. 多个ServerTypes匹配", fields: fields{ServerTypes: []string{"api", "cron"}}, args: args{tp: "api"}, want: true},
		{name: "4. 多个ServerTypes不匹配", fields: fields{ServerTypes: []string{"api", "cron"}}, args: args{tp: "xapi"}, want: false},
	}
	for _, tt := range tests {
		m := &global{
			ServerTypes: tt.fields.ServerTypes,
		}
		assert.Equal(t, tt.want, m.HasServerType(tt.args.tp), tt.name)
	}
}

func Test_parsePath(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name            string
		args            args
		wantPlatName    string
		wantSystemName  string
		wantServerTypes []string
		wantClusterName string
		wantErr         bool
	}{
		{name: "1. 正常格式-单个serverType", args: args{p: "/platName/sysName/serverType/clusterName"}, wantPlatName: "platName", wantSystemName: "sysName", wantServerTypes: []string{"serverType"}, wantClusterName: "clusterName", wantErr: false},
		{name: "2. 正常格式-多个serverType", args: args{p: "/platName/sysName/api-cron-mqc/clusterName"}, wantPlatName: "platName", wantSystemName: "sysName", wantServerTypes: []string{"api", "cron", "mqc"}, wantClusterName: "clusterName", wantErr: false},
		{name: "3. 正常格式-首尾多/", args: args{p: "/platName/sysName/serverType/clusterName/"}, wantPlatName: "platName", wantSystemName: "sysName", wantServerTypes: []string{"serverType"}, wantClusterName: "clusterName", wantErr: false},
		{name: "4. 正常格式-多个serverType-首尾多/", args: args{p: "/platName/sysName/api-cron-mqc/clusterName"}, wantPlatName: "platName", wantSystemName: "sysName", wantServerTypes: []string{"api", "cron", "mqc"}, wantClusterName: "clusterName", wantErr: false},
		{name: "5. 错误格式-分段不足", args: args{p: "/platName/sysName/api-cron-mqc"}, wantPlatName: "", wantSystemName: "", wantServerTypes: nil, wantClusterName: "", wantErr: true},
		{name: "6. 错误格式-分段超多", args: args{p: "/platName/sysName/api-cron-mqc/clusterName/xxx"}, wantPlatName: "", wantSystemName: "", wantServerTypes: nil, wantClusterName: "", wantErr: true},
	}
	for _, tt := range tests {

		gotPlatName, gotSystemName, gotServerTypes, gotClusterName, err := parsePath(tt.args.p)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.Equalf(t, tt.wantPlatName, gotPlatName, tt.name+"%s", "PlatName")
		assert.Equalf(t, tt.wantSystemName, gotSystemName, tt.name+"%s", "SystemName")
		assert.Equalf(t, tt.wantServerTypes, gotServerTypes, tt.name+"%s", "ServerTypes")
		assert.Equalf(t, tt.wantClusterName, gotClusterName, tt.name+"%s", "ClusterName")
	}
}

func Test_global_check(t *testing.T) {
	type fields struct {
		RegistryAddr    string
		PlatName        string
		SysName         string
		ServerTypes     []string
		ServerTypeNames string
		ClusterName     string
		Trace           string
		Name            string
		IsDebug         bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errMsg  string
		expect  CliFlagObject
		flagVal CliFlagObject
	}{

		{name: "1.Registry-不设置", fields: fields{PlatName: "p", ServerTypeNames: "api"}, wantErr: false, expect: CliFlagObject{RegistryAddr: "lm://."}, flagVal: CliFlagObject{}},
		{name: "2.Registry-设置", fields: fields{PlatName: "p", ServerTypeNames: "api"}, wantErr: false, expect: CliFlagObject{RegistryAddr: "zk://192.168.0.1"}, flagVal: CliFlagObject{RegistryAddr: "zk://192.168.0.1"}},

		{name: "3.PlatName-不设置", fields: fields{ServerTypeNames: "api"}, wantErr: true, errMsg: "平台名称不能为空", expect: CliFlagObject{PlatName: ""}, flagVal: CliFlagObject{}},
		{name: "4.PlatName-设置", fields: fields{ServerTypeNames: "api"}, wantErr: false, expect: CliFlagObject{PlatName: "PlatName"}, flagVal: CliFlagObject{PlatName: "PlatName"}},

		{name: "5.SysName-设置", fields: fields{PlatName: "p", ServerTypeNames: "api"}, wantErr: false, expect: CliFlagObject{SysName: "SysName"}, flagVal: CliFlagObject{SysName: "SysName"}},
		{name: "6.SysName-不为空", fields: fields{PlatName: "p", ServerTypeNames: "api"}, wantErr: false, expect: CliFlagObject{SysName: ""}, flagVal: CliFlagObject{SysName: ""}},

		{name: "7.ServerTypeNames-单个", fields: fields{PlatName: "p"}, wantErr: false, expect: CliFlagObject{ServerTypeNames: "api"}, flagVal: CliFlagObject{ServerTypeNames: "api"}},
		{name: "8.ServerTypeNames-多个", fields: fields{PlatName: "p"}, wantErr: false, expect: CliFlagObject{ServerTypeNames: "api-mqc"}, flagVal: CliFlagObject{ServerTypeNames: "api-mqc"}},
		{name: "9.ServerTypeNames-存在不包含", fields: fields{PlatName: "p"}, wantErr: false, expect: CliFlagObject{}, flagVal: CliFlagObject{ServerTypeNames: "api-mqc"}},

		{name: "10.Name-为空", fields: fields{Name: "", RegistryAddr: "RegistryAddr", PlatName: "PlatName", SysName: "SysName", ServerTypeNames: "api-mqc", ClusterName: "ClusterName"}, wantErr: false, expect: CliFlagObject{}, flagVal: CliFlagObject{RegistryAddr: "", Name: "", PlatName: "", SysName: "", ServerTypeNames: "", ClusterName: ""}},
		{name: "11.Name-不为空", fields: fields{Name: "/PlatName/SysName/api-mqc/ClusterName"}, wantErr: false, expect: CliFlagObject{RegistryAddr: "zk://192.168.0.1"}, flagVal: CliFlagObject{RegistryAddr: "zk://192.168.0.1"}},

		{name: "12.Trace不为空-在trace列表内", fields: fields{Trace: "cpu"}, wantErr: false, expect: CliFlagObject{RegistryAddr: "zk://192.168.0.1", PlatName: "PlatName", SysName: "SysName", ServerTypeNames: "api", ClusterName: "ClusterName"}, flagVal: CliFlagObject{RegistryAddr: "zk://192.168.0.1", PlatName: "PlatName", SysName: "SysName", ServerTypeNames: "api", ClusterName: "ClusterName"}},
		{name: "13.Trace不为空-不在trace列表内", fields: fields{Trace: "xpu"}, wantErr: true, errMsg: "trace名称只能是[cpu mem block mutex web]", expect: CliFlagObject{RegistryAddr: "zk://192.168.0.1", PlatName: "PlatName", SysName: "SysName", ServerTypeNames: "api", ClusterName: "ClusterName"}, flagVal: CliFlagObject{RegistryAddr: "zk://192.168.0.1", PlatName: "PlatName", SysName: "SysName", ServerTypeNames: "api", ClusterName: "ClusterName"}},

		{name: "14.IsDebug:true", fields: fields{IsDebug: true}, wantErr: false, expect: CliFlagObject{RegistryAddr: "zk://192.168.0.1", PlatName: "PlatName_debug", SysName: "SysName", ServerTypeNames: "api", ClusterName: "ClusterName"}, flagVal: CliFlagObject{RegistryAddr: "zk://192.168.0.1", PlatName: "PlatName", SysName: "SysName", ServerTypeNames: "api", ClusterName: "ClusterName"}},
		{name: "15.IsDebug:false", fields: fields{IsDebug: false}, wantErr: false, expect: CliFlagObject{RegistryAddr: "zk://192.168.0.1", PlatName: "PlatName", SysName: "SysName", ServerTypeNames: "api", ClusterName: "ClusterName"}, flagVal: CliFlagObject{RegistryAddr: "zk://192.168.0.1", PlatName: "PlatName", SysName: "SysName", ServerTypeNames: "api", ClusterName: "ClusterName"}},
	}

	ServerTypes = []string{"api", "mqc"}
	for _, tt := range tests {
		m := &global{
			RegistryAddr:    tt.fields.RegistryAddr,
			PlatName:        tt.fields.PlatName,
			SysName:         tt.fields.SysName,
			ServerTypes:     tt.fields.ServerTypes,
			ServerTypeNames: tt.fields.ServerTypeNames,
			ClusterName:     tt.fields.ClusterName,
			Name:            tt.fields.Name,
			Trace:           tt.fields.Trace,
		}
		//初始化测试用例参数

		FlagVal.RegistryAddr = tt.flagVal.RegistryAddr
		FlagVal.PlatName = tt.flagVal.PlatName
		FlagVal.SysName = tt.flagVal.SysName
		FlagVal.ServerTypeNames = tt.flagVal.ServerTypeNames
		FlagVal.ClusterName = tt.flagVal.ClusterName
		IsDebug = tt.fields.IsDebug

		//执行被测试方法
		err := m.check()

		assert.Equalf(t, tt.wantErr, err != nil, "%s:错误与预期不符 expect:%v,got:%v", tt.name, tt.wantErr, err)
		if tt.wantErr {
			assert.Equalf(t, tt.errMsg, err.Error(), "%s", tt.name)
		}

		if tt.expect.RegistryAddr != "" {
			assert.Equalf(t, tt.expect.RegistryAddr, m.GetRegistryAddr(), "%s:RegistryAddr与预期不匹配：expect:%s,got:%s", tt.name, tt.expect.RegistryAddr, m.GetRegistryAddr())
		}
		if tt.expect.PlatName != "" {
			assert.Equalf(t, tt.expect.PlatName, m.GetPlatName(), "%s:PlatName与预期不匹配：expect:%s,got:%s", tt.name, tt.expect.PlatName, m.GetPlatName())
		}
		if tt.expect.SysName != "" {
			assert.Equalf(t, tt.expect.SysName, m.GetSysName(), "%s:SysName与预期不匹配：expect:%s,got:%s", tt.name, tt.expect.SysName, m.GetSysName())
		}
		if tt.expect.ServerTypeNames != "" {
			assert.Equalf(t, tt.expect.ServerTypeNames, strings.Join(m.GetServerTypes(), "-"), "%s:ServerTypeNames与预期不匹配：expect:%s,got:%s", tt.name, tt.expect.ServerTypeNames, m.GetServerTypes())
		}
		if tt.expect.ClusterName != "" {
			assert.Equalf(t, tt.expect.ClusterName, m.GetClusterName(), "%s:ClusterName与预期不匹配：expect:%s,got:%s", tt.name, tt.expect.ClusterName, m.GetClusterName())
		}
	}
}
