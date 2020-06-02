package lmq

//Message 消息信息
type Message struct {
	Message string
	HasData bool
}

//Ack 确定消息
func (m *Message) Ack() error {
	return nil
}

//Nack 取消消息
func (m *Message) Nack() error {
	return nil
}

//GetMessage 获取消息
func (m *Message) GetMessage() string {
	return m.Message
}

//Has 是否有数据
func (m *Message) Has() bool {
	return m.HasData
}

//newMessage 创建消息
func newMessage(msg string) *Message {
	return &Message{Message: msg, HasData: len(msg) > 0}
}
