package confvars

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/db"
	dbmysql "github.com/micro-plat/hydra/conf/vars/db/mysql"
	dboracle "github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNewDB(t *testing.T) {
	tests := []struct {
		name string
		args map[string]map[string]interface{}
		want *Vardb
	}{
		{name: "初始化db对象", args: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}},
			want: &Vardb{vars: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}}}},
	}
	for _, tt := range tests {
		got := NewDB(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVardb_Oracle(t *testing.T) {
	type args struct {
		name string
		q    *dboracle.Oracle
	}
	tests := []struct {
		name   string
		fields *Vardb
		args   args
		want   *Vardb
	}{
		{name: "初始化对象", fields: NewDB(map[string]map[string]interface{}{}), args: args{name: "oracleDB", q: dboracle.New("connstr")},
			want: NewDB(map[string]map[string]interface{}{db.TypeNodeName: map[string]interface{}{"oracleDB": dboracle.New("connstr")}})},
	}
	for _, tt := range tests {
		got := tt.fields.Oracle(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVardb_MySQL(t *testing.T) {
	type args struct {
		name string
		q    *dbmysql.MySQL
	}
	tests := []struct {
		name   string
		fields *Vardb
		args   args
		want   *Vardb
	}{
		{name: "初始化对象", fields: NewDB(map[string]map[string]interface{}{}), args: args{name: "mysqlDB", q: dbmysql.New("connstr1")},
			want: NewDB(map[string]map[string]interface{}{db.TypeNodeName: map[string]interface{}{"mysqlDB": dbmysql.New("connstr1")}})},
	}
	for _, tt := range tests {
		got := tt.fields.MySQL(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVardb_Custom(t *testing.T) {
	type args struct {
		name string
		q    interface{}
	}
	tests := []struct {
		name   string
		fields *Vardb
		args   args
		want   *Vardb
	}{
		{name: "初始化对象", fields: NewDB(map[string]map[string]interface{}{}), args: args{name: "customer", q: map[string]interface{}{"sss": "sdfdsfsdf"}},
			want: NewDB(map[string]map[string]interface{}{db.TypeNodeName: map[string]interface{}{"customer": map[string]interface{}{"sss": "sdfdsfsdf"}}})},
	}
	for _, tt := range tests {
		got := tt.fields.Custom(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
