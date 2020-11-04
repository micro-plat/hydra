package creator

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/server/rpc"
	"github.com/micro-plat/hydra/services"
)

func Test_newRPC(t *testing.T) {
	type args struct {
		address     string
		fnGetRouter func(string) *services.ORouter
		opts        []rpc.Option
	}
	tests := []struct {
		name string
		args args
		want *rpcBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRPC(tt.args.address, tt.args.fnGetRouter, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRPC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rpcBuilder_Load(t *testing.T) {
	type fields struct {
		httpBuilder *httpBuilder
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &rpcBuilder{
				httpBuilder: tt.fields.httpBuilder,
			}
			b.Load()
		})
	}
}
