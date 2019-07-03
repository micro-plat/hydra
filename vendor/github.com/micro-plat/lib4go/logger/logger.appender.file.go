package logger

import (
	"bytes"
	"io"
	"sync"
	"time"

	"fmt"

	"github.com/micro-plat/lib4go/file"
)

//FileAppender 文件输出器
type FileAppender struct {
	name      string
	buffer    *bytes.Buffer
	lastWrite time.Time
	layout    *Appender
	file      io.WriteCloser
	ticker    *time.Ticker
	locker    sync.Mutex
	Level     int
}

//NewFileAppender 构建基于文件流的日志输出对象
func NewFileAppender(path string, layout *Appender) (fa *FileAppender, err error) {
	fa = &FileAppender{layout: layout}
	fa.Level = GetLevel(layout.Level)
	fa.buffer = bytes.NewBufferString("\n--------------------begin------------------------\n\n")
	intervalStr := layout.Interval
	if intervalStr == "" {
		intervalStr = "1s"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		err = fmt.Errorf("日志配置文件错误：%v", err)
		return
	}
	fa.ticker = time.NewTicker(interval)
	fa.file, err = file.CreateFile(path)
	if err != nil {
		return
	}
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *FileAppender) Write(event *LogEvent) {
	current := GetLevel(event.Level)
	if current < f.Level {
		return
	}
	f.locker.Lock()
	f.buffer.WriteString(event.Output)
	f.locker.Unlock()
	f.lastWrite = time.Now()
}

//Close 关闭当前appender
func (f *FileAppender) Close() {
	f.Level = ILevel_OFF
	f.ticker.Stop()
	f.locker.Lock()
	f.buffer.WriteString("\n---------------------end-------------------------\n")
	f.buffer.WriteTo(f.file)
	f.file.Close()
	f.locker.Unlock()
}

//writeTo 定时写入文件
func (f *FileAppender) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.locker.Lock()
				f.buffer.WriteTo(f.file)
				f.buffer.Reset()
				f.locker.Unlock()
			} else {
				break START
			}
		}
	}
}
