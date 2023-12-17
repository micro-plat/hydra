package logger

import "strings"

const (
	ILevel_ALL = iota
	ILevel_Debug
	ILevel_Info
	ILevel_Warn
	ILevel_Error
	ILevel_Fatal
	ILevel_OFF
)
const (
	SLevel_OFF   = "Off"
	SLevel_Info  = "Info"
	SLevel_Warn  = "Warn"
	SLevel_Error = "Error"
	SLevel_Fatal = "Fatal"
	SLevel_Debug = "Debug"
	SLevel_ALL   = "All"
)

const (
	appender_file   = "file"
	appender_stdout = "stdout"
)

var levelMap = map[string]int{
	SLevel_OFF:   ILevel_OFF,
	SLevel_Info:  ILevel_Info,
	SLevel_Warn:  ILevel_Warn,
	SLevel_Error: ILevel_Error,
	SLevel_Fatal: ILevel_Fatal,
	SLevel_Debug: ILevel_Debug,
	SLevel_ALL:   ILevel_ALL,
}

//GetLevel 获取日志等级编号
func GetLevel(name string) int {
	if len(name) > 0 {
		if l, ok := levelMap[strings.ToUpper(name[:1])+name[1:]]; ok {
			return l
		}
	}

	return ILevel_ALL
}

//LogWriter 提供Write函数的日志方法
type LogWriter func(content ...interface{})

func (l LogWriter) Write(p []byte) (n int, err error) {
	l(string(p))
	return len(p), nil
}
