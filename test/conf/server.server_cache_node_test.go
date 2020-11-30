package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNewCNode(t *testing.T) {
	type args struct {
		name  string
		mid   string
		index int
	}
	type wantInfo struct {
		name     string
		root     string
		path     string
		host     string
		port     string
		index    int
		serverID string
		mid      string
	}
	tests := []struct {
		name string
		args args
		want *wantInfo
	}{
		{name: "1. Conf-NewCNode-初始化默认对象", args: args{name: "1212_dd"}, want: &wantInfo{name: "1212_dd", root: "", path: "", host: "1212", port: "", index: 0, serverID: "dd", mid: ""}},
		{name: "2. Conf-NewCNode-设置域名默认对象", args: args{name: "baid.com_dd", mid: "midxx", index: 3}, want: &wantInfo{name: "baid.com_dd", root: "", path: "", host: "baid.com", port: "", index: 3, serverID: "dd", mid: "midxx"}},
		{name: "3. Conf-NewCNode-设置ip默认对象", args: args{name: "192.169.0.11:3333_dd", mid: "midxx", index: 3}, want: &wantInfo{name: "192.169.0.11:3333_dd", root: "", path: "", host: "192.169.0.11", port: "3333", index: 3, serverID: "dd", mid: "midxx"}},
	}
	for _, tt := range tests {
		got := server.NewCNode(tt.args.name, tt.args.mid, tt.args.index)
		assert.Equal(t, tt.want.name, got.GetName(), tt.name+",name")
		assert.Equal(t, tt.want.host, got.GetHost(), tt.name+",host")
		assert.Equal(t, tt.want.root, got.GetRoot(), tt.name+",root")
		assert.Equal(t, tt.want.path, got.GetPath(), tt.name+",path")
		assert.Equal(t, tt.want.port, got.GetPort(), tt.name+",port")
		assert.Equal(t, tt.want.index, got.GetIndex(), tt.name+",index")
		assert.Equal(t, tt.want.serverID, got.GetNodeID(), tt.name+",serverID")
		assert.Equal(t, tt.want.mid, got.GetMid(), tt.name+",mid")
	}
}

func TestCNode_IsAvailable(t *testing.T) {
	tests := []struct {
		name   string
		fields *server.CNode
		want   bool
	}{
		{name: "1. Conf-CNodeIsAvailable-初始化对象肯定返回true", fields: server.NewCNode("11_22", "", 0), want: true},
	}
	for _, tt := range tests {
		got := tt.fields.IsAvailable()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCNode_IsCurrent(t *testing.T) {
	tests := []struct {
		name   string
		fields *server.CNode
		want   bool
	}{
		{name: "1. Conf-CNodeIsCurrent-空mid", fields: server.NewCNode("11_22", "", 0), want: false},
		{name: "2. Conf-CNodeIsCurrent-midserverid不想等", fields: server.NewCNode("11_22", "33", 0), want: false},
		{name: "3. Conf-CNodeIsCurrent-midserverid想等", fields: server.NewCNode("11_22", "22", 0), want: true},
	}
	for _, tt := range tests {
		got := tt.fields.IsCurrent()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCNode_IsMaster(t *testing.T) {
	tests := []struct {
		name   string
		fields *server.CNode
		args   int
		want   bool
	}{
		{name: "1. Conf-CNodeIsMaster-入参小于节点index", fields: server.NewCNode("11_22", "", 2), args: 1, want: false},
		{name: "2. Conf-CNodeIsMaster-入参等于节点index", fields: server.NewCNode("11_22", "", 2), args: 2, want: false},
		{name: "3. Conf-CNodeIsMaster-入参大于节点index", fields: server.NewCNode("11_22", "", 2), args: 3, want: true},
	}
	for _, tt := range tests {
		got := tt.fields.IsMaster(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCNode_Clone(t *testing.T) {
	tests := []struct {
		name   string
		fields *server.CNode
		want   conf.ICNode
	}{
		{name: "1. Conf-CNodeClone-设置index1", fields: server.NewCNode("11_22", "", 2), want: server.NewCNode("11_22", "", 2)},
	}
	for _, tt := range tests {
		got := tt.fields.Clone()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCNode_GetIndex(t *testing.T) {
	tests := []struct {
		name   string
		fields *server.CNode
		want   int
	}{
		{name: "1. Conf-CNodeGetIndex-设置index1", fields: server.NewCNode("11_22", "", 0), want: 0},
		{name: "2. Conf-CNodeGetIndex-设置index2", fields: server.NewCNode("11_22", "", 1), want: 1},
		{name: "3. Conf-CNodeGetIndex-设置index3", fields: server.NewCNode("11_22", "", 2), want: 2},
	}
	for _, tt := range tests {
		got := tt.fields.GetIndex()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCNode_GetName(t *testing.T) {
	tests := []struct {
		name   string
		fields *server.CNode
		want   string
	}{
		{name: "1. Conf-CNodeGetName-设置name1", fields: server.NewCNode("11_22aa", "", 0), want: "11_22aa"},
		{name: "2. Conf-CNodeGetName-设置name2", fields: server.NewCNode("11_22bb", "", 1), want: "11_22bb"},
		{name: "3. Conf-CNodeGetName-设置name3", fields: server.NewCNode("11_22cc", "", 2), want: "11_22cc"},
	}
	for _, tt := range tests {
		got := tt.fields.GetName()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCNode_GetServerID(t *testing.T) {

	tests := []struct {
		name   string
		fields *server.CNode
		want   string
	}{
		{name: "1. Conf-CNodeGetServerID-设置ServerID1", fields: server.NewCNode("11_22aa", "", 0), want: "22aa"},
		{name: "2. Conf-CNodeGetServerID-设置ServerID2", fields: server.NewCNode("11_22bb", "", 1), want: "22bb"},
		{name: "3. Conf-CNodeGetServerID-设置ServerID3", fields: server.NewCNode("11_22cc", "", 2), want: "22cc"},
	}
	for _, tt := range tests {
		got := tt.fields.GetNodeID()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCNode_GetPort(t *testing.T) {
	tests := []struct {
		name   string
		fields *server.CNode
		want   string
	}{
		{name: "1. Conf-CNodeGetPort-设置Port1", fields: server.NewCNode("11_22aa", "", 0), want: ""},
		{name: "2. Conf-CNodeGetPort-设置Port2", fields: server.NewCNode("11:1212_22bb", "", 1), want: "1212"},
		{name: "3. Conf-CNodeGetPort-设置Port3", fields: server.NewCNode("11:3333_22cc", "", 2), want: "3333"},
	}
	for _, tt := range tests {
		got := tt.fields.GetPort()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCNode_GetHost(t *testing.T) {
	tests := []struct {
		name   string
		fields *server.CNode
		want   string
	}{
		{name: "1. Conf-CNodeGetHost-设置Host1", fields: server.NewCNode("11_22aa", "", 0), want: "11"},
		{name: "2. Conf-CNodeGetHost-设置Host2", fields: server.NewCNode("192.169.1.11:1212_22bb", "", 1), want: "192.169.1.11"},
		{name: "3. Conf-CNodeGetHost-设置Host3", fields: server.NewCNode("www.baidu.com_22cc", "", 2), want: "www.baidu.com"},
	}
	for _, tt := range tests {
		got := tt.fields.GetHost()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
