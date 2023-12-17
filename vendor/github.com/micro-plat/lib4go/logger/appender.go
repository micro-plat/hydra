package logger

import (
	"fmt"
	"sync"
)

//IAppender 定义appender接口
type IAppender interface {
	Write(*Layout, *LogEvent) error
	Close() error
}

type appenderWriter struct {
	appenders map[string]IAppender
	layouts   []*Layout
	lock      sync.RWMutex
}

func newAppenderWriter() *appenderWriter {
	return &appenderWriter{
		appenders: make(map[string]IAppender),
		layouts:   make([]*Layout, 0, 2),
	}
}

//AddAppender  添加appender
func (a *appenderWriter) AddAppender(typ string, i IAppender) {
	a.lock.Lock()
	defer a.lock.Unlock()
	if _, ok := a.appenders[typ]; ok {
		panic(fmt.Errorf("不能重复注册appender:%s", typ))
	}
	a.appenders[typ] = i
}

//RemoveAppender 移除某个Appender
func (a *appenderWriter) RemoveAppender(typ string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	delete(a.appenders, typ)
}

//AddLayout 添加layout配置
func (a *appenderWriter) AddLayout(layouts ...*Layout) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, layout := range layouts {
		if layout.Level == SLevel_OFF {
			continue
		}
		if _, ok := a.appenders[layout.Type]; ok {
			a.layouts = append(a.layouts, layout)
		}
	}
}

//ResetLayout 重置layout配置
func (a *appenderWriter) ResetLayout(layouts ...*Layout) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.layouts = make([]*Layout, 0, 2)
	for _, layout := range layouts {
		if layout.Level == SLevel_OFF {
			continue
		}
		if _, ok := a.appenders[layout.Type]; ok {
			a.layouts = append(a.layouts, layout)
		}
	}
}

//Log 记录日志信息
func (a *appenderWriter) Log(event *LogEvent) {
	defer func() {
		if err := recover(); err != nil {
			SysLog.Errorf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, layout := range a.layouts {
		if GetLevel(layout.Level) > GetLevel(event.Level) {
			continue
		}

		e := event.Event(layout.Layout)
		if apppender, ok := a.appenders[layout.Type]; ok {
			apppender.Write(layout, e)
			continue
		}
	}
}

//Close 关闭日志
func (a *appenderWriter) Close() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, v := range a.appenders {
		appender := v.(IAppender)
		appender.Close()
	}
	return nil

}

//默认appender写入器
var defWriter = newAppenderWriter()

//AddAppender 添加appender
func AddAppender(typ string, i IAppender) {
	defWriter.AddAppender(typ, i)
}

//RemoveAppender 移除Appender
func RemoveAppender(typ string) {
	defWriter.RemoveAppender(typ)
}

//RemoveStdoutAppender 移除stdout appender
func RemoveStdoutAppender() {
	defWriter.RemoveAppender("stdout")
}

//AddLayout 添加日志输出配置
func AddLayout(l ...*Layout) {
	defWriter.AddLayout(l...)
}

func logNow(event *LogEvent) {
	defWriter.Log(event)
}

//进行日志配置文件初始化
func init() {
	AddAppender("file", NewFileAppender())
	AddAppender("stdout", NewStudoutAppender())
}
