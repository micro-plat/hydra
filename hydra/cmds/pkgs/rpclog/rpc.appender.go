package rpclog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/jsons"
	"github.com/micro-plat/lib4go/logger"
	"github.com/zkfy/log"
)

//RPCAppender 文件输出器
type RPCAppender struct {
	output      *log.Logger
	name        string
	buffer      *bytes.Buffer
	lastWrite   time.Time
	layout      *logger.Appender
	ticker      *time.Ticker
	locker      sync.Mutex
	writer      io.WriteCloser
	intervalStr string
	Level       int
}

//NewRPCAppender 构建writer日志输出对象
func NewRPCAppender(writer io.WriteCloser, layout *logger.Appender) (fa *RPCAppender, err error) {
	fa = &RPCAppender{layout: layout, writer: writer}
	fa.output = log.New(fa.buffer, "", log.Llongcolor)
	fa.output.SetOutput(os.Stderr)
	fa.output.SetOutputLevel(log.Ldebug)
	fa.Level = logger.GetLevel(layout.Level)
	fa.intervalStr = layout.Interval
	fa.buffer = bytes.NewBufferString("")
	interval, err := time.ParseDuration(fa.intervalStr)
	if err != nil {
		err = fmt.Errorf("rpc日志的interval字段配置有误:%v", interval)
		return
	}
	fa.ticker = time.NewTicker(interval)
	go fa.writeTo()
	return
}

//Reset 重置保存周期
func (f *RPCAppender) Reset(intervalStr string, writer io.WriteCloser) error {
	f.locker.Lock()
	defer f.locker.Unlock()
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		err = fmt.Errorf("rpc日志的interval字段配置有误:%v", interval)
		return err
	}
	//f.ticker.Stop()
	f.ticker = time.NewTicker(interval)
	f.writer = writer
	return nil
}

//Write 写入日志
func (f *RPCAppender) Write(event *logger.LogEvent) {
	current := logger.GetLevel(event.Level)
	if current < f.Level {
		return
	}
	f.locker.Lock()
	f.buffer.WriteString(",")
	f.buffer.WriteString(jsons.Escape(event.Output))
	f.locker.Unlock()
	f.lastWrite = time.Now()
}

//Close 关闭当前appender
func (f *RPCAppender) Close() {
	f.Level = logger.ILevel_OFF
	f.ticker.Stop()
	f.locker.Lock()
	f.buffer.WriteTo(f.writer)
	f.writer.Close()
	f.locker.Unlock()
}

//writeTo 定时写入
func (f *RPCAppender) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.locker.Lock()
				_, err := f.buffer.WriteTo(f.writer)
				f.buffer.Reset()
				f.locker.Unlock()
				if err != nil {
					log.Error(err)
				}
			} else {
				break START
			}
		}
	}
}
