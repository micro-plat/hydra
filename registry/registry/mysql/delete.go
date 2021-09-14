package mysql

import (
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/sql"
	"github.com/micro-plat/lib4go/errs"
)

//Delete 删除节点
func (r *Mysql) Delete(path string) error {
	count, err := r.db.Execute(sql.Delete, map[string]interface{}{
		"path": path,
	})

	if err != nil || count < 1 {
		return errs.New("删除节点错误:%+v,count:%d", err, count)
	}

	//@todo notify
	r.notifyParentChange(path, 0)
	return nil
}
