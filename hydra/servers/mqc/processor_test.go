package mqc

import (
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

func TestNewProcessor(t *testing.T) {
	tests := []struct {
		name    string
		proto   string
		confRaw string
		wantP   *Processor
		wantErr string
	}{
		{name: "协议错误", proto: "proto", confRaw: "{}", wantErr: "构建mqc服务失败(proto:proto,raw:{}) mqc: 未知的协议类型 proto"},
		{name: "协议配置正确", proto: "redis", confRaw: `{"proto":"redis","addrs":["192.168.5.79:6379"]}`, wantP: &Processor{status: unstarted,
			closeChan: make(chan struct{}),
			startTime: time.Now(),
			queues:    cmap.New(4)}},
	}
	for _, tt := range tests {
		gotP, err := NewProcessor(tt.proto, tt.confRaw)
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, 4, len(gotP.Engine.Handlers), tt.name)
	}

}

func TestProcessor_Add(t *testing.T) {
	tests := []struct {
		name       string
		queueNames []string
		queues     []*queue.Queue
		wantErr    string
	}{
		{name: "添加消息队列", queues: []*queue.Queue{queue.NewQueue("queue1", "services1"), queue.NewQueue("queue2", "services2")}},
		{name: "再次添加消息队列", queues: []*queue.Queue{queue.NewQueue("queue1", "services1"), queue.NewQueue("queue3", "services3")}},
	}
	s, _ := NewProcessor("redis", `{"proto":"redis","addrs":["192.168.5.79:6379"]}`)
	for _, tt := range tests {
		err := s.Add(tt.queues...)
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, nil, err, tt.name)
		for _, v := range tt.queues {
			assert.Equal(t, s.queues.Items()[v.Queue], v, tt.name)
		}
	}
}

func TestProcessor_Remove(t *testing.T) {
	s, _ := NewProcessor("redis", `{"proto":"redis","addrs":["192.168.5.79:6379"]}`)
	//添加消息队列
	queues := []*queue.Queue{queue.NewQueue("queue1", "services1"), queue.NewQueue("queue2", "services2")}
	err := s.Add(queues...)
	assert.Equal(t, nil, err, "Add")
	//移除消息队列
	l := len(queues)
	for _, v := range queues {
		err := s.Remove(v)
		assert.Equal(t, nil, err, "Remove")
		l--
		assert.Equal(t, len(s.queues.Items()), l, "Remove")
	}

}

func TestProcessor_Resume(t *testing.T) {

	s, _ := NewProcessor("redis", `{"proto":"redis","addrs":["192.168.5.79:6379"]}`)
	//添加消息队列
	queues := []*queue.Queue{queue.NewQueue("queue1", "services1"), queue.NewQueue("queue2", "services2")}
	err := s.Add(queues...)
	assert.Equal(t, nil, err, "Add")

	//暂停
	got, err := s.Pause()
	assert.Equal(t, nil, err, "Pause")
	assert.Equal(t, true, got, "Pause")

	//暂停
	got, err = s.Pause()
	assert.Equal(t, nil, err, "Pause2")
	assert.Equal(t, false, got, "Pause2")

	//重启
	got, err = s.Resume()
	assert.Equal(t, nil, err, "Resume")
	assert.Equal(t, true, got, "Resume")
	time.Sleep(time.Second)
	got, err = s.Resume()
	assert.Equal(t, nil, err, "Resume2")
	assert.Equal(t, false, got, "Resume2")

	assert.Equal(t, len(queues), len(s.queues.Items()), "Resume2")

	//再次添加
	addQueues := []*queue.Queue{queue.NewQueue("queue3", "services3"), queue.NewQueue("queue4", "services4")}
	err = s.Add(addQueues...)
	assert.Equal(t, nil, err, "Add")
	assert.Equal(t, len(queues)+len(addQueues), len(s.queues.Items()), "Resume2")

	//暂停
	got, err = s.Pause()
	assert.Equal(t, nil, err, "Pause3")
	assert.Equal(t, true, got, "Pause3")

}

func TestProcessor_Close(t *testing.T) {
	s, _ := NewProcessor("redis", `{"proto":"redis","addrs":["192.168.5.79:6379"]}`)
	//添加消息队列
	queues := []*queue.Queue{queue.NewQueue("queue1", "services1"), queue.NewQueue("queue2", "services2")}
	err := s.Add(queues...)
	assert.Equal(t, nil, err, "Add")
	s.Close()
	assert.Equal(t, 0, len(s.queues.Items()), "Close")
	assert.Equal(t, true, s.done, "Close")
}
