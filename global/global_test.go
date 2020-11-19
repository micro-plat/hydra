package global

import (
	"fmt"
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
		{
			name: "测试未传入名称-nil",
			args: args{
				n: nil,
			},
			want: "global.test_04f81c71",
		},
		{
			name: "测试未传入名称-空数组",
			args: args{
				n: []string{},
			},
			want: "global.test_04f81c71",
		},
		{
			name: "测试有传入名称-短",
			args: args{
				n: []string{"name"},
			},
			want: "name_6a848303",
		},
		{
			name: "测试有传入名称-超过32",
			args: args{
				n: []string{"name123456789012345678901234567890"},
			},
			want: "name1234567890123456789_a0e7045b",
		},
	}
	for _, tt := range tests {
		m := &global{}
		assert.Equal(t, tt.want, m.GetLongAppName(tt.args.n...), tt.name)
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

		{
			name: "单个ServerTypes匹配",
			fields: fields{
				ServerTypes: []string{"api"},
			},
			args: args{tp: "api"},
			want: true,
		},
		{
			name: "单个ServerTypes不匹配",
			fields: fields{
				ServerTypes: []string{"api"},
			},
			args: args{tp: "xapi"},
			want: false,
		},
		{
			name: "多个ServerTypes不匹配",
			fields: fields{
				ServerTypes: []string{"api", "cron"},
			},
			args: args{tp: "xapi"},
			want: false,
		}, {
			name: "多个ServerTypes匹配",
			fields: fields{
				ServerTypes: []string{"api", "cron"},
			},
			args: args{tp: "api"},
			want: true,
		},
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
		{
			name: "正常格式-单个serverType",
			args: args{
				p: "/platName/sysName/serverType/clusterName",
			},
			wantPlatName:    "platName",
			wantSystemName:  "sysName",
			wantServerTypes: []string{"serverType"},
			wantClusterName: "clusterName",
			wantErr:         false,
		},
		{
			name: "正常格式-多个serverType",
			args: args{
				p: "/platName/sysName/api-cron-mqc/clusterName",
			},
			wantPlatName:    "platName",
			wantSystemName:  "sysName",
			wantServerTypes: []string{"api", "cron", "mqc"},
			wantClusterName: "clusterName",
			wantErr:         false,
		},
		{
			name: "正常格式-首尾多/",
			args: args{
				p: "/platName/sysName/serverType/clusterName/",
			},
			wantPlatName:    "platName",
			wantSystemName:  "sysName",
			wantServerTypes: []string{"serverType"},
			wantClusterName: "clusterName",
			wantErr:         false,
		},
		{
			name: "正常格式-多个serverType-首尾多/",
			args: args{
				p: "/platName/sysName/api-cron-mqc/clusterName",
			},
			wantPlatName:    "platName",
			wantSystemName:  "sysName",
			wantServerTypes: []string{"api", "cron", "mqc"},
			wantClusterName: "clusterName",
			wantErr:         false,
		},
		{
			name: "错误格式-分段不足",
			args: args{
				p: "/platName/sysName/api-cron-mqc",
			},
			wantPlatName:    "",
			wantSystemName:  "",
			wantServerTypes: nil,
			wantClusterName: "",
			wantErr:         true,
		},
		{
			name: "错误格式-分段超多",
			args: args{
				p: "/platName/sysName/api-cron-mqc/clusterName/xxx",
			},
			wantPlatName:    "",
			wantSystemName:  "",
			wantServerTypes: nil,
			wantClusterName: "",
			wantErr:         true,
		},

		// TODO: Add test cases.
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

type Assert func(t *testing.T, m *global, err error)
type Init func()

func Test_global_check(t *testing.T) {
	type fields struct {
		RegistryAddr    string
		PlatName        string
		SysName         string
		ServerTypes     []string
		ServerTypeNames string
		ClusterName     string
		Name            string
		Trace           string
		LocalConfName   string
	}
	tests := []struct {
		name         string
		fields       fields
		wantErr      bool
		errMsg       string
		assertResult Assert
		initParams   Init
	}{

		{
			name:    "RegistryAddr",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.Equalf(t, m.RegistryAddr, FlagVal.RegistryAddr, "RegistryAddr与指定值不相等，expect:%s,actual:%s", FlagVal.RegistryAddr, m.RegistryAddr)
			},
			initParams: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
			},
		},
		{
			name:    "PlatName",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.Equalf(t, m.PlatName, FlagVal.PlatName, "PlatName与指定值不相等，expect:%s,actual:%s", FlagVal.PlatName, m.PlatName)
			},
			initParams: func() {
				FlagVal.PlatName = "PlatName"
			},
		},
		{
			name:    "SysName-不为空",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.Equalf(t, m.SysName, FlagVal.SysName, "SysName与指定值不相等，expect:%s,actual:%s", FlagVal.SysName, m.SysName)
			},
			initParams: func() {
				FlagVal.SysName = "SysName"
			},
		},
		{
			name:    "SysName-为空",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.Equalf(t, m.SysName, AppName, "SysName与指定值不相等，expect:%s,actual:%s", AppName, m.SysName)
			},
			initParams: func() {
				FlagVal.SysName = ""
			},
		},
		{
			name:    "ServerTypeNames-单个",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.Equalf(t, m.ServerTypeNames, FlagVal.ServerTypeNames, "ServerTypeNames与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypeNames)
				assert.Equalf(t, 1, len(m.ServerTypes), "ServerTypes与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypes[0])
				assert.Equalf(t, "api", m.ServerTypes[0], "ServerTypes与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypes[0])
			},
			initParams: func() {
				ServerTypes = []string{"api"}
				FlagVal.ServerTypeNames = "api"
			},
		},
		{
			name:    "ServerTypeNames-多个",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.Equalf(t, m.ServerTypeNames, FlagVal.ServerTypeNames, "ServerTypeNames与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypeNames)
				assert.Equalf(t, 2, len(m.ServerTypes), "ServerTypes与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypes[0])
			},
			initParams: func() {
				ServerTypes = []string{"api", "mqc"}
				FlagVal.ServerTypeNames = "api-mqc"
			},
		},
		{
			name:    "ServerTypeNames-存在不包含",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				expectMsg := fmt.Sprintf("%s不支持，只能是%v", "mqc", ServerTypes)

				assert.IsNilf(t, false, err, "ServerTypeNames-存在不包含，不包含的用例不通过;expect:%s,actual:%s", expectMsg, err)
				assert.Equalf(t, expectMsg, err.Error(), "ServerTypeNames-存在不包含，不包含的用例不通过;expect:%s,actual:%s", expectMsg, err)
			},
			initParams: func() {
				ServerTypes = []string{"api"}
				FlagVal.ServerTypeNames = "api-mqc"
			},
		},
		{
			name: "Name-为空",
			fields: fields{
				Name:            "",
				RegistryAddr:    "RegistryAddr",
				PlatName:        "PlatName",
				SysName:         "SysName",
				ServerTypeNames: "api-mqc",
				ClusterName:     "ClusterName",
			},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {

				assert.Equalf(t, "RegistryAddr", m.RegistryAddr, "RegistryAddr 与指定值不相等，expect:%s,actual:%s", "RegistryAddr", m.RegistryAddr)

				assert.Equalf(t, "PlatName", m.PlatName, "PlatName 与指定值不相等，expect:%s,actual:%s", "PlatName", m.PlatName)

				assert.Equalf(t, "SysName", m.SysName, "SysName 与指定值不相等，expect:%s,actual:%s", "SysName", m.SysName)

				assert.Equalf(t, "api-mqc", m.ServerTypeNames, "ServerTypeNames 与指定值不相等，expect:%s,actual:%s", "api-mqc", m.ServerTypeNames)
				assert.Equalf(t, "ClusterName", m.ClusterName, "ClusterName 与指定值不相等，expect:%s,actual:%s", "ClusterName", m.ClusterName)

			},
			initParams: func() {
				FlagVal.RegistryAddr = ""
				FlagVal.Name = ""
				FlagVal.PlatName = ""
				FlagVal.SysName = ""
				FlagVal.ServerTypeNames = ""
				FlagVal.ClusterName = ""
			},
		},
		{
			name: "Name-不为空",
			fields: fields{
				Name: "/PlatName/SysName/api-mqc/ClusterName",
			},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {

				assert.Equalf(t, FlagVal.RegistryAddr, m.RegistryAddr, "RegistryAddr 与指定值不相等，expect:%s,actual:%s", FlagVal.RegistryAddr, m.RegistryAddr)

				assert.Equalf(t, "PlatName", m.PlatName, "PlatName 与指定值不相等，expect:%s,actual:%s", "PlatName", m.PlatName)

				assert.Equalf(t, "SysName", m.SysName, "SysName 与指定值不相等，expect:%s,actual:%s", "SysName", m.SysName)

				assert.Equalf(t, []string{"api", "mqc"}, m.ServerTypes, "ServerTypeNames 与指定值不相等，expect:%s,actual:%s", "api-mqc", m.ServerTypeNames)
				assert.Equalf(t, "ClusterName", m.ClusterName, "ClusterName 与指定值不相等，expect:%s,actual:%s", "ClusterName", m.ClusterName)

			},
			initParams: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
			},
		},
		{
			name: "Trace不为空-在trace列表内",
			fields: fields{
				Trace: "cpu",
			},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.IsNilf(t, true, err, "Trace在trace列表内，expect:%s,actual:%s", "nil", err)
			},
			initParams: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
				FlagVal.PlatName = "PlatName"
				FlagVal.SysName = "SysName"
				FlagVal.ServerTypeNames = "api"
				FlagVal.ClusterName = "ClusterName"
			},
		},
		{
			name: "Trace不为空-不在trace列表内",
			fields: fields{
				Trace: "xpu", //不在trace列表
			},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				expectMsg := fmt.Sprintf("trace名称只能是%v", traces)
				assert.IsNilf(t, false, err, "Trace不在trace列表内，expect:%s,actual:%s", expectMsg, err.Error())
				assert.Equalf(t, expectMsg, err.Error(), "Trace不在trace列表内，expect:%s,actual:%s", expectMsg, err.Error())

			},
			initParams: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
				FlagVal.PlatName = "PlatName"
				FlagVal.SysName = "SysName"
				FlagVal.ServerTypeNames = "api"
				FlagVal.ClusterName = "ClusterName"
			},
		},
		{
			name:    "IsDebug=true",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.NotEqualf(t, FlagVal.PlatName, m.PlatName, "PlatName与指定值不相等，expect:%s,actual:%s", FlagVal.PlatName, m.PlatName)
			},
			initParams: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
				FlagVal.PlatName = "PlatName"
				FlagVal.SysName = "SysName"
				FlagVal.ServerTypeNames = "api"
				FlagVal.ClusterName = "ClusterName"
				IsDebug = true
			},
		},
		{
			name:    "IsDebug=false",
			fields:  fields{},
			wantErr: false,
			assertResult: func(t *testing.T, m *global, err error) {
				assert.Equalf(t, FlagVal.PlatName, m.PlatName, "PlatName与指定值不相等，expect:%s,actual:%s", FlagVal.PlatName, m.PlatName)
			},
			initParams: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
				FlagVal.PlatName = "PlatName"
				FlagVal.SysName = "SysName"
				FlagVal.ServerTypeNames = "api"
				FlagVal.ClusterName = "ClusterName"
				IsDebug = false
			},
		},
	}
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
			LocalConfName:   tt.fields.LocalConfName,
		}
		//初始化测试用例参数
		tt.initParams()
		//执行被测试方法
		err := m.check()
		//检查测试结果
		tt.assertResult(t, m, err)
	}
}
