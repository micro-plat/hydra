package global

import (
	"sync"
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func Test_global_Close(t *testing.T) {
	Clear := func() {
		closers = nil
	}
	testcases := []struct {
		name      string
		closer    []*mockClose
		expectLen int
	}{
		{name: "1. 没有closer", closer: nil, expectLen: 0},
		{name: "2. 存在closer-单个", closer: []*mockClose{&mockClose{}}, expectLen: 1},
		{name: "3. 存在closer-多个", closer: []*mockClose{&mockClose{}, &mockClose{}}, expectLen: 2},
	}

	for _, tt := range testcases {
		Clear()
		m := &global{
			isClose: false,
			close:   make(chan struct{}),
		}
		if len(tt.closer) > 0 {
			for i := range tt.closer {
				m.AddCloser(tt.closer[i])
			}
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
	defer func() {
		if obj := recover(); obj != nil {
			assert.Equal(t, nil, obj, "1. 重复调用Close出现崩溃")
		}
	}()
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
		{name: "1. 正常的io.Closer", closer: []interface{}{&mockClose{}}, wantLen: 1},
		{name: "2. 正常的closeHandle", closer: []interface{}{closeHandle(func() (err error) { return })}, wantLen: 1},
		{name: "3. 正常的io.Closer+closeHandle", closer: []interface{}{&mockClose{}, closeHandle(func() (err error) { return })}, wantLen: 2},
		{name: "4. 函数签名不匹配io.Closer", closer: []interface{}{&mockCloseNotMatch{}}, wantLen: 0},
		{name: "5. 函数签名不匹配closeHandle", closer: []interface{}{func(x int) (err error) { return }}, wantLen: 0},
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
