package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/db/tpl"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
)

// executeBatch  批量执行SQL语句
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

// queryBatch  批量执行SQL语句
func queryBatch(db IDBExecuter, sqls []string, input map[string]interface{}) (QueryRows, error) {
	outputs := types.NewXMaps()
	if len(sqls) == 0 {
		return outputs, fmt.Errorf("未传入任何SQL语句")
	}

	for _, sql := range sqls {
		if len(sql) < 6 {
			return nil, fmt.Errorf("sql语句错误%s", sql)
		}
		prefix := strings.Trim(strings.TrimSpace(strings.TrimLeft(sql, "\n")), "\t")[:6]
		switch strings.ToUpper(prefix) {
		case "SELECT":
			coutput, err := db.Query(sql, input)
			if err != nil {
				return outputs, err
			}
			for _, v := range coutput {
				outputs.Append(v)
			}
		default:
			return outputs, fmt.Errorf("不支持的SQL语句，或SQL语句前包含有特殊字符:%s", sql)
		}
	}
	return outputs, nil
}

func insertBatch(db IBaseDB, tpl tpl.ITPLContext, sql string, inputs []map[string]interface{}) (row int64, err error) {
	// 不区分大小写查找VALUES位置
	lowerSql := strings.ToLower(sql)
	valuesIndex := strings.Index(lowerSql, "values")
	if valuesIndex == -1 {
		return 0, fmt.Errorf("sql格式错误，不包含values: %s", sql)
	}

	// 使用原始SQL分割，保留大小写
	insertPart, valuesPart := sql[:valuesIndex], sql[valuesIndex+6:] // "values"长度为6

	// 构建完整的SQL语句和参数列表
	var valueParts []string
	var argParts []interface{}

	for _, input := range inputs {
		query, args := tpl.GetSQLContext(valuesPart, input)
		valueParts = append(valueParts, query)
		argParts = append(argParts, args...)
	}

	// 执行批量插入
	fullSql := fmt.Sprintf("%s VALUES %s", insertPart, strings.Join(valueParts, ","))
	result, err := db.Execute(fullSql, argParts...)
	if err != nil {
		return 0, getDBError(err, fullSql, argParts)
	}
	return result, nil
}
func updateSave(db IBaseDB, tpl tpl.ITPLContext, sql string, inputs []map[string]interface{}) (row int64, err error) {
	totalRows, result := int64(0), int64(0)
	for _, input := range inputs {
		fullSql, args := tpl.GetSQLContext(sql, input)
		result, err = db.Execute(fullSql, args...)
		if err != nil {
			return 0, getDBError(err, fullSql, args)
		}
		totalRows += result
	}

	return totalRows, nil
}
func fetchRows(db IBaseDB, tpl tpl.ITPLContext, sql string, input map[string]interface{}) (*sql.Rows, error) {
	query, args := tpl.GetSQLContext(sql, input)
	data, err := db.FetchRows(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	return data, err
}
func query(db IBaseDB, tpl tpl.ITPLContext, sql string, input map[string]interface{}) (data QueryRows, err error) {
	query, args := tpl.GetSQLContext(sql, input)
	data, err = db.Query(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	return
}
func scalar(db IBaseDB, tpl tpl.ITPLContext, sql string, input map[string]interface{}) (data interface{}, err error) {
	query, args := tpl.GetSQLContext(sql, input)
	result, err := db.Query(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	if result.Len() == 0 || result.Get(0).IsEmpty() {
		return nil, nil
	}
	data, _ = result.Get(0).Get(result.Get(0).Keys()[0])
	return
}
func executes(db IBaseDB, tpl tpl.ITPLContext, sql string, input map[string]interface{}) (insertID int64, row int64, err error) {
	query, args := tpl.GetSQLContext(sql, input)
	insertID, row, err = db.Executes(query, args...)
	if err != nil {
		return 0, 0, getDBError(err, query, args)
	}
	return
}
func execute(db IBaseDB, tpl tpl.ITPLContext, sql string, input map[string]interface{}) (row int64, err error) {
	query, args := tpl.GetSQLContext(sql, input)
	row, err = db.Execute(query, args...)
	if err != nil {
		return 0, getDBError(err, query, args)
	}
	return
}
func executeSP(db IBaseDB, tpl tpl.ITPLContext, procName string, input map[string]interface{}, output ...interface{}) (row int64, err error) {
	query, args := tpl.GetSPContext(procName, input)
	ni := append(args, output...)
	row, err = db.Execute(query, ni...)
	if err != nil {
		return 0, getDBError(err, query, args)
	}
	return
}

func resolveFullRows(rows *sql.Rows) (dataRows QueryRows, err error) {
	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := types.XMaps{}

	// 逐行解析
	for rows.Next() {
		currentRow, err := resolveRow(rows, columns)
		if err != nil {
			return nil, err
		}
		// 将当前行添加到结果切片
		result = append(result, currentRow)
	}

	// 检查错误
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
func resolveRow(rows *sql.Rows, columns []string) (map[string]interface{}, error) {
	// 使用 sql.NullString 来处理可能的 NULL 值
	values := make([]sql.NullString, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// 扫描数据
	if err := rows.Scan(scanArgs...); err != nil {
		return nil, err
	}

	// 构建当前行的map
	currentRow := make(map[string]interface{})
	for i, col := range columns {
		if values[i].Valid {
			currentRow[col] = values[i].String
		} else {
			currentRow[col] = "" // 将 NULL 值转换为空字符串
		}
	}

	return currentRow, nil
}
func getDBError(err error, query string, args []interface{}) error {
	return fmt.Errorf("%w(sql:%s,args:%+v)", err, query, args)
}
