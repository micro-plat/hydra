package creator

import (
	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/conf/vars/db/mysql"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
)

//Vardb db配置
type Vardb struct {
	vars vars
}

//NewDB 构建db配置
func NewDB(internal map[string]map[string]interface{}) *Vardb {
	return &Vardb{
		vars: internal,
	}
}

//Oracle 添加oracle
func (c *Vardb) Oracle(nodeName string, uName string, pwd string, tnsName string, opts ...db.Option) vars {
	return c.Custom(nodeName, oracle.NewBy(uName, pwd, tnsName, opts...))
}

//OracleByConnStr 添加oracle
func (c *Vardb) OracleByConnStr(nodeName string, connStr string, opts ...db.Option) vars {
	return c.Custom(nodeName, oracle.New(connStr, opts...))
}

//MySQL 添加MySQL
func (c *Vardb) MySQL(nodeName string, uName string, pwd string, serverIP string, dbName string, opts ...db.Option) vars {
	return c.Custom(nodeName, mysql.NewBy(uName, pwd, serverIP, dbName, opts...))
}

//MySQLByConnStr 添加MySQLByConnStr
func (c *Vardb) MySQLByConnStr(nodeName string, connStr string, opts ...db.Option) vars {
	return c.Custom(nodeName, mysql.New(connStr, opts...))
}

//Custom 自定义数据库配置
func (c *Vardb) Custom(nodeName string, q interface{}) vars {
	if _, ok := c.vars[db.TypeNodeName]; !ok {
		c.vars[db.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[db.TypeNodeName][nodeName] = q
	return c.vars
}
