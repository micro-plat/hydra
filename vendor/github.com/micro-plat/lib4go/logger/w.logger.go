package logger

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"bytes"
)

//Logger 日志对象
type Logger struct {
	index        int64
	names        string
	sessions     string
	cacheSession bool
	tags         map[string]string
	isPause      bool
	DoPrint      func(content ...interface{})
	DoPrintf     func(format string, content ...interface{})
}

var loggerEventChan chan *LogEvent
var loggerCloserChan chan *Logger
var loggerPool *sync.Pool
var closeChan chan struct{}
var onceClose sync.Once
var done bool

func init() {
	loggerPool = &sync.Pool{
		New: func() interface{} {
			return New("")
		},
	}
	closeChan = make(chan struct{})
	loggerEventChan = make(chan *LogEvent, 2000)
	loggerCloserChan = make(chan *Logger, 1000)
	go loopDoLog()

}

//AddWriteThread 添加count个写线程用于并发写日志
func AddWriteThread(count int) {
	for i := 0; i < count; i++ {
		go loopDoLog()
	}
}

//New 根据一个或多个日志名称构建日志对象，该日志对象具有新的session id系统不会缓存该日志组件
func New(names string, tags ...string) (logger *Logger) {
	initConf()
	logger = &Logger{index: 100}
	logger.names = names
	logger.sessions = getSession()
	logger.DoPrint = logger.Info
	logger.DoPrintf = logger.Infof
	logger.tags = make(map[string]string)
	if len(tags) > 0 && len(tags) != 2 {
		panic(fmt.Sprintf("日志输入参数错误，扩展参数必须成对出现:%s,%v", names, tags))
	}
	for i := 0; i < len(tags)-1; i++ {
		logger.tags[tags[i]] = tags[i+1]
	}
	return logger
}

//CreateAndCache 使用当前session id创建logger,并将session id 缓存到当前协程
func CreateAndCache(name string, sessionID string, tags ...string) (logger *Logger) {
	log := GetSession(name, sessionID, tags...)
	cacheSession(sessionID)
	log.cacheSession = true
	return log
}

//GetSession 根据日志名称及session获取日志组件
func GetSession(name string, sessionID string, tags ...string) (logger *Logger) {
	logger = loggerPool.Get().(*Logger)
	logger.names = name
	logger.sessions = sessionID
	logger.tags = make(map[string]string)
	if len(tags) > 0 && len(tags) != 2 {
		panic(fmt.Sprintf("日志输入参数错误，扩展参数必须成对出现:%s,%v", name, tags))
	}
	for i := 0; i < len(tags)-1; i++ {
		logger.tags[tags[i]] = tags[i+1]
	}
	return logger
}

//SetName 设置日志名称
func (logger *Logger) SetName(name string) {
	logger.names = name
}

//Close 关闭当前日志组件
func (logger *Logger) Close() {
	select {
	case loggerCloserChan <- logger:
	default:
		if logger.cacheSession {
			removeSession()
		}
		loggerPool.Put(logger)
	}
}

//Pause 暂停记录
func (logger *Logger) Pause() {
	logger.isPause = true
}

//Resume 恢复记录
func (logger *Logger) Resume() {
	logger.isPause = false
	initConf()
}

//GetSessionID 获取当前日志的session id
func (logger *Logger) GetSessionID() string {
	if len(logger.sessions) > 0 {
		return logger.sessions
	}
	return ""
}

//Debug 输出debug日志
func (logger *Logger) Debug(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.log(SLevel_Debug, content...)
}

//Debugf 输出debug日志
func (logger *Logger) Debugf(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.logfmt(format, SLevel_Debug, content...)
}

//Info 输出info日志
func (logger *Logger) Info(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.log(SLevel_Info, content...)
}

//Infof 输出info日志
func (logger *Logger) Infof(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.logfmt(format, SLevel_Info, content...)
}

//Warn 输出info日志
func (logger *Logger) Warn(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.log(SLevel_Warn, content...)
}

//Warnf 输出info日志
func (logger *Logger) Warnf(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.logfmt(format, SLevel_Warn, content...)
}

//Error 输出Error日志
func (logger *Logger) Error(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.log(SLevel_Error, content...)
}

//Errorf 输出Errorf日志
func (logger *Logger) Errorf(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.logfmt(format, SLevel_Error, content...)
}

//Fatal 输出Fatal日志
func (logger *Logger) Fatal(content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.log(SLevel_Fatal, content...)
	Close()
	os.Exit(99)
}

//Fatalf 输出Fatalf日志
func (logger *Logger) Fatalf(format string, content ...interface{}) {
	if logger.isPause || globalPause {
		return
	}
	logger.logfmt(format, SLevel_Fatal, content...)
	Close()
	os.Exit(99)
}

//Fatalln 输出Fatal日志
func (logger *Logger) Fatalln(content ...interface{}) {
	logger.Fatal(content...)
}

//Print 输出info日志
func (logger *Logger) Print(content ...interface{}) {
	if logger.DoPrint == nil {
		return
	}
	logger.DoPrint(content...)
}

//Printf 输出info日志
func (logger *Logger) Printf(format string, content ...interface{}) {
	if logger == nil || logger.DoPrintf == nil {
		return
	}
	logger.DoPrintf(format, content...)
}

//Println 输出info日志
func (logger *Logger) Println(content ...interface{}) {
	logger.Print(content...)

}
func (logger *Logger) logfmt(f string, level string, content ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			SysLog.Errorf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	if !done {
		event := NewLogEvent(logger.names, level, logger.sessions, fmt.Sprintf(f, content...), logger.tags, atomic.AddInt64(&logger.index, 1))
		loggerEventChan <- event
	}
}
func (logger *Logger) log(level string, content ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			SysLog.Errorf("[Recovery] panic recovered:\n%s\n%s", err, getStack())
		}
	}()
	if !done {
		event := NewLogEvent(logger.names, level, logger.sessions, getString(content...), logger.tags, atomic.AddInt64(&logger.index, 1))
		loggerEventChan <- event
	}
}
func loopDoLog() {
	for {
		select {
		case logger := <-loggerCloserChan:
			loggerPool.Put(logger)
		case v, ok := <-loggerEventChan:
			if !ok {
				onceClose.Do(func() {
					close(closeChan)
				})
				return
			}
			logNow(v)
			v.Close()
		}
	}
}
func getString(c ...interface{}) string {
	if len(c) == 1 {
		return fmt.Sprintf("%v", c[0])
	}
	var buf bytes.Buffer
	for i := 0; i < len(c); i++ {
		buf.WriteString(fmt.Sprint(c[i]))
		if i != len(c)-1 {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

//Close 关闭所有日志组件
func Close() {
	if done {
		return
	}
	done = true
	close(loggerEventChan)
	<-closeChan
	logNow(GetEndWriteEvent())
	defWriter.Close()
}
