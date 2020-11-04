package confvars

import (
	"github.com/micro-plat/hydra/conf/vars/db"
	dbmysql "github.com/micro-plat/hydra/conf/vars/db/mysql"
	dboracle "github.com/micro-plat/hydra/conf/vars/db/oracle"
)

type Vardb struct {
	vars vars
}

func NewDB(confVars map[string]map[string]interface{}) *Vardb{
	return &Vardb{
		vars:confVars,
	}
}

func (c *Vardb) Oracle(name string, q *dboracle.Oracle) *Vardb {
	return c.Custom(name, q)
}

func (c *Vardb) MySQL(name string, q *dbmysql.MySQL) *Vardb {
	return c.Custom(name, q)
}

func (c *Vardb) Custom(name string, q interface{}) *Vardb {
	if _, ok := c.vars[db.TypeNodeName]; !ok {
		c.vars[db.TypeNodeName] = make(map[string]interface{})
	}
	c.vars[db.TypeNodeName][name] = q
	return c
}
