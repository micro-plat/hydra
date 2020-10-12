package skywalking

import (
	"reflect"
	"testing"
	"time"
	"unsafe"

	grpcreporter "github.com/SkyAPM/go2sky/reporter"
	conf "github.com/micro-plat/hydra/conf/vars/apm"
	"github.com/micro-plat/hydra/context/apm"
)

func Test_buildOptions(t *testing.T) {
	type args struct {
		config *conf.APM
	}
	tests := []struct {
		name     string
		args     args
		wantOpts []grpcreporter.GRPCReporterOption
		wantErr  bool
	}{
		{name: "1", args: args{config: &conf.APM{}}, wantOpts: []grpcreporter.GRPCReporterOption{}, wantErr: false},
		{name: "2", args: args{config: &conf.APM{ReportCheckInterval: 1}}, wantOpts: []grpcreporter.GRPCReporterOption{
			grpcreporter.WithCheckInterval(time.Second * time.Duration(1)),
		}, wantErr: false},
		{name: "3", args: args{config: &conf.APM{
			InstanceProps: map[string]string{"a": "b"},
		}}, wantOpts: []grpcreporter.GRPCReporterOption{grpcreporter.WithInstanceProps(map[string]string{"a": "b"})}, wantErr: false},
		{name: "4", args: args{config: &conf.APM{
			MaxSendQueueSize: 500000,
		}}, wantOpts: []grpcreporter.GRPCReporterOption{grpcreporter.WithMaxSendQueueSize(500000)}, wantErr: false},
		{name: "5", args: args{config: &conf.APM{
			Credentials: &conf.Credential{CertFile: "filename", ServerName: "servername"}, //@todo 正确的文件地址
		}}, wantOpts: []grpcreporter.GRPCReporterOption{}, wantErr: true},
		{name: "6", args: args{config: &conf.APM{
			AuthenticationKey: "key",
		}}, wantOpts: []grpcreporter.GRPCReporterOption{grpcreporter.WithAuthentication("key")}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOpts, err := buildOptions(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sh := (*reflect.SliceHeader)(unsafe.Pointer(&gotOpts))
			sh1 := (*reflect.SliceHeader)(unsafe.Pointer(&tt.wantOpts))
			// fmt.Printf("n1 Data:Ox%x,Len:%d,Cap:%d\n", sh.Data, sh.Len, sh.Cap)
			// fmt.Printf("n2 Data:Ox%x,Len:%d,Cap:%d\n", sh1.Data, sh1.Len, sh1.Cap)
			if sh1.Len != sh.Len && sh1.Cap != sh.Cap {
				t.Errorf("NewReporter() Data:Ox%x,Len:%d,Cap:%d\n", sh.Data, sh.Len, sh.Cap)
				t.Errorf("want Data:Ox%x,Len:%d,Cap:%d\n", sh1.Data, sh1.Len, sh1.Cap)
			}
		})
	}
}

func TestNewReporter(t *testing.T) {
	type args struct {
		serverAddr string
		config     *conf.APM
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{serverAddr: "serverAddr", config: &conf.APM{}}, wantErr: false},
		{name: "2", args: args{serverAddr: "serverAddr", config: &conf.APM{
			ReportCheckInterval: 1,
			InstanceProps:       map[string]string{"a": "b"},
			MaxSendQueueSize:    500000,
			AuthenticationKey:   "key",
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReporter, err := NewReporter(tt.args.serverAddr, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReporter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, ok := gotReporter.(apm.Reporter); !ok {
				t.Error("NewReporter() didn't return apm.Reporter")
			}
		})
	}
}
