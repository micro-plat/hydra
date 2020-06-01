package stomp

import (
	"testing"

	"github.com/micro-plat/lib4go/ut"
)

var consumerQueue = "queue_test"
var consumerMsg = "msg_test"
var consumerTimeOut = 10

// TestNewStompConsumer 测试创建一个消费者对象
func TestNewStompConsumer(t *testing.T) {
	consumer, err := NewStompConsumer(address)
	ut.Expect(t, err, nil)
	ut.Refute(t, consumer, nil)

}

// TestConsumerConnect 测试消费者对象连接到服务器
func TestConsumerConnect(t *testing.T) {
	// 正常连接到服务器
	consumer, err := NewStompConsumer(address)
	ut.Expect(t, err, nil)
	err = consumer.ConnectOnce()
	ut.Expect(t, err, nil)

	// 端口错误
	addr := "192.168.0.165:80"
	consumer, err = NewStompConsumer(addr)
	ut.Expect(t, err, nil)

	err = consumer.ConnectOnce()
	ut.Refute(t, err, nil)

	// ip地址格式错误
	addr = "168.165:61613"
	consumer, err = NewStompConsumer(addr)
	ut.Expect(t, err, nil)
	err = consumer.ConnectOnce()
	ut.Refute(t, err, nil)
}
