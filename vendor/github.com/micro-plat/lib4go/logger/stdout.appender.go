package logger

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"sync"

	"github.com/zkfy/log"
)

//StdoutAppender 标准输出器
type StdoutAppender struct {
	name       string
	lastWrite  time.Time
	output     *log.Logger
	buffer     *bytes.Buffer
	ticker     *time.Ticker
	interval   time.Duration
	lock       sync.Mutex
	onceNotify sync.Once
	notifyChan chan struct{}
}

//NewStudoutAppender 构建基于文件流的日志输出对象
func NewStudoutAppender() (fa *StdoutAppender) {
	fa = &StdoutAppender{interval: time.Millisecond * 200, notifyChan: make(chan struct{})}
	fa.buffer = bytes.NewBufferString("")
	fa.output = log.New(fa.buffer, "", log.Llongcolor)
	fa.output.SetOutputLevel(log.Ldebug)
	fa.ticker = time.NewTicker(fa.interval)
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *StdoutAppender) Write(layout *Layout, event *LogEvent) error {
	if event.IsClose() {
		f.onceNotify.Do(func() {
			close(f.notifyChan)
		})
		return nil
	}
	current := GetLevel(event.Level)
	if GetLevel(layout.Level) > GetLevel(event.Level) {
		return nil
	}

	f.lastWrite = time.Now()
	f.lock.Lock()
	defer f.lock.Unlock()
	switch current {
	case ILevel_Debug:
		f.output.Debug(event.Output)
	case ILevel_Info:
		f.output.Info(event.Output)
	case ILevel_Warn:
		f.output.Warn(event.Output)
	case ILevel_Error:
		f.output.Error(event.Output)
	case ILevel_Fatal:
		f.output.Output("", log.Lfatal, 1, fmt.Sprintln(event.Output))
	}
	return nil
}

//writeTo 定时写入文件
func (f *StdoutAppender) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.lock.Lock()
				f.buffer.WriteTo(os.Stdout)
				f.buffer.Reset()
				f.lock.Unlock()
			} else {
				break START
			}
		}
	}

}

//Close 关闭当前appender
func (f *StdoutAppender) Close() error {
	select {
	case <-f.notifyChan:
	case <-time.After(time.Second):
	}
	f.ticker.Stop()
	f.lock.Lock()
	defer f.lock.Unlock()
	f.buffer.WriteTo(os.Stdout)
	return nil
}
