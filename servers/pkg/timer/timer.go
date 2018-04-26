package timer

import (
	"sync"
	"time"

	"github.com/zkfy/cron"
)

type Timer struct {
	schedule  cron.Schedule
	mailboxes []chan struct{}
	done      bool
	closeChan chan struct{}
	sync      sync.Once
}

func NewTimer(c string) (m *Timer, err error) {
	m = &Timer{}
	m.schedule, err = cron.ParseStandard(c)
	m.mailboxes = make([]chan struct{}, 0, 2)
	m.closeChan = make(chan struct{})
	if err != nil {
		return nil, err
	}
	return m, nil
}
func (m *Timer) Subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	m.mailboxes = append(m.mailboxes, ch)
	return ch
}
func (m *Timer) Start() {
	go m.run()
}
func (m *Timer) Close() {
	m.sync.Do(func() {
		close(m.closeChan)
	})
	m.done = true
}

func (m *Timer) run() {
	var intervalTicker int64
LOOP:
	for {
		select {
		case <-m.closeChan:
			break LOOP
		case <-time.After(time.Second):
			if m.done {
				break LOOP
			}
			now := time.Now()
			if intervalTicker > now.Unix() {
				break
			}
			go m.notify()
			intervalTicker = m.schedule.Next(now).Unix()
		}
	}
}
func (m *Timer) notify() {
	for _, c := range m.mailboxes {
		c <- struct{}{}
	}
}
