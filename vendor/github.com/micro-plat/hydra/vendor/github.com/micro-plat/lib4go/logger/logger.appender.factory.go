package logger

import (
	"fmt"
	"strings"
)

type loggerAppenderFactory struct {
}

//MakeAppender 根据appender配置及日志信息生成appender对象
func (f *loggerAppenderFactory) MakeAppender(l *Appender, event *LogEvent) (IAppender, error) {
	switch strings.ToLower(l.Type) {
	case appender_file:
		return NewFileAppender(f.MakeUniq(l, event), l)
	case appender_stdout:
		return NewStudoutAppender(f.MakeUniq(l, event), l)
	default:
		if value, b := registedFactory.Get(l.Type); b {
			fa := value.(LoggerAppenderFactory)
			return fa.MakeAppender(l, event)
		}
	}
	return nil, fmt.Errorf("不支持的日志类型:%s", l.Type)
}

//MakeUniq 根据appender配置及日志信息生成appender唯一编号
func (f *loggerAppenderFactory) MakeUniq(l *Appender, event *LogEvent) string {
	switch strings.ToLower(l.Type) {
	case appender_file:
		return transform(l.Path, event)
	default:
		if value, b := registedFactory.Get(l.Type); b {
			fa := value.(LoggerAppenderFactory)
			return fa.MakeUniq(l, event)
		}
	}
	return l.Type
}
