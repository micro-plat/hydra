package cron

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/lib4go/assert"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		tasks   []*task.Task
		wantT   *Server
		wantErr bool
	}{
		{name: "1. 初始化空对象", tasks: []*task.Task{}, wantT: &Server{Processor: NewProcessor()}, wantErr: true},
		{name: "2. 初始化单个错误任务对象", tasks: []*task.Task{task.NewTask("错误", "/cron/serve1")}, wantT: &Server{Processor: NewProcessor()}, wantErr: false},
		{name: "3. 初始化单个正确任务对象", tasks: []*task.Task{task.NewTask("@every 10s", "/cron/serve1")}, wantT: &Server{Processor: NewProcessor()}, wantErr: true},
		{name: "4. 初始化多个错误任务对象", tasks: []*task.Task{task.NewTask("错误", "/cron/serve1")}, wantT: &Server{Processor: NewProcessor()}, wantErr: false},
		{name: "5. 初始化多个正确任务对象", tasks: []*task.Task{task.NewTask("@every 10s", "/cron/serve1"), task.NewTask("@every 10s", "/cron/serve2")}, wantT: &Server{Processor: NewProcessor()}, wantErr: true},
	}
	for _, tt := range tests {
		if len(tt.tasks) > 0 {
			tt.wantT.Add(tt.tasks...)
		}

		gotT, err := NewServer(tt.tasks)
		assert.Equal(t, tt.wantErr, err == nil, tt.name, err)
		if tt.wantErr {
			assert.Equalf(t, tt.wantT.span, gotT.span, tt.name+",span")
			assert.Equalf(t, tt.wantT.length, gotT.length, tt.name+",length")
		}
	}
}

func TestServer_Shutdown(t *testing.T) {
	tests := []struct {
		name      string
		Processor *Processor
		running   bool
		addr      string
		want      bool
	}{
		{name: "1. 关闭服务器-running值西相同", Processor: NewProcessor(), running: false, addr: "", want: false},
		{name: "2. 关闭服务器-running值西不同", Processor: NewProcessor(), running: true, addr: "", want: false},
	}
	for _, tt := range tests {
		s := &Server{
			Processor: tt.Processor,
			running:   tt.running,
			addr:      tt.addr,
		}
		s.Shutdown()
		assert.Equal(t, tt.want, s.running, tt.name)
	}
}

func TestServer_Pause(t *testing.T) {
	tests := []struct {
		name      string
		Processor *Processor
		running   bool
		addr      string
		want      bool
		wantErr   bool
	}{
		{name: "1. 暂停服务器-关闭时暂停", Processor: NewProcessor(), running: false, addr: "", want: true, wantErr: true},
		{name: "2. 暂停服务器-启动时暂停", Processor: NewProcessor(), running: true, addr: "", want: true, wantErr: true},
	}
	for _, tt := range tests {
		s := &Server{
			Processor: tt.Processor,
			addr:      tt.addr,
		}

		if !tt.running {
			s.Shutdown()
		}
		got, err := s.Pause()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestServer_Resume(t *testing.T) {
	tests := []struct {
		name      string
		Processor *Processor
		running   bool
		addr      string
		want      bool
		wantErr   bool
	}{
		{name: "1. 恢复服务器-关闭时恢复", Processor: NewProcessor(), running: false, addr: "", want: true, wantErr: true},
		{name: "2. 恢复服务器-启动时恢复", Processor: NewProcessor(), running: true, addr: "", want: true, wantErr: true},
	}
	for _, tt := range tests {
		s := &Server{
			Processor: tt.Processor,
			running:   tt.running,
			addr:      tt.addr,
		}
		if !tt.running {
			s.Shutdown()
		}
		got, err := s.Resume()
		assert.Equal(t, tt.wantErr, err == nil, tt.name+",err")
		assert.Equal(t, tt.want, got, tt.name)
	}
}
