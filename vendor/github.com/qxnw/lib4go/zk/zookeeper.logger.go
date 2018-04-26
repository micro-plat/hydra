package zk

import "github.com/qxnw/lib4go/logger"

type zkLogger struct {
	logger *logger.Logger
}

func (l *zkLogger) Printf(f string, c ...interface{}) {
	if l.logger != nil {
		l.logger.Printf(f, c...)
	}
}
