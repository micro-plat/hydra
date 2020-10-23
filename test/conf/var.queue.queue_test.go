package conf

import (
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/vars/queue"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/hydra/test/mocks"
)

func TestQueueNew(t *testing.T) {

	tests := []struct {
		name  string
		proto string
		raw   []byte
		want  *queue.Queue
	}{
		{
			name:  "测试新增",
			proto: "redis",
			raw:   []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
			want: &queue.Queue{
				Proto: "redis",
				Raw:   []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
			},
		},
	}
	for _, tt := range tests {
		got := queue.New(tt.proto, tt.raw)
		assert.Equal(t, tt.want.Proto, got.Proto, tt.name)
		assert.Equal(t, tt.want.Raw, got.Raw, tt.name)

	}
}

func TestQueueGetConf(t *testing.T) {
	type args struct {
		cnfData []byte
		version int32
		tp      string
		name    string
	}
	tests := []struct {
		name     string
		args     args
		want     *queue.Queue
		IsNilErr bool
	}{
		{
			name: "测试-var中无该配置",
			args: args{
				cnfData: []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
				tp:      "queue",
				name:    "xredis",
			},
			want:     nil,
			IsNilErr: false,
		},
		{
			name: "测试-var中有配置",
			args: args{
				cnfData: []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
				tp:      "queue",
				name:    "redis",
			},
			want: &queue.Queue{
				Proto: "redis",
				Raw:   []byte(`{"proto":"redis","addrs":["192.168.5.79:6379"]}`),
			},
			IsNilErr: true,
		},

		// TODO: Add test cases.
	}

	for _, tt := range tests {
		rawCnf, err := conf.NewRawConfByJson(tt.args.cnfData, tt.args.version)
		fmt.Println(tt.name)
		cnf := &mocks.MockVarConf{
			Version: tt.args.version,
			ConfData: map[string]map[string]*conf.RawConf{
				"queue": map[string]*conf.RawConf{
					"redis": rawCnf,
				},
			},
		}

		got, err := queue.GetConf(cnf, tt.args.tp, tt.args.name)
		//fmt.Println("queue.GetConf:", got, err)
		assert.IsNil(t, tt.IsNilErr, err, tt.name)

		assert.Equal(t, tt.want, got, tt.name)
	}
}
