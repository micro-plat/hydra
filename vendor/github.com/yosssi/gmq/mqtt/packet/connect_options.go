package packet

import (
	"errors"

	"github.com/yosssi/gmq/mqtt"
)

// Error values
var (
	ErrClientIDExceedsMaxStringsLen    = errors.New("the length of the Client Identifier exceeds the maximum strings length")
	ErrUserNameExceedsMaxStringsLen    = errors.New("the length of the User Name exceeds the maximum strings length")
	ErrPasswordExceedsMaxStringsLen    = errors.New("the length of the Password exceeds the maximum strings length")
	ErrWillTopicExceedsMaxStringsLen   = errors.New("the length of the Will Topic exceeds the maximum strings length")
	ErrWillMessageExceedsMaxStringsLen = errors.New("the length of the Will Message exceeds the maximum strings length")
	ErrInvalidClientIDCleanSession     = errors.New("the Clean Session must be true if the Client Identifier is zero-byte")
	ErrInvalidClientIDPassword         = errors.New("the Password must be zero-byte if the Client Identifier is zero-byte")
	ErrInvalidWillTopicMessage         = errors.New("the Will Topic (Message) must not be zero-byte if the Will Message (Topic) is not zero-byte")
	ErrInvalidWillQoS                  = errors.New("the Will QoS is invalid")
	ErrInvalidWillTopicMessageQoS      = errors.New("the Will QoS must be zero if both the Will Topic and the Will Message are zero-byte")
	ErrInvalidWillTopicMessageRetain   = errors.New("the Will Retain must be false if both the Will Topic and the Will Message are zero-byte")
)

// CONNECTOptions represents options for a CONNECT Packet.
type CONNECTOptions struct {
	// ClientID is the Client Identifier of the payload.
	ClientID []byte
	// UserName is the User Name of the payload.
	UserName []byte
	// Password is the Password of the payload.
	Password []byte
	// CleanSession is the Clean Session of the variable header.
	CleanSession bool
	// KeepAlive is the Keep Alive of the variable header.
	KeepAlive uint16
	// WillTopic is the Will Topic of the payload.
	WillTopic []byte
	// WillMessage is the Will Message of the payload.
	WillMessage []byte
	// WillQoS is the Will QoS of the variable header.
	WillQoS byte
	// WillRetain is the Will Retain of the variable header.
	WillRetain bool
}

func (opts *CONNECTOptions) validate() error {
	// Check the length of the Client Identifier.
	if len(opts.ClientID) > maxStringsLen {
		return ErrClientIDExceedsMaxStringsLen
	}

	// Check the combination of the Client Identifier and the Clean Session.
	if len(opts.ClientID) == 0 && !opts.CleanSession {
		return ErrInvalidClientIDCleanSession
	}

	// Check the length of the User Name.
	if len(opts.UserName) > maxStringsLen {
		return ErrUserNameExceedsMaxStringsLen
	}

	// Check the length of the Password.
	if len(opts.Password) > maxStringsLen {
		return ErrPasswordExceedsMaxStringsLen
	}

	// Check the combination of the Client Identifier and the Password.
	if len(opts.UserName) == 0 && len(opts.Password) > 0 {
		return ErrInvalidClientIDPassword
	}

	// Check the length of the Will Topic.
	if len(opts.WillTopic) > maxStringsLen {
		return ErrWillTopicExceedsMaxStringsLen
	}

	// Check the length of the Will Message.
	if len(opts.WillMessage) > maxStringsLen {
		return ErrWillMessageExceedsMaxStringsLen
	}

	// Check the combination of the Will Topic and the Will Message.
	if (len(opts.WillTopic) > 0 && len(opts.WillMessage) == 0) || (len(opts.WillTopic) == 0 && len(opts.WillMessage) > 0) {
		return ErrInvalidWillTopicMessage
	}

	// Check the Will QoS.
	if !mqtt.ValidQoS(opts.WillQoS) {
		return ErrInvalidWillQoS
	}

	// Check the combination of the Will Topic, the Will Message and the Will QoS.
	if len(opts.WillTopic) == 0 && len(opts.WillMessage) == 0 && opts.WillQoS != mqtt.QoS0 {
		return ErrInvalidWillTopicMessageQoS
	}

	// Check the combination of the Will Topic, the Will Message and the Will Retain.
	if len(opts.WillTopic) == 0 && len(opts.WillMessage) == 0 && opts.WillRetain {
		return ErrInvalidWillTopicMessageRetain
	}

	return nil
}
