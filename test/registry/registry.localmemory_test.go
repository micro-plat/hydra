package registry

import (
	"sort"
	"testing"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/micro-plat/hydra/test/assert"
)

var cases = []struct {
	name  string
	path  string
	value string
}{
	{name: "一段路径", path: "hydra", value: "1"},
	{name: "路径中有数字", path: "1231222", value: "2"},
	{name: "数字在前", path: "123hydra", value: "3"},
	{name: "路径中有特殊字符", path: "1232hydra#$%", value: "4"},
	{name: "路径中只有特殊字符", path: "#$%", value: "5"},
	{name: "路径中有数字", path: "/123123", value: "6"},
	{name: "带/线", path: "/hydra#$%xee", value: "7"},
	{name: "前后/", path: "/hydra#$%/", value: "8"},
	{name: "以段以上路径", path: "/hydra/abc/", value: "18"},
	{name: "多段有数字", path: "/hydra/454/", value: "17"},
	{name: "多段有特殊字符", path: "/hydra/#$#%/", value: "189"},
	{name: "较长路径", path: "/hydraabcefgjijklmnopqrstuvwxyz", value: "1445"},
	{name: "较长分段", path: "/hydra/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/", value: "1255"},
}

func TestCreateTempNode(t *testing.T) {

	lm := localmemory.NewLocalMemory()

	//创建节点
	for _, c := range cases {
		err := lm.CreateTempNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)
	}

	//检查节点值是否正确，是否有被覆盖等
	for _, c := range cases {
		data, v, err := lm.GetValue(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.NotEqual(t, v, int32(0), c.name)
		assert.Equal(t, string(data), c.value, c.name)
	}
}
func TestUpdateNode(t *testing.T) {
	cases := []struct {
		name   string
		path   string
		value  string
		nvalue string
	}{
		{name: "更新为数字", path: "/hydra", value: "1", nvalue: "2333"},
		{name: "更新为字符", path: "1231222", value: "2", nvalue: "sdfd"},
		{name: "更新为中文", path: "123hydra", value: "3", nvalue: "研发"},
		{name: "更新为特殊字符", path: "1232hydra#$%", value: "4", nvalue: "研发12312@#@"},
		{name: "更新为json", path: "#$%", value: "5", nvalue: `{"abc":"ef",age:[10,20]}`},
		{name: "更新为xml", path: "/hydra/apiserver/api/conf", value: "5", nvalue: `<xml><node id="abc"/></xml>`},
	}
	lm := localmemory.NewLocalMemory()

	//创建节点,更新节点
	for _, c := range cases {
		err := lm.CreateTempNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)

		err = lm.Update(c.path, c.nvalue)
		assert.Equal(t, nil, err, c.name)
	}

	//检查节点值是否正确
	for _, c := range cases {
		data, v, err := lm.GetValue(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.NotEqual(t, v, int32(0), c.name)
		assert.Equal(t, string(data), c.nvalue, c.name)
	}
}

func TestExists(t *testing.T) {

	lm := localmemory.NewLocalMemory()

	for _, c := range cases {

		//节点不存在
		exists := false
		b, err := lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, b, exists, c.name)

		//创建节点
		err = lm.CreateTempNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)

		//节点应存在
		exists = true
		b, err = lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, b, exists, c.name)
	}
}
func TestDelete(t *testing.T) {

	lm := localmemory.NewLocalMemory()
	exists := false
	for _, c := range cases {
		//创建节点
		err := lm.CreateTempNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)

		//删除节点
		err = lm.Delete(c.path)
		assert.Equal(t, nil, err, c.name)

		//是否存在
		b, err := lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, b, exists, c.name)
	}
}

func TestChildren(t *testing.T) {

	lm := localmemory.NewLocalMemory()
	cases := []struct {
		name     string
		path     string
		children []string
		value    string
	}{
		{name: "一个", path: "/hydra1", value: "1", children: []string{"efg"}},
		{name: "多个", path: "/hydra2", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},
		{name: "空", path: "/hydra3", value: "1"},
	}

	for _, c := range cases {

		//创建节点
		for _, ch := range c.children {
			err := lm.CreateTempNode(registry.Join(c.path, ch), c.value)
			assert.Equal(t, nil, err, c.name)
		}

		//获取子节点
		paths, v, err := lm.GetChildren(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.NotEqual(t, v, 0, c.name)
		assert.Equal(t, len(paths), len(c.children), paths)

		if len(c.children) == 0 {
			continue
		}

		//排序列表
		sort.Strings(paths)
		sort.Strings(c.children)
		assert.Equal(t, paths, c.children, c.name)
	}
}

func TestWatchValue(t *testing.T) {
	lm := localmemory.NewLocalMemory()
	cases := []struct {
		name   string
		path   string
		value  string
		nvalue string
	}{
		{name: "一个", path: "/hydra1", value: "1", nvalue: "2"},
		{name: "一个", path: "/hydra2", value: "2", nvalue: "234"},
	}
	for _, c := range cases {

		//创建临时节点
		err := lm.CreateTempNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)

		//监控值变化
		notify, err := lm.WatchValue(c.path)
		assert.Equal(t, nil, err, c.name)

		//此时值未变化不应收到通知
		select {
		case <-notify:
			t.Error("测试未通过")
		default:
		}

		//更新值
		err = lm.Update(c.path, c.nvalue)
		assert.Equal(t, nil, err, c.name)

		//应收到值变化通知
		select {
		case v := <-notify:
			value, version := v.GetValue()
			assert.NotEqual(t, version, int32(0), c.name)
			assert.Equal(t, c.nvalue, string(value), c.name)
		default:
			t.Error("测试未通过")
		}

	}

}
func TestWatchChildren(t *testing.T) {
	lm := localmemory.NewLocalMemory()
	cases := []struct {
		name     string
		path     string
		children []string
		value    string
	}{
		{name: "一个", path: "/hydra1", value: "1", children: []string{"efg"}},
		{name: "多个", path: "/hydra2", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},
	}

	for _, c := range cases {

		//监控父节点
		notify, err := lm.WatchChildren(c.path)
		assert.Equal(t, nil, err, c.name)

		//创建节点
		for _, ch := range c.children {
			err := lm.CreateTempNode(registry.Join(c.path, ch), c.value)
			assert.Equal(t, nil, err, c.name)
		}

		//应收到值变化通知
		select {
		case v := <-notify:
			assert.Equal(t, v.GetPath(), c.path)
		default:
			t.Error("测试未通过")
		}
	}
}
