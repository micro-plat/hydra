package cron

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/micro-plat/hydra/conf/server/router"
	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/hydra/servers/pkg/adapter"
	"github.com/micro-plat/hydra/hydra/servers/pkg/middleware"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/utility"
)

const (
	unstarted = 1
	pause     = 2
	running   = 4
)

//Processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type Processor struct {
	//*dispatcher.Engine
	lock      sync.Mutex
	done      bool
	closeChan chan struct{}
	length    int
	index     int
	span      time.Duration
	slots     []cmap.ConcurrentMap //time slots
	startTime time.Time
	metric    *middleware.Metric
	status    int
	engine    *adapter.DispatcherEngine
}

//NewProcessor 创建processor
func NewProcessor(routers ...*router.Router) (p *Processor) {
	p = &Processor{
		status:    unstarted,
		closeChan: make(chan struct{}),
		span:      time.Second,
		length:    60,
		startTime: time.Now(),
		metric:    middleware.NewMetric(),
	}
	p.engine = adapter.NewDispatcherEngine(CRON)

	p.engine.Use(middleware.Recovery(true))
	p.engine.Use(p.metric.Handle())
	p.engine.Use(middleware.Logging())
	p.engine.Use(middleware.Recovery())

	p.engine.Use(middleware.Trace()) //跟踪信息
	p.engine.Use(middlewares...)

	p.addRouter(routers...)

	p.slots = make([]cmap.ConcurrentMap, p.length)
	for i := 0; i < p.length; i++ {
		p.slots[i] = cmap.New(2)
	}

	return p
}

func (s *Processor) addRouter(routers ...*router.Router) {
	s.engine.Handles(routers, middleware.ExecuteHandler())
}

//Start 所有任务
func (s *Processor) Start() error {
START:
	for {
		select {
		case <-s.closeChan:
			break START
		case <-time.After(s.span):
			s.execute()
		}
	}
	return nil
}

//Add 添加任务
func (s *Processor) Add(ts ...*task.Task) (err error) {
	for _, t := range ts {
		if t.Disable {
			s.Remove(t.GetUNQ())
			continue
		}
		task, err := NewCronTask(t)
		if err != nil {
			return fmt.Errorf("构建cron.task失败:%v", err)
		}

		if _, _, err := s.add(task); err != nil {
			return err
		}
	}
	return

}

func (s *Processor) add(task *CronTask) (offset int, round int, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.done {
		return -1, -1, nil
	}
	now := time.Now()
	nextTime := task.NextTime(now)
	if nextTime.Sub(now) < 0 {
		return -1, -1, errors.New("next time less than now.1")
	}
	offset, round = s.getOffset(now, nextTime)
	task.Round.Update(round)
	s.slots[offset].Set(utility.GetGUID(), task)
	return
}

//Remove 移除服务
func (s *Processor) Remove(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, slot := range s.slots {
		slot.RemoveIterCb(func(k string, value interface{}) bool {
			task := value.(*CronTask)
			task.Disable = true
			return task.GetName() == name
		})
	}
}

//Pause 暂停所有任务
func (s *Processor) Pause() (bool, error) {
	if s.status != pause {
		s.status = pause
		return true, nil
	}
	return false, nil
}

//Resume 恢复所有任务
func (s *Processor) Resume() (bool, error) {
	if s.status != running {
		s.status = running
		return true, nil
	}
	return false, nil
}

//Close 退出
func (s *Processor) Close() {
	defer s.metric.Stop()
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.done {
		s.done = true
		close(s.closeChan)
	}
}

//TaskCount 获取当前启用的Task数量
func (s *Processor) TaskCount() int {
	count := 0
	for i := range s.slots {
		for item := range s.slots[i].IterBuffered() {
			if !item.Val.(*CronTask).Disable {
				count++
			}
		}
	}
	return count
}

//-------------------------------------内部处理------------------------------------

func (s *Processor) getOffset(now time.Time, next time.Time) (pos int, circle int) {
	d := next.Sub(now) //剩余时间
	delaySeconds := int(math.Ceil(float64(d) / float64(1e9)))
	circle = int(delaySeconds) / s.length
	pos = int(s.index+delaySeconds) % s.length
	if pos == s.index { //offset与当前index相同时，应减少一环
		circle--
	}
	return
}

func (s *Processor) execute() {
	s.startTime = time.Now()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.index = (s.index + 1) % s.length
	current := s.slots[s.index]
	current.RemoveIterCb(func(k string, value interface{}) bool {
		task := value.(*CronTask)
		task.Round.Reduce()
		if task.Round.Get() < 0 {
			if task.Round.Get() == -1 { //所有环数已扣减完成
				go s.handle(task)
			}
			return true
		}
		return false
	})
}
func (s *Processor) handle(task *CronTask) error {
	if s.done || task.Disable {
		return nil
	}
	if s.status == running {
		task.Counter.Increase()
		s.engine.HandleRequest(task) //触发服务引擎进行业务处理
	}
	if task.IsImmediately() {
		return nil
	}
	_, _, err := s.add(task)
	return err

}
