package registry

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/hydra/registry/registry/filesystem"
	"github.com/micro-plat/hydra/test/assert"
	r "github.com/micro-plat/lib4go/registry"
	"github.com/micro-plat/lib4go/types"
)

func checkData(r registry.IRegistry, path string, data map[string]string) error {
	bd, v, err := r.GetValue(path)
	if err != nil {
		return err
	}

	if v <= 0 {
		return fmt.Errorf("节点版本号错误，path：%s", path)
	}

	str := string(bd)
	if _, ok := data[path]; ok {
		if !strings.EqualFold(str, data[path]) {
			return fmt.Errorf("数据不合法,str:%s, data[path]:%s", str, data[path])
		}
	} else if !(str == "" || str == "{}") {
		return fmt.Errorf("数据不合法1,str:%s", str)
	}

	paths, v, err := r.GetChildren(path)
	fmt.Println("xxxxxxxxxx:", paths)
	if err != nil {
		return err
	}

	if v <= 0 {
		return fmt.Errorf("节点版本号错误1，path：%s", path)
	}

	if len(paths) == 0 {
		return nil
	}

	for _, xpath := range paths {
		xpath = registry.Join(path, xpath)
		return checkData(r, xpath, data)
	}

	return nil
}

func getPaths(path string) []string {
	nodes := strings.Split(strings.Trim(path, "/"), "/")
	len := len(nodes)
	paths := make([]string, 0, len)
	for i := 0; i < len; i++ {
		npath := "/" + strings.Join(nodes[:i+1], "/")
		paths = append(paths, npath)
	}
	return paths
}

func TestFSCreateTempNode(t *testing.T) {

	//构建所有注册中心
	fs, _ := filesystem.NewFileSystem(".")
	var fscases = []struct {
		name  string
		path  string
		value string
	}{
		{name: "1.1 FSCreateTemp-一段路径-字母", path: "hydrax", value: "1"},
		{name: "1.2 FSCreateTemp-一段路径-数字", path: "1231222", value: "2"},
		{name: "1.3 FSCreateTemp-一段路径-字母数字混合", path: "123hydra3", value: "3"},
		{name: "1.4 FSCreateTemp-一段路径-含特殊字符", path: "1232hydra#$%", value: "4"},
		{name: "1.5 FSCreateTemp-一段路径-全特殊字符", path: "#$%", value: "5"},
		{name: "1.6 FSCreateTemp-一段路径-有前/", path: "/123123", value: "6"},
		{name: "1.7 FSCreateTemp-一段路径-有后/", path: "hydra#$%xee/", value: "7"},
		{name: "1.8 FSCreateTemp-一段路径-有前后/", path: "/hydra#$%/", value: "8"},
		{name: "1.9 FSCreateTemp-一段路径-长路径", path: "/hydraabcefgjijkfsnopqrstuvwxyz", value: "1445"},

		{name: "2.1 FSCreateTemp-二段路径-以段以上路径", path: "/hydra1/abc/", value: "18"},
		{name: "2.2 FSCreateTemp-二段路径-多段有数字", path: "/hydra2/454/", value: "17"},
		{name: "2.3 FSCreateTemp-二段路径-多段有特殊字符", path: "/hydra3/#$#%/", value: "189"},
		{name: "2.4 FSCreateTemp-二段路径-有后/", path: "hydra4/abc/", value: "181"},
		{name: "2.5 FSCreateTemp-二段路径-前后/", path: "/hydra5/454/", value: "173"},
		{name: "2.6 FSCreateTemp-二段路径-前/", path: "/hydra6/#$#%", value: "189x"},

		{name: "3.1 FSCreateTemp-多段-较长分段", path: "/hydra11/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/xxx", value: "1255"},
		{name: "3.2 FSCreateTemp-多段-较长分段1", path: "hydra22/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/yyy", value: "12225"},
	}

	//按注册中心进行测试
	//创建节点
	for _, c := range fscases {
		err := fs.CreateTempNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)
		mp := map[string]string{c.path: c.value}
		err = checkData(fs, c.path, mp)
		assert.Equal(t, nil, err, c.name)
	}

	fs.Close()

	fs, _ = filesystem.NewFileSystem(".")
	//检查临时节点是否删除
	for _, c := range fscases {
		ok, _ := fs.Exists(c.path)
		assert.Equal(t, false, ok, "临时节点没有删除掉")
	}

	//删除文件
	for _, c := range fscases {
		paths := getPaths(c.path)
		if len(paths) > 0 {
			fs.Delete(paths[0])
		} else {
			fs.Delete(c.path)
		}
	}

	fs.Close()
}

func TestFSCreateSeqNode(t *testing.T) {

	//构建所有注册中心
	fs, _ := filesystem.NewFileSystem(".")
	var fscases = []struct {
		name  string
		path  string
		value string
	}{
		{name: "1.1 FSCreateSeq-一段路径-字母", path: "hydray", value: "1"},
		{name: "1.2 FSCreateSeq-一段路径-数字", path: "1231452", value: "2"},
		{name: "1.3 FSCreateSeq-一段路径-字母数字混合", path: "123hssydra3", value: "3"},
		{name: "1.4 FSCreateSeq-一段路径-含特殊字符", path: "1232hyffdra#$%", value: "4"},
		{name: "1.5 FSCreateSeq-一段路径-全特殊字符", path: "#@@$%", value: "5"},
		{name: "1.6 FSCreateSeq-一段路径-有前/", path: "/12315623", value: "6"},
		{name: "1.7 FSCreateSeq-一段路径-有后/", path: "hydddra#$%xee/", value: "7"},
		{name: "1.8 FSCreateSeq-一段路径-有前后/", path: "/hydggra#$%/", value: "8"},
		{name: "1.9 FSCreateSeq-一段路径-长路径", path: "/hydradfdfabcefgjijkfsnopqrstuvwxyz", value: "1445"},

		{name: "2.1 FSCreateSeq-二段路径-以段以上路径", path: "/hydraa/abc/", value: "18"},
		{name: "2.2 FSCreateSeq-二段路径-多段有数字", path: "/hydrab/454/", value: "17"},
		{name: "2.3 FSCreateSeq-二段路径-多段有特殊字符", path: "/hydrac/#$#%/", value: "189"},
		{name: "2.4 FSCreateSeq-二段路径-有后/", path: "hydrad/abc/", value: "181"},
		{name: "2.5 FSCreateSeq-二段路径-前后/", path: "/hydrae/454/", value: "173"},
		{name: "2.6 FSCreateSeq-二段路径-前/", path: "/hydraf/#$#%", value: "189x"},

		{name: "3.1 FSCreateSeq-多段-较长分段", path: "/hydrag/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/xxx", value: "1255"},
		{name: "3.2 FSCreateSeq-多段-较长分段1", path: "hydrah/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/yyy", value: "12225"},
	}

	//按注册中心进行测试
	//创建节点
	rpaths := []string{}
	for _, c := range fscases {
		rpath, err := fs.CreateSeqNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)
		ca := strings.Split(rpath, "_")
		assert.Equal(t, 2, len(ca), c.name, "返回的路径不合法")
		assert.Equal(t, true, types.GetInt(ca[1]) > 0, c.name, "返回的路径不合法1")
		mp := map[string]string{rpath: c.value}
		err = checkData(fs, rpath, mp)
		assert.Equal(t, nil, err, c.name)
		rpaths = append(rpaths, rpath)
	}

	fs.Close()

	fs, _ = filesystem.NewFileSystem(".")
	//检查临时节点是否删除
	for _, p := range rpaths {
		ok, _ := fs.Exists(p)
		assert.Equal(t, false, ok, "seq临时节点没有删除掉")
	}

	//删除文件
	for _, p := range rpaths {
		paths := getPaths(p)
		if len(paths) > 0 {
			fs.Delete(paths[0])
		} else {
			fs.Delete(p)
		}
	}

	fs.Close()
}

func TestFSCreatePersistentNode(t *testing.T) {

	//构建所有注册中心
	fs, _ := filesystem.NewFileSystem(".")
	var fscases = []struct {
		name  string
		path  string
		value string
	}{
		{name: "1.1 FSCreatePersistent-一段路径-字母", path: "hydrawy", value: "1"},
		{name: "1.2 FSCreatePersistent-一段路径-数字", path: "12314532", value: "2"},
		{name: "1.3 FSCreatePersistent-一段路径-字母数字混合", path: "123hsdfsydra3", value: "3"},
		{name: "1.4 FSCreatePersistent-一段路径-含特殊字符", path: "1232hdfyffdra#$%", value: "4"},
		{name: "1.5 FSCreatePersistent-一段路径-全特殊字符", path: "#@@%%$%", value: "5"},
		{name: "1.6 FSCreatePersistent-一段路径-有前/", path: "/1231567623", value: "6"},
		{name: "1.7 FSCreatePersistent-一段路径-有后/", path: "hydddrha#$%xee/", value: "7"},
		{name: "1.8 FSCreatePersistent-一段路径-有前后/", path: "/hydgghra#$%/", value: "8"},
		{name: "1.9 FSCreatePersistent-一段路径-长路径", path: "/hydradfddfabcefgjijkfsnopqrstuvwxyz", value: "1445"},

		{name: "2.1 FSCreatePersistent-二段路径-以段以上路径", path: "/hydraa1/abcd/", value: "18"},
		{name: "2.2 FSCreatePersistent-二段路径-多段有数字", path: "/hydrab2/454/", value: "17"},
		{name: "2.3 FSCreatePersistent-二段路径-多段有特殊字符", path: "/hydrac3/#$#%/", value: "189"},
		{name: "2.4 FSCreatePersistent-二段路径-有后/", path: "hydrad4/abc/", value: "181"},
		{name: "2.5 FSCreatePersistent-二段路径-前后/", path: "/hydrae5/454/", value: "173"},
		{name: "2.6 FSCreatePersistent-二段路径-前/", path: "/hydraf6/#$#%", value: "189x"},

		{name: "3.1 FSCreatePersistent-多段-较长分段", path: "/hydrag7/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/xxx", value: "1255"},
		{name: "3.2 FSCreatePersistent-多段-较长分段1", path: "hydrah8/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z/yyy", value: "12225"},
	}

	//按注册中心进行测试
	//创建节点
	for _, c := range fscases {
		err := fs.CreatePersistentNode(c.path, c.value)
		assert.Equal(t, nil, err, c.name)
		mp := map[string]string{c.path: c.value}
		err = checkData(fs, c.path, mp)
		assert.Equal(t, nil, err, c.name)
	}

	fs.Close()

	fs, _ = filesystem.NewFileSystem(".")
	//检查永久节点是否删除
	for _, p := range fscases {
		ok, _ := fs.Exists(p.path)
		assert.Equal(t, true, ok, "永久节点被删除掉")
	}

	//删除文件
	for _, p := range fscases {
		paths := getPaths(p.path)
		if len(paths) > 0 {
			fs.Delete(paths[0])
		} else {
			fs.Delete(p.path)
		}
	}
	fs.Close()
}

func TestFSUpdateNode(t *testing.T) {
	fscases := []struct {
		name   string
		isE    bool
		path   string
		value  string
		nvalue string
	}{
		{name: "1.1 FSUpdate-节点不存在-更新数据", isE: false, path: "/hydrauodate1", value: "1", nvalue: "2333"},

		{name: "2.1 FSUpdate-节点存在-数据为空-更新为空", isE: true, path: "/hydrauodate2", value: "1", nvalue: "2333"},
		{name: "2.1 FSUpdate-节点存在-数据为空-更新为数字", isE: true, path: "/hydrauodate3", value: "1", nvalue: "2333"},
		{name: "2.1 FSUpdate-节点存在-数据为空-更新为字符", isE: true, path: "sdsd343434", value: "2", nvalue: "sdfd"},
		{name: "2.1 FSUpdate-节点存在-数据为空-更新为中文", isE: true, path: "123hydraddd", value: "3", nvalue: "研发"},
		{name: "2.1 FSUpdate-节点存在-数据为空-更新为特殊字符", isE: true, path: "1232hydrqqa#$%", value: "4", nvalue: "研发12312@#@"},
		{name: "2.1 FSUpdate-节点存在-数据为空-更新为json", isE: true, path: "#$@×&#(%", value: "5", nvalue: `{"abc":"ef",age:[10,20]}`},
		{name: "2.1 FSUpdate-节点存在-数据为空-更新为xml", isE: true, path: "/hydrauodate/apiserver/api/conf", value: "5", nvalue: `<xml><node id="abc"/></xml>`},

		{name: "3.1 FSUpdate-节点存在-数据存在-更新为空", isE: true, path: "/hydrauodate5", value: "1", nvalue: ""},
		{name: "3.1 FSUpdate-节点存在-数据存在-更新为数字", isE: true, path: "/hydrauodate6", value: "1", nvalue: "2333"},
		{name: "3.1 FSUpdate-节点存在-数据存在-更新为字符", isE: true, path: "561233gdfg", value: "2", nvalue: "sdfd"},
		{name: "3.1 FSUpdate-节点存在-数据存在-更新为中文", isE: true, path: "123hydra212", value: "3", nvalue: "研发"},
		{name: "3.1 FSUpdate-节点存在-数据存在-更新为特殊字符", isE: true, path: "1232hydra#**&$%", value: "4", nvalue: "研发12312@#@"},
		{name: "3.1 FSUpdate-节点存在-数据存在-更新为json", isE: true, path: "#$8787%", value: "5", nvalue: `{"abc":"ef",age:[10,20]}`},
		{name: "3.1 FSUpdate-节点存在-数据存在-更新为xml", isE: true, path: "/hydrauodate9/apiserver/api/conf", value: "5", nvalue: `<xml><node id="abc"/></xml>`},
	}
	//构建所有注册中心
	fs, _ := filesystem.NewFileSystem(".")

	//按注册中心进行测试
	//创建节点,更新节点
	for _, c := range fscases {
		if c.isE {
			err := fs.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)

			data, v, err := fs.GetValue(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.Equal(t, c.value, string(data), c.name)
			assert.Equal(t, true, v > 0, c.name)
		}

		err := fs.Update(c.path, c.nvalue)
		if !c.isE {
			assert.Equal(t, true, strings.Contains(err.Error(), "不存在"), c.name)
		} else {
			assert.Equal(t, nil, err, c.name)
		}
	}

	//检查节点值是否正确
	for _, c := range fscases {
		data, v, err := fs.GetValue(c.path)
		if !c.isE {
			assert.Equal(t, true, strings.Contains(err.Error(), "不存在"), c.name)
		} else {
			assert.Equal(t, nil, err, c.name)
			assert.NotEqual(t, v, int32(0), c.name)
			assert.Equal(t, string(data), c.nvalue, c.name)
		}
	}

	//删除文件
	for _, p := range fscases {
		paths := getPaths(p.path)
		if len(paths) > 0 {
			fs.Delete(paths[0])
		} else {
			fs.Delete(p.path)
		}
	}

	//关闭注册中心
	fs.Close()
}

func TestFSExists(t *testing.T) {

	//构建所有注册中心
	fs, _ := filesystem.NewFileSystem(".")

	var fscases = []struct {
		name  string
		ctype string
		path  string
		value string
	}{
		{name: "1.1 FSEXists-永久节点-一段路径-有后/", ctype: "1", path: "hydddrha#$%xee/", value: "7"},
		{name: "1.2 FSEXists-永久节点-一段路径-有前后/", ctype: "1", path: "/hydgghra#$%/", value: "8"},
		{name: "1.3 FSEXists-永久节点-一段路径-长路径", ctype: "1", path: "/hydradfddfabcefgjijkfsnopqrstuvwxyz", value: "1445"},
		{name: "1.4 FSEXists-永久节点-二段路径-有后/", ctype: "1", path: "hydrad4/abc/", value: "181"},
		{name: "1.5 FSEXists-永久节点-二段路径-前后/", ctype: "1", path: "/hydrae5/454/", value: "173"},
		{name: "1.6 FSEXists-永久节点-二段路径-前/", ctype: "1", path: "/hydraf6/#$#%", value: "189x"},

		{name: "2.1 FSEXists-tmp节点-一段路径-有后/", ctype: "2", path: "hydddrha#$%xee1/", value: "7"},
		{name: "2.2 FSEXists-tmp节点-一段路径-有前后/", ctype: "2", path: "/hydgghra#$%1/", value: "8"},
		{name: "2.3 FSEXists-tmp节点-一段路径-长路径", ctype: "2", path: "/hydradfddfabcefgjijkfsnopqrstuvwxy1z", value: "1445"},
		{name: "2.4 FSEXists-tmp节点-二段路径-有后/", ctype: "2", path: "hydrad41/abc1/", value: "181"},
		{name: "2.5 FSEXists-tmp节点-二段路径-前后/", ctype: "2", path: "/hydrae51/454/", value: "173"},
		{name: "2.6 FSEXists-tmp节点-二段路径-前/", ctype: "2", path: "/hydraf61/#$#%", value: "189x"},

		{name: "3.1 FSEXists-seq节点-一段路径-有后/", ctype: "3", path: "hydddrha#$%xee2/", value: "7"},
		{name: "3.2 FSEXists-seq节点-一段路径-有前后/", ctype: "3", path: "/hydgghra#$%2/", value: "8"},
		{name: "3.3 FSEXists-seq节点-一段路径-长路径", ctype: "3", path: "/hydra2dfddfabcefgjijkfsnopqrstuvwxyz", value: "1445"},
		{name: "3.4 FSEXists-seq节点-二段路径-有后/", ctype: "3", path: "hydrad42/abc/", value: "181"},
		{name: "3.5 FSEXists-seq节点-二段路径-前后/", ctype: "3", path: "/hydrae52/454/", value: "173"},
		{name: "3.6 FSEXists-seq节点-二段路径-前/", ctype: "3", path: "/hydraf62/#$#%", value: "189x"},
	}

	for _, c := range fscases {
		//节点不存在
		exists := false
		b, err := fs.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, b, exists, c.name)

		//创建节点
		switch c.ctype {
		case "1":
			err = fs.CreatePersistentNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
		case "2":
			err = fs.CreateTempNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
		case "3":
			rpath, err := fs.CreateSeqNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
			c.path = rpath
		default:
			assert.Equal(t, true, false, c.name, "用列的ctype类型错误，只能是1.2.3")
		}

		//节点应存在
		exists = true
		b, err = fs.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, b, exists, c.name)
	}

	//删除文件
	for _, p := range fscases {
		paths := getPaths(p.path)
		if len(paths) > 0 {
			fs.Delete(paths[0])
		} else {
			fs.Delete(p.path)
		}
	}

	fs.Close()
}

func TestFsDelete(t *testing.T) {

	//构建所有注册中心
	fs, _ := filesystem.NewFileSystem(".")

	var fscases = []struct {
		name  string
		isE   bool
		ctype string
		path  string
		value string
	}{
		{name: "1.1 FsDelete-节点不存在-删除数据", isE: false, path: "/hydraufgggte1", value: "1"},
		{name: "1.2 FsDelete-永久节点删除-二段路径-有后/", isE: true, ctype: "1", path: "hyeead4/abc/", value: "181"},
		{name: "1.3 FsDelete-永久节点删除-二段路径-前后/", isE: true, ctype: "1", path: "/hyrrae5/454/", value: "173"},
		{name: "1.4 FsDelete-永久节点删除-二段路径-前/", isE: true, ctype: "1", path: "/hydrtt6/#$#%", value: "189x"},

		{name: "2.1 FsDelete-节点不存在-删除数据", isE: false, path: "/hydrauodate1", value: "1"},
		{name: "2.2 FsDelete-tmp节点删除-二段路径-有后/", isE: true, ctype: "2", path: "hydhff41/abc1/", value: "181"},
		{name: "2.3 FsDelete-tmp节点删除-二段路径-前后/", isE: true, ctype: "2", path: "/hyhjhe51/454/", value: "173"},
		{name: "2.4 FsDelete-tmp节点删除-二段路径-前/", isE: true, ctype: "2", path: "/hydrkjk1/#$#%", value: "189x"},

		{name: "3.1 FsDelete-节点不存在-删除数据", isE: false, path: "/hydrauodate1", value: "1"},
		{name: "3.2 FsDelete-seq节点删除-二段路径-有后/", isE: true, ctype: "3", path: "hydsds42/abc/", value: "181"},
		{name: "3.3 FsDelete-seq节点删除-二段路径-前后/", isE: true, ctype: "3", path: "/hyuuue52/454/", value: "173"},
		{name: "3.4 FsDelete-seq节点删除-二段路径-前/", isE: true, ctype: "3", path: "/hydrbbb2/#$#%", value: "189x"},
	}

	//按注册中心进行测试
	exists := false
	for _, c := range fscases {
		//创建节点
		if c.isE {
			//创建节点
			switch c.ctype {
			case "1":
				err := fs.CreatePersistentNode(c.path, c.value)
				assert.Equal(t, nil, err, c.name)
			case "2":
				err := fs.CreateTempNode(c.path, c.value)
				assert.Equal(t, nil, err, c.name)
			case "3":
				rpath, err := fs.CreateSeqNode(c.path, c.value)
				assert.Equal(t, nil, err, c.name)
				c.path = rpath
			default:
				assert.Equal(t, true, false, c.name, "用列的ctype类型错误，只能是1.2.3")
			}

			//判断数据是否添加成功
			b, err := fs.Exists(c.path)
			assert.Equal(t, nil, err, c.name)
			assert.Equal(t, true, b, c.name)
		}

		//删除节点
		err := fs.Delete(c.path)
		assert.Equal(t, nil, err, c.name)

		//是否存在
		b, err := fs.Exists(c.path)
		assert.Equal(t, nil, err, c.name)
		assert.Equal(t, b, exists, c.name)
	}

	//删除文件
	for _, p := range fscases {
		paths := getPaths(p.path)
		if len(paths) > 0 {
			fs.Delete(paths[0])
		} else {
			fs.Delete(p.path)
		}
	}

	fs.Close()
}

func TestFSChildren(t *testing.T) {

	//构建所有注册中心
	fs, _ := filesystem.NewFileSystem(".")

	//按注册中心进行测试
	fscases := []struct {
		name     string
		path     string
		children []string
		value    string
	}{
		{name: "1.1 单级目录-无子节点", path: "/hydra3", value: "1"},
		{name: "1.2 单级目录-一个子节点", path: "/hydra1", value: "1", children: []string{"efg"}},
		{name: "1.3 单级目录-多个子节点", path: "/hydra2", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},

		{name: "2.1 二级目录-无子节点", path: "/hydra3/x1", value: "1"},
		{name: "2.2 二级目录-一个子节点", path: "/hydra1/x2", value: "1", children: []string{"efg"}},
		{name: "2.3 二级目录-多个子节点", path: "/hydra2/x3", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},

		{name: "3.1 多级目录-无子节点", path: "/hydra3/x1/xx", value: "1"},
		{name: "3.2 多级目录-一个子节点", path: "/hydra1/x2/xd/cd", value: "1", children: []string{"efg"}},
		{name: "3.3 多级目录-多个子节点", path: "/hydra2/x3/x/c/v/b/n", value: "1", children: []string{"abc", "efg", "efss", "12", "!@#"}},
	}

	for _, c := range fscases {

		//创建节点
		for _, ch := range c.children {
			err := fs.CreateTempNode(registry.Join(c.path, ch), c.value)
			assert.Equal(t, nil, err, c.name)
		}

		//获取子节点
		paths, v, err := fs.GetChildren(c.path)
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

	//删除文件
	for _, p := range fscases {
		paths := getPaths(p.path)
		if len(paths) > 0 {
			fs.Delete(paths[0])
		} else {
			fs.Delete(p.path)
		}
	}

	fs.Close()
}

func TestFSWatchValue(t *testing.T) {
	//构建所有注册中心
	fs, _ := filesystem.NewFileSystem(".")
	//开启文件系统监听
	fs.Start()

	//按注册中心进行测试
	fscases := []struct {
		name   string
		path   string
		isE    bool
		isD    bool
		value  string
		nvalue string
	}{
		// {name: "1. FSWatch-删除节点", isE: true, isD: true, path: "/hydr12a1", value: "1", nvalue: ""},
		// {name: "2. FSWatch-新增节点", isE: false, isD: false, path: "/hyd33ra1", value: "", nvalue: "2"},
		// {name: "3. FSWatch-修改节点数据,有数据变更", isE: true, isD: false, path: "/hyd434ra2", value: "2", nvalue: "234"},
		// {name: "4. FSWatch-修改节点数据,空->有", isE: true, isD: false, path: "/hydr545a2", value: "", nvalue: "234"},
		// {name: "5. FSWatch-修改节点数据,有->空", isE: true, isD: false, path: "/hy565dra2", value: "2", nvalue: ""},
		// {name: "6. FSWatch-修改节点数据，值不变", isE: true, isD: false, path: "/hyd678ra2", value: "234", nvalue: "234"},
	}
	for _, c := range fscases {

		//创建节点
		if c.isE {
			err := fs.CreatePersistentNode(c.path, c.value)
			assert.Equal(t, nil, err, c.name)
		}

		//监控值变化
		notify, err := fs.WatchValue(c.path)
		assert.Equal(t, nil, err, c.name)
		//此时值未变化不应收到通知
		go func(c chan r.ValueWatcher, name, nvalue string, isD bool) {
			select {
			case v := <-c:
				value, _ := v.GetValue()
				if isD {
					assert.NotEqual(t, true, strings.Contains(v.GetError().Error(), "文件发生变化"), name, v.GetError())
				} else {
					// assert.Equal(t, true, version > 0, name, version)
					assert.Equal(t, nvalue, string(value), name, nvalue, string(value))
				}
			case <-time.After(2 * time.Second):
				assert.Equal(t, false, true, name, "通知没有及时回来，测试未通过")
			}
		}(notify, c.name, c.nvalue, c.isD)

		if !c.isE {
			err := fs.CreatePersistentNode(c.path, c.nvalue)
			assert.Equal(t, nil, err, c.name)
		} else {
			if c.isD {
				err := fs.Delete(c.path)
				assert.Equal(t, nil, err, c.name)
			} else {
				err := fs.Update(c.path, c.nvalue)
				assert.Equal(t, nil, err, c.name)
			}
		}
	}

	time.Sleep(5 * time.Second)

	//删除文件
	for _, p := range fscases {
		paths := getPaths(p.path)
		if len(paths) > 0 {
			fs.Delete(paths[0])
		} else {
			fs.Delete(p.path)
		}
	}

	fs.Close()
}
