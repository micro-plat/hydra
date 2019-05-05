package packet

import (
	"errors"

	"github.com/yosssi/gmq/mqtt"
)

// Error values
var (
	ErrNoTopicFilter                   = errors.New("the Topic Filter must be specified")
	ErrTopicFilterExceedsMaxStringsLen = errors.New("the length of the Topic Filter exceeds the maximum strings length")
)

// SubReq represents subscription request.
type SubReq struct {
	// TopicFilter is the Topic Filter of the Subscription.
	TopicFilter []byte
	// QoS is the requsting QoS.
	QoS byte
}

// validate validates the subscription request.
func (s *SubReq) validate() error {
	// Check the length of the Topic Filter.
	l := len(s.TopicFilter)

	if l == 0 {
		return ErrNoTopicFilter
	}

	if l > maxStringsLen {
		return ErrTopicFilterExceedsMaxStringsLen
	}

	// Check the QoS.
	if !mqtt.ValidQoS(s.QoS) {
		return ErrInvalidQoS
	}

	return nil
}
