/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 09:36:32
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-22 11:52:04
 */
package dbr

import (
	"fmt"
	"strings"
)

//GetValue 获取节点值
func (r *DBR) GetValue(path string) (data []byte, version int32, err error) {
	datas, err := r.db.Query(r.sqltexture.getValue, newInput(path))
	if err != nil {
		return nil, 0, err
	}
	if datas.IsEmpty() {
		return nil, 0, fmt.Errorf("数据不存在")
	}
	return []byte(datas.Get(0).GetString(FieldValue)), datas.Get(0).GetInt32(FieldDataVersion), nil
}

//GetChildren 获取所有子节点
func (r *DBR) GetChildren(path string) (paths []string, version int32, err error) {
	datas, err := r.db.Query(r.sqltexture.getChildren, newInput(path))
	if err != nil {
		return nil, 0, err
	}
	paths = make([]string, 0, len(datas))
	for _, p := range datas {
		if p.GetString(FieldPath) == path {
			continue
		}
		np := strings.TrimPrefix(p.GetString(FieldPath), path)
		paths = append(paths, strings.Trim(strings.Trim(np, path), "/"))
	}
	return paths, datas.Get(0).GetInt32(FieldDataVersion), nil
}
