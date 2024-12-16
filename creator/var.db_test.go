package creator

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/db"
	dbmysql "github.com/micro-plat/hydra/conf/vars/db/mysql"
	dboracle "github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/lib4go/assert"
)

func TestNewDB(t *testing.T) {
	tests := []struct {
		name string
		args map[string]map[string]interface{}
		want *Vardb
	}{
		{name: "1. 初始化db对象", args: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}},
			want: &Vardb{vars: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}}}},
	}
	for _, tt := range tests {
		got := NewDB(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVardb_Oracle(t *testing.T) {
	type args struct {
		name    string
		connStr string
	}
	tests := []struct {
		name   string
		fields *Vardb
		args   args
		want   vars
	}{
		{name: "1. 初始化Oracle对象", fields: NewDB(map[string]map[string]interface{}{}), args: args{name: "oracleDB", connStr: "connstr"},
			want: map[string]map[string]interface{}{db.TypeNodeName: map[string]interface{}{"oracleDB": dboracle.New("connstr")}}},
	}
	for _, tt := range tests {
		got := tt.fields.OracleByConnStr(tt.args.name, tt.args.connStr)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVardb_MySQL(t *testing.T) {
	type args struct {
		name    string
		connStr string
	}
	tests := []struct {
		name   string
		fields *Vardb
		args   args
		want   vars
	}{
		{name: "1. 初始化MySQL对象", fields: NewDB(map[string]map[string]interface{}{}), args: args{name: "mysqlDB", connStr: "connstr1"},
			want: map[string]map[string]interface{}{db.TypeNodeName: map[string]interface{}{"mysqlDB": dbmysql.New("connstr1")}}},
	}
	for _, tt := range tests {
		got := tt.fields.MySQLByConnStr(tt.args.name, tt.args.connStr)
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
		repeat *args
		want   vars
	}{
		{name: "1. 初始化空子定义对象", fields: NewDB(map[string]map[string]interface{}{}), args: args{name: "", q: map[string]interface{}{}},
			want: map[string]map[string]interface{}{db.TypeNodeName: map[string]interface{}{"": map[string]interface{}{}}}},
		{name: "2. 初始化自定义对象", fields: NewDB(map[string]map[string]interface{}{}), args: args{name: "customer", q: map[string]interface{}{"sss": "sdfdsfsdf"}},
			want: map[string]map[string]interface{}{db.TypeNodeName: map[string]interface{}{"customer": map[string]interface{}{"sss": "sdfdsfsdf"}}}},
		{name: "3. 重复初始化自定义对象", fields: NewDB(map[string]map[string]interface{}{}),
			args:   args{name: "customer", q: map[string]interface{}{"sss": "sdfdsfsdf"}},
			repeat: &args{name: "customer", q: map[string]interface{}{"xxx": "54dfdff"}},
			want:   map[string]map[string]interface{}{db.TypeNodeName: map[string]interface{}{"customer": map[string]interface{}{"xxx": "54dfdff"}}}},
	}
	for _, tt := range tests {
		got := tt.fields.Custom(tt.args.name, tt.args.q)
		if tt.repeat != nil {
			got = tt.fields.Custom(tt.repeat.name, tt.repeat.q)
		}
		assert.Equal(t, tt.want, got, tt.name)
	}
}
