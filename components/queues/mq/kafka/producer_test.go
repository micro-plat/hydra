package kafka

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/lib4go/assert"
)

func TestProducer(t *testing.T) {
	producer, err := NewProducerByRaw(`{"addrs":["xlh-kafka01:30021","xlh-kafka02:30022","xlh-kafka03:30023"],"write_timeout":300}`)
	assert.Equal(t, nil, err, err)
	err = producer.Push("kakfa-exporter", "hello world32")
	assert.Equal(t, nil, err, err)
}
func TestConsumer(t *testing.T) {
	consumer, err := NewConsumerByRaw(`{"addrs":["xlh-kafka01:30021","xlh-kafka02:30022","xlh-kafka03:30023"],"write_timeout":300,"group":"c4"}`)
	assert.Equal(t, nil, err, err)
	err = consumer.Connect()
	assert.Equal(t, nil, err, err)
	cc := 0
	err = consumer.Consume("kakfa-exporter", 1, func(data mq.IMQCMessage) {
		msg := data.GetMessage()
		cc++
		// data.Ack()

		fmt.Println("msg:", cc, string(msg))
	})

	assert.Equal(t, nil, err, err)
	time.Sleep(time.Second * 5)
	consumer.UnConsume("kakfa-exporter")
	assert.Equal(t, 1, 2)

}
