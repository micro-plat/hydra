package xmq

import (
	"bytes"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"strings"

	"github.com/micro-plat/lib4go/encoding"
	"github.com/micro-plat/lib4go/jsons"
	"github.com/micro-plat/lib4go/security/md5"
)

var xmqSEQId int64 = 10000

//Message 消息体
type Message struct {
	CMD       int      `json:"cmd"`  //0发送
	Mode      int      `json:"mod"`  //模式是否需要回复:0:需要 1:不需要
	Timeout   int      `json:"time"` //超时时长
	QueueName string   `json:"name"` //队列名称
	SEQ       int64    `json:"seq"`  //请滶序列号
	Timestmap int64    `json:"ts"`   //时间戳，当前秒数
	Data      []string `json:"data"` //消息数据
	Sign      string   `json:"sign"` //签名串
	signKey   string
}

const defaultSignKey = "sdfsdfsdfdsf"

//newHeartBit 构建消息信息
func newHeartBit() *Message {

	r := &Message{
		CMD:       99,
		Mode:      1,
		Timestmap: time.Now().Unix(),
		signKey:   defaultSignKey,
	}
	id := atomic.AddInt64(&xmqSEQId, 1)
	r.SEQ, _ = strconv.ParseInt(fmt.Sprintf("%d", id), 10, 0)
	return r
}

//newMessage 构建消息信息
func newMessage(queueName string, msg string, timeout int) *Message {

	r := &Message{
		CMD:       0,
		Mode:      1,
		QueueName: queueName,
		Data:      []string{msg},
		Timeout:   timeout * 1000,
		Timestmap: time.Now().Unix(),
		signKey:   defaultSignKey,
	}
	id := atomic.AddInt64(&xmqSEQId, 1)
	r.SEQ, _ = strconv.ParseInt(fmt.Sprintf("%d", id), 10, 0)
	return r
}

//Make 构建消息
func (x *Message) Make() (string, error) {
	buff := &bytes.Buffer{}
	buff.WriteString(strconv.Itoa(x.CMD))
	buff.WriteString(fmt.Sprint(x.SEQ))
	buff.WriteString(fmt.Sprint(x.Timestmap))
	buff.WriteString(fmt.Sprint(x.signKey))
	gbkValue, err := encoding.DecodeBytes(buff.Bytes(), "gbk")
	if err != nil {
		return "", err
	}
	x.Sign = strings.ToUpper(md5.EncryptBytes(gbkValue))
	r, err := jsons.Marshal(x)
	if err != nil {
		return "", err
	}
	return string(r) + "\n", nil
}
