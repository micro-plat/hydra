package logs

import (
	"os"

	"github.com/zkfy/log"
)

//Logger 日志组件
type Logger struct {
	*log.Logger
}

//Log 日志
var Log = New()

//New 日志组件
func New() *Logger {
	l := &Logger{
		Logger: log.New(os.Stdout, "", log.Llongcolor),
	}
	l.SetOutputLevel(log.Ldebug)
	return l
}
