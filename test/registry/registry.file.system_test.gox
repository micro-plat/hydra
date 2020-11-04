package registry

import (
	"testing"

	"github.com/micro-plat/hydra/registry/registry/filesystem"
	"github.com/micro-plat/hydra/test/assert"
)

//  /hydra/apiserver/api/test/conf:{"address":":51001"}
//  /hydra/apiserver/api/test/conf/connection_max:5000
//  /hydra/apiserver/api/test/conf/enabled:true
//  /hydra/apiserver/api/test/conf/port:[8001,8001,8002]
//  /hydra/apiserver/api/test/conf/router:"router"
//  /hydra/apiserver/cron/test/conf:{}
//  /hydra/apiserver/cron/test/conf/tesk:{"task":[]}
//  /hydra/apiserver/servers/test/conf/alpha:{"dc":"eqdc10","ip":"10.0.0.1"}
//  /hydra/apiserver/servers/test/conf/beta:{"dc":"eqdc10","ip":"10.0.0.2"}

func TestNewFileSystem(t *testing.T) {
	tests := []struct {
		name        string
		platName    string
		systemName  string
		clusterName string
		path        string
		wantErr     bool
	}{
		{name: "错误路径", platName: "hydra", systemName: "apiserver", clusterName: "test", path: "./a.toml", wantErr: true},
		{name: "正确路径", platName: "hydra", systemName: "apiserver", clusterName: "test", path: "./registry.file.system_test.toml", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := filesystem.NewFileSystem(tt.platName, tt.systemName, tt.clusterName, tt.path)
			assert.Equal(t, tt.wantErr, err != nil, tt.name)
		})
	}
}

func Test_fileSystem_Exists(t *testing.T) {

	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{name: "不存在的节点", path: "/a/!@#$%^&*/c", want: false, wantErr: false},
		//{name: "存在的父节点", path: "/hydra/apiserver/api/test", want: true},
		{name: "存在的子节点", path: "/hydra/apiserver/api/test/conf", want: true},
	}
	l, err := filesystem.NewFileSystem("hydra", "apiserver", "test", "./registry.file.system_test.toml")
	assert.Equal(t, false, err != nil, "构建对象")

	for _, tt := range tests {
		got, err := l.Exists(tt.path)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_fileSystem_GetValue(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantData    []byte
		wantVersion int32
		wantErr     bool
	}{
		{name: "不存在的节点", path: "/a/!@#$%^&*/c", wantErr: true},
		//{name: "存在的父节点", path: "/hydra/apiserver/api/test",},
		{name: "存在的子节点", path: "/hydra/apiserver/api/test/conf", wantData: []byte(`{"address":":51001"}`)},
	}
	l, err := filesystem.NewFileSystem("hydra", "apiserver", "test", "./registry.file.system_test.toml")
	assert.Equal(t, false, err != nil, "构建对象")

	for _, tt := range tests {
		data, version, err := l.GetValue(tt.path)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.Equal(t, tt.wantVersion, version, tt.name)
		assert.Equal(t, tt.wantData, data, tt.name)
	}
}

func Test_fileSystem_Update(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		data    string
		wantErr bool
	}{
		{name: "不存在的节点", path: "/a/!@#$%^&*/c", data: "data", wantErr: true},
		//{name: "存在的父节点", path: "/hydra/apiserver/api/test",},
		{name: "存在的子节点", path: "/hydra/apiserver/api/test/conf", data: "data"},
	}
	l, err := filesystem.NewFileSystem("hydra", "apiserver", "test", "./registry.file.system_test.toml")
	assert.Equal(t, false, err != nil, "构建对象")

	for _, tt := range tests {
		err := l.Update(tt.path, tt.data)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			continue
		}
		data, _, err := l.GetValue(tt.path)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.data, string(data), tt.name)
	}
}

func Test_fileSystem_GetChild(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		version int32
		paths   []string
		wantErr bool
	}{
		{name: "不存在的节点", path: "/a/!@#$%^&*/c", paths: []string{}},
		{name: "存在的父节点", path: "/hydra/apiserver/api", paths: []string{}},
	}
	l, err := filesystem.NewFileSystem("hydra", "apiserver", "test", "./registry.file.system_test.toml")
	assert.Equal(t, false, err != nil, "构建对象")

	for _, tt := range tests {
		paths, version, err := l.GetChildren(tt.path)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		assert.Equal(t, tt.version, version, tt.name)
		assert.Equal(t, tt.paths, paths, tt.name)
	}
}

func Test_fileSystem_Delete(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{name: "不存在的节点", path: "/a/!@#$%^&*/c"},
		{name: "存在的子节点", path: "/hydra/apiserver/api/test/conf"},
	}
	l, err := filesystem.NewFileSystem("hydra", "apiserver", "test", "./registry.file.system_test.toml")
	assert.Equal(t, false, err != nil, "构建对象")

	for _, tt := range tests {
		err := l.Delete(tt.path)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
	}
}

func Test_fileSystem_CreateTempNode(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		data    string
		version int32
		wantErr bool
	}{
		{name: "不存在的节点", path: "/a/!@#$%^&*/c", data: "data"},
		{name: "存在的子节点", path: "/hydra/apiserver/api/test/conf", data: "data"},
	}
	l, err := filesystem.NewFileSystem("hydra", "apiserver", "test", "./registry.file.system_test.toml")
	assert.Equal(t, false, err != nil, "构建对象")

	for _, tt := range tests {
		err := l.CreateTempNode(tt.path, tt.data)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		data, version, err := l.GetValue(tt.path)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.version, version, tt.name)
		assert.Equal(t, tt.data, string(data), tt.name)
	}
}

func Test_fileSystem_CreateSeqNode(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		data    string
		version int32
		wantErr bool
	}{
		{name: "不存在的节点", path: "/a/!@#$%^&*/c", data: "data"},
		{name: "存在的子节点", path: "/hydra/apiserver/api/test/conf", data: "data"},
	}
	l, err := filesystem.NewFileSystem("hydra", "apiserver", "test", "./registry.file.system_test.toml")
	assert.Equal(t, false, err != nil, "构建对象")

	for _, tt := range tests {
		rpath, err := l.CreateSeqNode(tt.path, tt.data)
		assert.Equal(t, tt.wantErr, err != nil, tt.name)
		data, version, err := l.GetValue(rpath)
		assert.Equal(t, false, err != nil, tt.name)
		assert.Equal(t, tt.version, version, tt.name)
		assert.Equal(t, tt.data, string(data), tt.name)
	}
}
