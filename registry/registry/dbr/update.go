package dbr

import (
	"fmt"

	"github.com/micro-plat/lib4go/errs"
)

//Update 更新节点值
func (r *DBR) Update(path string, data string) (err error) {

	//获取原数据
	datas, err := r.db.Query(getValue, newInput(path))
	if err != nil {
		return err
	}
	if datas.IsEmpty() {
		return fmt.Errorf("数据不存在")
	}

	count, err := r.db.Execute(update, newInputByUpdate(path, data, datas.Get(0).GetInt32("data_version")))

	if err != nil || count < 1 {
		return errs.New("更新节点错误:%+v,count:%d", err, count)
	}

	//通知变更
	r.notifyValueChange(path, data, datas.Get(0).GetInt32("data_version"))
	return nil
}
