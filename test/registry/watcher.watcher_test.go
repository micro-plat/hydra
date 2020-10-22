package registry

import (
	"testing"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewCArgsByChange(t *testing.T) {
	confObj := mocks.NewConf() //构建对象
	confObj.API(":8080")
	apiconf := confObj.GetAPIConf() //初始化参数
	c := apiconf.GetMainConf()      //获取配置
	r := c.GetRegistry()

	tests := []struct {
		name     string
		op       int
		deep     int
		parent   string
		children []string
		v        int32
		r        registry.IRegistry
		want     *watcher.ChildChangeArgs
	}{
		{name: "构建子节点变动参数,父节点不为空", op: 1, deep: 1, parent: "a/b/c/!@#$%^&*", children: []string{"children1", "children2"},
			v: 1, r: r, want: &watcher.ChildChangeArgs{
				OP:       1,
				Registry: r,
				Parent:   "a/b/c/!@#$%^&*",
				Version:  1,
				Children: []string{"children1", "children2"},
				Deep:     1,
				Name:     "!@#$%^&*",
			}},
		{name: "构建子节点变动参数,父节点为空", op: 1, deep: 1, parent: "", children: []string{"children1", "children2"},
			v: 1, r: r, want: &watcher.ChildChangeArgs{
				OP:       1,
				Registry: r,
				Parent:   "",
				Version:  1,
				Children: []string{"children1", "children2"},
				Deep:     1,
				Name:     "",
			}},
	}

	for _, tt := range tests {
		got := watcher.NewCArgsByChange(tt.op, tt.deep, tt.parent, tt.children, tt.v, tt.r)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestValueChangeArgs_IsConf(t *testing.T) {
	tests := []struct {
		name string
		n    *watcher.ValueChangeArgs
		want bool
	}{
		{name: "conf根节点", n: &watcher.ValueChangeArgs{Path: "/conf"}, want: true},
		{name: "非conf根节点", n: &watcher.ValueChangeArgs{Path: "conf/"}, want: false},
		{name: "conf子节点", n: &watcher.ValueChangeArgs{Path: "/conf/!@#%%^&*"}, want: true},
		{name: "非conf子根节点", n: &watcher.ValueChangeArgs{Path: "conf/a/b"}, want: false},
	}
	for _, tt := range tests {
		got := tt.n.IsConf()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestValueChangeArgs_IsVarRoot(t *testing.T) {
	tests := []struct {
		name string
		n    *watcher.ValueChangeArgs
		want bool
	}{
		{name: "var根节点", n: &watcher.ValueChangeArgs{Path: "/var"}, want: true},
		{name: "非var根节点", n: &watcher.ValueChangeArgs{Path: "var/"}, want: false},
		{name: "var子节点", n: &watcher.ValueChangeArgs{Path: "/var/!@#%%^&*"}, want: true},
		{name: "非var子根节点", n: &watcher.ValueChangeArgs{Path: "var/a/b"}, want: false},
	}
	for _, tt := range tests {
		got := tt.n.IsVarRoot()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
