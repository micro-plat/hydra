package skywalking

import (
	"testing"

	"github.com/micro-plat/hydra/context/apm"
)

func TestNewTracer(t *testing.T) {
	type args struct {
		service string
		opts    []apm.TracerOption
	}
	c, _ := New("instance", `{"check_interval":1,"max_send_queue_size":500000,"instance_props":{"":""},"authentication_key":""}`)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "1", args: args{service: "service1", opts: []apm.TracerOption{
			WithReporter(c.reporter), WithInstance(c.instance),
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTracer, err := NewTracer(tt.args.service, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTracer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, ok := gotTracer.(apm.Tracer); !ok {
				t.Error("NewTracer() doesn't return an apm.Tracer")
			}
		})
	}
}
