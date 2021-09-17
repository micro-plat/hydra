package dbr

import (
	"fmt"

	"github.com/micro-plat/lib4go/errs"
)

//Delete 删除节点
func (r *DBR) Delete(path string) error {
	count, err := r.db.Execute(r.sqltexture.delete, newInput(path))
	if err != nil || count < 1 {
		return errs.New("删除节点发生错误(%s)%w", path, errs.GetDBError(count, err))
	}
	r.notifyParentChange(path, 0)
	return nil
}

//clear 清除节点
func (r *DBR) clear(path string) {
	_, err := r.db.Execute(r.sqltexture.clear, newInput(path))
	fmt.Println(err)
	return
}
