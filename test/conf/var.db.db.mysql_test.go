package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/conf/vars/db/mysql"
	"github.com/micro-plat/hydra/test/assert"
)

var connectStr = "root:xxxxxx@tcp(192.168.0.1:3306)/dbname?charset=utf8"

func TestDBMysqlNew(t *testing.T) {
	tests := []struct {
		name       string
		connString string
		opts       []db.Option
		want       *db.DB
	}{
		{name: "1. Conf-DBMysqlNew-测试新增-无OPTION", connString: connectStr, want: &db.DB{Provider: "mysql", ConnString: connectStr, MaxOpen: 10, MaxIdle: 3, LifeTime: 600}},
		{name: "2. Conf-DBMysqlNew-测试新增-WithConnect", connString: connectStr, opts: []db.Option{db.WithConnect(11, 22, 33)}, want: &db.DB{Provider: "mysql", ConnString: connectStr, MaxOpen: 11, MaxIdle: 22, LifeTime: 33}},
	}
	for _, tt := range tests {
		got := mysql.New(tt.connString, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestDBMysqlNewBy(t *testing.T) {

	tests := []struct {
		name     string
		uName    string
		pwd      string
		serverIP string
		dbName   string
		opts     []db.Option
		want     *db.DB
	}{
		{name: "1. Conf-DBMysqlNewBy-测试新增-无OPTION", uName: "root", pwd: "xxxxxx", serverIP: "192.168.0.1:3306", dbName: "dbname", want: &db.DB{Provider: "mysql", ConnString: connectStr, MaxOpen: 10, MaxIdle: 3, LifeTime: 600}},
		{name: "2. Conf-DBMysqlNewBy-测试新增-WithConnect", uName: "root", pwd: "xxxxxx", serverIP: "192.168.0.1:3306", dbName: "dbname", opts: []db.Option{db.WithConnect(11, 22, 33)}, want: &db.DB{Provider: "mysql", ConnString: connectStr, MaxOpen: 11, MaxIdle: 22, LifeTime: 33}},
	}
	for _, tt := range tests {
		got := mysql.NewBy(tt.uName, tt.pwd, tt.serverIP, tt.dbName, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
