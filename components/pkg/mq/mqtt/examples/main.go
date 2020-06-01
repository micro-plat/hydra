package main

import (
	"fmt"
	"time"

	"github.com/micro-plat/lib4go/mq"
	"github.com/micro-plat/lib4go/mq/mqtt"
	pmqtt "github.com/micro-plat/lib4go/queue/mqtt"
)

func main() {

	consumer, err := mqtt.NewConsumer("", mq.WithRaw(`
	{
				"proto":"mqtt",
				"address":"222.209.84.37:11883",
				"userName":"mosquittouser",
				"password":"abc123$"
		}
	`))
	if err != nil {
		fmt.Println("consumer.err:", err)
		return
	}

	if err := consumer.Connect(); err != nil {
		fmt.Println("connect.err:", err)
		return
	}

	consumer.Consume("device.request1", 1, func(m mq.IMessage) {
		fmt.Println("recv:", m.GetMessage())
	})
	fmt.Println("success")

	time.Sleep(time.Second)
	publisher, err := pmqtt.New([]string{}, `{
				"proto":"mqtt",
				"address":"222.209.84.37:11883",
				"userName":"mosquittouser",
				"password":"abc123$"
		}`)
	if err != nil {
		fmt.Println("publisher.err:", err)
		return
	}
	for {
		fmt.Println("send.message:abc")
		err = publisher.Push("device.request1", `abc`)
		if err != nil {
			fmt.Println("push.err:", err)
		}
		time.Sleep(time.Second * 2)
	}

	time.Sleep(time.Hour)

}

// func main() {
// 	cc := client.New(&client.Options{
// 		ErrorHandler: func(err error) {
// 			fmt.Println("err:", err)
// 		},
// 	})

// 	if err := cc.Connect(&client.ConnectOptions{
// 		Network:   "tcp",
// 		Address:   "222.209.84.37:11883",
// 		UserName:  []byte("mosquittouser"),
// 		Password:  []byte("abc123$"),
// 		ClientID:  []byte("hydra-client5665"),
// 		KeepAlive: 3,
// 	}); err != nil {
// 		fmt.Println("err1:", err)
// 	}
// 	fmt.Println("conn.success")

// 	err := cc.Subscribe(&client.SubscribeOptions{
// 		SubReqs: []*client.SubReq{
// 			&client.SubReq{
// 				TopicFilter: []byte("device.request1"),
// 				QoS:         mqtt.QoS0,
// 				Handler: func(topicName, message []byte) {
// 					fmt.Println("recv:", string(topicName), string(message))
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	//----------------------publisher------------
// 	cc1 := client.New(&client.Options{
// 		ErrorHandler: func(err error) {
// 			fmt.Println("err:", err)
// 		},
// 	})

// 	if err := cc1.Connect(&client.ConnectOptions{
// 		Network: "tcp",
// 		Address: "222.209.84.37:11883",
// 		//UserName:  []byte("mosquittouser"),
// 		//Password:  []byte("abc123$"),
// 		ClientID:  []byte("hydra-client2222"),
// 		KeepAlive: 3,
// 	}); err != nil {
// 		fmt.Println("err1:", err)
// 	}

// 	for {
// 		err = cc1.Publish(&client.PublishOptions{
// 			QoS:       mqtt.QoS0,
// 			TopicName: []byte("device.request1"),
// 			Message:   []byte(`{"id":100}`),
// 		})
// 		if err != nil {
// 			fmt.Println("send.err:", err)
// 		}
// 		time.Sleep(time.Second)

// 	}

// 	time.Sleep(time.Second * 100000)
// 	fmt.Println("sub.success")

// }
