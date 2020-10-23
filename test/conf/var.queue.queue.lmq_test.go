package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/queue/lmq"
	"github.com/micro-plat/hydra/test/assert"
)

func TestQueueLMQNew(t *testing.T) {
	tests := []struct {
		name string
		want *lmq.LMQ
	}{
		{
			name: "测试新增 ",
			want: &lmq.LMQ{
				Proto: "lmq",
				Raw:   nil,
			},
		},
	}
	for _, tt := range tests {
		got := lmq.New()
		assert.Equal(t, tt.want, got, tt.name)
	}
}
