package creator

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/micro-plat/hydra/conf/server/acl/whitelist"

	"github.com/micro-plat/hydra/conf/server"
	varpub "github.com/micro-plat/hydra/conf/vars"
	"github.com/micro-plat/hydra/conf/vars/http"
	"github.com/micro-plat/hydra/conf/vars/rpc"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	_ "github.com/micro-plat/hydra/registry/registry/filesystem"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	_ "github.com/micro-plat/hydra/registry/registry/zookeeper"
	"github.com/micro-plat/hydra/test/assert"
)

func Test_conf_Pub(t *testing.T) {
	type fields struct {
		data    map[string]iCustomerBuilder
		olddata map[string]iCustomerBuilder
		vars    map[string]map[string]interface{}
		oldvars map[string]map[string]interface{}
	}
	type args struct {
		platName     string
		systemName   string
		clusterName  string
		registryAddr string
		cover        bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		isExsit bool
		wantErr bool
	}{
		//文件系统注册的分支没有测试  因为关系到toml文件发布的问题,暂时没有实现  所以不测试
		{name: "1. 发布时,注册中心地址错误", fields: fields{}, args: args{registryAddr: "errdata:"}, isExsit: false, wantErr: true},
		{name: "2. 发布时,地址正确,空对象,不覆盖", fields: fields{data: map[string]iCustomerBuilder{"api": CustomerBuilder{"main": "123456", "testvar1": "22222"}},
			vars: map[string]map[string]interface{}{"db": map[string]interface{}{"dcc": "545454"}, "cache1": map[string]interface{}{"dccsss": "5454"}}},
			args:    args{registryAddr: "lm://.", platName: "platName1", systemName: "systemName1", clusterName: "clusterName1", cover: false},
			isExsit: false, wantErr: false},
		{name: "3. 发布时,地址正确,空对象,覆盖", fields: fields{data: map[string]iCustomerBuilder{"api": CustomerBuilder{"main": "123456", "testvar1": "22222"}},
			vars: map[string]map[string]interface{}{"db": map[string]interface{}{"dcc": "545454"}, "cache1": map[string]interface{}{"dccsss": "5454"}}},
			args:    args{registryAddr: "lm://.", platName: "platName2", systemName: "systemName2", clusterName: "clusterName2", cover: true},
			isExsit: false, wantErr: false},
		{name: "4. 发布时,地址正确,实体对象,不覆盖", fields: fields{data: map[string]iCustomerBuilder{"api": CustomerBuilder{"main": "123456", "testvar1": "22222"}},
			vars:    map[string]map[string]interface{}{"db": map[string]interface{}{"dcc": "545454"}, "cache1": map[string]interface{}{"dccsss": "5454"}},
			olddata: map[string]iCustomerBuilder{"api": CustomerBuilder{"main": "{}", "testvar1": "{}"}},
			oldvars: map[string]map[string]interface{}{"db": map[string]interface{}{"dcc": "{}"}, "cache1": map[string]interface{}{"dccsss": "{}"}}},
			args:    args{registryAddr: "lm://.", platName: "platName3", systemName: "systemName3", clusterName: "clusterName3", cover: false},
			isExsit: true, wantErr: true},
		{name: "5. 发布时,地址正确,实体对象,覆盖", fields: fields{data: map[string]iCustomerBuilder{"api": CustomerBuilder{"main": "123456", "testvar1": "22222"}},
			vars:    map[string]map[string]interface{}{"db": map[string]interface{}{"dcc": "545454"}, "cache1": map[string]interface{}{"dccsss": "5454"}},
			olddata: map[string]iCustomerBuilder{"api": CustomerBuilder{"main": "{}", "testvar1": "{}"}},
			oldvars: map[string]map[string]interface{}{"db": map[string]interface{}{"dcc": "{}"}, "cache1": map[string]interface{}{"dccsss": "{}"}}},
			args:    args{registryAddr: "lm://.", platName: "platName3", systemName: "systemName3", clusterName: "clusterName3", cover: true},
			isExsit: true, wantErr: false},
	}

	global.Def.ServerTypes = []string{}
	for _, tt := range tests {
		c := &conf{}
		if tt.isExsit {
			c.data = tt.fields.olddata
			c.vars = tt.fields.oldvars
			err := c.Pub(tt.args.platName, tt.args.systemName, tt.args.clusterName, tt.args.registryAddr, true)
			assert.Equal(t, false, err != nil, tt.name+",err")
		}

		c.data = tt.fields.data
		c.vars = tt.fields.vars
		err := c.Pub(tt.args.platName, tt.args.systemName, tt.args.clusterName, tt.args.registryAddr, tt.args.cover)
		assert.Equal(t, tt.wantErr, err != nil, tt.name+",err1")

		rgt, err := registry.GetRegistry("lm://.", global.Def.Log())
		assert.Equal(t, true, err == nil, tt.name+",err2")
		if !tt.isExsit || tt.args.cover {
			for tp, subs := range c.data {
				pub := server.NewServerPub(tt.args.platName, tt.args.systemName, tp, tt.args.clusterName)
				data, _, err := rgt.GetValue(pub.GetServerPath())
				assert.Equalf(t, true, err == nil, tt.name+",err3", err)
				data1, _ := json.Marshal(subs.Map()["main"])
				assert.Equal(t, string(data), string(data1)[1:len(string(data1))-1], tt.name+",data")
				for name, value := range subs.Map() {
					if name == "main" {
						continue
					}

					data, _, err = rgt.GetValue(pub.GetSubConfPath(name))
					assert.Equal(t, true, err == nil, tt.name+",err4")
					data1, _ = json.Marshal(value)
					assert.Equal(t, string(data), string(data1)[1:len(string(data1))-1], tt.name+",data1")
				}
			}

			for tp, subs := range c.vars {
				pub := varpub.NewVarPub(tt.args.platName)
				for k, v := range subs {
					data, _, err := rgt.GetValue(pub.GetVarPath(tp, k))
					assert.Equal(t, true, err == nil, tt.name+",err5")
					data1, _ := json.Marshal(v)
					assert.Equal(t, string(data), string(data1)[1:len(string(data1))-1], tt.name+",data2")
				}
			}
		} else {
			for tp, subs := range tt.fields.olddata {
				pub := server.NewServerPub(tt.args.platName, tt.args.systemName, tp, tt.args.clusterName)
				data, _, err := rgt.GetValue(pub.GetServerPath())
				assert.Equal(t, true, err == nil, tt.name+",err6")
				data1, _ := json.Marshal(subs.Map()["main"])
				assert.Equal(t, string(data), string(data1)[1:len(string(data1))-1], tt.name+",data3")
				for name, value := range subs.Map() {
					if name == "main" {
						continue
					}

					data, _, err = rgt.GetValue(pub.GetSubConfPath(name))
					assert.Equal(t, true, err == nil, tt.name+",err7")
					data1, _ = json.Marshal(value)
					assert.Equal(t, string(data), string(data1)[1:len(string(data1))-1], tt.name+",data4")
				}
			}

			for tp, subs := range tt.fields.oldvars {
				pub := varpub.NewVarPub(tt.args.platName)
				for k, v := range subs {
					data, _, err := rgt.GetValue(pub.GetVarPath(tp, k))
					assert.Equal(t, true, err == nil, tt.name+",err8")
					data1, _ := json.Marshal(v)
					assert.Equal(t, string(data), string(data1)[1:len(string(data1))-1], tt.name+",data5")
				}
			}
		}
	}
}

func Test_conf_Pub1(t *testing.T) {
	data := map[string]string{
		"/platname3/systemname3/api/clustername3/conf":                `{"address":":8585","status":"start"}`,
		"/platname3/systemname3/api/clustername3/conf/acl/white.list": `{"disable":true}`,
		"/platname3/systemname3/api/clustername3/conf/static":         `{"dir":"./src","exclude":["/view/","/views/","/web/",".exe",".so"],"homePage":"index.html","rewriters":["/","index.htm","default.html","default.htm"]}`,
		"/platname3/var/http/httpclient":                              `{"connectionTimeout":10,"requestTimeout":10,"certs":null,"ca":"","proxy":"","keepAlive":true,"trace":false}`,
		"/platname3/var/rpc/rpcclinent":                               `{"connectionTimeout":20,"log":"","sortPrefix":"","tls":null,"balancer":"localfirst"}`,
	}
	Conf.API(":8585").WhiteList(whitelist.WithDisable()).Static()
	Conf.Vars().HTTP("httpclient", http.WithConnTimeout(10), http.WithKeepalive(true))
	Conf.Vars().RPC("rpcclinent", rpc.WithConnectionTimeout(20), rpc.WithLocalFirst())
	global.Def.ServerTypes = []string{"api"}
	Conf.Load()
	tests := []struct {
		name      string
		regstType string
		wantErr   bool
	}{
		{name: "1. 发布到fs系统", regstType: "fs://.", wantErr: true},
		{name: "2. 发布到lm系统", regstType: "lm://.", wantErr: true},
		{name: "3. 发布到zk系统", regstType: "zk://192.168.0.101", wantErr: true},
	}

	for _, tt := range tests {
		err := Conf.Pub("platname3", "systemname3", "clustername3", tt.regstType, true)
		assert.Equal(t, tt.wantErr, err == nil, "发布异常", err)
		r, err := registry.GetRegistry(tt.regstType, global.Def.Log())
		assert.Equal(t, true, err == nil, "获取注册中心异常", err)
		err = checkData(r, registry.Join("platname3"), data)
		assert.Equal(t, true, err == nil, "获取注册中心异常", err)
		r.Delete(registry.Join("platname3"))
	}
}

func checkData(r registry.IRegistry, path string, data map[string]string) error {
	bd, _, err := r.GetValue(path)
	if err != nil {
		return err
	}

	str := string(bd)
	if _, ok := data[path]; ok {
		if !strings.EqualFold(str, data[path]) {
			return fmt.Errorf("数据不合法,str:%s, data[path]:%s", str, data[path])
		}
	} else if !(str == "" || str == "{}") {
		return fmt.Errorf("数据不合法1,str:%s", str)
	}

	paths, _, err := r.GetChildren(path)
	if err != nil {
		return err
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

func Test_publish(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		v       interface{}
		isExsit bool
		cover   bool
		wantErr bool
	}{
		{name: "1. 空对象,不覆盖数据", path: "/path/x1/y1", v: `{"testdata":"1"}`, isExsit: false, cover: false, wantErr: false},
		{name: "2. 空对象,覆盖数据", path: "/path/x2/y2", v: `{"testdata":"2"}`, isExsit: false, cover: true, wantErr: false},
		{name: "3. 实体对象,不覆盖数据", path: "/path/x3/y3", v: `{"testdata":"3"}`, isExsit: true, cover: false, wantErr: true},
		{name: "4. 实体对象,覆盖数据", path: "/path/x4/y4", v: `{"testdata":"4"}`, isExsit: true, cover: true, wantErr: false},
	}
	for _, tt := range tests {
		rgt, err := registry.GetRegistry("lm://.", global.Def.Log())
		assert.Equal(t, true, err == nil, "注册中心初始化失败")
		if tt.isExsit {
			err := rgt.CreatePersistentNode(tt.path, "{}")
			assert.Equal(t, true, err == nil, "创建初始化节点失败")
		}
		err = publish(rgt, tt.path, tt.v, tt.cover)
		assert.Equal(t, tt.wantErr, err != nil, tt.name+",err")

		data, _, err := rgt.GetValue(tt.path)
		assert.Equal(t, true, err == nil, "获取新的节点数据失败")
		if !tt.isExsit {
			assert.Equal(t, tt.v, string(data), tt.name+",1")
		} else {
			if tt.cover {
				assert.Equal(t, tt.v, string(data), tt.name+",2")
			} else {
				assert.Equal(t, "{}", string(data), tt.name+",3")
			}
		}
	}
}

func Test_deleteAll(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		subList []string
		wantErr bool
	}{
		{name: "1. 节点不存在", path: "/path1", subList: []string{}, wantErr: false},
		{name: "2. 单级节点存在,删除所有节点", path: "/path1", subList: []string{"/path1/cx1", "/path1/cx2"}, wantErr: false},
		{name: "3. 多级节点存在,删除所有节点", path: "/path1", subList: []string{"/path1/cx1", "/path1/cx2", "/path1/cx1/cc", "/path1/cx1/cc/xx"}, wantErr: false},
	}
	for _, tt := range tests {
		rgt, err := registry.GetRegistry("lm://.", global.Def.Log())
		assert.Equal(t, true, err == nil, "注册中心初始化失败")
		for _, str := range tt.subList {
			rgt.CreatePersistentNode(str, "{}")
		}
		err = deleteAll(rgt, tt.path)
		assert.Equal(t, tt.wantErr, err != nil, tt.name+",err")

		got, err := getAllPath(rgt, tt.path)
		assert.Equal(t, false, err != nil, tt.name+",err1")
		assert.Equal(t, []string{tt.path}, got, tt.name+",err")
	}
}

func Test_getAllPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		subList []string
		want    []string
		wantErr bool
	}{
		{name: "1. 无子级节点,获取所有的路径", path: "/path1", subList: []string{}, want: []string{"/path1"}, wantErr: false},
		{name: "2. 有单级子节点节点,获取所有的路径", path: "/path1", subList: []string{"/path1/cx1", "/path1/cx2"}, want: []string{"/path1/cx1", "/path1/cx2", "/path1"}, wantErr: false},
		{name: "3. 有多级子节点节点,获取所有的路径", path: "/path1", subList: []string{"/path1/cx1", "/path1/cx2", "/path1/cx1/cc", "/path1/cx1/cc/xx"}, want: []string{"/path1/cx1/cc/xx", "/path1/cx1/cc", "/path1/cx1", "/path1/cx2", "/path1"}, wantErr: false},
	}
	for _, tt := range tests {
		rgt, err := registry.GetRegistry("lm://.", global.Def.Log())
		assert.Equal(t, true, err == nil, "注册中心初始化失败")
		for _, str := range tt.subList {
			rgt.CreatePersistentNode(str, "{}")
		}
		got, err := getAllPath(rgt, tt.path)

		assert.Equal(t, tt.wantErr, err != nil, tt.name+",err")

		sort.Strings(tt.want)
		sort.Strings(got)
		assert.Equal(t, tt.want, got, tt.name+",value")
	}
}

type testss struct {
	XX string `json:"xx"`
}

func Test_getJSON(t *testing.T) {
	buff, _ := json.Marshal(map[string]string{"xx": "cc"})
	tests := []struct {
		name      string
		args      interface{}
		wantValue string
		wantErr   bool
	}{
		{name: "1. 参数是字符串", args: "string", wantValue: "string", wantErr: false},
		{name: "2. 参数是map", args: map[string]string{"xx": "cc"}, wantValue: string(buff), wantErr: false},
		{name: "3. 参数是struct", args: testss{XX: "cc"}, wantValue: string(buff), wantErr: false},
		{name: "3. 参数是prt", args: &testss{XX: "cc"}, wantValue: string(buff), wantErr: false},
		{name: "4. 参数是int", args: 1, wantValue: "1", wantErr: false},
		{name: "5. 参数是float", args: 1.5, wantValue: "1.5", wantErr: false},
		{name: "6. 参数是byte", args: []byte("d"), wantValue: `"ZA=="`, wantErr: false},
	}
	for _, tt := range tests {
		got, err := getJSON(tt.args)
		assert.Equal(t, tt.wantErr, err != nil, tt.name+",err")
		assert.Equal(t, tt.wantValue, got, tt.name+",value")
	}
}
