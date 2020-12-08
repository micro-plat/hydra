package global

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"
)

func TestIsLocal(t *testing.T) {
	type args struct {
		proto string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "1.本地-lm", args: args{proto: ProtoLM}, want: true},
		{name: "2.本地-lmq", args: args{proto: ProtoLMQ}, want: true},
		{name: "3.非本地-http", args: args{proto: ProtoHTTP}, want: false},
	}
	for _, tt := range tests {
		got := IsLocal(tt.args.proto)
		assert.Equal(t, tt.want, got, tt.name)
	}
}

func TestParseProto(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name      string
		args      args
		wantProto string
		wantAddr  string
		isNilErr  bool
	}{

		{name: "1. 正常的地址-单个地址", args: args{address: "zk://192.168.0.1"}, wantProto: "zk", wantAddr: "192.168.0.1", isNilErr: true},
		{name: "2. 正常的地址-多个地址", args: args{address: "zk://192.168.0.1,192.168.0.2"}, wantProto: "zk", wantAddr: "192.168.0.1,192.168.0.2", isNilErr: true},
		{name: "3. 错误地址-多个://", args: args{address: "zk://192.168.0.1://192.168.0.2"}, wantProto: "", wantAddr: "", isNilErr: false},
		{name: "4. 错误地址-无协议", args: args{address: "://192.168.0.1"}, wantProto: "", wantAddr: "", isNilErr: false},
		{name: "5. 错误地址-无地址", args: args{address: "zk://"}, wantProto: "", wantAddr: "", isNilErr: false},
	}

	for _, tt := range tests {
		gotproto, gotAddr, err := ParseProto(tt.args.address)

		assert.Equal(t, tt.wantProto, gotproto, tt.name)
		assert.Equal(t, tt.wantAddr, gotAddr, tt.name)
		assert.IsNil(t, tt.isNilErr, err, tt.name)
	}
}

func TestIsProto(t *testing.T) {
	type args struct {
		addr  string
		proto string
	}
	tests := []struct {
		name     string
		args     args
		wantAddr string
		wantIs   bool
	}{
		{name: "1. IsProto-匹配的proto", args: args{addr: "zk://192.168.0.1", proto: "zk"}, wantAddr: "192.168.0.1", wantIs: true},
		{name: "2. IsProto-不匹配的proto", args: args{addr: "zk://192.168.0.1", proto: "lm"}, wantAddr: "192.168.0.1", wantIs: false},
	}
	for _, tt := range tests {
		gotAddr, gotIs := IsProto(tt.args.addr, tt.args.proto)
		assert.Equal(t, tt.wantAddr, gotAddr, tt.name)
		assert.Equal(t, tt.wantIs, gotIs, tt.name)
	}
}
