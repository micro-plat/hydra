package db

import (
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

//executeBatch  批量执行SQL语句
func executeBatch(db IDBExecuter, sqls []string, input map[string]interface{}) (QueryRows, error) {
	output := types.NewXMaps()
	if len(sqls) == 0 {
		return output, fmt.Errorf("未传入任何SQL语句")
	}
	ninput := types.XMap(input)
	output.Append(ninput)
	for i, sql := range sqls {

		if len(sql) < 6 {
			return nil, fmt.Errorf("sql语句错误%s", sql)
		}
		prefix := strings.Trim(strings.TrimSpace(strings.TrimLeft(sql, "\n")), "\t")[:6]
		switch strings.ToUpper(prefix) {
		case "SELECT":
			coutput, err := db.Query(sql, ninput.ToMap())
			if err != nil {
				return output, err
			}
			if coutput.Len() == 0 {
				return output, fmt.Errorf("%s数据不存在%w input:%+v", sql, errs.ErrNotExist, types.Sprint(ninput.ToMap()))
			}
			ninput.Merge(coutput.Get(0))
			if i == len(sqls)-1 && coutput.Len() > 1 {
				return coutput, nil
			}
		case "UPDATE", "INSERT":
			rows, err := db.Execute(sql, ninput.ToMap())
			if err != nil {
				return output, err
			}
			if rows == 0 {
				return output, fmt.Errorf("%s数据修改失败%w input:%+v", sql, errs.ErrNotExist, types.Sprint(ninput.ToMap()))
			}
		default:
			return output, fmt.Errorf("不支持的SQL语句，或SQL语句前包含有特殊字符:%s", sql)
		}
	}
	return output, nil
}
