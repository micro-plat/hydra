package main

import (
	"github.com/micro-plat/hydra/context"

	"github.com/micro-plat/hydra/component"
	"github.com/micro-plat/hydra/hydra"
)

func main() {
	app := hydra.NewApp(
		hydra.WithPlatName("install"),
		hydra.WithSystemName("demo"),
		hydra.WithServerTypes("api"),
		//	hydra.WithRegistry("zk://192.168.0.107"),
		//	hydra.WithClusterName("test"),
		hydra.WithDebug())
	app.Conf.API.SetMainConf(`{"address":"#address"}`)
	app.Conf.Plat.SetVarConf("db", "db", `{			
			"provider":"ora",
			"connString":"sso/123456@orcl136",
			"maxOpen":10,
			"maxIdle":1,
			"lifeTime":10		
	}`)

	app.Conf.API.Installer(func(c component.IContainer) error {
		db, err := c.GetDB()
		if err != nil {
			return err
		}
		_, _, _, err = db.Execute(`create table sso_role_menu123(
	id number(20) not null,
  sys_id number(20) not null,
  role_id number(20) not null,
  menu_id number(20) not null,
  enable number(1) default 0 not null,
  create_time date default sysdate not null,
  sortrank number(20) not null
  )`, map[string]interface{}{})
		return err
	})

	app.Micro("/hello", (component.ServiceFunc)(helloWorld))
	app.Start()
}

func helloWorld(ctx *context.Context) (r interface{}) {
	return "hello world"
}
