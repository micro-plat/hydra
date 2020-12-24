package creator

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/cache/cacheredis"
	"github.com/micro-plat/hydra/conf/vars/db/mysql"
	"github.com/micro-plat/hydra/conf/vars/http"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/conf/vars/rpc"
	"github.com/micro-plat/lib4go/assert"
)

func Test_vars_DB(t *testing.T) {
	tests := []struct {
		name string
		v    vars
		want vars
	}{
		{name: "1. 初始化mysqldb对象", v: vars{}.DB().MySQLByConnStr("mysql", "contion"), want: NewDB(vars{}).MySQLByConnStr("mysql", "contion")},
		{name: "2. 初始化oracledb对象", v: vars{}.DB().OracleByConnStr("oracle", "contion"), want: NewDB(vars{}).OracleByConnStr("oracle", "contion")},
		{name: "3. 初始化cutomdb对象", v: vars{}.DB().Custom("customer", mysql.New("contion")), want: NewDB(vars{}).Custom("customer", mysql.New("contion"))},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.v, tt.name)
	}
}

func Test_vars_Cache(t *testing.T) {
	tests := []struct {
		name string
		v    vars
		want vars
	}{
		{name: "1. 初始化redis对象", v: vars{"redis": map[string]interface{}{"redis": redis.New("192.196.0.1")}}.Cache().Redis("redisname", "", cacheredis.WithConfigName("redis")),
			want: NewCache(vars{"redis": map[string]interface{}{"redis": redis.New("192.196.0.1")}}).Redis("redisname", "", cacheredis.WithConfigName("redis"))},
		{name: "2. 初始化gocache对象", v: vars{}.Cache().GoCache("GoCache"), want: NewCache(vars{}).GoCache("GoCache")},
		{name: "3. 初始化Memcache对象", v: vars{}.Cache().Memcache("Memcache", ""), want: NewCache(vars{}).Memcache("Memcache", "")},
		{name: "4. 初始化customer对象", v: vars{}.Cache().Custom("Custom", "sdsdsd"), want: NewCache(vars{}).Custom("Custom", "sdsdsd")},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.v, tt.name)
	}
}

func Test_vars_Queue(t *testing.T) {
	tests := []struct {
		name string
		v    vars
		want vars
	}{
		{name: "1. 初始化redis对象", v: vars{"redis": map[string]interface{}{"redis": redis.New("192.196.0.1")}}.Queue().Redis("redis", "", queueredis.WithConfigName("redis")),
			want: NewQueue(vars{"redis": map[string]interface{}{"redis": redis.New("192.196.0.1")}}).Redis("redis", "", queueredis.WithConfigName("redis"))},
		{name: "2. 初始化MQTT对象", v: vars{}.Queue().MQTT("mqtt", "192.168.0.101:8017"), want: NewQueue(vars{}).MQTT("mqtt", "192.168.0.101:8017")},
		{name: "3. 初始化LMQ对象", v: vars{}.Queue().LMQ("lmq"), want: NewQueue(vars{}).LMQ("lmq")},
		{name: "4. 初始化Custom对象", v: vars{}.Queue().Custom("custom", "dddddd"), want: NewQueue(vars{}).Custom("custom", "dddddd")},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.v, tt.name)
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
		{name: "1. 节点不存在,设置默认rlog", v: vars{}, args: args{service: "service1", opts: []rlog.Option{}}, want: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: rlog.New("service1")}}},
		{name: "2. 节点不存在,设置自定义rlog", v: vars{}, args: args{service: "service1", opts: []rlog.Option{rlog.WithEnable(), rlog.WithAll()}}, want: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: rlog.New("service1", rlog.WithEnable(), rlog.WithAll())}}},
		{name: "3. 节点存在,设置默认rlog", v: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: "xx"}}, args: args{service: "service1", opts: []rlog.Option{}}, want: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: rlog.New("service1")}}},
		{name: "4. 节点存在,设置自定义rlog", v: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: "xx"}}, args: args{service: "service1", opts: []rlog.Option{rlog.WithEnable(), rlog.WithAll()}}, want: vars{rlog.TypeNodeName: map[string]interface{}{rlog.LogName: rlog.New("service1", rlog.WithEnable(), rlog.WithAll())}}},
	}
	for _, tt := range tests {
		got := tt.v.RLog(tt.args.service, tt.args.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_vars_RPC(t *testing.T) {
	type args struct {
		service string
		opts    []rpc.Option
	}
	tests := []struct {
		name string
		v    vars
		args args
		want vars
	}{
		{name: "1. 节点不存在,设置默认rpc", v: vars{}, args: args{service: "service1", opts: []rpc.Option{}}, want: vars{rpc.RPCTypeNode: map[string]interface{}{"service1": rpc.New()}}},
		{name: "2. 节点不存在,设置自定义rpc", v: vars{}, args: args{service: "service1", opts: []rpc.Option{rpc.WithLocalFirst(), rpc.WithConnectionTimeout(10)}}, want: vars{rpc.RPCTypeNode: map[string]interface{}{"service1": rpc.New(rpc.WithLocalFirst(), rpc.WithConnectionTimeout(10))}}},
		{name: "3. 节点存在,设置默认rpc", v: vars{rpc.RPCTypeNode: map[string]interface{}{"service1": "xx"}}, args: args{service: "service1", opts: []rpc.Option{}}, want: vars{rpc.RPCTypeNode: map[string]interface{}{"service1": rpc.New()}}},
		{name: "4. 节点存在,设置自定义rpc", v: vars{rpc.RPCTypeNode: map[string]interface{}{"service1": "xx"}}, args: args{service: "service1", opts: []rpc.Option{rpc.WithLocalFirst(), rpc.WithConnectionTimeout(10)}}, want: vars{rpc.RPCTypeNode: map[string]interface{}{"service1": rpc.New(rpc.WithLocalFirst(), rpc.WithConnectionTimeout(10))}}},
	}
	for _, tt := range tests {
		got := tt.v.RPC(tt.args.service, tt.args.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_vars_HTTP(t *testing.T) {
	type args struct {
		service string
		opts    []http.Option
	}
	tests := []struct {
		name string
		v    vars
		args args
		want vars
	}{
		{name: "1. 节点不存在,设置默认Http", v: vars{}, args: args{service: http.HttpNameNode, opts: []http.Option{}}, want: vars{http.HttpTypeNode: map[string]interface{}{http.HttpNameNode: http.New()}}},
		{name: "2. 节点不存在,设置自定义http", v: vars{}, args: args{service: http.HttpNameNode, opts: []http.Option{http.WithConnTimeout(10), http.WithRaw([]byte(`{"t1":"t2"}`))}}, want: vars{http.HttpTypeNode: map[string]interface{}{http.HttpNameNode: http.New(http.WithConnTimeout(10), http.WithRaw([]byte(`{"t1":"t2"}`)))}}},
		{name: "3. 节点存在,设置默认http", v: vars{http.HttpTypeNode: map[string]interface{}{http.HttpNameNode: "xx"}}, args: args{service: http.HttpNameNode, opts: []http.Option{}}, want: vars{http.HttpTypeNode: map[string]interface{}{http.HttpNameNode: http.New()}}},
		{name: "4. 节点存在,设置自定义http", v: vars{http.HttpTypeNode: map[string]interface{}{http.HttpNameNode: "xx"}}, args: args{service: http.HttpNameNode, opts: []http.Option{http.WithConnTimeout(20), http.WithProxy("192.168.141.2")}}, want: vars{http.HttpTypeNode: map[string]interface{}{http.HttpNameNode: http.New(http.WithConnTimeout(20), http.WithProxy("192.168.141.2"))}}},
	}
	for _, tt := range tests {
		got := tt.v.HTTP(tt.args.service, tt.args.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_vars_Redis(t *testing.T) {
	type args struct {
		name string
		addr string
		opts []redis.Option
	}
	tests := []struct {
		name string
		v    vars
		args args
		want vars
	}{
		{name: "1. 节点不存在,设置默认redis", v: vars{}, args: args{name: redis.TypeNodeName, addr: "", opts: nil}, want: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: redis.New("")}}},
		{name: "2. 节点不存在,设置自定义redis", v: vars{}, args: args{name: redis.TypeNodeName, addr: "192.168.0.101:8888", opts: []redis.Option{redis.WithPoolSize(10), redis.WithDbIndex(1)}}, want: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: redis.New("192.168.0.101:8888", redis.WithPoolSize(10), redis.WithDbIndex(1))}}},
		{name: "3. 节点存在,设置默认redis", v: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: "xx"}}, args: args{name: redis.TypeNodeName, addr: "", opts: nil}, want: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: redis.New("")}}},
		{name: "4. 节点存在,设置自定义redis", v: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: "xx"}}, args: args{name: redis.TypeNodeName, addr: "192.168.0.101:8888", opts: []redis.Option{redis.WithPoolSize(10), redis.WithDbIndex(1)}}, want: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: redis.New("192.168.0.101:8888", redis.WithPoolSize(10), redis.WithDbIndex(1))}}},
	}
	for _, tt := range tests {
		got := tt.v.Redis(tt.args.name, tt.args.addr, tt.args.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
