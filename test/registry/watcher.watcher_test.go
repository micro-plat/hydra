package registry

import (
	"testing"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/watcher"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewCArgsByChange(t *testing.T) {
	confObj := mocks.NewConfBy("hydra_rgst_watcher_test", "rgtwatchertest") //构建对象
	r := confObj.Registry

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
		{name: "1. 构建子节点变动参数,父节点不为空,无后/", op: 1, deep: 1, parent: "a/b/c/!@#$%^&*", children: []string{"children1", "children2"}, v: 1, r: r, want: &watcher.ChildChangeArgs{OP: 1, Registry: r, Parent: "a/b/c/!@#$%^&*", Version: 1, Children: []string{"children1", "children2"}, Deep: 1, Name: "!@#$%^&*"}},
		{name: "2. 构建子节点变动参数,父节点不为空,有后/", op: 1, deep: 1, parent: "a/b/c/!@#$%^&*/", children: []string{"children1", "children2"}, v: 1, r: r, want: &watcher.ChildChangeArgs{OP: 1, Registry: r, Parent: "a/b/c/!@#$%^&*/", Version: 1, Children: []string{"children1", "children2"}, Deep: 1, Name: "!@#$%^&*"}},
		{name: "3. 构建子节点变动参数,父节点为空", op: 1, deep: 1, parent: "", children: []string{"children1", "children2"}, v: 1, r: r, want: &watcher.ChildChangeArgs{OP: 1, Registry: r, Parent: "", Version: 1, Children: []string{"children1", "children2"}, Deep: 1, Name: ""}},
	}

	for _, tt := range tests {
		got := watcher.NewCArgsByChange(tt.op, tt.deep, tt.parent, tt.children, tt.v, tt.r)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
