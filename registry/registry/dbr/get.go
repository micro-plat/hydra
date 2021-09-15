package dbr

import (
	"fmt"
	"strings"
)

//GetValue 获取节点值
func (r *DBR) GetValue(path string) (data []byte, version int32, err error) {
	datas, err := r.db.Query(getValue, map[string]interface{}{
		"path": path,
	})
	if err != nil {
		return nil, 0, err
	}
	if datas.IsEmpty() {
		return nil, 0, fmt.Errorf("数据不存在")
	}
	return []byte(datas.Get(0).GetString("value")), datas.Get(0).GetInt32("data_version"), nil
}

//GetChildren 获取所有子节点
func (r *DBR) GetChildren(path string) (paths []string, version int32, err error) {
	datas, err := r.db.Query(getChildren, map[string]interface{}{
		"path": path,
	})
	if err != nil {
		return nil, 0, err
	}
	paths = make([]string, 0, len(datas))
	for _, p := range datas {
		np := strings.TrimPrefix(p.GetString("path"), path)
		paths = append(paths, strings.Trim(strings.Trim(np, path), "/"))
	}
	return paths, datas.Get(0).GetInt32("data_version"), nil
}
