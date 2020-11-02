package global

import (
	"io"
	"sync"
)

type closeHandle func() error

func (c closeHandle) Close() error {
	return c()
}

var closers = make([]io.Closer, 0, 4)
var closerLock sync.Mutex

//Close 关闭全局应用
func (m *global) Close() {
	closerLock.Lock()
	defer closerLock.Unlock()
	if m.isClose {
		return
	}
	m.isClose = true
	close(m.close)
	for _, c := range closers {
		c.Close()
	}
}
func (m *global) AddCloser(f interface{}) {
	closerLock.Lock()
	defer closerLock.Unlock()

	if v, ok := f.(io.Closer); ok {
		closers = append(closers, v)
		return
	}
	if v, ok := f.(closeHandle); ok {
		closers = append(closers, v)
		return
	}
	if v, ok := f.(func() error); ok {
		closers = append(closers, closeHandle(v))
		return
	}
}
