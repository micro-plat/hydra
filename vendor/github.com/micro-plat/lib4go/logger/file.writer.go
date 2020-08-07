package logger

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/file"
)

//writer 文件输出器
type writer struct {
	name       string
	buffer     *bytes.Buffer
	lastWrite  time.Time
	layout     *Layout
	interval   time.Duration
	file       io.WriteCloser
	ticker     *time.Ticker
	lock       sync.Mutex
	onceNotify sync.Once
	notifyChan chan struct{}
	Level      int
}

//newwriter 构建基于文件流的日志输出对象
func newWriter(path string, layout *Layout) (fa *writer, err error) {
	fa = &writer{layout: layout, interval: time.Second, notifyChan: make(chan struct{})}
	fa.Level = GetLevel(layout.Level)
	fa.buffer = bytes.NewBufferString("\n--------------------begin------------------------\n\n")
	fa.ticker = time.NewTicker(fa.interval)
	fa.file, err = file.CreateFile(path)
	if err != nil {
		return
	}
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *writer) Write(event *LogEvent) {
	if event.IsClose() {
		f.onceNotify.Do(func() {
			close(f.notifyChan)
		})
		return
	}
	if f.Level > GetLevel(event.Level) {
		return
	}
	f.lock.Lock()
	defer f.lock.Unlock()

	f.buffer.WriteString(event.Output)
	f.lastWrite = time.Now()
}

//Close 关闭当前appender
func (f *writer) Close() {

	//等待日志被关闭
	select {
	case <-f.notifyChan:
	case <-time.After(time.Second):
	}

	f.lock.Lock()
	defer f.lock.Unlock()
	f.ticker.Stop()
	f.buffer.WriteString("\n---------------------end-------------------------\n")
	f.buffer.WriteTo(f.file)
	f.file.Close()
}

//writeTo 定时写入文件
func (f *writer) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.lock.Lock()
				f.buffer.WriteTo(f.file)
				f.buffer.Reset()
				f.lock.Unlock()
			} else {
				break START
			}
		}
	}
}
