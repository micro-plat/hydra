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
		{name: "1", args: args{}, wantErr: true},
		{name: "2", args: args{opts: []varredis.Option{varredis.WithRaw(`{"addrs":["192.168.0.111:6379","192.168.0.112:6379"]}`)}}, wantErr: false},
		//{name: "3", args: args{opts: []Option{WithRaw(`{"addrs":["192.168.0.111:6379","192.168.0.116:6379"],"master":"master"}`)}}, wantErr: false}, //redis 的sentinel没有启动
		{name: "4", args: args{opts: []varredis.Option{varredis.WithRaw(`{"addrs":["192.168.0.111:6379"]}`)}}, wantErr: false},
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
