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

type XMQMessage struct {
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

//NewXMQHeartBit 构建消息信息
func NewXMQHeartBit() *XMQMessage {

	r := &XMQMessage{
		CMD:       99,
		Mode:      1,
		Timestmap: time.Now().Unix(),
		signKey:   "sdfsdfsdfdsf",
	}
	id := atomic.AddInt64(&xmqSEQId, 1)
	r.SEQ, _ = strconv.ParseInt(fmt.Sprintf("%d", id), 10, 0)
	//r.SEQ = r.Timestmap
	return r
}

//NewXMQMessage 构建消息信息
func NewXMQMessage(queueName string, msg string, timeout int) *XMQMessage {

	r := &XMQMessage{
		CMD:       0,
		Mode:      1,
		QueueName: queueName,
		Data:      []string{msg},
		Timeout:   timeout * 1000,
		Timestmap: time.Now().Unix(),
		signKey:   "sdfsdfsdfdsf",
	}
	id := atomic.AddInt64(&xmqSEQId, 1)
	r.SEQ, _ = strconv.ParseInt(fmt.Sprintf("%d", id), 10, 0)
	// r.SEQ = r.Timestmap
	return r
}

//MakeMessage 构建消息
func (x *XMQMessage) MakeMessage() (string, error) {
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
