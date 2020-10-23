package conf

import (
	"testing"

	"github.com/micro-plat/hydra/conf/vars/queue"
	"github.com/micro-plat/hydra/conf/vars/queue/mqtt"
	"github.com/micro-plat/hydra/test/assert"
)

func TestQueueMQTTNew(t *testing.T) {

	tests := []struct {
		name    string
		address string
		opts    []mqtt.Option
		want    *mqtt.MQTT
	}{
		{
			name:    "测试新增-无option",
			address: "192.168.5.79:6379",
			want: &mqtt.MQTT{
				Queue: &queue.Queue{
					Proto: "mqtt",
				},
				Address: "192.168.5.79:6379",
			},
		},
		{
			name:    "测试新增-WithUP",
			address: "192.168.5.79:6379",
			opts: []mqtt.Option{
				mqtt.WithUP("name", "pwd"),
			},
			want: &mqtt.MQTT{
				Queue: &queue.Queue{
					Proto: "mqtt",
				},
				Address:  "192.168.5.79:6379",
				UserName: "name",
				Password: "pwd",
			},
		},
		{
			name:    "测试新增-WithUP",
			address: "192.168.5.79:6379",
			opts: []mqtt.Option{
				mqtt.WithUP("name", "pwd"),
				mqtt.WithCert("cert"),
			},
			want: &mqtt.MQTT{
				Queue: &queue.Queue{
					Proto: "mqtt",
				},
				Address:  "192.168.5.79:6379",
				UserName: "name",
				Password: "pwd",
				Cert:     "cert",
			},
		},
	}
	for _, tt := range tests {
		got := mqtt.New(tt.address, tt.opts...)
		assert.Equal(t, tt.want, got, tt.name)
	}
}
