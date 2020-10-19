package global

import (
	"sync"
	"testing"

	"github.com/micro-plat/hydra/test/assert"
)

func Test_global_Close(t *testing.T) {
	Clear := func() {
		closers = nil
	}
	testcases := []struct {
		name      string
		closer    *mockClose
		expectLen int
	}{
		{name: "存在closer", closer: &mockClose{}, expectLen: 1},
		{name: "没有closer", closer: nil, expectLen: 0},
	}

	for _, tt := range testcases {
		Clear()
		m := &global{
			isClose: false,
			close:   make(chan struct{}),
		}
		if tt.closer != nil {
			m.AddCloser(tt.closer)
		}
		m.Close()

		assert.Equalf(t, len(closers), tt.expectLen, tt.name)
	}

}

func Test_global_Close_repeat(t *testing.T) {
	//测试close 重复调用
	m := &global{
		isClose: false,
		close:   make(chan struct{}),
	}
	m.Close() //第一次
	m.Close() //第二次
}

type mockClose struct {
	isClose bool
}

func (c *mockClose) Close() error {
	c.isClose = true
	return nil
}

type mockCloseNotMatch struct{}

func (mockCloseNotMatch) XClose() error {
	return nil
}

func Test_global_AddCloser(t *testing.T) {
	lock := sync.Mutex{}
	Clear := func() {
		closers = nil
	}

	testcases := []struct {
		name    string
		closer  []interface{}
		wantLen int
	}{
		{name: "正常的io.Closer", closer: []interface{}{&mockClose{}}, wantLen: 1},
		{name: "正常的closeHandle", closer: []interface{}{closeHandle(func() (err error) { return })}, wantLen: 1},
		{name: "正常的io.Closer+closeHandle", closer: []interface{}{&mockClose{}, closeHandle(func() (err error) { return })}, wantLen: 2},
		{name: "函数签名不匹配io.Closer", closer: []interface{}{&mockCloseNotMatch{}}, wantLen: 0},
		{name: "函数签名不匹配closeHandle", closer: []interface{}{func(x int) (err error) { return }}, wantLen: 0},
	}

	for _, tt := range testcases {
		lock.Lock()
		Clear()
		m := &global{}
		for _, c := range tt.closer {
			m.AddCloser(c)
		}
		assert.Equalf(t, tt.wantLen, len(closers), tt.name)
		lock.Unlock()
	}
}
