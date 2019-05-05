package logger

import "time"
import "sync"

var eventPool *sync.Pool

func init() {
	eventPool = &sync.Pool{
		New: func() interface{} {
			return &LogEvent{}
		},
	}
}

type LogEvent struct {
	Level   string
	Now     time.Time
	Name    string
	Session string
	Content string
	Output  string
	Index   int64
	Tags    map[string]string
}

func NewLogEvent(name string, level string, session string, content string, tags map[string]string, index int64) *LogEvent {
	e := eventPool.Get().(*LogEvent)
	e.Now = time.Now()
	e.Level = level
	e.Name = name
	e.Session = session
	e.Content = content
	e.Tags = tags
	e.Index = index
	return e
}
func (l *LogEvent) Close() {
	eventPool.Put(l)
}
