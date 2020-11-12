package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestNewEmptyQueues(t *testing.T) {
	tests := []struct {
		name string
		want *queue.Queues
	}{
		{name: "初始化空对象", want: &queue.Queues{Queues: make([]*queue.Queue, 0)}},
	}
	for _, tt := range tests {
		got := queue.NewEmptyQueues()
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestNewQueues(t *testing.T) {
	tests := []struct {
		name string
		args []*queue.Queue
		want *queue.Queues
	}{
		{name: "初始化空对象", args: []*queue.Queue{}, want: &queue.Queues{Queues: make([]*queue.Queue, 0)}},
		{name: "初始化空对象,如参nil", args: nil, want: &queue.Queues{Queues: make([]*queue.Queue, 0)}},
		{name: "初始化单个队列对象", args: []*queue.Queue{queue.NewQueue("queue1", "service1")},
			want: &queue.Queues{Queues: []*queue.Queue{queue.NewQueue("queue1", "service1")}}},
		{name: "初始化多个队列对象", args: []*queue.Queue{queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2")},
			want: &queue.Queues{Queues: []*queue.Queue{queue.NewQueue("queue1", "service1"), queue.NewQueue("queue2", "service2")}}},
		{name: "初始化单个队列对象", args: []*queue.Queue{queue.NewQueueByConcurrency("queue1", "service1", 1)},
			want: &queue.Queues{Queues: []*queue.Queue{queue.NewQueueByConcurrency("queue1", "service1", 1)}}},
		{name: "初始化多个队列对象", args: []*queue.Queue{queue.NewQueueByConcurrency("queue1", "service1", 1), queue.NewQueueByConcurrency("queue2", "service2", 1)},
			want: &queue.Queues{Queues: []*queue.Queue{queue.NewQueueByConcurrency("queue1", "service1", 1), queue.NewQueueByConcurrency("queue2", "service2", 1)}}},
	}
	for _, tt := range tests {
		got := queue.NewQueues(tt.args...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestQueues_Append(t *testing.T) {
	tests := []struct {
		name   string
		fields *queue.Queues
		args   []*queue.Queue
		want   *queue.Queues
	}{
		{name: "添加空对象", fields: queue.NewEmptyQueues(), args: []*queue.Queue{}, want: queue.NewEmptyQueues()},
		{name: "添加nil对象", fields: queue.NewEmptyQueues(), args: nil, want: queue.NewEmptyQueues()},
		{name: "添加单个队列对象", fields: queue.NewEmptyQueues(), args: []*queue.Queue{queue.NewQueue("queue1", "service1")}, want: queue.NewQueues(queue.NewQueue("queue1", "service1"))},
		{name: "添加单个队列对象c", fields: queue.NewEmptyQueues(), args: []*queue.Queue{queue.NewQueueByConcurrency("queue1", "service1", 1)}, want: queue.NewQueues(queue.NewQueueByConcurrency("queue1", "service1", 1))},
		{name: "添加多个个队列对象", fields: queue.NewEmptyQueues(), args: []*queue.Queue{queue.NewQueueByConcurrency("queue1", "service1", 1), queue.NewQueue("queue2", "service2")}, want: queue.NewQueues(queue.NewQueueByConcurrency("queue1", "service1", 1), queue.NewQueue("queue2", "service2"))},
	}
	for _, tt := range tests {
		got, _ := tt.fields.Append(tt.args...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestQueuesGetConf(t *testing.T) {

	type test struct {
		name    string
		cnf     conf.IServerConf
		want    *queue.Queues
		wantErr bool
	}

	conf := mocks.NewConfBy("hydra", "queuetest")
	confB := conf.MQC("redis://192.168.0.101")
	test1 := test{name: "queues节点不存在", cnf: conf.GetMQCConf().GetServerConf(), want: queue.NewEmptyQueues(), wantErr: false}
	queueObj, err := queue.GetConf(test1.cnf)
	assert.Equal(t, test1.wantErr, (err != nil), test1.name)
	assert.Equal(t, test1.want, queueObj, test1.name)

	confB.Queue(queue.NewQueue("队列", "service1"))
	test2 := test{name: "queues节点存在,数据错误", cnf: conf.GetMQCConf().GetServerConf(), want: queue.NewEmptyQueues(), wantErr: false}
	queueObj, err = queue.GetConf(test2.cnf)
	assert.Equal(t, test2.wantErr, (err != nil), test2.name+",err")
	assert.Equal(t, test2.want, queueObj, test2.name+",obj")

}
