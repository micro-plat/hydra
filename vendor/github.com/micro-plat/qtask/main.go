package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro-plat/hydra/conf"
	_ "github.com/micro-plat/hydra/hydra"
	"github.com/micro-plat/hydra/registry"
	ldb "github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/qtask/qtask/db"
	"github.com/micro-plat/zkcli/rgsts"

	_ "github.com/zkfy/go-oci8"
)

func main() {

	//处理输入参数
	defer logger.Close()
	logger := logger.New("qtask")
	if len(os.Args) < 3 {
		logger.Error("命令错误，请使用 ’qtask [注册中心地址] [平台名称]‘ 注册中心连接串(proto://host)，平台名称(根据平台名称获取数据库配置串)")
		return
	}
	zkAddr := os.Args[1]
	platName := strings.Trim(os.Args[2], "/")
	dbName := "db"
	if len(os.Args) > 3 {
		dbName = os.Args[3]
	}

	//构建注册中心参数
	registry, err := registry.NewRegistryWithAddress(zkAddr, logger)
	if err != nil {
		rgsts.Log.Error(err)
		return
	}
	buff, _, err := registry.GetValue(fmt.Sprintf("/%s/var/db/%s", platName, dbName))
	if err != nil {
		logger.Error(err)
		return
	}

	//构建数据库对象
	xdb, err := createDB(buff)
	if err != nil {
		logger.Error(err)
		return
	}

	//创建数据库
	err = db.CreateDB(xdb)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("数据表创建成功")
}

func createDB(buff []byte) (ldb.IDB, error) {
	var dbConf conf.DBConf
	if err := json.Unmarshal(buff, &dbConf); err != nil {
		return nil, err
	}
	if b, err := govalidator.ValidateStruct(&dbConf); !b {
		return nil, err
	}
	return ldb.NewDB(dbConf.Provider,
		dbConf.ConnString,
		dbConf.MaxOpen,
		dbConf.MaxIdle,
		dbConf.LefeTime)
}
