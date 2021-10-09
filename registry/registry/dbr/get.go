/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 09:36:32
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-26 14:57:48
 */
package dbr

import (
	"strings"
)

//GetValue 获取节点值
func (r *DBR) GetValue(path string) (data []byte, version int32, err error) {
	datas, err := r.db.Query(r.sqltexture.getValue, newInput(path))
	if err != nil {
		return nil, 0, err
	}
	if datas.IsEmpty() {
		return []byte(""), 0, nil
	}
	return []byte(datas.Get(0).GetString(FieldValue)), datas.Get(0).GetInt32(FieldDataVersion), nil
}

//GetChildren 获取所有子节点
func (r *DBR) GetChildren(path string) (paths []string, version int32, err error) {
	datas, err := r.db.Query(r.sqltexture.getChildren, newInput(path))
	if err != nil {
		return nil, 0, err
	}
	if len(datas) == 0 {
		return []string{}, 0, nil
	}
	paths = make([]string, 0, len(datas)-1)
	cache := map[string]bool{}
	for _, p := range datas {
		if p.GetString(FieldPath) == path {
			version = p.GetInt32(FieldDataVersion)
			continue
		}
		np := strings.TrimPrefix(p.GetString(FieldPath), path)
		np = strings.TrimLeft(np, "/")
		arryPath := strings.Split(np, "/")
		if len(arryPath) > 0 {
			if len(arryPath[0]) > 0 {
				p := arryPath[0]
				if ok, _ := cache[p]; ok {
					continue
				}
				cache[p] = true
				paths = append(paths, p)
			}
		}
	}
	return paths, version, nil
}
