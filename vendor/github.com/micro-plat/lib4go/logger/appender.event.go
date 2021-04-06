package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/net"
)

var eventPool *sync.Pool
var appName string = filepath.Base(os.Args[0])

//EndWriteEvent 关闭日志事件
var EndWriteEvent = &LogEvent{isCloseEvent: true, Level: SLevel_ALL}

//GetEndWriteEvent 获取EndWriteEvent
func GetEndWriteEvent() *LogEvent {
	event := *EndWriteEvent
	event.Now = time.Now()
	return &event
}

func init() {
	eventPool = &sync.Pool{
		New: func() interface{} {
			return &LogEvent{}
		},
	}
}

//LogEvent 日志信息
type LogEvent struct {
	Level        string
	Now          time.Time
	Name         string
	Session      string
	Content      string
	Output       string
	Index        int64
	Tags         map[string]string
	isCloseEvent bool
}

//NewLogEvent 构建日志事件
func NewLogEvent(name string, level string, session string, content string, tags map[string]string, index int64) *LogEvent {
	e := eventPool.Get().(*LogEvent)
	e.Now = time.Now()
	e.Level = level
	e.Name = name
	e.Session = session
	e.Content = content
	e.Tags = tags
	e.Index = index
	e.isCloseEvent = false
	return e
}

//Event 获取转换后的日志事件
func (e *LogEvent) Event(format string) *LogEvent {
	e.Output = e.Transform(format)
	return e
}

var ip = net.GetLocalIPAddress()

//Transform 翻译模板串
func (e *LogEvent) Transform(tpl string) (result string) {
	word, _ := regexp.Compile(`%\w+`)
	isJson := (strings.HasPrefix(tpl, "[") || strings.HasPrefix(tpl, "{")) && (strings.HasSuffix(tpl, "]") || strings.HasSuffix(tpl, "}"))
	//@变量, 将数据放入params中
	result = word.ReplaceAllStringFunc(tpl, func(s string) string {
		key := s[1:]
		switch key {
		case "app":
			return appName
		case "session":
			return e.Session
		case "date":
			return e.Now.Format("20060102")
		case "datetime":
			return e.Now.Format("2006/01/02 15:04:05")
		case "yy":
			return e.Now.Format("2006")
		case "mm":
			return e.Now.Format("01")
		case "dd":
			return e.Now.Format("02")
		case "hh":
			return e.Now.Format("15")
		case "mi":
			return e.Now.Format("04")
		case "ss":
			return e.Now.Format("05")
		case "ms":
			return fmt.Sprintf("%06d", e.Now.Nanosecond()/1e3)
		case "level":
			return strings.ToLower(e.Level)
		case "l":
			if len(e.Level) > 0 {
				return strings.ToLower(e.Level)[:1]
			}
			return ""
		case "name":
			return e.Name
		case "pid":
			return fmt.Sprintf("%d", os.Getpid())
		case "n":
			return "\n"
		case "caller":
			return getCaller(8)
		case "content":
			if isJson {
				buff, err := json.Marshal(e.Content)
				if err != nil {
					return e.Content
				}
				if len(buff) > 2 {
					return string(string(buff[1 : len(buff)-1]))
				}
			}
			return e.Content
		case "index":
			return fmt.Sprintf("%d", e.Index)
		case "ip":
			return ip
		default:
			v, ok := e.Tags[key]
			if ok {
				return v
			}
			return ""
		}
	})
	return
}

//IsClose 是否是关闭事件
func (e *LogEvent) IsClose() bool {
	return e.isCloseEvent
}

//Close 关闭回收日志
func (e *LogEvent) Close() {
	eventPool.Put(e)
}

func getCaller(index int) string {
	defer recover()
	_, file, line, ok := runtime.Caller(index)
	if ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return ""
}
