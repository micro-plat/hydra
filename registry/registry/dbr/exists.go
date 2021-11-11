package dbr

import (
	"github.com/micro-plat/lib4go/types"
)

//Exists 检查节点是否存在
func (r *DBR) Exists(path string) (bool, error) {
	data, err := r.db.Scalar(r.sqltexture.exists, newInput(path))
	return types.GetInt(data) > 0, err
}
