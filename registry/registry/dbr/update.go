/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 09:36:32
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-23 16:14:33
 */
package dbr

import (
	"fmt"

	"github.com/micro-plat/lib4go/errs"
)

//Update 更新节点值
func (r *DBR) Update(path string, data string) (err error) {

	//获取原数据
	datas, err := r.db.Query(r.sqltexture.getValue, newInput(path))
	if err != nil {
		return err
	}
	if datas.IsEmpty() {
		return fmt.Errorf("数据不存在")
	}

	count, err := r.db.Execute(r.sqltexture.update, newInputByUpdate(path, data, datas.Get(0).GetInt32(FieldDataVersion)))
	if err != nil || count < 1 {
		return errs.New("更新节点错误:%+v,count:%d", err, count)
	}

	//通知变更
	r.notifyValueChange(path, data, datas.Get(0).GetInt32(FieldDataVersion))
	return nil
}
