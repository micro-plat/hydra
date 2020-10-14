package global

import (
	"fmt"
	"reflect"
	"testing"
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
			name: "测试有传入名称-短",
			args: args{
				n: []string{"name"},
			},
			want: "com_micro-plat_hydra_global_name",
		},
		{
			name: "测试有传入名称-超过32",
			args: args{
				n: []string{"name123456789012345678901234567890"},
			},
			want: "me123456789012345678901234567890",
		},
		{
			name: "测试未传入名称",
			args: args{
				n: []string{},
			},
			want: "ro-plat_hydra_global_global.test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &global{}
			if got := m.GetLongAppName(tt.args.n...); got != tt.want {
				t.Errorf("global.GetLongAppName() = %v, want %v", got, tt.want)
			}
		})
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
		t.Run(tt.name, func(t *testing.T) {
			m := &global{
				ServerTypes: tt.fields.ServerTypes,
			}
			if got := m.HasServerType(tt.args.tp); got != tt.want {
				t.Errorf("global.HasServerType() = %v, want %v", got, tt.want)
			}
		})
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
		t.Run(tt.name, func(t *testing.T) {
			gotPlatName, gotSystemName, gotServerTypes, gotClusterName, err := parsePath(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPlatName != tt.wantPlatName {
				t.Errorf("parsePath() gotPlatName = %v, want %v", gotPlatName, tt.wantPlatName)
			}
			if gotSystemName != tt.wantSystemName {
				t.Errorf("parsePath() gotSystemName = %v, want %v", gotSystemName, tt.wantSystemName)
			}
			if !reflect.DeepEqual(gotServerTypes, tt.wantServerTypes) {
				t.Errorf("parsePath() gotServerTypes = %v, want %v", gotServerTypes, tt.wantServerTypes)
			}
			if gotClusterName != tt.wantClusterName {
				t.Errorf("parsePath() gotClusterName = %v, want %v", gotClusterName, tt.wantClusterName)
			}
		})
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
		name    string
		fields  fields
		wantErr bool
		errMsg  string
		assert  Assert
		init    Init
	}{

		{
			name:    "RegistryAddr",
			fields:  fields{},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				if m.RegistryAddr != FlagVal.RegistryAddr {
					t.Errorf("RegistryAddr与指定值不相等，expect:%s,actual:%s", FlagVal.RegistryAddr, m.RegistryAddr)
				}
			},
			init: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
			},
		},
		{
			name:    "PlatName",
			fields:  fields{},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				if m.PlatName != FlagVal.PlatName {
					t.Errorf("PlatName与指定值不相等，expect:%s,actual:%s", FlagVal.PlatName, m.PlatName)
				}
			},
			init: func() {
				FlagVal.PlatName = "PlatName"
			},
		},
		{
			name:    "SysName-不为空",
			fields:  fields{},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				if m.SysName != FlagVal.SysName {
					t.Errorf("SysName与指定值不相等，expect:%s,actual:%s", FlagVal.SysName, m.SysName)
				}
			},
			init: func() {
				FlagVal.SysName = "SysName"
			},
		},
		{
			name:    "SysName-为空",
			fields:  fields{},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				if m.SysName != AppName {
					t.Errorf("SysName与指定值不相等，expect:%s,actual:%s", AppName, m.SysName)
				}
			},
			init: func() {
				FlagVal.SysName = ""
			},
		},
		{
			name:    "ServerTypeNames-单个",
			fields:  fields{},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				if m.ServerTypeNames != FlagVal.ServerTypeNames {
					t.Errorf("ServerTypeNames与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypeNames)
				}
				if len(m.ServerTypes) != 1 || m.ServerTypes[0] != "api" {
					t.Errorf("ServerTypes与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypes[0])
				}

			},
			init: func() {
				ServerTypes = []string{"api"}
				FlagVal.ServerTypeNames = "api"
			},
		},
		{
			name:    "ServerTypeNames-多个",
			fields:  fields{},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				if m.ServerTypeNames != FlagVal.ServerTypeNames {
					t.Errorf("ServerTypeNames与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypeNames)
				}
				if len(m.ServerTypes) != 2 {
					t.Errorf("ServerTypes与指定值不相等，expect:%s,actual:%s", FlagVal.ServerTypeNames, m.ServerTypes)
				}
			},
			init: func() {
				ServerTypes = []string{"api", "mqc"}
				FlagVal.ServerTypeNames = "api-mqc"
			},
		},
		{
			name:    "ServerTypeNames-存在不包含",
			fields:  fields{},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				expectMsg := fmt.Sprintf("%s不支持，只能是%v", "mqc", ServerTypes)
				if err == nil || err.Error() != expectMsg {
					t.Errorf("ServerTypeNames-存在不包含，不包含的用例不通过;expect:%s,actual:%s", expectMsg, err)
				}
			},
			init: func() {
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
			assert: func(t *testing.T, m *global, err error) {
				if m.RegistryAddr != "RegistryAddr" {
					t.Errorf("RegistryAddr 与指定值不相等，expect:%s,actual:%s", "RegistryAddr", m.RegistryAddr)
				}
				if m.PlatName != "PlatName" {
					t.Errorf("PlatName 与指定值不相等，expect:%s,actual:%s", "PlatName", m.PlatName)
				}
				if m.SysName != "SysName" {
					t.Errorf("SysName 与指定值不相等，expect:%s,actual:%s", "SysName", m.SysName)
				}
				if m.ServerTypeNames != "api-mqc" {
					t.Errorf("ServerTypeNames 与指定值不相等，expect:%s,actual:%s", "api-mqc", m.ServerTypeNames)
				}
				if m.ClusterName != "ClusterName" {
					t.Errorf("ClusterName 与指定值不相等，expect:%s,actual:%s", "ClusterName", m.ClusterName)
				}
			},
			init: func() {
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
				Name: "/platName/sysName/api-mqc/clusterName",
			},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				if m.RegistryAddr != FlagVal.RegistryAddr {
					t.Errorf("RegistryAddr 与指定值不相等，expect:%s,actual:%s", FlagVal.RegistryAddr, m.RegistryAddr)
				}
				if m.PlatName != "platName" {
					t.Errorf("PlatName 与指定值不相等，expect:%s,actual:%s", "platName", m.PlatName)
				}
				if m.SysName != "sysName" {
					t.Errorf("SysName 与指定值不相等，expect:%s,actual:%s", "sysName", m.SysName)
				}
				if len(m.ServerTypes) != 2 {
					t.Errorf("ServerTypeNames 与指定值不相等，expect:%s,actual:%s", "api-mqc", m.ServerTypeNames)
				}
				if m.ClusterName != "clusterName" {
					t.Errorf("ClusterName 与指定值不相等，expect:%s,actual:%s", "clusterName", m.ClusterName)
				}
			},
			init: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
			},
		},
		{
			name: "Trace不为空-在trace列表内",
			fields: fields{
				Trace: "cpu",
			},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				if err != nil {
					t.Errorf("Trace在trace列表内，expect:%s,actual:%s", "nil", err.Error())
				}
			},
			init: func() {
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
				Trace: "xcpu",
			},
			wantErr: false,
			assert: func(t *testing.T, m *global, err error) {
				expectMsg := fmt.Sprintf("trace名称只能是%v", traces)
				if err == nil || err.Error() != expectMsg {
					t.Errorf("Trace不在trace列表内，expect:%s,actual:%s", expectMsg, err.Error())
				}
			},
			init: func() {
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
			assert: func(t *testing.T, m *global, err error) {
				if m.PlatName == FlagVal.PlatName {
					t.Errorf("PlatName与指定值不相等，expect:%s,actual:%s", FlagVal.PlatName, m.PlatName)
				}
			},
			init: func() {
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
			assert: func(t *testing.T, m *global, err error) {
				if m.PlatName != FlagVal.PlatName {
					t.Errorf("PlatName与指定值不相等，expect:%s,actual:%s", FlagVal.PlatName, m.PlatName)
				}
			},
			init: func() {
				FlagVal.RegistryAddr = "zk://192.168.0.1"
				FlagVal.PlatName = "PlatName"
				FlagVal.SysName = "SysName"
				FlagVal.ServerTypeNames = "api"
				FlagVal.ClusterName = "ClusterName"
				IsDebug = false
			},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			tt.init()
			err := m.check()
			tt.assert(t, m, err)
		})
	}
}
