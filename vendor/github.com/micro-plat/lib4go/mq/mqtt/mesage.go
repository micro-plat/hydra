package mqtt

//Message reids消息
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
func (m *Message) Write(b []byte) (int, error) {
	m.Message = string(b)
	m.HasData = len(b) > 0
	return len(b), nil
}

//NewMessage 创建消息
func NewMessage() *Message {
	return &Message{}
}
