package cron

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNewServer(t *testing.T) {
	type args struct {
		tasks []*task.Task
	}
	tests := []struct {
		name    string
		args    args
		wantT   *Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := NewServer(tt.args.tasks...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("NewServer() = %v, want %v", gotT, tt.wantT)
			}
		})
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
		{name: "相同值修改", Processor: NewProcessor(), running: false, addr: "", want: false},
		{name: "不同值修改", Processor: NewProcessor(), running: true, addr: "", want: false},
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
		{name: "关闭时暂停", Processor: NewProcessor(), running: false, addr: "", want: true, wantErr: true},
		{name: "启动时暂停", Processor: NewProcessor(), running: true, addr: "", want: true, wantErr: true},
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
		{name: "关闭时恢复", Processor: NewProcessor(), running: false, addr: "", want: true, wantErr: true},
		{name: "启动时恢复", Processor: NewProcessor(), running: true, addr: "", want: true, wantErr: true},
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
