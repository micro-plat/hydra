package dbr

import (
	"github.com/micro-plat/lib4go/errs"
)

//Delete 删除节点
func (r *DBR) Delete(path string) error {
	count, err := r.db.Execute(delete, newInput(path))
	if err != nil || count < 1 {
		return errs.New("删除节点错误:%+v,count:%d", err, count)
	}
	r.notifyParentChange(path, 0)
	return nil
}
