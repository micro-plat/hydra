package rlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/hydra/components"
	"github.com/micro-plat/hydra/components/rpcs/rpc"
	"github.com/micro-plat/hydra/conf/vars/rlog"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/hydra/registry"
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
	f.buffer.WriteString(jsonEscape(event.Output))
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
		f.writeNow(true)
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
				if err := f.writeNow(false); err != nil {
					logger.SysLog.Errorf("未正确写入日志:%v", err)
				}
			} else {
				break START
			}
		}
	}
}

//writeNow 将数据写入远程请求
func (f *RPCAppender) writeNow(b bool) (err error) {
	if f.buffer.Len() == 0 {
		return nil
	}
	write := func(p []byte) (n int, err error) {
		p[0] = byte('[')
		p = append(p, byte(']'))
		var buff bytes.Buffer
		if err := json.Compact(&buff, []byte(p)); err != nil {
			err = fmt.Errorf("json.compact.err:%v", err)
			return 0, err
		}

		_, err = components.Def.RPC().GetRegularRPC().Request(
			f.service,
			buff.String(),
			rpc.WithHeader("Plat", global.Def.PlatName),
			rpc.WithHeader("System", global.Def.SysName))
		if err != nil {
			return 0, fmt.Errorf("rlog写入日志失败 %s %w", f.service, err)
		}
		return len(p) - 1, nil
	}
	if !b {
		f.lock.Lock()
		defer f.lock.Unlock()
	}
	if _, err := f.buffer.WriteTo(writeHandler(write)); err != nil {
		f.buffer.Reset()
		return err
	}
	f.buffer.Reset()
	return nil
}

type writeHandler func(p []byte) (n int, err error)

func (w writeHandler) Write(p []byte) (n int, err error) {
	return w(p)
}

//Registry 注册日志组件
func Registry(platName string, addr string) error {

	//初始化注册中心
	registry, err := registry.GetRegistry(addr, global.Def.Log())
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

// Escape 把编码 \\u0026，\\u003c，\\u003e 替换为 &,<,>
func jsonEscape(input string) string {
	r := strings.Replace(input, "\\u0026", "&", -1)
	r = strings.Replace(r, "\\u003c", "<", -1)
	r = strings.Replace(r, "\\u003e", ">", -1)
	r = strings.Replace(r, "\n", "<br />", -1)
	r = strings.Replace(r, "\r", "", -1)
	r = strings.Replace(r, "\t", "    ", -1)
	r = strings.TrimRight(strings.TrimLeft(r, "\""), "\"")
	return strings.Replace(r, `"`, `\"`, -1)
}
