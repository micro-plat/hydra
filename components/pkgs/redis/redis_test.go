package redis

import (
	"testing"

	varredis "github.com/micro-plat/hydra/conf/vars/redis"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []varredis.Option
	}
	tests := []struct {
		name    string
		args    args
		wantR   *Client
		wantErr bool
	}{
		{name: "1. 初始化空redis对象", args: args{}, wantErr: true},
		{name: "2. 多ip初始化redis对象", args: args{opts: []varredis.Option{varredis.WithRaw(`{"addrs":["192.168.0.111:6379","192.168.0.112:6379"]}`)}}, wantErr: false},
		{name: "3. 单ip初始化redis对象", args: args{opts: []varredis.Option{varredis.WithRaw(`{"addrs":["192.168.0.111:6379"]}`)}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewByOpts(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// if !reflect.DeepEqual(gotR, tt.wantR) {
			// 	t.Errorf("New() = %v, want %v", gotR, tt.wantR)
			// }
		})
	}
}
