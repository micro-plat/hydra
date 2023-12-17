package logger

import (
	"bytes"
	"runtime"

	"github.com/micro-plat/lib4go/concurrent/cmap"
	"github.com/micro-plat/lib4go/utility"
)

var sessions = cmap.New(16)

//CreateSession create logger session
func CreateSession() string {
	return utility.GetGUID()[0:9]
}

//获取当前协程的session id 没有则自动创建
func getSession() string {
	v, ok := sessions.Get(getGoroutineID())
	if ok {
		return v.(string)
	}
	return CreateSession()
}

func cacheSession(session string) {
	sessions.Set(getGoroutineID(), session)
}
func removeSession() {
	sessions.Remove(getGoroutineID())
}

// GetGoroutineID 获取goroutine id
func getGoroutineID() string {
	b := make([]byte, 64)
	b = b[:32]
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}
