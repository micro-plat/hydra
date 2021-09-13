package mysql

import (
	"fmt"

	"github.com/micro-plat/hydra/registry/registry/mysql/internal/sql"
	"github.com/micro-plat/lib4go/errs"
)

//Update 更新节点值
func (r *Mysql) Update(path string, data string) (err error) {

	//获取原数据
	datas, err := r.db.Query(sql.GetValue, map[string]interface{}{
		"path": path,
	})
	if err != nil {
		return err
	}
	if datas.IsEmpty() {
		return fmt.Errorf("数据不存在")
	}

	//解析并判断节点类型
	ovalue, err := newValueByJSON(datas.Get(0).GetString("value"))
	if err != nil {
		return err
	}

	//构建新对象，并修改
	value := newValue(data, ovalue.IsTemp)

	count, err := r.db.Execute(sql.Update, map[string]interface{}{
		"path":  path,
		"value": value.String(),
	})

	if err != nil || count < 1 {
		return errs.New("更新节点错误:%+v,count:%d", err, count)
	}

	//通知变更
	r.notifyValueChange(path, value)
	return nil
}
