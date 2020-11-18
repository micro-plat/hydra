package mqc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/micro-plat/hydra/components/queues/mq/redis"
	"github.com/micro-plat/hydra/conf/server/queue"
	"github.com/micro-plat/hydra/test/assert"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name      string
		queueName string
		service   string
		message   string
		hasData   bool
		errStr    string
	}{
		{name: "队列数据格式错误", queueName: "queue", service: "service", message: `message`, hasData: true, errStr: "队列queue中存放的数据不是有效的json:message invalid character 'm' looking for beginning of value"},
		{name: "添加队列数据", queueName: "queue", service: "service", message: `{"key":"value","__header__":{"Content-Type":"application/json"}}`, hasData: true},
	}
	for _, tt := range tests {
		gotR, err := NewRequest(queue.NewQueue(tt.queueName, tt.service), &redis.RedisMessage{Message: tt.message, HasData: tt.hasData})
		if tt.errStr != "" {
			assert.Equal(t, tt.errStr, err.Error(), tt.name)
			continue
		}
		assert.Equal(t, tt.queueName, gotR.GetName(), tt.name)
		assert.Equal(t, tt.service, gotR.GetService(), tt.name)
		assert.Equal(t, DefMethod, gotR.GetMethod(), tt.name)
		f := make(map[string]interface{})
		json.Unmarshal([]byte(tt.message), &f)
		header := make(map[string]string)
		for n, m := range f["__header__"].(map[string]interface{}) {
			header[n] = fmt.Sprint(m)
		}
		header["Client-IP"] = "127.0.0.1"
		assert.Equal(t, header, gotR.GetHeader(), tt.name)
		f["__body_"] = tt.message
		assert.Equal(t, f, gotR.GetForm(), tt.name)
	}
}
