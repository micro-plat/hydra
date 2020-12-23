package logs

import (
	"os"

	"github.com/micro-plat/lib4go/logger"
	"github.com/zkfy/log"
)

//Log 日志
var Log logger.ILogging = New()

//Logger 日志组件
type Logger struct {
	*log.Logger
}

//New 日志组件
func New() *Logger {
	l := &Logger{
		Logger: log.New(os.Stdout, "", log.Llongcolor),
	}
	l.SetOutputLevel(log.Ldebug)
	return l
}
