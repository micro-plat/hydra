package mqc

import (
	"testing"
	"time"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/hydra/test/assert"
	"github.com/micro-plat/lib4go/concurrent/cmap"
)

func TestNewProcessor(t *testing.T) {
	tests := []struct {
		name    string
		proto   string
		confRaw string
		wantP   *Processor
		wantErr string
	}{
		{name: "协议错误", proto: "proto", confRaw: "{}", wantErr: "构建mqc服务失败(proto:proto,raw:{}) mqc: 未知的协议类型 proto"},
		{name: "协议配置正确", proto: "redis", confRaw: `{"proto":"redis","addrs":["192.168.5.79:6379"]}`, wantP: &Processor{status: unstarted,
			closeChan: make(chan struct{}),
			startTime: time.Now(),
			queues:    cmap.New(4)}},
	}
	for _, tt := range tests {
		gotP, err := NewProcessor(tt.proto, tt.confRaw)
		if tt.wantErr != "" {
			assert.Equal(t, tt.wantErr, err.Error(), tt.name)
			continue
		}
		tt.wantP.customer, _ = mq.NewMQC(tt.proto, tt.confRaw)
		assert.Equal(t, tt.wantP, gotP, tt.name)

		//
		assert.Equal(t, 4, len(gotP.Engine.Handlers), tt.name)
	}

}
