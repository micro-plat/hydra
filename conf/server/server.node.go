package server

import (
	"strings"

	"github.com/micro-plat/hydra/conf"
)

var _ conf.ICluster = ClusterNodes{}

//ClusterNodes 集群节点信息
type ClusterNodes []*CNode

//Current 获取当前节点
func (c ClusterNodes) Current() conf.ICNode {
	for _, v := range c {
		if v.IsCurrent() {
			return v
		}
	}
	return nil
}

//Iter 迭代所有集群节点
func (c ClusterNodes) Iter(f func(conf.ICNode) bool) {
	for _, c := range c {
		if !f(c) {
			return
		}
	}
}

//Clone 克隆当前对象
func (c ClusterNodes) Clone() conf.ICluster {
	nodes := make([]*CNode, 0, len(c))
	for _, m := range c {
		nodes = append(nodes, m.Clone())
	}
	return ClusterNodes(nodes)
}

//CNode 集群节点
type CNode struct {
	name    string
	root    string
	path    string
	host    string
	index   int
	clustID string
	mid     string
}

//NewCNode 构建集群节点信息
func NewCNode(name string, mid string, index int) *CNode {
	items := strings.Split(name, "_")
	return &CNode{
		name:    name,
		mid:     mid,
		index:   index,
		host:    items[0],
		clustID: items[1],
	}
}

//GetHost 获取服务器信息
func (c *CNode) GetHost() string {
	return c.host
}

//GetClusterID 获取节点编号
func (c *CNode) GetClusterID() string {
	return c.clustID
}

//GetName 获取当前服务名称
func (c *CNode) GetName() string {
	return c.name
}

//IsCurrent 是否是当前服务
func (c *CNode) IsCurrent() bool {
	return c.clustID == c.mid
}

//GetIndex 获取当前节点的索引编号
func (c *CNode) GetIndex() int {
	return c.index
}

//Clone 克隆当前对象
func (c *CNode) Clone() *CNode {
	node := *c
	return &node
}
