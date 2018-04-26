package cron

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/servers/pkg/dispatcher"
	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/redis"
	"github.com/micro-plat/lib4go/utility"
)

//Processor 任务处理程序
type Processor struct {
	*dispatcher.Dispatcher
	lock         sync.Mutex
	once         sync.Once
	done         bool
	closeChan    chan struct{}
	length       int
	index        int
	span         time.Duration
	slots        []cmap.ConcurrentMap //time slots
	startTime    time.Time
	isPause      bool
	redisSetting string
	redisClient  *redis.Client
	historyNode  string
}

//NewProcessor 创建processor
func NewProcessor(redisSetting string, historyNode string) (p *Processor, err error) {
	p = &Processor{
		Dispatcher:   dispatcher.New(),
		closeChan:    make(chan struct{}),
		span:         time.Second,
		length:       60,
		startTime:    time.Now(),
		redisSetting: redisSetting,
		historyNode:  historyNode,
	}
	p.slots = make([]cmap.ConcurrentMap, p.length, p.length)
	for i := 0; i < p.length; i++ {
		p.slots[i] = cmap.New(2)
	}
	if redisSetting != "" {
		p.redisClient, err = redis.NewClientByJSON(string(redisSetting))
	}
	return p, nil
}

//Start 启动cron timer
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

func (s *Processor) execute() {
	s.startTime = time.Now()
	s.lock.Lock()
	defer s.lock.Unlock()
	s.index = (s.index + 1) % s.length
	current := s.slots[s.index]
	current.RemoveIterCb(func(k string, value interface{}) bool {
		task := value.(iCronTask)
		task.ReduceRound(1)
		if task.GetRound() <= 0 {
			go s.handle(task)
			return true
		}
		return false
	})
}
func (s *Processor) handle(task iCronTask) error {
	if s.done {
		return nil
	}
	if !s.isPause {
		task.AddExecuted()
		rw, err := s.Dispatcher.HandleRequest(task)
		if err != nil {
			task.Errorf("%s执行出错:%v", task.GetName(), err)
		}
		task.SetResult(rw.Status(), rw.Data())
		s.saveHistory(task)

	}
	_, _, err := s.Add(task, false)
	if err != nil {
		fmt.Println(err)
	}
	return err

}

//Add 添加任务
func (s *Processor) Add(task iCronTask, r bool) (offset int, round int, err error) {
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
	if offset < 0 || round < 0 {
		return -1, -1, errors.New("next time less than now.2")
	}
	task.SetRound(round)
	s.slots[offset].Set(utility.GetGUID(), task)
	if r {
		s.Dispatcher.Handle(task.GetMethod(), task.GetService(), task.GetHandler().(dispatcher.HandlerFunc))
	}
	return
}
func (s *Processor) getOffset(now time.Time, next time.Time) (pos int, circle int) {
	d := next.Sub(now) //剩余时间
	delaySeconds := int(d/1e9) + 1
	intervalSeconds := int(s.span.Seconds())
	circle = int(delaySeconds / intervalSeconds / s.length)
	pos = int(s.index+delaySeconds/intervalSeconds) % s.length
	return
}

//Pause 暂停所有任务
func (s *Processor) Pause() error {
	s.isPause = true
	return nil
}

//Resume 恢复所有任务
func (s *Processor) Resume() error {
	s.isPause = false
	return nil
}

//Close 退出
func (s *Processor) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.done = true
	s.once.Do(func() {
		close(s.closeChan)
	})
}
