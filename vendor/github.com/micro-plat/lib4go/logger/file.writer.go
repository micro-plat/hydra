package logger

import (
	"bufio"
	"io"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/file"
)

//writer 文件输出器
type writer struct {
	name       string
	writer     *bufio.Writer
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

//newwriter 构建基于文件流的日志输出对象,使用带缓冲区的文件写入，缓存区达到4K或每隔3秒写入一次文件。
func newWriter(path string, layout *Layout) (fa *writer, err error) {
	fa = &writer{layout: layout, interval: time.Second * 3, notifyChan: make(chan struct{})}
	fa.file, err = file.CreateFile(path)
	if err != nil {
		return
	}
	fa.Level = GetLevel(layout.Level)
	fa.ticker = time.NewTicker(fa.interval)
	fa.writer = bufio.NewWriterSize(fa.file, 4096)
	fa.writer.WriteString("\n--------------------begin------------------------\n\n")
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *writer) Write(event *LogEvent) {
	if event.IsClose() {
		f.Close()
		return
	}
	if f.Level > GetLevel(event.Level) {
		return
	}
	f.lock.Lock()
	defer f.lock.Unlock()

	f.writer.WriteString(event.Output)
	f.lastWrite = time.Now()
}

//Close 关闭当前appender
func (f *writer) Close() {
	f.onceNotify.Do(func() {
		close(f.notifyChan)
		f.ticker.Stop()
	})
}

//writeTo 定时写入文件
func (f *writer) writeTo() {
START:
	for {
		select {
		case <-f.notifyChan:
			f.lock.Lock()
			f.writer.WriteString("\n---------------------end-------------------------\n")
			f.writer.Flush()
			f.file.Close()
			f.lock.Unlock()
			break START
		case <-f.ticker.C:
			f.lock.Lock()
			if err := f.writer.Flush(); err != nil {
				SysLog.Error("file.write.err:", err)
			}
			f.lock.Unlock()
		}
	}
}
