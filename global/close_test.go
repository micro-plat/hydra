package global

import (
	"sync"
	"testing"
)

func Test_global_Close(t *testing.T) {
	type fields struct {
		isClose bool
		close   chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "test1-重复调用",
			fields: fields{
				isClose: true,
				close:   make(chan struct{}),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &global{
				isClose: tt.fields.isClose,
				close:   tt.fields.close,
			}
			m.Close()
			m.Close()
		})
	}
}

type mockClose struct{}

func (mockClose) Close() error {
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

	t.Run("正常的io.Closer", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		Clear()
		m := &global{}
		m.AddCloser(&mockClose{})
		t.Log("-----", len(closers))
		if len(closers) != 1 {
			t.Errorf("正常的io.Closer未正常添加进去;expect:%d,actual:%d", 1, len(closers))
		}
	})

	t.Run("正常的closeHandle", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		closers = nil
		m := &global{}
		m.AddCloser(closeHandle(func() (err error) { return }))
		t.Log("-----", len(closers))
		if len(closers) != 1 {
			t.Errorf("正常的closeHandle未正常添加进去;expect:%d,actual:%d", 1, len(closers))
		}
	})
	t.Run("正常的io.Closer+closeHandle", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		closers = nil
		m := &global{}
		m.AddCloser(closeHandle(func() (err error) { return }))
		m.AddCloser(&mockClose{})
		t.Log("-----", len(closers))
		if len(closers) != 2 {
			t.Errorf("正常的io.Closer+closeHandle未正常添加进去;expect:%d,actual:%d", 2, len(closers))
		}
	})
	t.Run("函数签名不匹配io.Closer和closeHandle", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		closers = nil
		m := &global{}
		m.AddCloser(&mockCloseNotMatch{})
		m.AddCloser(func(x int) (err error) { return })
		t.Log("-----", len(closers))
		if len(closers) != 0 {
			t.Errorf("函数签名不匹配io.Closer和closeHandle;expect:%d,actual:%d", 0, len(closers))
		}
	})
}
