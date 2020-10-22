package pub

import (
	"sync"
	"testing"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/lib4go/logger"
)

func TestPublisher_WatchClusterChange(t *testing.T) {
	type fields struct {
		c          conf.IMainConf
		log        logger.ILogging
		serverNode string
		serverName string
		lock       sync.Mutex
		closeChan  chan struct{}
		watchChan  chan struct{}
		pubs       map[string]string
		done       bool
	}
	type args struct {
		notify func(isMaster bool, sharding int, total int)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Publisher{
				c:          tt.fields.c,
				log:        tt.fields.log,
				serverNode: tt.fields.serverNode,
				serverName: tt.fields.serverName,
				lock:       tt.fields.lock,
				closeChan:  tt.fields.closeChan,
				watchChan:  tt.fields.watchChan,
				pubs:       tt.fields.pubs,
				done:       tt.fields.done,
			}
			if err := p.WatchClusterChange(tt.args.notify); (err != nil) != tt.wantErr {
				t.Errorf("Publisher.WatchClusterChange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
