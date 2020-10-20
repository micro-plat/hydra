package global

import (
	"testing"

	"github.com/micro-plat/hydra/test/assert"
)

func TestGetHostPort(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name     string
		args     args
		wantHost string
		wantPort string
		IsNilErr bool
	}{
		{
			name: "正常的IP+端口",
			args: args{
				addr: "192.168.0.1:9090",
			},
			wantHost: "192.168.0.1",
			wantPort: "9090",
			IsNilErr: true,
		},
		{
			name: "无IP+端口",
			args: args{
				addr: ":9090",
			},
			wantHost: "0.0.0.0",
			wantPort: "9090",
			IsNilErr: true,
		},
		{
			name: "只有端口",
			args: args{
				addr: "8080",
			},
			wantHost: "0.0.0.0",
			wantPort: "8080",
			IsNilErr: true,
		},
		{
			name: "80端口",
			args: args{
				addr: "80",
			},
			wantHost: "",
			wantPort: "",
			IsNilErr: false,
		},
		{
			name: "错误的端口",
			args: args{
				addr: ":aa",
			},
			wantHost: "",
			wantPort: "",
			IsNilErr: false,
		},
	}
	for _, tt := range tests {
		gotHost, gotPort, err := GetHostPort(tt.args.addr)

		//t.Log(gotHost, gotPort, err, tt.IsNilErr)

		assert.Equal(t, tt.wantHost, gotHost, tt.name)
		assert.Equal(t, tt.wantPort, gotPort, tt.name)
		assert.IsNil(t, tt.IsNilErr, err, tt.name)
	}
}

func TestWithIPMask(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name     string
		args     args
		wantMask string
	}{
		{
			name: "设置mask",
			args: args{
				val: "192.168",
			},
			wantMask: "192.168",
		},
		{
			name: "设置mask-空",
			args: args{
				val: "",
			},
			wantMask: "",
		},
	}
	for _, tt := range tests {
		WithIPMask(tt.args.val)
		assert.Equal(t, tt.wantMask, mask, tt.name)
	}
}
