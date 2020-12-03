package registry

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/registry/redis"
	"github.com/micro-plat/hydra/test/assert"
	r "github.com/micro-plat/lib4go/registry"
	"github.com/micro-plat/lib4go/types"
)

func getRegistry() (*redis.Redis, error) {
	return redis.NewRedisBy("", "", []string{"192.168.106.204:6379"}, 0, 10)
}

func TestRedisCreateTempNode(t *testing.T) {

	//构建注册中心
	lm, err := getRegistry()
	assert.Equal(t, nil, err)
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
	for _, c := range cases {
		err := lm.CreateTempNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)

		b, err := lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, true, b, c.name)

		d, _, err := lm.GetValue(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, c.value, string(d))

		err = lm.Delete(c.path)
		assert.Equal(t, nil, err, c.name)

		b, err = lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, false, b, c.name)
	}
}

func TestRedisCreateSeqNode(t *testing.T) {

	//构建注册中心
	lm, err := getRegistry()
	assert.Equal(t, nil, err)
	var cases = []struct {
		name  string
		path  string
		value string
	}{
		{name: "1.1 LMCreateSeq-一段路径-字母", path: "hydrax", value: "1"},
		{name: "1.2 LMCreateSeq-一段路径-数字", path: "1231222", value: "2"},
		{name: "1.3 LMCreateSeq-一段路径-字母数字混合", path: "123hydra3", value: "3"},
		{name: "1.4 LMCreateSeq-一段路径-含特殊字符", path: "1232hydra#$%", value: "4"},
		{name: "1.5 LMCreateSeq-一段路径-全特殊字符", path: "#$%", value: "5"},
		{name: "1.6 LMCreateSeq-一段路径-有前/", path: "/123123", value: "6"},
		{name: "1.7 LMCreateSeq-一段路径-有后/", path: "hydra#$%xee/", value: "7"},
		{name: "1.8 LMCreateSeq-一段路径-有前后/", path: "/hydra#$%/", value: "8"},
		{name: "1.9 LMCreateSeq-一段路径-长路径", path: "/hydraabcefgjijkfsnopqrstuvwxyz", value: "1445"},

		{name: "2.1 LMCreateSeq-二段路径-以段以上路径", path: "/hydra1/abc/", value: "18"},
		{name: "2.2 LMCreateSeq-二段路径-多段有数字", path: "/hydra2/454/", value: "17"},
		{name: "2.3 LMCreateSeq-二段路径-多段有特殊字符", path: "/hydra3/#$#%/", value: "189"},
		{name: "2.4 LMCreateSeq-二段路径-有后/", path: "hydra4/abc/", value: "181"},
		{name: "2.5 LMCreateSeq-二段路径-前后/", path: "/hydra5/454/", value: "173"},
		{name: "2.6 LMCreateSeq-二段路径-前/", path: "/hydra6/#$#%", value: "189x"},

		{name: "3.1 LMCreateSeq-多段-较长分段", path: "/hydra11/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/xxx", value: "1255"},
		{name: "3.2 LMCreateSeq-多段-较长分段1", path: "hydra22/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/yyy", value: "12225"},
	}

	//按注册中心进行测试
	for _, c := range cases {
		rpath, err := lm.CreateSeqNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)
		ca := strings.Split(rpath, "_")
		assert.Equal(t, 2, len(ca), c.name, "返回的路径不合法")
		assert.Equal(t, true, types.GetInt(ca[1]) > 0, c.name, "返回的路径不合法")
		c.path = rpath

		b, err := lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, true, b, c.name)

		d, _, err := lm.GetValue(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, c.value, string(d))

		err = lm.Delete(rpath)
		assert.Equal(t, nil, err, c.name)

		b, err = lm.Exists(rpath)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, false, b, c.name)
	}
}

func TestRedisCreatePersistentNode(t *testing.T) {

	//构建注册中心
	lm, err := getRegistry()
	assert.Equal(t, nil, err)
	var cases = []struct {
		name  string
		path  string
		value string
	}{
		{name: "1.1 LMCreatePersistent-一段路径-字母", path: "hydrax", value: "1"},
		{name: "1.2 LMCreatePersistent-一段路径-数字", path: "1231222", value: "2"},
		{name: "1.3 LMCreatePersistent-一段路径-字母数字混合", path: "123hydra3", value: "3"},
		{name: "1.4 LMCreatePersistent-一段路径-含特殊字符", path: "1232hydra#$%", value: "4"},
		{name: "1.5 LMCreatePersistent-一段路径-全特殊字符", path: "#$%", value: "5"},
		{name: "1.6 LMCreatePersistent-一段路径-有前/", path: "/123123", value: "6"},
		{name: "1.7 LMCreatePersistent-一段路径-有后/", path: "hydra#$%xee/", value: "7"},
		{name: "1.8 LMCreatePersistent-一段路径-有前后/", path: "/hydra#$%/", value: "8"},
		{name: "1.9 LMCreatePersistent-一段路径-长路径", path: "/hydraabcefgjijkfsnopqrstuvwxyz", value: "1445"},

		{name: "2.1 LMCreatePersistent-二段路径-以段以上路径", path: "/hydra1/abc/", value: "18"},
		{name: "2.2 LMCreatePersistent-二段路径-多段有数字", path: "/hydra2/454/", value: "17"},
		{name: "2.3 LMCreatePersistent-二段路径-多段有特殊字符", path: "/hydra3/#$#%/", value: "189"},
		{name: "2.4 LMCreatePersistent-二段路径-有后/", path: "hydra4/abc/", value: "181"},
		{name: "2.5 LMCreatePersistent-二段路径-前后/", path: "/hydra5/454/", value: "173"},
		{name: "2.6 LMCreatePersistent-二段路径-前/", path: "/hydra6/#$#%", value: "189x"},

		{name: "3.1 LMCreatePersistent-多段-较长分段", path: "/hydra11/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/xxx", value: "1255"},
		{name: "3.2 LMCreatePersistent-多段-较长分段1", path: "hydra22/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/yyy", value: "12225"},
	}

	//按注册中心进行测试
	for _, c := range cases {
		err := lm.CreatePersistentNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)

		b, err := lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, true, b, c.name)

		d, _, err := lm.GetValue(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, c.value, string(d))

		err = lm.Delete(c.path)
		assert.Equal(t, nil, err, c.name)

		b, err = lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, false, b, c.name)

	}
}

func TestRedisUpdateNode(t *testing.T) {
	lm, err := getRegistry()
	assert.Equal(t, nil, err)
	cases := []struct {
		name   string
		path   string
		isE    bool
		value  string
		nvalue string
	}{
		{name: "1.1 LMUpdate-节点不存在-更新数据", isE: false, path: "/hydrauodate1", value: "1", nvalue: "2333"},

		{name: "2.1 LMUpdate-节点存在-数据为空-更新为空", isE: true, path: "/hydrauodate2", value: "1", nvalue: "2333"},
		{name: "2.1 LMUpdate-节点存在-数据为空-更新为数字", isE: true, path: "/hydrauodate3", value: "1", nvalue: "2333"},
		{name: "2.1 LMUpdate-节点存在-数据为空-更新为字符", isE: true, path: "sdsd343434", value: "2", nvalue: "sdfd"},
		{name: "2.1 LMUpdate-节点存在-数据为空-更新为中文", isE: true, path: "123hydraddd", value: "3", nvalue: "研发"},
		{name: "2.1 LMUpdate-节点存在-数据为空-更新为特殊字符", isE: true, path: "1232hydrqqa#$%", value: "4", nvalue: "研发12312@#@"},
		{name: "2.1 LMUpdate-节点存在-数据为空-更新为json", isE: true, path: "#$@×&#(%", value: "5", nvalue: `{"abc":"ef",age:[10,20]}`},
		{name: "2.1 LMUpdate-节点存在-数据为空-更新为xml", isE: true, path: "/hydrauodate/apiserver/api/conf", value: "5", nvalue: `<xml><node id="abc"/></xml>`},

		{name: "3.1 LMUpdate-节点存在-数据存在-更新为空", isE: true, path: "/hydrauodate5", value: "1", nvalue: ""},
		{name: "3.1 LMUpdate-节点存在-数据存在-更新为数字", isE: true, path: "/hydrauodate6", value: "1", nvalue: "2333"},
		{name: "3.1 LMUpdate-节点存在-数据存在-更新为字符", isE: true, path: "561233gdfg", value: "2", nvalue: "sdfd"},
		{name: "3.1 LMUpdate-节点存在-数据存在-更新为中文", isE: true, path: "123hydra212", value: "3", nvalue: "研发"},
		{name: "3.1 LMUpdate-节点存在-数据存在-更新为特殊字符", isE: true, path: "1232hydra#**&$%", value: "4", nvalue: "研发12312@#@"},
		{name: "3.1 LMUpdate-节点存在-数据存在-更新为json", isE: true, path: "#$8787%", value: "5", nvalue: `{"abc":"ef",age:[10,20]}`},
		{name: "3.1 LMUpdate-节点存在-数据存在-更新为xml", isE: true, path: "/hydrauodate9/apiserver/api/conf", value: "5", nvalue: `<xml><node id="abc"/></xml>`},
	}
	//创建节点,更新节点
	for _, c := range cases {
		if c.isE {
			err := lm.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)

			data, v, err := lm.GetValue(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.Equal(t, c.value, string(data), c.name)
			assert.Equal(t, true, v > 0, c.name)
		}

		err := lm.Update(c.path, c.nvalue)
		if !c.isE {
			assert.Equal(t, true, strings.Contains(err.Error(), "不存在"), err.Error())
		} else {
			assert.Equal(t, nil, err, c.name)
		}
	}

	//检查节点值是否正确
	for _, c := range cases {
		data, v, err := lm.GetValue(c.path)
		if !c.isE {
			assert.Equal(t, true, strings.Contains(err.Error(), "不存在"), c.name+err.Error())
		} else {
			assert.Equal(t, nil, err, c.name)
			assert.NotEqual(t, v, int32(0), c.name)
			assert.Equal(t, string(data), c.nvalue, c.name)
		}
	}
}

func TestRedisExists(t *testing.T) {

	//构建所有注册中心
	lm, err := getRegistry()
	assert.Equal(t, nil, err)

	var cases = []struct {
		name  string
		ctype string
		path  string
		value string
	}{
		{name: "1.1 LMEXists-永久节点-一段路径-有后/", ctype: "1", path: "hydddrha#$%xee/", value: "7"},
		{name: "1.2 LMEXists-永久节点-一段路径-有前后/", ctype: "1", path: "/hydgghra#$%/", value: "8"},
		{name: "1.3 LMEXists-永久节点-一段路径-长路径", ctype: "1", path: "/hydradfddfabcefgjijkfsnopqrstuvwxyz", value: "1445"},
		{name: "1.4 LMEXists-永久节点-二段路径-有后/", ctype: "1", path: "hydrad4/abc/", value: "181"},
		{name: "1.5 LMEXists-永久节点-二段路径-前后/", ctype: "1", path: "/hydrae5/454/", value: "173"},
		{name: "1.6 LMEXists-永久节点-二段路径-前/", ctype: "1", path: "/hydraf6/#$#%", value: "189x"},

		{name: "2.1 LMEXists-tmp节点-一段路径-有后/", ctype: "2", path: "hydddrha#$%xee1/", value: "7"},
		{name: "2.2 LMEXists-tmp节点-一段路径-有前后/", ctype: "2", path: "/hydgghra#$%1/", value: "8"},
		{name: "2.3 LMEXists-tmp节点-一段路径-长路径", ctype: "2", path: "/hydradfddfabcefgjijkfsnopqrstuvwxy1z", value: "1445"},
		{name: "2.4 LMEXists-tmp节点-二段路径-有后/", ctype: "2", path: "hydrad41/abc1/", value: "181"},
		{name: "2.5 LMEXists-tmp节点-二段路径-前后/", ctype: "2", path: "/hydrae51/454/", value: "173"},
		{name: "2.6 LMEXists-tmp节点-二段路径-前/", ctype: "2", path: "/hydraf61/#$#%", value: "189x"},

		{name: "3.1 LMEXists-seq节点-一段路径-有后/", ctype: "3", path: "hydddrha#$%xee2/", value: "7"},
		{name: "3.2 LMEXists-seq节点-一段路径-有前后/", ctype: "3", path: "/hydgghra#$%2/", value: "8"},
		{name: "3.3 LMEXists-seq节点-一段路径-长路径", ctype: "3", path: "/hydra2dfddfabcefgjijkLMnopqrstuvwxyz", value: "1445"},
		{name: "3.4 LMEXists-seq节点-二段路径-有后/", ctype: "3", path: "hydrad42/abc/", value: "181"},
		{name: "3.5 LMEXists-seq节点-二段路径-前后/", ctype: "3", path: "/hydrae52/454/", value: "173"},
		{name: "3.6 LMEXists-seq节点-二段路径-前/", ctype: "3", path: "/hydraf62/#$#%", value: "189x"},
	}

	//按注册中心进行测试
	for _, c := range cases {

		exists := false
		b, err := lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, exists, b, c.name)

		//创建节点
		switch c.ctype {
		case "1":
			err = lm.CreatePersistentNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
		case "2":
			err = lm.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
		case "3":
			rpath, err := lm.CreateSeqNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
			c.path = rpath
		default:
			assert.Equal(t, true, false, c.name, "用列的ctype类型错误，只能是1.2.3")
		}

		//节点应存在
		exists = true
		b, err = lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, b, exists, c.name)

		err = lm.Delete(c.path)
		assert.Equal(t, nil, err, c.name)

	}
}
func TestRedisDelete(t *testing.T) {

	//构建所有注册中心
	lm, err := getRegistry()
	assert.Equal(t, nil, err)
	var cases = []struct {
		name  string
		isE   bool
		ctype string
		path  string
		value string
	}{
		{name: "1.1 LMDelete-节点不存在-删除数据", isE: false, path: "/hydraufgggte1", value: "1"},
		{name: "1.2 LMDelete-永久节点删除-二段路径-有后/", isE: true, ctype: "1", path: "hyeead4/abc/", value: "181"},
		{name: "1.3 LMDelete-永久节点删除-二段路径-前后/", isE: true, ctype: "1", path: "/hyrrae5/454/", value: "173"},
		{name: "1.4 LMDelete-永久节点删除-二段路径-前/", isE: true, ctype: "1", path: "/hydrtt6/#$#%", value: "189x"},

		{name: "2.1 LMDelete-节点不存在-删除数据", isE: false, path: "/hydrauodate1", value: "1"},
		{name: "2.2 LMDelete-tmp节点删除-二段路径-有后/", isE: true, ctype: "2", path: "hydhff41/abc1/", value: "181"},
		{name: "2.3 LMDelete-tmp节点删除-二段路径-前后/", isE: true, ctype: "2", path: "/hyhjhe51/454/", value: "173"},
		{name: "2.4 LMDelete-tmp节点删除-二段路径-前/", isE: true, ctype: "2", path: "/hydrkjk1/#$#%", value: "189x"},

		{name: "3.1 LMDelete-节点不存在-删除数据", isE: false, path: "/hydrauodate1", value: "1"},
		{name: "3.2 LMDelete-seq节点删除-二段路径-有后/", isE: true, ctype: "3", path: "hydsds42/abc/", value: "181"},
		{name: "3.3 LMDelete-seq节点删除-二段路径-前后/", isE: true, ctype: "3", path: "/hyuuue52/454/", value: "173"},
		{name: "3.4 LMDelete-seq节点删除-二段路径-前/", isE: true, ctype: "3", path: "/hydrbbb2/#$#%", value: "189x"},
	}

	//按注册中心进行测试
	exists := false
	for _, c := range cases {
		//创建节点

		if c.isE {
			//创建节点
			switch c.ctype {
			case "1":
				err := lm.CreatePersistentNode(c.path, c.value)
				assert.Equal(t, nil, err, c.name)
			case "2":
				err := lm.CreateTempNode(c.path, c.value)
				assert.Equal(t, nil, err, c.name)
			case "3":
				rpath, err := lm.CreateSeqNode(c.path, c.value)
				assert.Equal(t, nil, err, c.name)
				c.path = rpath
			default:
				assert.Equal(t, true, false, c.name, "用列的ctype类型错误，只能是1.2.3")
			}

			//判断数据是否添加成功
			b, err := lm.Exists(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.Equal(t, true, b, c.name)
		}

		//删除节点
		err := lm.Delete(c.path)
		assert.Equal(t, nil, err, c.name)

		//是否存在
		b, err := lm.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, b, exists, c.name)
	}
}

func TestRedisChildren(t *testing.T) {

	lm, err := getRegistry()
	assert.Equal(t, nil, err)
	cases := []struct {
		name     string
		path     string
		children []string
		value    string
	}{
		{name: "1.1 单级目录-无子节点", path: "/hydra3", value: "1"},
		{name: "1.2 单级目录-一个子节点", path: "/hydra1", value: "1", children: []string{"efg"}},
		{name: "1.3 单级目录-一个子节点,节点名包含", path: "/hydr", value: "1", children: []string{"abc"}},
		{name: "1.4 单级目录-多个子节点", path: "/hydra2", value: "1", children: []string{"abcccc", "efg", "efss", "12", "!@#"}},
		{name: "1.5 单级目录-多个子节点,子节点名包含", path: "/hydr", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},

		{name: "2.1 二级目录-无子节点", path: "/hydra3/x1", value: "1"},
		{name: "2.2 二级目录-一个子节点", path: "/hydra1/x2", value: "1", children: []string{"efg"}},
		{name: "2.3 二级目录-多个子节点", path: "/hydra2/x3", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},

		{name: "3.1 多级目录-无子节点", path: "/hydra3/x1/xx", value: "1"},
		{name: "3.2 多级目录-一个子节点", path: "/hydra1/x2/xd/cd", value: "1", children: []string{"efg"}},
		{name: "3.3 多级目录-多个子节点", path: "/hydra2/x3/x/c/v/b/n", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},
	}

	for _, c := range cases {

		//创建节点
		for _, ch := range c.children {
			err := lm.CreateTempNode(registry.Join(c.path, ch), c.value)
			assert.Equal(t, nil, err, c.name)
		}

		if len(c.children) == 0 {
			continue
		}

		//获取子节点
		paths, v, err := lm.GetChildren(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.NotEqual(t, v, 0, c.name)
		assert.Equal(t, len(c.children), len(paths), paths)

		//排序列表
		sort.Strings(paths)
		sort.Strings(c.children)
		assert.Equal(t, c.children, paths, c.name)

		for _, ch := range c.children {
			err := lm.Delete(registry.Join(c.path, ch))
			assert.Equal(t, nil, err, c.name)
		}

	}
}

func TestRedisWatchValue(t *testing.T) {
	//构建所有注册中心
	lm, err := getRegistry()
	assert.Equal(t, nil, err)

	//按注册中心进行测试
	cases := []struct {
		name   string
		path   string
		isE    bool
		isD    bool
		value  string
		nvalue string
	}{
		{name: "1. LMWatch-删除节点", isE: true, isD: true, path: "/hydr12a1", value: "1", nvalue: "1"},
		{name: "2. LMWatch-新增节点", isE: false, isD: false, path: "/hyd33ra1", value: "", nvalue: "2"},
		{name: "3. LMWatch-修改节点数据,有数据变更", isE: true, isD: false, path: "/hyd434ra2", value: "2", nvalue: "234"},
		{name: "4. LMWatch-修改节点数据,空->有", isE: true, isD: false, path: "/hydr545a2", value: "", nvalue: "234"},
		{name: "5. LMWatch-修改节点数据,有->空", isE: true, isD: false, path: "/hy565dra2", value: "2", nvalue: ""},
		{name: "6. LMWatch-修改节点数据，值不变", isE: true, isD: false, path: "/hyd678ra2", value: "234", nvalue: "234"},
	}
	for _, c := range cases {
		//创建节点
		if c.isE {
			err := lm.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
		}

		//监控值变化
		notify, err := lm.WatchValue(c.path)
		assert.Equal(t, nil, err, c.name)
		//此时值未变化不应收到通知
		go func(c chan r.ValueWatcher, name, nvalue string, isD, isE bool) {
			select {
			case v := <-c:
				if !isE {
					assert.Equal(t, true, true, name, "新增是没有值通知的，所以不合法")
				}

				assert.Equal(t, nil, v.GetError(), name)

				value, version := v.GetValue()
				assert.Equal(t, true, version > 0, name, version)
				assert.Equal(t, nvalue, string(value), name, nvalue, string(value))
			case <-time.After(30 * time.Second):
				//如果是新增   是没有通知的
				if isE {
					assert.Equal(t, false, true, name, "通知没有及时回来，测试未通过--")
				}
			}
		}(notify, c.name, c.nvalue, c.isD, c.isE)

		if !c.isE {
			err := lm.CreateTempNode(c.path, c.nvalue)
			assert.Equal(t, nil, err, c.name)
		} else {
			if c.isD {
				err := lm.Delete(c.path)
				assert.Equal(t, nil, err, c.name)
			} else {
				err := lm.Update(c.path, c.nvalue)
				assert.Equal(t, nil, err, c.name)
			}
		}
		lm.Delete(c.path)
	}

	time.Sleep(time.Second)
}

func TestRedisWatchChildren(t *testing.T) {
	lm, err := getRegistry()
	assert.Equal(t, nil, err)

	//按注册中心进行测试
	cases := []struct {
		name     string
		path     string
		isE      bool
		isD      bool
		children []string
		value    string
		nvalue   string
	}{
		{name: "1. LMWatch-删除节点", isE: true, isD: true, path: "/hydr12a1", value: "1", nvalue: "1", children: []string{"efg"}},
		{name: "2. LMWatch-新增节点", isE: false, isD: false, path: "/hyd33ra1", value: "", nvalue: "2", children: []string{"efg"}},
	}

	for _, c := range cases {
		//创建节点
		if c.isE {
			err := lm.CreateTempNode(registry.Join(c.path, c.children[0]), c.value)
			assert.Equal(t, nil, err, c.name)
		}

		//监控父节点
		notify, err := lm.WatchChildren(c.path)
		assert.Equal(t, nil, err, c.name)
		//此时值未变化不应收到通知
		go func(c chan r.ChildrenWatcher, name, path string, paths []string, isD, isE bool) {
			select {
			case v := <-c:
				if isE && !isD {
					assert.Equal(t, true, true, name, "修改数据不会有父级节点通知，所以不合法")
				}

				cPath, cVersion := v.GetValue()
				assert.NotEqual(t, int32(0), cVersion, name)
				assert.Equal(t, []string{registry.Join(path, paths[0])}, cPath, name)
				assert.Equal(t, v.GetPath(), path, name)
			case <-time.After(20 * time.Second):
				//如果是新增   是没有通知的
				if !(isE && !isD) {
					assert.Equal(t, false, true, name, "通知没有及时回来，测试未通过")
				}
			}
		}(notify, c.name, c.path, c.children, c.isD, c.isE)

		if !c.isE {
			err := lm.CreatePersistentNode(c.path, c.nvalue)
			assert.Equal(t, nil, err, c.name)
		} else {
			if c.isD {
				err := lm.Delete(c.path)
				assert.Equal(t, nil, err, c.name)
			} else {
				err := lm.Update(c.path, c.nvalue)
				assert.Equal(t, nil, err, c.name)
			}
		}
		err = lm.Delete(registry.Join(c.path, c.children[0]))
		assert.Equal(t, nil, err, c.name)

		err = lm.Delete(c.path)
		assert.Equal(t, nil, err, c.name)
	}
}
