package logger

import "github.com/micro-plat/lib4go/concurrent/cmap"

type LoggerAppenderFactory interface {
	GetType() string
	MakeAppender(l *Appender, event *LogEvent) (IAppender, error)
	MakeUniq(l *Appender, event *LogEvent) string
}

var registedFactory cmap.ConcurrentMap

func init() {
	registedFactory = cmap.New(2)
}
func RegistryFactory(f LoggerAppenderFactory, appender *Appender) {
	registedFactory.SetIfAbsent(f.GetType(), f)
	manager.append(appender)
}
func UnRegistryFactory(tp string) {
	registedFactory.Remove(tp)
	manager.remote(tp)
}
