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
	name      string
	lastWrite time.Time
	layout    *Appender
	output    *log.Logger
	buffer    *bytes.Buffer
	ticker    *time.Ticker
	unq       string
	Level     int
	locker    sync.Mutex
}

//NewStudoutAppender 构建基于文件流的日志输出对象
func NewStudoutAppender(unq string, layout *Appender) (fa *StdoutAppender, err error) {
	fa = &StdoutAppender{layout: layout, unq: unq}
	fa.Level = GetLevel(layout.Level)
	fa.buffer = bytes.NewBufferString("")
	fa.output = log.New(fa.buffer, "", log.Llongcolor)
	intervalStr := layout.Interval
	if intervalStr == "" {
		intervalStr = "200ms"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		err = fmt.Errorf("日志配置文件错误：%v", err)
		return
	}
	fa.ticker = time.NewTicker(interval)
	fa.output.SetOutputLevel(log.Ldebug)
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *StdoutAppender) Write(event *LogEvent) {
	current := GetLevel(event.Level)
	if current < f.Level {
		return
	}
	f.lastWrite = time.Now()
	f.locker.Lock()
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
		f.output.Fatal(event.Output)
	}
	f.locker.Unlock()
}

//writeTo 定时写入文件
func (f *StdoutAppender) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.locker.Lock()
				f.buffer.WriteTo(os.Stdout)
				f.buffer.Reset()
				f.locker.Unlock()
			} else {
				break START
			}
		}
	}
}

//Close 关闭当前appender
func (f *StdoutAppender) Close() {
	f.Level = ILevel_OFF
	f.ticker.Stop()
	f.locker.Lock()
	f.buffer.WriteTo(os.Stdout)
	f.locker.Unlock()
}
