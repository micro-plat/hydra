package logger

import "fmt"

//SysLog 系统日志
var SysLog = newSysLogger()

type sysLogger struct {
	appender IAppender
	layout   *Layout
}

func newSysLogger() *sysLogger {
	return &sysLogger{
		appender: NewStudoutAppender(),
		layout:   &Layout{Layout: "[%datetime.%ms][%l][%session]%content", Level: SLevel_ALL},
	}
}
func (s *sysLogger) Error(content ...interface{}) {
	s.Log(NewLogEvent("sys", SLevel_Error, CreateSession(), fmt.Sprint(content...), nil, 0).Event(s.layout.Layout))
}
func (s *sysLogger) Errorf(f string, content ...interface{}) {
	s.Log(NewLogEvent("sys", SLevel_Error, CreateSession(), fmt.Sprintf(f, content...), nil, 0).Event(s.layout.Layout))
}
func (s *sysLogger) Log(event *LogEvent) {
	s.appender.Write(s.layout, event)
}
