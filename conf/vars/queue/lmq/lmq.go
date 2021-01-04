package lmq

import "github.com/micro-plat/hydra/conf/vars/queue"

//LMQ 本地队列
type LMQ = queue.Queue

//New 构建mqtt配置
func New() *LMQ {
	return &LMQ{Proto: "lmq"}
}

//MQ LMQ地址
const MQ = "lmq://."
