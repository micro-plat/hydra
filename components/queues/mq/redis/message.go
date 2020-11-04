package redis

import (
	"github.com/go-redis/redis"
)

//RedisMessage reids消息
type RedisMessage struct {
	Message string
	HasData bool
}

//Ack 确定消息
func (m *RedisMessage) Ack() error {
	return nil
}

//Nack 取消消息
func (m *RedisMessage) Nack() error {
	return nil
}

//GetMessage 获取消息
func (m *RedisMessage) GetMessage() string {
	return m.Message
}

//Has 是否有数据
func (m *RedisMessage) Has() bool {
	return m.HasData
}

//NewRedisMessage 创建消息
func NewRedisMessage(cmd *redis.StringSliceCmd) *RedisMessage {
	msg, err := cmd.Result()
	hasData := err == nil && len(msg) > 0
	ndata := ""
	if hasData {
		ndata = msg[len(msg)-1]
	}
	return &RedisMessage{Message: ndata, HasData: hasData}
}
