package internal

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/queue"
	queuelmq "github.com/micro-plat/hydra/conf/vars/queue/lmq"
	queuemqtt "github.com/micro-plat/hydra/conf/vars/queue/mqtt"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	"github.com/micro-plat/hydra/conf/vars/redis"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/types"
)

func TestNewQueue(t *testing.T) {
	tests := []struct {
		name string
		args map[string]map[string]interface{}
		want *Varqueue
	}{
		{name: "1. 初始化queue对象", args: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}},
			want: &Varqueue{vars: map[string]map[string]interface{}{"main": map[string]interface{}{"test1": "123456"}}}},
	}
	for _, tt := range tests {
		got := NewQueue(tt.args)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarqueue_Redis(t *testing.T) {
	newRedis := queueredis.New(queueredis.WithConfigName("address"))
	newRedis.Redis = redis.New([]string{"192.196.0.1"})
	type args struct {
		name string
		q    *queueredis.Redis
	}
	tests := []struct {
		name    string
		fields  *Varqueue
		args    args
		want    *Varqueue
		wantErr string
	}{
		{name: "1. configname是空", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "redis", q: queueredis.New(queueredis.WithAddrs("address"))},
			want: NewQueue(map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"redis": queueredis.New(queueredis.WithAddrs("address"))}})},
		{name: "2. configname不为空,无redis节点", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "redis", q: queueredis.New(queueredis.WithConfigName("address"))},
			want: nil, wantErr: "请确认已配置/var/redis"},
		{name: "3. configname不为空,redis节点存在,无configname节点", fields: NewQueue(map[string]map[string]interface{}{"redis": map[string]interface{}{}}), args: args{name: "redis", q: queueredis.New(queueredis.WithConfigName("address"))},
			want: nil, wantErr: "请确认已配置/var/redis/address"},
		{name: "4. configname不为空,节点存在,configname节点", fields: NewQueue(map[string]map[string]interface{}{"redis": map[string]interface{}{"address": redis.New([]string{"192.196.0.1"})}}),
			args: args{name: "redis", q: queueredis.New(queueredis.WithConfigName("address"))},
			want: NewQueue(map[string]map[string]interface{}{"redis": map[string]interface{}{"address": redis.New([]string{"192.196.0.1"})},
				queue.TypeNodeName: map[string]interface{}{"redis": newRedis}}), wantErr: ""},
	}
	for _, tt := range tests {
		func() {
			defer func() {
				e := recover()
				if e != nil {
					assert.Equal(t, tt.wantErr, types.GetString(e), tt.name+",err")
				}
			}()

			got := tt.fields.Redis(tt.args.name, tt.args.q)
			assert.Equal(t, tt.want, got, tt.name)
		}()

	}
}

func TestVarqueue_MQTT(t *testing.T) {
	type args struct {
		name string
		q    *queuemqtt.MQTT
	}
	tests := []struct {
		name   string
		fields *Varqueue
		args   args
		want   *Varqueue
	}{
		{name: "1. 初始化MQTT对象", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "mqttQueue", q: queuemqtt.New("address1")},
			want: NewQueue(map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"mqttQueue": queuemqtt.New("address1")}})},
	}
	for _, tt := range tests {
		got := tt.fields.MQTT(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarqueue_LMQ(t *testing.T) {
	type args struct {
		name string
		q    *queuelmq.LMQ
	}
	tests := []struct {
		name   string
		fields *Varqueue
		args   args
		want   *Varqueue
	}{
		{name: "1. 初始化LMQ对象", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "lmqQueue", q: queuelmq.New()},
			want: NewQueue(map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"lmqQueue": queuelmq.New()}})},
	}
	for _, tt := range tests {
		got := tt.fields.LMQ(tt.args.name, tt.args.q)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarqueue_Custom(t *testing.T) {
	type args struct {
		name string
		q    interface{}
	}
	tests := []struct {
		name   string
		fields *Varqueue
		args   args
		repeat *args
		want   *Varqueue
	}{
		{name: "1. 初始化自定义空对象", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "", q: map[string]interface{}{}},
			want: NewQueue(map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"": map[string]interface{}{}}})},
		{name: "2. 初始化自定义对象", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "cusomeQueue", q: map[string]interface{}{"test": "123456"}},
			want: NewQueue(map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"cusomeQueue": map[string]interface{}{"test": "123456"}}})},
		{name: "3. 重复初始化自定义对象", fields: NewQueue(map[string]map[string]interface{}{}),
			args:   args{name: "cusomeQueue", q: map[string]interface{}{"test": "123456"}},
			repeat: &args{name: "cusomeQueue", q: map[string]interface{}{"dddd": "989898"}},
			want:   NewQueue(map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"cusomeQueue": map[string]interface{}{"dddd": "989898"}}})},
	}
	for _, tt := range tests {
		got := tt.fields.Custom(tt.args.name, tt.args.q)
		if tt.repeat != nil {
			got = tt.fields.Custom(tt.repeat.name, tt.repeat.q)
		}
		assert.Equal(t, tt.want, got, tt.name)
	}
}
