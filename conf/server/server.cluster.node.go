package server

import (
	"strings"

	"github.com/micro-plat/hydra/conf"
)

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

//IsBefore 判断当前节点是否是前i个
func (c *CNode) IsBefore(i int) bool {
	return c.index < i
}

//Clone 克隆当前对象
func (c *CNode) Clone() conf.ICNode {
	node := *c
	return &node
}
