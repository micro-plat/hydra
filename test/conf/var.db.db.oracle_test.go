package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/db"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/hydra/test/assert"
)

func TestDBOracleNew(t *testing.T) {

	tests := []struct {
		name       string
		connString string
		opts       []db.Option
		want       *db.DB
	}{
		{name: "1. Conf-DBOracleNew-测试新增-无OPTION", connString: "zhjy/123456@orcl136", want: &db.DB{Provider: "oracle", ConnString: "zhjy/123456@orcl136", MaxOpen: 10, MaxIdle: 3, LifeTime: 600}},
		{name: "2. Conf-DBOracleNew-测试新增-WithConnect", connString: "zhjy/123456@orcl136", opts: []db.Option{db.WithConnect(11, 22, 33)}, want: &db.DB{Provider: "oracle", ConnString: "zhjy/123456@orcl136", MaxOpen: 11, MaxIdle: 22, LifeTime: 33}},
	}
	for _, tt := range tests {
		got := oracle.New(tt.connString, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestDBOracleNewBy(t *testing.T) {

	tests := []struct {
		name    string
		uName   string
		pwd     string
		tnsName string
		opts    []db.Option
		want    *db.DB
	}{
		{name: "1. Conf-DBOracleNewBy-测试新增-无OPTION", uName: "zhjy", pwd: "123456", tnsName: "orcl136", want: &db.DB{Provider: "oracle", ConnString: "zhjy/123456@orcl136", MaxOpen: 10, MaxIdle: 3, LifeTime: 600}},
		{name: "2. Conf-DBOracleNewBy-测试新增-WithConnect", uName: "zhjy", pwd: "123456", tnsName: "orcl136", opts: []db.Option{db.WithConnect(11, 22, 33)}, want: &db.DB{Provider: "oracle", ConnString: "zhjy/123456@orcl136", MaxOpen: 11, MaxIdle: 22, LifeTime: 33}},
	}
	for _, tt := range tests {
		got := oracle.NewBy(tt.uName, tt.pwd, tt.tnsName, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
