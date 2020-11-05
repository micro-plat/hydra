package conf

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	_ "github.com/micro-plat/hydra/registry/registry/localmemory"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/types"
)

//@todo 注册中心的cluster的需要验证通知功能

func xTest_NewCluster(t *testing.T) {
	//初始化注册中心
	rgt, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")

	clusterName := "cluster1"
	pub := server.NewServerPub("platName", "sysName", "serverType", "cluster1")
	rootPath := pub.GetServerPubPath(clusterName)
	gotS, err := server.NewCluster(pub, rgt, clusterName)
	assert.Equal(t, true, err == nil, "初始化集群对象失败")
	assert.Equal(t, 0, gotS.Len(), "集群数量不正确0")

	path := rootPath + "/132456_111"
	err = rgt.CreateTempNode(path, "1")
	time.Sleep(1 * time.Second)
	assert.Equal(t, true, err == nil, "创建临时节点异常")
	assert.Equal(t, &server.CNode{}, gotS.Current(), "不能存在当前配置节点")
	assert.Equal(t, 1, gotS.Len(), "集群数量不正确1")

	path1 := rootPath + "/132456_" + pub.GetServerID()
	err = rgt.CreateTempNode(path1, "2")
	time.Sleep(1 * time.Second)
	assert.Equal(t, true, err == nil, "创建临时节点异常")
	assert.Equal(t, pub.GetServerID(), gotS.Current().GetNodeID(), "不能存在当前配置节点1")
	assert.Equal(t, 2, gotS.Len(), "集群数量不正确2")

	err = rgt.Delete(path1)
	assert.Equal(t, true, err == nil, "删除节点异常")
	time.Sleep(1 * time.Second)
	assert.Equal(t, &server.CNode{}, gotS.Current(), "不能存在当前配置节点2")
	assert.Equal(t, 1, gotS.Len(), "集群数量不正确3")

	var addCount int64 = 0
	var reduceCount int64 = 0
	reduceCh := make(chan string, 100)

	for x := 0; x < 10; x++ {
		for i := 0; i < 10; i++ {
			go func() {
				pathX := rootPath + "/55555_" + types.GetString(time.Now().Nanosecond())
				err = rgt.CreateTempNode(pathX, pathX)
				assert.Equal(t, true, err == nil, "for创建临时节点异常")
				nid := atomic.AddInt64(&addCount, 1)
				if nid%3 == 0 {
					reduceCh <- pathX
				}
			}()
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("len(reduceCh):", len(reduceCh))
	lenCount := len(reduceCh)
	for x := 0; x < lenCount; x++ {
		if len(reduceCh) == 0 {
			break
		}
		v := <-reduceCh
		go func() {
			err = rgt.Delete(v)
			assert.Equal(t, true, err == nil, "删除节点异常")
			atomic.AddInt64(&reduceCount, 1)
		}()
	}

	time.Sleep(10 * time.Second)
	assert.Equal(t, addCount-reduceCount+1, int64(gotS.Len()), "集群数量不正确n")
}

func xTestCluster_Current(t *testing.T) {
	rgt, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")

	clusterName := "cluster2"
	pub := server.NewServerPub("platName", "sysName", "serverType", clusterName)
	rootPath := pub.GetServerPubPath(clusterName)
	gotS, err := server.NewCluster(pub, rgt, clusterName)
	assert.Equal(t, true, err == nil, "初始化集群对象失败")
	assert.Equal(t, &server.CNode{}, gotS.Current(), "不能存在当前配置节点")

	path1 := rootPath + "/132456_" + pub.GetServerID()
	err = rgt.CreateTempNode(path1, "2")
	time.Sleep(1 * time.Second)
	assert.Equal(t, true, err == nil, "创建临时节点异常")
	assert.Equal(t, pub.GetServerID(), gotS.Current().GetNodeID(), "不能存在当前配置节点1")

	err = rgt.Delete(path1)
	assert.Equal(t, true, err == nil, "删除节点异常")
	time.Sleep(1 * time.Second)
	assert.Equal(t, &server.CNode{}, gotS.Current(), "不能存在当前配置节点2")
}

func TestCluster_GetType(t *testing.T) {
	clusterName := "cluster3"
	rgt, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")

	pub := server.NewServerPub("platName", "sysName", "serverType", clusterName)
	obj1, err := server.NewCluster(pub, rgt, clusterName)
	assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")

	pub1 := server.NewServerPub("platName", "sysName", "xxxx", clusterName)
	obj2, err := server.NewCluster(pub1, rgt, clusterName)
	assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")

	tests := []struct {
		name   string
		fields *server.Cluster
		want   string
	}{
		{name: "实体对象", fields: obj1, want: "serverType"},
		{name: "实体对象1", fields: obj2, want: "xxxx"},
	}
	for _, tt := range tests {
		got := tt.fields.GetServerType()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestCluster_Next(t *testing.T) {

	clusterName := "cluster4"
	rgt, err := registry.NewRegistry("lm://.", global.Def.Log())
	assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")

	pub := server.NewServerPub("platName", "sysName", "serverType", clusterName)
	obj1, err := server.NewCluster(pub, rgt, clusterName)
	assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")
	got, got1 := obj1.Next()
	assert.Equal(t, nil, got, "空对象返回node不为空")
	assert.Equal(t, false, got1, "空对象获取结果失败")

	clusterName1 := "cluster5"
	rootPath := pub.GetServerPubPath(clusterName1)
	nodeMap := map[string]string{
		"132456_111": "1",
		"132456_222": "2",
		"132456_333": "3",
		"132456_444": "4",
	}

	for path, data := range nodeMap {
		err = rgt.CreateTempNode(rootPath+"/"+path, data)
		assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")
	}

	time.Sleep(1 * time.Second)
	obj2, err := server.NewCluster(pub, rgt, clusterName1)
	assert.Equal(t, true, err == nil, "初始化集群对象,获取注册中心对象失败")
	resMap := []string{}
	for i := 0; i < 4; i++ {
		got, got1 = obj2.Next()
		if got1 {
			resMap = append(resMap, got.GetName())
		} else {
			break
		}
	}

	assert.Equal(t, len(resMap), len(nodeMap), "数据异常")
	for _, str := range resMap {
		_, ok := nodeMap[str]
		if ok {
			delete(nodeMap, str)
		}
	}
	assert.Equal(t, 0, len(nodeMap), "数据异常1")
}
