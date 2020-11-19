package container

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_newVers(t *testing.T) {
	l := newVers()
	if !reflect.DeepEqual(l, &vers{keys: make(map[string]*ver)}) {
		t.Error("newVers() didn't return *vers")
	}
}

func Test_vers_Add(t *testing.T) {
	type args struct {
		tp   string
		name string
		key  string
	}
	s := newVers()
	tests := []struct {
		name string
		args args
	}{
		{name: "1", args: args{tp: "test", name: "test1", key: "123"}},
		{name: "2", args: args{tp: "test", name: "test2", key: "123"}},
		{name: "3", args: args{tp: "test", name: "test1", key: "321"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Add(tt.args.tp, tt.args.name, tt.args.key)
			tps := fmt.Sprintf("%s-%s", tt.args.tp, tt.args.name)
			if ver, ok := s.keys[tps]; !ok {
				t.Errorf("vers.keys[%s] doesn't exist", tps)
			} else {
				if ver.current != tt.args.key {
					t.Errorf("vers.keys[%s].current isn't right", tps)
				}
				if ver.keys[len(ver.keys)-1] != tt.args.key {
					t.Errorf("vers.keys[%s].keys isn't right", tps)
				}
			}
		})
	}
}

func Test_vers_Remove(t *testing.T) {
	type data struct {
		tp   string
		name string
		key  string
	}

	type args struct {
		f func(key string) bool
	}

	s := newVers()
	f := func(key string) bool {
		t.Log("移除key", key)
		return true
	}

	tests := []struct {
		name string
		data data
		args args
	}{
		{name: "1", data: data{tp: "test", name: "test1", key: "123"}, args: args{f: f}},
		{name: "2", data: data{tp: "test", name: "test2", key: "123"}, args: args{f: f}},
		{name: "3", data: data{tp: "test", name: "test1", key: "321"}, args: args{f: f}},
	}

	for _, tt := range tests {
		s.Add(tt.data.tp, tt.data.name, tt.data.key)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.Remove(tt.args.f)
		})
	}

	for k, v := range s.keys {
		t.Logf("当前剩余key %s:%+v", k, v.keys)
	}
}
