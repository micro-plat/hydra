package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// 测试的时候给这个变量赋值，可以进行回调
var testCallBack func(error)

func sysLoggerInfo(content ...interface{}) {
	sysLoggerWrite(SLevel_Info, fmt.Sprint(content...))
}
func sysLoggerError(content ...interface{}) {
	sysLoggerWrite(SLevel_Error, fmt.Sprint(content...))
}

func sysLoggerWrite(level string, content interface{}) {
	if strings.EqualFold(level, "") {
		level = "All"
	}

	e := &LogEvent{}
	e.Now = time.Now()
	e.Level = level
	e.Name = "sys"
	e.Session = CreateSession()
	e.Content = fmt.Sprintf("%v", content)
	e.Output = "%content%n"
	os.Stderr.WriteString(transform(e.Output, e))

	// 测试时候回调
	if testCallBack != nil {
		testCallBack(content.(error))
	}
}
