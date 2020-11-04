package creator

import (
	"reflect"
	"testing"

	"github.com/micro-plat/hydra/conf/server/mqc"
	"github.com/micro-plat/hydra/conf/server/queue"
)

func Test_newMQC(t *testing.T) {
	type args struct {
		addr string
		opts []mqc.Option
	}
	tests := []struct {
		name string
		args args
		want *mqcBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMQC(tt.args.addr, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMQC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mqcBuilder_Load(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &mqcBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
			}
			b.Load()
		})
	}
}

func Test_mqcBuilder_Queue(t *testing.T) {
	type fields struct {
		CustomerBuilder CustomerBuilder
	}
	type args struct {
		mq []*queue.Queue
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mqcBuilder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &mqcBuilder{
				CustomerBuilder: tt.fields.CustomerBuilder,
			}
			if got := b.Queue(tt.args.mq...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mqcBuilder.Queue() = %v, want %v", got, tt.want)
			}
		})
	}
}
