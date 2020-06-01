package stomp

import "testing"
import "github.com/micro-plat/lib4go/ut"
import "time"

var address = "192.168.0.165:61613"

func TestStompProducer1(t *testing.T) {
	producer, err := NewStompProducer(address)
	ut.Expect(t, err, nil)
	err = producer.connectOnce()
	ut.Expect(t, err, nil)
}

func TestStompProducer2(t *testing.T) {
	producer, err := NewStompProducer(address)
	ut.Expect(t, err, nil)

	err = producer.connectOnce()
	ut.Expect(t, err, nil)

	err = producer.Send("hydra_test01", "hello", time.Second*60)
	ut.Expect(t, err, nil)
	time.Sleep(time.Millisecond * 100)
	ut.Expect(t, len(producer.backupMsg), 0)
	ut.Expect(t, len(producer.messages), 0)
}

func TestStompProducer3(t *testing.T) {
	addr := "192.168.0.165:61612"
	producer, err := NewStompProducer(addr)
	ut.Expect(t, err, nil)

	err = producer.connectOnce()
	ut.Refute(t, err, nil)

	err = producer.Send("hydra_test01", "hello", time.Second*60)
	ut.Expect(t, err, nil)
	time.Sleep(time.Millisecond * 100)
	ut.Expect(t, len(producer.backupMsg), 0)
	ut.Expect(t, len(producer.messages), 1)

	//重新连接
	producer.address = address
	err = producer.connectOnce()
	ut.Expect(t, err, nil)
	time.Sleep(time.Millisecond * 100)
	ut.Expect(t, len(producer.backupMsg), 0)
	ut.Expect(t, len(producer.messages), 0)
}
func TestStompProducer4(t *testing.T) {
	producer, err := NewStompProducer(address)
	ut.Expect(t, err, nil)

	err = producer.connectOnce()
	ut.Expect(t, err, nil)

	producer.Close()
	err = producer.Send("hydra_test01", "hello", time.Second*60)
	ut.Refute(t, err, nil)
}
