package rlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/jsons"
	"github.com/micro-plat/lib4go/logger"
)

//RPCAppender 文件输出器
type RPCAppender struct {
	buffer     *bytes.Buffer
	lastWrite  time.Time
	ticker     *time.Ticker
	lock       sync.Mutex
	interval   time.Duration
	onceNotify sync.Once
	service    string
	notifyChan chan struct{}
	onceClose  sync.Once
}

//NewRPCAppender 构建writer日志输出对象
func NewRPCAppender(service string) (fa *RPCAppender) {
	fa = &RPCAppender{interval: time.Second, notifyChan: make(chan struct{})}
	fa.service = service
	fa.buffer = bytes.NewBufferString("")
	fa.ticker = time.NewTicker(fa.interval)
	go fa.writeTo()
	return
}

//Write 写入日志
func (f *RPCAppender) Write(layout *logger.Layout, event *logger.LogEvent) error {
	if event.IsClose() {
		f.onceNotify.Do(func() {
			close(f.notifyChan)
		})
		return nil
	}

	if logger.GetLevel(layout.Level) > logger.GetLevel(event.Level) {
		return nil
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.buffer.WriteString(",")
	f.buffer.WriteString(jsons.Escape(event.Output))
	f.lastWrite = time.Now()
	return nil
}

//Close 关闭当前appender
func (f *RPCAppender) Close() error {
	f.onceClose.Do(func() {
		select {
		case <-f.notifyChan:
		case <-time.After(time.Second):
		}
		f.ticker.Stop()
		f.lock.Lock()
		defer f.lock.Unlock()
		f.buffer.WriteTo(writeHandler(f.writeNow))
	})

	return nil

}

//writeTo 定时写入
func (f *RPCAppender) writeTo() {
START:
	for {
		select {
		case _, ok := <-f.ticker.C:
			if ok {
				f.lock.Lock()
				_, err := f.buffer.WriteTo(writeHandler(f.writeNow))
				f.buffer.Reset()
				f.lock.Unlock()
				if err != nil {
					logger.SysLog.Error(err)
				}
			} else {
				break START
			}
		}
	}
}

//writeNow 将数据写入远程请求
func (f *RPCAppender) writeNow(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	p[0] = byte('[')
	p = append(p, byte(']'))
	var buff bytes.Buffer
	if err := json.Compact(&buff, []byte(p)); err != nil {
		err = fmt.Errorf("json.compact.err:%v", err)
		return 0, err
	}
	_, err = components.Def.RPC().GetRegularRPC().Request(context.GetContextWithDefault(), f.service, buff.Bytes())
	if err != nil {
		return 0, fmt.Errorf("rlog写入日志失败 %s %w", f.service, err)
	}
	return len(p) - 1, nil
}

type writeHandler func(p []byte) (n int, err error)

func (w writeHandler) Write(p []byte) (n int, err error) {
	return w(p)
}

//Registry 注册日志组件
func Registry(platName string, addr string) error {

	//初始化注册中心
	registry, err := registry.NewRegistry(addr, global.Def.Log())
	if err != nil {
		return err
	}

	//获取远程配置
	layout, err := rlog.GetConfByAddr(registry, platName)
	if err != nil {
		return err
	}
	if layout.Disable {
		return nil
	}
	//注册日志组件
	logger.AddAppender(rlog.LogName, NewRPCAppender(layout.Service))
	logger.AddLayout(layout.ToLoggerLayout())
	return nil
}
