package container

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_newHistories(t *testing.T) {
	l := newHistories()
	if !reflect.DeepEqual(l, &histories{records: make(map[string]*history)}) {
		t.Error("newHistories() didn't return *histories")
	}
}

func Test_histories_Add(t *testing.T) {
	type args struct {
		tp   string
		name string
		key  string
	}
	s := newHistories()
	tests := []struct {
		name string
		args args
	}{
		{name: "1. VerAdd-添加数据", args: args{tp: "test", name: "test1", key: "123"}},
		{name: "2. VerAdd-添加数据", args: args{tp: "test", name: "test2", key: "123"}},
		{name: "3. VerAdd-添加数据", args: args{tp: "test", name: "test1", key: "321"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Add(fmt.Sprintf("%s_%s", tt.args.tp, tt.args.name), tt.args.key)
			tps := fmt.Sprintf("%s_%s", tt.args.tp, tt.args.name)
			if ver, ok := s.records[tps]; !ok {
				t.Errorf("histories.keys[%s] doesn't exist", tps)
			} else {
				if ver.current != tt.args.key {
					t.Errorf("histories.keys[%s].current isn't right", tps)
				}
				if ver.keys[len(ver.keys)-1] != tt.args.key {
					t.Errorf("histories.keys[%s].keys isn't right", tps)
				}
			}
		})
	}
}

func Test_histories_Remove(t *testing.T) {
	type data struct {
		tp   string
		name string
		key  string
	}

	type args struct {
		f func(key string) bool
	}

	s := newHistories()
	f := func(key string) bool {
		t.Log("移除key", key)
		return true
	}

	tests := []struct {
		name string
		data data
		args args
	}{
		{name: "1. historiesRemove-移除数据-数据存在", data: data{tp: "test", name: "test1", key: "123"}, args: args{f: f}},
		{name: "2. historiesRemove-移除数据-数据不存在", data: data{tp: "test", name: "test2", key: "123"}, args: args{f: f}},
		{name: "3. historiesRemove-移除数据-数据存在1", data: data{tp: "test", name: "test1", key: "321"}, args: args{f: f}},
	}

	for _, tt := range tests {
		s.Add(fmt.Sprintf("%s_%s", tt.data.tp, tt.data.name), tt.data.key)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Remove(tt.args.f)
		})
	}

	for k, v := range s.records {
		t.Logf("当前剩余key %s:%+v", k, v.keys)
	}
}
