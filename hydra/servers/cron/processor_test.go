package cron

import (
	"testing"
	"time"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

func TestNewProcessor(t *testing.T) {
	tests := []struct {
		name    string
		handles []middleware.Handler
		wantP   *Processor
	}{
		{name: "1. cron-NewProcessor-设置空任务服务", handles: []middleware.Handler{}, wantP: &Processor{status: unstarted, span: time.Second, length: 60, slots: cMap(60)}},
		{name: "2. cron-NewProcessor-设置有任务服务", handles: []middleware.Handler{func(middleware.IMiddleContext) {}},
			wantP: &Processor{status: unstarted, span: time.Second, length: 60, slots: cMap(60)}},
	}
	for _, tt := range tests {
		middlewares = tt.handles
		gotP := NewProcessor()
		assert.Equalf(t, 4+len(tt.handles), len(gotP.Engine.RouterGroup.Handlers), tt.name+",中间件数量")
		assert.Equalf(t, tt.wantP.slots, gotP.slots, tt.name+",slots")
		assert.Equalf(t, tt.wantP.span, gotP.span, tt.name+",span")
		assert.Equalf(t, tt.wantP.length, gotP.length, tt.name+",length")
		middlewares = []middleware.Handler{}
	}
}

func cMap(lenth int) []cmap.ConcurrentMap {
	slots := make([]cmap.ConcurrentMap, lenth, lenth)
	for i := 0; i < lenth; i++ {
		slots[i] = cmap.New(2)
	}
	return slots
}

func TestProcessor_Add(t *testing.T) {
	tests := []struct {
		name    string
		ts      []*task.Task
		count   int
		wantErr bool
	}{
		{name: "1. cron-ProcessorAdd-添加空列表", ts: []*task.Task{}, count: 0, wantErr: true},
		{name: "2. cron-ProcessorAdd-添加disable列表", ts: []*task.Task{task.NewTask("@every 10s", "/cron/serve1", task.WithDisable())}, count: 0, wantErr: true},
		{name: "3. cron-ProcessorAdd-添加错误的列表", ts: []*task.Task{task.NewTask("错误", "/cron/serve1")}, count: 0, wantErr: false},
		{name: "4. cron-ProcessorAdd-添加重复的列表", ts: []*task.Task{task.NewTask("@every 10s", "/cron/serve1"), task.NewTask("@every 10h", "/cron/serve1")}, count: 1, wantErr: true},
		{name: "5. cron-ProcessorAdd-添加重复+disable的列表", ts: []*task.Task{task.NewTask("@every 10s", "/cron/serve1"), task.NewTask("@every 10h", "/cron/serve1"), task.NewTask("@every 10s", "/cron/serve1", task.WithDisable())}, count: 1, wantErr: true},
	}
	for _, tt := range tests {
		s := NewProcessor()
		err := s.Add(tt.ts...)
		assert.Equalf(t, tt.wantErr, err == nil, tt.name, err)
		assert.Equalf(t, 4+tt.count, len(s.Engine.RouterGroup.Handlers)+len(s.Engine.Routes()), tt.name+",服务数量")
	}
}

func TestProcessor_Remove(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		ts    []*task.Task
		count int
		args  args
	}{
		{name: "1. cron-ProcessorRemove-删除不存在的服务", ts: []*task.Task{task.NewTask("@every 10s", "/cron/serve1")}, count: 1, args: args{name: "/cron/serve2"}},
		{name: "2. cron-ProcessorRemove-删除存在的服务,没有禁用", ts: []*task.Task{task.NewTask("@every 10s", "/cron/serve1")}, count: 1, args: args{name: "/cron/serve1"}},
		{name: "3. cron-ProcessorRemove-删除存在的服务,禁用", ts: []*task.Task{task.NewTask("@every 10s", "/cron/serve1", task.WithDisable())}, count: 0, args: args{name: "/cron/serve1"}},
	}
	for _, tt := range tests {
		s := NewProcessor()
		err := s.Add(tt.ts...)
		assert.Equalf(t, nil, err, tt.name, err)
		s.Remove(tt.args.name)
		count := 0
		for _, slot := range s.slots {
			count += slot.Count()
		}
		assert.Equalf(t, tt.count, count, tt.name+",剩下的服务数量")
	}
}

func TestProcessor_Pause(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		want    bool
		wantErr bool
	}{
		{name: "1. cron-ProcessorPause-初始化为启用", status: 4, want: true, wantErr: true},
		{name: "2. cron-ProcessorPause-初始化为带启用", status: 1, want: true, wantErr: true},
		{name: "3. cron-ProcessorPause-初始化为禁用", status: 2, want: false, wantErr: true},
	}
	for _, tt := range tests {
		s := NewProcessor()
		s.status = tt.status
		got, err := s.Pause()
		assert.Equalf(t, tt.wantErr, err == nil, tt.name, err)
		assert.Equalf(t, tt.want, got, tt.name, got)
	}
}

func TestProcessor_Resume(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		want    bool
		wantErr bool
	}{
		{name: "1. cron-ProcessorResume-初始化为启用", status: 4, want: false, wantErr: true},
		{name: "2. cron-ProcessorResume-初始化为带启用", status: 1, want: true, wantErr: true},
		{name: "3. cron-ProcessorResume-初始化为禁用", status: 2, want: true, wantErr: true},
	}
	for _, tt := range tests {
		s := NewProcessor()
		s.status = tt.status
		got, err := s.Resume()
		assert.Equalf(t, tt.wantErr, err == nil, tt.name, err)
		assert.Equalf(t, tt.want, got, tt.name, got)
	}
}

func TestProcessor_getOffset(t *testing.T) {

	tests := []struct {
		name       string
		args       *task.Task
		wantPos    int
		wantCircle int
	}{
		{name: "1. 根据任务cron-获取任务分片和位置-1s", args: task.NewTask("@every 1s", "/server1"), wantPos: 1, wantCircle: 0},
		{name: "2. 根据任务cron-获取任务分片和位置-59s", args: task.NewTask("@every 59s", "/server1"), wantPos: 59, wantCircle: 0},
		{name: "3. 根据任务cron-获取任务分片和位置-60s", args: task.NewTask("@every 60s", "/server1"), wantPos: 0, wantCircle: 1},
		{name: "4. 根据任务cron-获取任务分片和位置-120s", args: task.NewTask("@every 120s", "/server1"), wantPos: 0, wantCircle: 2},
	}
	for _, tt := range tests {
		s := NewProcessor()
		taska, _ := NewCronTask(tt.args)
		now := time.Now()
		next := taska.NextTime(now)
		gotPos, gotCircle := s.getOffset(now, next)
		assert.Equalf(t, tt.wantPos, gotPos, tt.name)
		assert.Equalf(t, tt.wantCircle, gotCircle, tt.name)
		// fmt.Println("--", gotPos, "--", gotCircle)
	}
}
