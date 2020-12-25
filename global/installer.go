package global

type installer struct {
	//DB 数据库安装配置
	DB *db
}

//Installer 安装程序
var Installer = &installer{
	DB: &db{sqls: make([]string, 0, 1), handlers: make([]func() error, 0, 1)},
}
