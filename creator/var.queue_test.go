package creator

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/queue"
	queuelmq "github.com/micro-plat/hydra/conf/vars/queue/lmq"
	queuemqtt "github.com/micro-plat/hydra/conf/vars/queue/mqtt"
	"github.com/micro-plat/hydra/conf/vars/queue/queueredis"
	"github.com/micro-plat/lib4go/assert"
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
	type args struct {
		name string
		addr string
		opts []queueredis.Option
	}
	tests := []struct {
		name    string
		fields  *Varqueue
		args    args
		want    vars
		wantErr string
	}{
		{name: "1. configname是空", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "redis", addr: "address"},
			want: map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"redis": queueredis.New("address")}}},
	}
	for _, tt := range tests {
		func() {
			// defer func() {
			// 	e := recover()
			// 	if e != nil {
			// 		assert.Equal(t, tt.wantErr, types.GetString(e), tt.name+",err")
			// 	}
			// }()

			got := tt.fields.Redis(tt.args.name, tt.args.addr, tt.args.opts...)
			assert.Equal(t, tt.want, got, tt.name)
		}()

	}
}

func TestVarqueue_MQTT(t *testing.T) {
	type args struct {
		name string
		addr string
		opts queuemqtt.Option
	}
	tests := []struct {
		name   string
		fields *Varqueue
		args   args
		want   vars
	}{
		{name: "1. 初始化MQTT对象", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "mqttQueue", addr: "address1"},
			want: map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"mqttQueue": queuemqtt.New("address1")}}},
	}
	for _, tt := range tests {
		got := tt.fields.MQTT(tt.args.name, tt.args.addr)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestVarqueue_LMQ(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields *Varqueue
		args   args
		want   vars
	}{
		{name: "1. 初始化LMQ对象", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "lmqQueue"},
			want: map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"lmqQueue": queuelmq.New()}}},
	}
	for _, tt := range tests {
		got := tt.fields.LMQ(tt.args.name)
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
		want   vars
	}{
		{name: "1. 初始化自定义空对象", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "", q: map[string]interface{}{}},
			want: map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"": map[string]interface{}{}}}},
		{name: "2. 初始化自定义对象", fields: NewQueue(map[string]map[string]interface{}{}), args: args{name: "cusomeQueue", q: map[string]interface{}{"test": "123456"}},
			want: map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"cusomeQueue": map[string]interface{}{"test": "123456"}}}},
		{name: "3. 重复初始化自定义对象", fields: NewQueue(map[string]map[string]interface{}{}),
			args:   args{name: "cusomeQueue", q: map[string]interface{}{"test": "123456"}},
			repeat: &args{name: "cusomeQueue", q: map[string]interface{}{"dddd": "989898"}},
			want:   map[string]map[string]interface{}{queue.TypeNodeName: map[string]interface{}{"cusomeQueue": map[string]interface{}{"dddd": "989898"}}}},
	}
	for _, tt := range tests {
		got := tt.fields.Custom(tt.args.name, tt.args.q)
		if tt.repeat != nil {
			got = tt.fields.Custom(tt.repeat.name, tt.repeat.q)
		}
		assert.Equal(t, tt.want, got, tt.name)
	}
}
