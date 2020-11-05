package creator

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/db/mysql"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"

	"github.com/micro-plat/hydra/test/assert"

	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/creator/confvars"
)

func Test_vars_DB(t *testing.T) {
	tests := []struct {
		name string
		v    *confvars.Vardb
		want *confvars.Vardb
	}{
		{name: "初始化mysqldb对象", v: vars{}.DB().MySQL("name", mysql.New("contion")), want: confvars.NewDB(vars{}).MySQL("name", mysql.New("contion"))},
		{name: "初始化oracledb对象", v: vars{}.DB().Oracle("name", oracle.New("contion")), want: confvars.NewDB(vars{}).MySQL("name", oracle.New("contion"))},
		{name: "初始化cutomdb对象", v: vars{}.DB().MySQL("name", mysql.New("contion")), want: confvars.NewDB(vars{}).MySQL("name", mysql.New("contion"))},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.v, tt.name)
	}
}

func Test_vars_Cache(t *testing.T) {
	tests := []struct {
		name string
		v    vars
		want *confvars.Varcache
	}{
		{name: "初始化cache对象", v: vars{}, want: confvars.NewCache(vars{})},
	}
	for _, tt := range tests {
		got := tt.v.Cache()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_vars_Queue(t *testing.T) {
	tests := []struct {
		name string
		v    vars
		want *confvars.Varqueue
	}{
		{name: "初始化queue对象", v: vars{}, want: confvars.NewQueue(vars{})},
	}
	for _, tt := range tests {
		got := tt.v.Queue()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_vars_RLog(t *testing.T) {
	type args struct {
		service string
		opts    []rlog.Option
	}
	tests := []struct {
		name string
		v    vars
		args args
		want vars
	}{
		{name: "节点不存在,设置rlog", v: vars{}, args: args{service: "service1", opts: []rlog.Option{}}, want: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: rlog.New("service1")}}},
		{name: "节点存在,设置rlog", v: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: "xx"}}, args: args{service: "service1", opts: []rlog.Option{}}, want: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: rlog.New("service1")}}},
	}
	for _, tt := range tests {
		got := tt.v.RLog(tt.args.service, tt.args.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_vars_Redis(t *testing.T) {
	type args struct {
		name string
		opts *redis.Redis
	}
	tests := []struct {
		name string
		v    vars
		args args
		want vars
	}{
		{name: "节点不存在,设置redis", v: vars{}, args: args{name: "service1", opts: &redis.Redis{}}, want: vars{redis.TypeNodeName: map[string]interface{}{"service1": &redis.Redis{}}}},
		{name: "节点存在,设置redis", v: vars{redis.TypeNodeName: map[string]interface{}{"service1": "xx"}}, args: args{name: "service1", opts: redis.New([]string{"192.168.0.101"})}, want: vars{redis.TypeNodeName: map[string]interface{}{"service1": redis.New([]string{"192.168.0.101"})}}},
	}
	for _, tt := range tests {
		got := tt.v.Redis(tt.args.name, tt.args.opts)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
