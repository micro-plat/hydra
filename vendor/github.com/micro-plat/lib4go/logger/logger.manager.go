package logger

import (
	"errors"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/concurrent/cmap"
)

var isOpen bool

type ILoggerAppenderFactory interface {
	MakeAppender(*Appender, *LogEvent) (IAppender, error)
	MakeUniq(*Appender, *LogEvent) string
}

type loggerManager struct {
	appenders cmap.ConcurrentMap
	factory   ILoggerAppenderFactory
	configs   []*Appender
	ticker    *time.Ticker
	lock      sync.RWMutex
	isClose   bool
}
type appenderEntity struct {
	appender IAppender
	last     time.Time
}

func newLoggerManager() (m *loggerManager, err error) {
	m = &loggerManager{isClose: false}
	m.factory = &loggerAppenderFactory{}
	m.appenders = cmap.New(2)
	m.configs = ReadConfig()
	isOpen = len(m.configs) > 0
	if isOpen {
		m.ticker = time.NewTicker(time.Second * 300)
		go m.clearUp()
		return m, nil
	}
	return nil, errors.New("未启动日志")
}
func (a *loggerManager) append(app *Appender) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, v := range a.configs {
		if v.Type == app.Type {
			return
		}
	}
	a.configs = append(a.configs, app)
}
func (a *loggerManager) remote(t string) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for i, v := range a.configs {
		if v.Type == t {
			a.configs = append(a.configs[:i], a.configs[i+1])
		}
	}
}

// Log 将日志内容写入appender, 如果appender不存在则创建
// callBack回调函数,如果不需要传nil
func (a *loggerManager) Log(event *LogEvent) {
	if a.isClose {
		return
	}
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, config := range a.configs {
		uniq := a.factory.MakeUniq(config, event)
		_, currentAppender, err := a.appenders.SetIfAbsentCb(uniq, func(p ...interface{}) (entity interface{}, err error) {
			l := p[0].(*Appender)
			app, err := a.factory.MakeAppender(l, event)
			if err != nil {
				return
			}
			entity = &appenderEntity{appender: app, last: time.Now()}
			return
		}, config)
		if err == nil {
			capp := currentAppender.(*appenderEntity)
			a.write(capp, config.Layout, event)
		} else {
			sysLoggerError(err)
		}
	}
}
func (a *loggerManager) write(capp *appenderEntity, format string, event *LogEvent) {
	defer func() {
		if r := recover(); r != nil {
			sysLoggerError(r)
		}
	}()
	capp.last = time.Now()
	event.Output = transform(format, event)
	capp.appender.Write(event)
}
func (a *loggerManager) clearUp() {
START:
	for {
		select {
		case _, ok := <-a.ticker.C:
			if ok {
				count := a.appenders.RemoveIterCb(func(key string, v interface{}) bool {
					apd := v.(*appenderEntity)
					if time.Now().Sub(apd.last).Seconds() > 10 {
						apd.appender.Close()
						return true
					}
					return false
				})
				if count > 0 {
					//sysLoggerInfo("已移除:", count)
				}
			} else {
				break START
			}
		}
	}
}

func (a *loggerManager) Close() {
	a.isClose = true
	if a.ticker != nil {
		a.ticker.Stop()
	}

	a.appenders.RemoveIterCb(func(key string, v interface{}) bool {
		apd := v.(*appenderEntity)
		apd.appender.Close()
		return true
	})
}
