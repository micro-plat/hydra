package conf

import "strings"

import "github.com/micro-plat/hydra/registry"

//CNode 集群节点
type CNode struct {
	name    string
	root    string
	path    string
	host    string
	clustID string
	cid     string
}

//NewCNode 构建集群节点信息
func NewCNode(root string, name string, cid string) *CNode {
	items := strings.Split(name, "_")
	return &CNode{
		name:    name,
		root:    root,
		cid:     cid,
		path:    registry.Join(root, name),
		host:    items[0],
		clustID: items[1],
	}
}

//GetFullPath 获取完整路径信息
func (c *CNode) GetFullPath() string {
	return c.path
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
	return c.clustID == c.cid
}
