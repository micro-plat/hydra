package global

import "sync"

var closers = make([]func() error, 0, 4)
var closerLock sync.Mutex

//Close 关闭全局应用
func (m *global) Close() {
	m.isClose = true
	close(m.close)
	closerLock.Lock()
	defer closerLock.Unlock()
	for _, c := range closers {
		c()
	}

}
func (m *global) AddCloser(f func() error) {
	closerLock.Lock()
	defer closerLock.Unlock()
	closers = append(closers, f)
}
