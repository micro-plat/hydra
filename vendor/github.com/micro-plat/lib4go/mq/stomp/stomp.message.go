package stomp

import s "github.com/zkfy/stompngo"

//StompMessage stomp消息
type StompMessage struct {
	s       *StompConsumer
	msg     *s.Message
	Message string
}

//Ack 确定消息
func (m *StompMessage) Ack() error {
	return m.s.conn.Ack(m.msg.Headers)
}

//Nack 取消消息
func (m *StompMessage) Nack() error {
	h := s.Headers{"message-id", m.msg.Headers.Value("message-id"),
		"subscription", m.msg.Headers.Value("subscription"),
		"id", m.msg.Headers.Value("message-id")}
	return m.s.conn.Nack(h)
}

//GetMessage 获取消息
func (m *StompMessage) GetMessage() string {
	return m.Message
}

//Has 是否报含数据
func (m *StompMessage) Has() bool {
	return len(m.msg.Headers) > 0 && m.msg.Headers.Value("connection_read_error") != "EOF"
}

//NewStompMessage 创建消息
func NewStompMessage(s *StompConsumer, msg *s.Message) *StompMessage {
	return &StompMessage{s: s, msg: msg, Message: string(msg.Body)}
}
