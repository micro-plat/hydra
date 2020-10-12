package skywalking

import (
	"testing"

	ctxapm "github.com/micro-plat/hydra/context/apm"
)

func TestNew(t *testing.T) {
	raw := `{"check_interval":1,"max_send_queue_size":500000,"instance_props":{"":""},"authentication_key":""}`
	instance := "test"

	gotM, err := New(instance, raw)
	if (err != nil) != false {
		t.Errorf("New() error = %v, wantErr %v", err, false)
		return
	}
	if _, ok := gotM.reporter.(ctxapm.Reporter); !ok || gotM.instance != instance {
		t.Error("New() didn't return right")
	}
}
