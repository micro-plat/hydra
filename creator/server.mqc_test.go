package creator

import (
	"testing"

	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/lib4go/assert"

	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
)

func Test_newMQC(t *testing.T) {
	tests := []struct {
		name   string
		addr   string
		opts   []mqc.Option
		repeat []mqc.Option
		want   *mqcBuilder
	}{
		{name: "1. 初始化默认mqc对象", addr: "redis://192.168.0.101", opts: []mqc.Option{}, want: &mqcBuilder{BaseBuilder: map[string]interface{}{"main": mqc.New("redis://192.168.0.101")}}},
		{name: "2. 初始化自定义mqc对象", addr: "redis://192.168.0.101", opts: []mqc.Option{mqc.WithDisable(), mqc.WithMasterSlave()}, want: &mqcBuilder{BaseBuilder: map[string]interface{}{"main": mqc.New("redis://192.168.0.101", mqc.WithDisable(), mqc.WithMasterSlave())}}},
		{name: "3. 重复初始化自定义mqc对象", addr: "redis://192.168.0.101",
			opts: []mqc.Option{mqc.WithDisable(), mqc.WithMasterSlave()}, repeat: []mqc.Option{mqc.WithEnable(), mqc.WithTrace(), mqc.WithSharding(1)},
			want: &mqcBuilder{BaseBuilder: map[string]interface{}{"main": mqc.New("redis://192.168.0.101", mqc.WithEnable(), mqc.WithTrace(), mqc.WithSharding(1))}}},
	}
	for _, tt := range tests {
		got := newMQC(tt.addr, tt.opts...)
		if tt.repeat != nil {
			got = newMQC(tt.addr, tt.repeat...)
		}
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func Test_mqcBuilder_Load(t *testing.T) {

	tests := []struct {
		name   string
		fields *mqcBuilder
		args   map[string]string
		want   *mqcBuilder
	}{
		{name: "1. 原本没有节点,添加空节点", fields: &mqcBuilder{BaseBuilder: make(map[string]interface{})}, args: map[string]string{},
			want: &mqcBuilder{BaseBuilder: map[string]interface{}{"queue": queue.NewEmptyQueues()}}},
		{name: "2. 原本有节点,添加空节点", fields: &mqcBuilder{BaseBuilder: map[string]interface{}{"queue": queue.NewQueues(queue.NewQueue("queue1", "service1"))}}, args: map[string]string{},
			want: &mqcBuilder{BaseBuilder: map[string]interface{}{"queue": queue.NewQueues(queue.NewQueue("queue1", "service1"))}}},
		//用例三  由于services.MQC.Add该函数是延迟执行,所以无法测试效果
		// {name: "原本没有节点,添加节点", fields: &mqcBuilder{BaseBuilder: make(map[string]interface{})}, args: map[string]string{"queue2": "service2"},
		// 	want: &mqcBuilder{BaseBuilder: map[string]interface{}{"queue": queue.NewQueues(queue.NewQueue("queue2", "service2"))}}},
	}
	for _, tt := range tests {
		for k, v := range tt.args {
			services.MQC.Add(k, v)
		}
		tt.fields.Load()
		assert.Equal(t, tt.want, tt.fields, tt.name)
	}
}

func Test_mqcBuilder_Queue(t *testing.T) {
	tests := []struct {
		name   string
		fields mqcBuilder
		args   []*queue.Queue
		want   *mqcBuilder
	}{
		// 由于该函数是延迟执行,所以无法测试效果
	}
	for _, tt := range tests {
		got := tt.fields.Queue(tt.args...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
