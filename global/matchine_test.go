package global

import (
	"testing"

	"github.com/micro-plat/lib4go/assert"
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
		{name: "1.正常的IP+端口", args: args{addr: "192.168.0.1:9090"}, wantHost: "192.168.0.1", wantPort: "9090", IsNilErr: true},
		{name: "2.无IP+端口", args: args{addr: ":9090"}, wantHost: "0.0.0.0", wantPort: "9090", IsNilErr: true},
		{name: "3.只有端口", args: args{addr: "8080"}, wantHost: "0.0.0.0", wantPort: "8080", IsNilErr: true},
		{name: "4.80端口", args: args{addr: "80"}, wantHost: "", wantPort: "", IsNilErr: false},
		{name: "5.错误的端口", args: args{addr: ":aa"}, wantHost: "", wantPort: "", IsNilErr: false},
	}
	for _, tt := range tests {
		gotHost, gotPort, err := GetHostPort(tt.args.addr)
		if tt.IsNilErr {
			assert.Equal(t, tt.wantHost, gotHost, tt.name)
			assert.Equal(t, tt.wantPort, gotPort, tt.name)
			continue
		}
		assert.NotNil(t, err, tt.name)
	}
}
