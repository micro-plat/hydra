package oracle

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/vars/db"
)

//Oracle oracle数据库连接信息
type Oracle = db.DB

//New 构建oracle连接信息
func New(connString string, opts ...db.Option) *Oracle {
	return db.New("oracle", connString, opts...)
}

//NewBy 构建oracle连接信息
func NewBy(uName string, pwd string, tnsName string, opts ...db.Option) *Oracle {
	return New(fmt.Sprintf("%s/%s@%s", uName, pwd, tnsName), opts...)
}
