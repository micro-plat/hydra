package creator

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/cache/cacheredis"
	gocache "github.com/micro-plat/hydra/conf/vars/cache/gocache"
	memcached "github.com/micro-plat/hydra/conf/vars/cache/memcached"
	"github.com/micro-plat/hydra/conf/vars/db/mysql"
	"github.com/micro-plat/hydra/conf/vars/db/oracle"
	"github.com/micro-plat/hydra/conf/vars/http"
	queuelmq "github.com/micro-plat/hydra/conf/vars/queue/lmq"
	queuemqtt "github.com/micro-plat/hydra/conf/vars/queue/mqtt"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/conf/vars/rpc"
	"github.com/micro-plat/hydra/creator/internal"
	"github.com/micro-plat/lib4go/assert"
)

func Test_vars_DB(t *testing.T) {
	tests := []struct {
		name string
		v    *internal.Vardb
		want *internal.Vardb
	}{
		{name: "1. 初始化mysqldb对象", v: vars{}.DB().MySQL("mysql", mysql.New("contion")), want: internal.NewDB(vars{}).MySQL("mysql", mysql.New("contion"))},
		{name: "2. 初始化oracledb对象", v: vars{}.DB().Oracle("oracle", oracle.New("contion")), want: internal.NewDB(vars{}).MySQL("oracle", oracle.New("contion"))},
		{name: "3. 初始化cutomdb对象", v: vars{}.DB().Custom("customer", mysql.New("contion")), want: internal.NewDB(vars{}).Custom("customer", mysql.New("contion"))},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.v, tt.name)
	}
}

func Test_vars_Cache(t *testing.T) {
	tests := []struct {
		name string
		v    *internal.Varcache
		want *internal.Varcache
	}{
		{name: "1. 初始化redis对象", v: vars{"redis": map[string]interface{}{"redis": redis.New([]string{"192.196.0.1"})}}.Cache().Redis("redisname", cacheredis.New(cacheredis.WithConfigName("redis"))),
			want: internal.NewCache(vars{"redis": map[string]interface{}{"redis": redis.New([]string{"192.196.0.1"})}}).Redis("redisname", cacheredis.New(cacheredis.WithConfigName("redis")))},
		{name: "2. 初始化gocache对象", v: vars{}.Cache().GoCache("GoCache", gocache.New()), want: internal.NewCache(vars{}).GoCache("GoCache", gocache.New())},
		{name: "3. 初始化Memcache对象", v: vars{}.Cache().Memcache("Memcache", memcached.New()), want: internal.NewCache(vars{}).Memcache("Memcache", memcached.New())},
		{name: "4. 初始化customer对象", v: vars{}.Cache().Custom("Custom", "sdsdsd"), want: internal.NewCache(vars{}).Custom("Custom", "sdsdsd")},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.v, tt.name)
	}
}

func Test_vars_Queue(t *testing.T) {
	tests := []struct {
		name string
		v    *internal.Varqueue
		want *internal.Varqueue
	}{
		{name: "1. 初始化redis对象", v: vars{"redis": map[string]interface{}{"redis": redis.New([]string{"192.196.0.1"})}}.Queue().Redis("redis", queueredis.New(queueredis.WithConfigName("redis"))),
			want: internal.NewQueue(vars{"redis": map[string]interface{}{"redis": redis.New([]string{"192.196.0.1"})}}).Redis("redis", queueredis.New(queueredis.WithConfigName("redis")))},
		{name: "2. 初始化MQTT对象", v: vars{}.Queue().MQTT("mqtt", queuemqtt.New("192.168.0.101:8017")), want: internal.NewQueue(vars{}).MQTT("mqtt", queuemqtt.New("192.168.0.101:8017"))},
		{name: "3. 初始化LMQ对象", v: vars{}.Queue().LMQ("lmq", queuelmq.New()), want: internal.NewQueue(vars{}).LMQ("lmq", queuelmq.New())},
		{name: "4. 初始化Custom对象", v: vars{}.Queue().Custom("custom", "dddddd"), want: internal.NewQueue(vars{}).Custom("custom", "dddddd")},
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
		opts *redis.Redis
	}
	tests := []struct {
		name string
		v    vars
		args args
		want vars
	}{
		{name: "1. 节点不存在,设置默认redis", v: vars{}, args: args{name: redis.TypeNodeName, opts: redis.New([]string{})}, want: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: redis.New([]string{})}}},
		{name: "2. 节点不存在,设置自定义redis", v: vars{}, args: args{name: redis.TypeNodeName, opts: redis.New([]string{"192.168.0.101:8888"}, redis.WithPoolSize(10), redis.WithDbIndex(1))}, want: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: redis.New([]string{"192.168.0.101:8888"}, redis.WithPoolSize(10), redis.WithDbIndex(1))}}},
		{name: "3. 节点存在,设置默认redis", v: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: "xx"}}, args: args{name: redis.TypeNodeName, opts: redis.New([]string{})}, want: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: redis.New([]string{})}}},
		{name: "4. 节点存在,设置自定义redis", v: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: "xx"}}, args: args{name: redis.TypeNodeName, opts: redis.New([]string{"192.168.0.101:8888"}, redis.WithPoolSize(10), redis.WithDbIndex(1))}, want: vars{redis.TypeNodeName: map[string]interface{}{redis.TypeNodeName: redis.New([]string{"192.168.0.101:8888"}, redis.WithPoolSize(10), redis.WithDbIndex(1))}}},
	}
	for _, tt := range tests {
		got := tt.v.Redis(tt.args.name, tt.args.opts)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
