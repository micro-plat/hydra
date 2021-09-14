package mysql

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/registry/registry/mysql/internal/sql"
)

//GetValue 获取节点值
func (r *Mysql) GetValue(path string) (data []byte, version int32, err error) {
	datas, err := r.db.Query(sql.GetValue, map[string]interface{}{
		"path": path,
	})
	if err != nil {
		return nil, 0, err
	}
	if datas.IsEmpty() {
		return nil, 0, fmt.Errorf("数据不存在")
	}
	json := datas.Get(0)
	value, err := newValueByJSON(json.GetString("value"))
	if err != nil {
		return nil, 0, err
	}
	return []byte(value.Data), value.Version, nil
}

//GetChildren 获取所有子节点
func (r *Mysql) GetChildren(path string) (paths []string, version int32, err error) {
	datas, err := r.db.Query(sql.GetChildren, map[string]interface{}{
		"path": fmt.Sprintf("%s/", path),
	})
	if err != nil {
		return nil, 0, err
	}

	paths = make([]string, 0, len(datas))
	cache := map[string]bool{}

	for _, p := range datas {
		cpath := strings.TrimPrefix(p.GetString("path"), path)
		children := strings.Split(cpath, "/")[1] //取第一段
		if ok := cache[children]; ok {
			continue
		}
		cache[children] = true
		paths = append(paths, children)
	}

	return
}
