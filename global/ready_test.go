package global

import (
	"errors"
	"sync"
	"testing"
)

var lock sync.Mutex

func TestOnReady(t *testing.T) {
	t.Run("参数类型为func()-isReady=false,noerror", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		isReady = false
		funcs = nil
		OnReady(func() {})
		if len(funcs) != 0 {
			t.Errorf("参数类型为func()-isReady=false,noerror;expect:%d,actual:%d", 0, len(funcs))
		}
	})

	t.Run("参数类型为func()-isReady=true,noerror", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		isReady = true
		funcs = nil
		OnReady(func() {})
		if len(funcs) != 1 {
			t.Errorf("参数类型为func()-isReady=true,noerror;expect:%d,actual:%d", 1, len(funcs))
		}
	})

	t.Run("参数类型为func()error-isReady=false,noerror", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		isReady = false
		funcs = nil
		OnReady(func() error { return nil })
		if len(funcs) != 1 {
			t.Errorf("参数类型为func()-isReady=false,noerror;expect:%d,actual:%d", 1, len(funcs))
		}
	})

	t.Run("参数类型为func()error-isReady=true,noerror", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		isReady = true
		funcs = nil
		OnReady(func() error { return nil })
		if len(funcs) != 0 {
			t.Errorf("参数类型为func()error-isReady=true,noerror;expect:%d,actual:%d", 0, len(funcs))
		}
	})

	t.Run("参数类型为func()error-isReady=false,error", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		isReady = false
		funcs = nil
		OnReady(func() error { return errors.New("error") })
		if len(funcs) != 1 {
			t.Errorf("参数类型为func()-isReady=false,error;expect:%d,actual:%d", 1, len(funcs))
		}
	})

	t.Run("参数类型为func()error-isReady=true,error", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		isReady = true
		funcs = nil
		defer func() {
			if obj := recover(); obj == nil {
				t.Errorf("参数类型为func()error-isReady=true,error;expect:%s,actual:%s", "error", "null")
			}
		}()
		OnReady(func() error { return errors.New("error") })
	})

	t.Run("参数类型为不匹配", func(t *testing.T) {
		lock.Lock()
		defer lock.Unlock()
		isReady = true
		funcs = nil
		defer func() {
			if obj := recover(); obj == nil {
				t.Errorf("参数类型为func()error-isReady=true,error;expect:%s,actual:%s", "error", "null")
			}
		}()
		OnReady(func(x int) error { return errors.New("error") })
	})
}
