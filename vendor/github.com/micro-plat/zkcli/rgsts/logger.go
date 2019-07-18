package rgsts

import (
	"os"

	"github.com/micro-plat/lib4go/logger"

	"github.com/zkfy/log"
)

var _ logger.ILogging = &Logger{}

type Logger struct {
	*log.Logger
}

func newLogger() *Logger {
	l := &Logger{
		Logger: log.New(os.Stdout, "", log.Llongcolor),
	}
	l.SetOutputLevel(log.Ldebug)
	return l
}
