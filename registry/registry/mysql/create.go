package mysql

import (
	"fmt"
	"sync/atomic"

	"github.com/micro-plat/hydra/registry/registry/mysql/internal/sql"
	"github.com/micro-plat/lib4go/errs"
)

//CreatePersistentNode 创建永久节点
func (r *Mysql) CreatePersistentNode(path string, data string) (err error) {
	value := newValue(data, false)
	count, err := r.db.Execute(sql.CreatePersistentNode, map[string]interface{}{
		"path":  path,
		"value": value.String(),
	})
	if err != nil || count < 1 {
		return errs.New("创建节点错误:%+v,count:%d", err, count)
	}
	r.notifyParentChange(path, value.Version)
	return nil
}

//CreateTempNode 创建临时节点
func (r *Mysql) CreateTempNode(path string, data string) (err error) {
	value := newValue(data, true)
	count, err := r.db.Execute(sql.CreatePersistentNode, map[string]interface{}{
		"path":  path,
		"value": value.String(),
	})
	if err != nil || count < 1 {
		return errs.New("创建节点错误:%+v,count:%d", err, count)
	}
	r.tmpNodes.Set(path, 0)
	r.notifyParentChange(path, value.Version)
	return nil
}

//CreateSeqNode 创建序列节点 @todo
func (r *Mysql) CreateSeqNode(path string, data string) (rpath string, err error) {
	nid := atomic.AddInt32(&r.seqValue, 1)
	rpath = fmt.Sprintf("%s_%d", path, nid)
	err = r.CreatePersistentNode(rpath, data)
	return rpath, err
}
