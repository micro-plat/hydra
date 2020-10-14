package global

import (
	"testing"
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
		{
			name: "本地-lm",
			args: args{proto: ProtoLM},
			want: true,
		},
		{
			name: "本地-lmq",
			args: args{proto: ProtoLMQ},
			want: true,
		},
		{
			name: "非本地-http",
			args: args{proto: ProtoHTTP},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLocal(tt.args.proto); got != tt.want {
				t.Errorf("IsLocal() = %v, want %v", got, tt.want)
			}
		})
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
		wantErr   bool
	}{

		{
			name: "正常的地址-单个地址",
			args: args{
				address: "zk://192.168.0.1",
			},
			wantProto: "zk",
			wantAddr:  "192.168.0.1",
			wantErr:   false,
		},
		{
			name: "正常的地址-多个地址",
			args: args{
				address: "zk://192.168.0.1,192.168.0.2",
			},
			wantProto: "zk",
			wantAddr:  "192.168.0.1,192.168.0.2",
			wantErr:   false,
		},
		{
			name: "错误地址-多个://",
			args: args{
				address: "zk://192.168.0.1://192.168.0.2",
			},
			wantProto: "",
			wantAddr:  "",
			wantErr:   true,
		},
		{
			name: "错误地址-无协议",
			args: args{
				address: "://192.168.0.1",
			},
			wantProto: "",
			wantAddr:  "",
			wantErr:   true,
		},
		{
			name: "错误地址-无地址",
			args: args{
				address: "zk://",
			},
			wantProto: "",
			wantAddr:  "",
			wantErr:   true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotproto, gotAddr, err := ParseProto(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotproto != tt.wantProto {
				t.Errorf("ParseProto() gotproto = %v, wantProto %v", gotproto, tt.wantProto)
			}
			if gotAddr != tt.wantAddr {
				t.Errorf("ParseProto() gotAddr = %v, wantAddr %v", gotAddr, tt.wantAddr)
			}
		})
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
		{
			name: "匹配的proto",
			args: args{
				addr:  "zk://192.168.0.1",
				proto: "zk",
			},
			wantAddr: "192.168.0.1",
			wantIs:   true,
		},
		{
			name: "不匹配的proto",
			args: args{
				addr:  "zk://192.168.0.1",
				proto: "lm",
			},
			wantAddr: "192.168.0.1",
			wantIs:   false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddr, gotIs := IsProto(tt.args.addr, tt.args.proto)
			if gotAddr != tt.wantAddr {
				t.Errorf("IsProto() gotAddr = %v, want %v", gotAddr, tt.wantAddr)
			}
			if gotIs != tt.wantIs {
				t.Errorf("IsProto() gotIs = %v, want %v", gotIs, tt.wantIs)
			}
		})
	}
}
