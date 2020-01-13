package lmq


//LMQMessage reids消息
type LMQMessage struct {
	Message string
	HasData bool
}

//Ack 确定消息
func (m *LMQMessage) Ack() error {
	return nil
}

//Nack 取消消息
func (m *LMQMessage) Nack() error {
	return nil
}

//GetMessage 获取消息
func (m *LMQMessage) GetMessage() string {
	return m.Message
}

//Has 是否有数据
func (m *LMQMessage) Has() bool {
	return m.HasData
}

//NewLMQMessage 创建消息
func NewLMQMessage(msg string) *LMQMessage {
	return &LMQMessage{Message: msg, HasData: len(msg) > 0}
}
