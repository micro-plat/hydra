package creator

import (
	"testing"

	"github.com/micro-plat/hydra/services"
	"github.com/micro-plat/hydra/test/assert"

	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
)

func Test_newMQC(t *testing.T) {
	tests := []struct {
		name string
		addr string
		opts []mqc.Option
		want *mqcBuilder
	}{
		{name: "初始化默认对象", addr: "redis://192.168.0.101", opts: []mqc.Option{}, want: &mqcBuilder{CustomerBuilder: map[string]interface{}{"main": mqc.New("redis://192.168.0.101")}}},
		{name: "初始化实体对象", addr: "redis://192.168.0.101", opts: []mqc.Option{mqc.WithDisable(), mqc.WithMasterSlave()}, want: &mqcBuilder{CustomerBuilder: map[string]interface{}{"main": mqc.New("redis://192.168.0.101", mqc.WithDisable(), mqc.WithMasterSlave())}}},
	}
	for _, tt := range tests {
		got := newMQC(tt.addr, tt.opts...)
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
		{name: "原本没有节点,添加空节点", fields: &mqcBuilder{CustomerBuilder: make(map[string]interface{})}, args: map[string]string{},
			want: &mqcBuilder{CustomerBuilder: map[string]interface{}{"queue": queue.NewEmptyQueues()}}},
		{name: "原本有节点,添加空节点", fields: &mqcBuilder{CustomerBuilder: map[string]interface{}{"queue": queue.NewQueues(queue.NewQueue("queue1", "service1"))}}, args: map[string]string{},
			want: &mqcBuilder{CustomerBuilder: map[string]interface{}{"queue": queue.NewQueues(queue.NewQueue("queue1", "service1"))}}},
		//用例三  由于services.MQC.Add该函数是延迟执行,所以无法测试效果
		{name: "原本没有节点,添加节点", fields: &mqcBuilder{CustomerBuilder: make(map[string]interface{})}, args: map[string]string{"queue2": "service2"},
			want: &mqcBuilder{CustomerBuilder: map[string]interface{}{"queue": queue.NewQueues(queue.NewQueue("queue2", "service2"))}}},
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
