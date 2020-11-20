package cron

import (
	"testing"

	"github.com/micro-plat/hydra/conf/server/task"
	"github.com/micro-plat/hydra/hydra/servers/cron"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		tasks   []*task.Task
		wantT   *cron.Server
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotT, err := cron.NewServer(tt.tasks...)
		assert.Equalf(t, tt.wantErr, err == nil, tt.name, err)
		assert.Equalf(t, tt.wantT.Processor, gotT.Processor, tt.name, "Processor")
	}
}
