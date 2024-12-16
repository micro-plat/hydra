package kafka

import (
	"sync"

	"github.com/IBM/sarama"
	"github.com/micro-plat/hydra/components/queues/mq"
	"github.com/micro-plat/lib4go/logger"
)

// KafkaMessage reids消息
type KafkaMessage struct {
	Message *sarama.ConsumerMessage
	session sarama.ConsumerGroupSession
}

// NewKafkaMessage 创建消息
func NewKafkaMessage(msg *sarama.ConsumerMessage, session sarama.ConsumerGroupSession) *KafkaMessage {
	return &KafkaMessage{Message: msg, session: session}
}

// Ack 确定消息
func (m *KafkaMessage) Ack() error {
	m.session.MarkMessage(m.Message, "")
	return nil
}

// Nack 取消消息
func (m *KafkaMessage) Nack() error {
	return nil
}

// GetMessage 获取消息
func (m *KafkaMessage) GetMessage() string {
	return string(m.Message.Value)
}

// Has 是否有数据
func (m *KafkaMessage) Has() bool {
	return len(string(m.Message.Value)) > 0
}

type ConsumeHandler struct {
	ready        chan bool
	log          logger.ILogger
	handle       func(mq.IMQCMessage)
	nconcurrency int
}

func NewConsumeHandler(log logger.ILogger, handle func(mq.IMQCMessage), nconcurrency int) *ConsumeHandler {
	return &ConsumeHandler{
		log:          log,
		handle:       handle,
		ready:        make(chan bool),
		nconcurrency: nconcurrency,
	}

}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ConsumeHandler) Setup(s sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumeHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a ConsumeHandler loop of ConsumeHandlerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (h *ConsumeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	wg := &sync.WaitGroup{}
	wg.Add(h.nconcurrency)
	for i := 0; i < h.nconcurrency; i++ {
		go func() {
			defer wg.Done()
		START:
			for {
				select {
				case message, ok := <-claim.Messages():
					if !ok {
						break START
					}
					go h.handle(NewKafkaMessage(message, session))
				case <-session.Context().Done():
					break START

				}
			}

		}()
	}
	wg.Wait()
	return nil
}
