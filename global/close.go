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
	m.isClose = true
	close(m.close)
	closerLock.Lock()
	defer closerLock.Unlock()
	for _, c := range closers {
		c.Close()
	}

}
func (m *global) AddCloser(f interface{}) {
	closerLock.Lock()
	defer closerLock.Unlock()
	switch t := f.(type) {
	case io.Closer:
	case closeHandle:
		closers = append(closers, t)
	default:
	}
}
