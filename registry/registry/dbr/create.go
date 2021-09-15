package dbr

import (
	"fmt"
	"sync/atomic"

	"github.com/micro-plat/lib4go/errs"
)

//CreatePersistentNode 创建永久节点
func (r *DBR) CreatePersistentNode(path string, data string) (err error) {
	count, err := r.db.Execute(createNode, map[string]interface{}{
		"path":         path,
		"value":        data,
		"data_version": 1,
		"acl_version":  1,
		"temp":         1,
	})
	if err != nil || count < 1 {
		return errs.New("创建节点错误:%+v,count:%d", err, count)
	}
	r.notifyParentChange(path, 1)
	return nil
}

//CreateTempNode 创建临时节点
func (r *DBR) CreateTempNode(path string, data string) (err error) {
	count, err := r.db.Execute(createNode, map[string]interface{}{
		"path":         path,
		"value":        data,
		"data_version": 1,
		"acl_version":  1,
		"temp":         0,
	})
	if err != nil || count < 1 {
		return errs.New("创建节点错误:%+v,count:%d", err, count)
	}
	r.tmpNodes.Set(path, 0)
	r.notifyParentChange(path, 1)
	return nil
}

//CreateSeqNode 创建序列节点 @todo
func (r *DBR) CreateSeqNode(path string, data string) (rpath string, err error) {
	nid := atomic.AddInt32(&r.seqValue, 1)
	rpath = fmt.Sprintf("%s_%d", path, nid)
	err = r.CreateTempNode(rpath, data)
	return rpath, err
}
