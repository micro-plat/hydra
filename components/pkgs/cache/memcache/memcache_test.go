package memcache

import (
	"reflect"
	"testing"
)

func TestClient_Set(t *testing.T) {
	c, _ := New(WithAddress("192.168.106.58:11211"))
	type args struct {
		key       string
		value     string
		expiresAt int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{key: "hydra_test_key1", value: "100", expiresAt: 300}, wantErr: false},
		{name: "2", args: args{key: "hydra_test_key2", value: "value", expiresAt: 300}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Set(tt.args.key, tt.args.value, tt.args.expiresAt); (err != nil) != tt.wantErr {
				t.Errorf("Client.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//测试用例前.先执行Set的测试用例TestClient_Set
func TestClient_Get(t *testing.T) {
	c, _ := New(WithAddress("192.168.106.58:11211"))
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "1", args: args{key: ""}, wantErr: true},
		{name: "2", args: args{key: "xxxxxx"}, wantErr: true},
		{name: "3", args: args{key: "hydra_test_key1"}, want: "100", wantErr: false},
		{name: "4", args: args{key: "hydra_test_key2"}, want: "value", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

//测试用例前.先执行Set的测试用例TestClient_Set
func TestClient_Gets(t *testing.T) {
	c, _ := New(WithAddress("192.168.106.58:11211"))

	type args struct {
		key []string
	}
	tests := []struct {
		name    string
		args    args
		wantR   []string
		wantErr bool
	}{
		{name: "1", args: args{key: []string{"a", "", "hydra_test_key1", "hydra_test_key2"}}, wantR: []string{"", "", "100", "value"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := c.Gets(tt.args.key...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Gets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Client.Gets() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

//测试用例前.先执行Set的测试用例TestClient_Set
func TestClient_Add(t *testing.T) {
	c, _ := New(WithAddress("192.168.106.58:11211"))

	type args struct {
		key       string
		value     string
		expiresAt int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{key: "", value: "", expiresAt: 60}, wantErr: true},
		{name: "1", args: args{key: "hydra_test_key1", value: "value", expiresAt: 60}, wantErr: true},
		{name: "1", args: args{key: "hydra_test_key3", value: "value", expiresAt: 60}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Add(tt.args.key, tt.args.value, tt.args.expiresAt); (err != nil) != tt.wantErr {
				t.Errorf("Client.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//测试用例前.先执行Set的测试用例TestClient_Set
func TestClient_Decrement(t *testing.T) {
	c, _ := New(WithAddress("192.168.106.58:11211"))

	type args struct {
		key   string
		delta int64
	}
	tests := []struct {
		name    string
		args    args
		wantN   int64
		wantErr bool
	}{
		{name: "1", args: args{key: "hydra_test", delta: 100}, wantErr: true},
		{name: "2", args: args{key: "hydra_test_key2", delta: 100}, wantErr: true},
		{name: "3", args: args{key: "hydra_test_key1", delta: 100}, wantN: 0, wantErr: false},
		{name: "4", args: args{key: "hydra_test_key1", delta: 100}, wantN: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := c.Decrement(tt.args.key, tt.args.delta)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Decrement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Client.Decrement() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

//测试用例前.先执行Set的测试用例TestClient_Set
func TestClient_Increment(t *testing.T) {
	c, _ := New(WithAddress("192.168.106.58:11211"))

	type args struct {
		key   string
		delta int64
	}
	tests := []struct {
		name    string
		args    args
		wantN   int64
		wantErr bool
	}{
		{name: "1", args: args{key: "hydra_test", delta: 100}, wantErr: true},
		{name: "2", args: args{key: "hydra_test_key2", delta: 100}, wantErr: true},
		{name: "4", args: args{key: "hydra_test_key1", delta: -300}, wantErr: true},
		{name: "3", args: args{key: "hydra_test_key1", delta: 100}, wantN: 200, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := c.Increment(tt.args.key, tt.args.delta)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Increment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Client.Increment() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
