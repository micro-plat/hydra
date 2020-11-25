package registry

import (
	"sort"
	"testing"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/micro-plat/hydra/test/assert"
	r "github.com/micro-plat/lib4go/registry"
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

func createRegistry() []registry.IRegistry {
	rgs := make([]registry.IRegistry, 0, 1)
	rgs = append(rgs, localmemory.NewLocalMemory())
	return rgs

}

func TestCreateTempNode(t *testing.T) {

	//构建注册中心
	lm := localmemory.NewLocalMemory()
	var cases = []struct {
		name  string
		path  string
		value string
	}{
		{name: "1.1 LMCreateTemp-一段路径-字母", path: "hydrax", value: "1"},
		{name: "1.2 LMCreateTemp-一段路径-数字", path: "1231222", value: "2"},
		{name: "1.3 LMCreateTemp-一段路径-字母数字混合", path: "123hydra3", value: "3"},
		{name: "1.4 LMCreateTemp-一段路径-含特殊字符", path: "1232hydra#$%", value: "4"},
		{name: "1.5 LMCreateTemp-一段路径-全特殊字符", path: "#$%", value: "5"},
		{name: "1.6 LMCreateTemp-一段路径-有前/", path: "/123123", value: "6"},
		{name: "1.7 LMCreateTemp-一段路径-有后/", path: "hydra#$%xee/", value: "7"},
		{name: "1.8 LMCreateTemp-一段路径-有前后/", path: "/hydra#$%/", value: "8"},
		{name: "1.9 LMCreateTemp-一段路径-长路径", path: "/hydraabcefgjijkfsnopqrstuvwxyz", value: "1445"},

		{name: "2.1 LMCreateTemp-二段路径-以段以上路径", path: "/hydra1/abc/", value: "18"},
		{name: "2.2 LMCreateTemp-二段路径-多段有数字", path: "/hydra2/454/", value: "17"},
		{name: "2.3 LMCreateTemp-二段路径-多段有特殊字符", path: "/hydra3/#$#%/", value: "189"},
		{name: "2.4 LMCreateTemp-二段路径-有后/", path: "hydra4/abc/", value: "181"},
		{name: "2.5 LMCreateTemp-二段路径-前后/", path: "/hydra5/454/", value: "173"},
		{name: "2.6 LMCreateTemp-二段路径-前/", path: "/hydra6/#$#%", value: "189x"},

		{name: "3.1 LMCreateTemp-多段-较长分段", path: "/hydra11/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/xxx", value: "1255"},
		{name: "3.2 LMCreateTemp-多段-较长分段1", path: "hydra22/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/yyy", value: "12225"},
	}

	//按注册中心进行测试
	//创建节点
	for _, c := range cases {
		err := lm.CreateTempNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)
		mp := map[string]string{c.path: c.value}
		err = checkData(lm, c.path, mp)
		assert.Equal(t, nil, err, c.name)
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
	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, lm := range rgs {
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
}

func TestExists(t *testing.T) {

	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, lm := range rgs {
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
}
func TestDelete(t *testing.T) {

	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, lm := range rgs {
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
}

func TestChildren(t *testing.T) {

	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, lm := range rgs {
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
}

func TestWatchValue(t *testing.T) {
	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, lm := range rgs {
		cases := []struct {
			name   string
			path   string
			value  string
			nvalue string
		}{
			{name: "一个", path: "/hydra1", value: "1", nvalue: "2"},
			{name: "一个", path: "/hydra1/hydra1-1", value: "1", nvalue: "2"},
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
			go func(c chan r.ValueWatcher, name, nvalue string) {
				select {
				case v := <-c:
					value, version := v.GetValue()
					assert.NotEqual(t, version, int32(0), name)
					assert.Equal(t, nvalue, string(value), name)
				case <-time.After(time.Second):
					t.Error("测试未通过")
				}
			}(notify, c.name, c.nvalue)
		}

		//更新值
		for _, c := range cases {
			err := lm.Update(c.path, c.nvalue)
			assert.Equal(t, nil, err, c.name)
		}

		time.Sleep(time.Second)
	}
}

func TestWatchChildren(t *testing.T) {
	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, lm := range rgs {
		cases := []struct {
			name     string
			path     string
			children []string
			values   []string
		}{
			{name: "一个", path: "/hydra1", values: []string{"1", "2"}, children: []string{"efg"}},
			{name: "多个", path: "/hydra2", values: []string{"1", "3"}, children: []string{"abc", "efg", "efss", "12", "!@#"}},
		}

		for _, c := range cases {

			for _, value := range c.values {

				//监控父节点
				notify, err := lm.WatchChildren(c.path)
				assert.Equal(t, nil, err, c.name)

				//创建节点
				for _, ch := range c.children {
					err := lm.CreateTempNode(registry.Join(c.path, ch), value)
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
	}
}
func TestWatchChildrenForDelete(t *testing.T) {
	//构建所有注册中心
	rgs := createRegistry()

	//按注册中心进行测试
	for _, lm := range rgs {
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

			//创建节点
			for _, ch := range c.children {
				err := lm.CreateTempNode(registry.Join(c.path, ch), c.value)
				assert.Equal(t, nil, err, c.name)
			}

			//监控父节点
			notify, err := lm.WatchChildren(c.path)
			assert.Equal(t, nil, err, c.name)

			//删除
			for _, ch := range c.children {
				err := lm.Delete(registry.Join(c.path, ch))
				assert.Equal(t, nil, err, c.name)
				//应收到值变化通知
				select {
				case v := <-notify:
					assert.Equal(t, c.path, v.GetPath(), c.name)
					cPath, cVersion := v.GetValue()
					assert.NotEqual(t, int32(0), cVersion, c.name)
					assert.Equal(t, []string{registry.Join(c.path, ch)}, cPath, c.name)
					notify, err = lm.WatchChildren(c.path)
					assert.Equal(t, nil, err, c.name)
				default:
					t.Error("测试未通过", c.name)
				}
			}
		}
	}
}
