package mysql

import (
	"github.com/micro-plat/hydra/registry/registry/mysql/internal/sql"
	"github.com/micro-plat/lib4go/types"
)

//Exists 检查节点是否存在
func (r *Mysql) Exists(path string) (bool, error) {
	data, err := r.db.Scalar(sql.Exists, map[string]interface{}{
		"path": path,
	})
	return types.GetInt(data) > 0, err
}
