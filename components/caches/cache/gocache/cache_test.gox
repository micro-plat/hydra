package gocache

import (
	"testing"
)

func TestNew(t *testing.T) {
	_, err := NewByOpts()
	if (err != nil) != false {
		t.Errorf("New() error = %v, wantErr %v", err, false)
		return
	}
}

func TestClient_Add(t *testing.T) {
	c, _ := NewByOpts()
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
		{name: "1", args: args{key: "key1", value: "1", expiresAt: 60}, wantErr: false},
		{name: "2", args: args{key: "key1", value: "1", expiresAt: 60}, wantErr: true},
		{name: "3", args: args{key: "key2", value: "1", expiresAt: 60}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Add(tt.args.key, tt.args.value, tt.args.expiresAt); (err != nil) != tt.wantErr {
				t.Errorf("Client.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	type args struct {
		key string
	}
	c, _ := NewByOpts()
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "1", args: args{key: "key1"}, want: "1", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Set(tt.args.key, "1", 10)
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
