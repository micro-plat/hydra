package packet

import (
	"bytes"
	"errors"

	"github.com/yosssi/gmq/mqtt"
)

// Wildcard characters
const (
	wildcardMulti  = "#"
	wildcardSingle = "+"
	wildcards      = wildcardMulti + wildcardSingle
)

// Error values
var (
	ErrInvalidQoS                    = errors.New("the QoS is invalid")
	ErrTopicNameExceedsMaxStringsLen = errors.New("the length of the Topic Name exceeds the maximum strings length")
	ErrTopicNameContainsWildcards    = errors.New("the Topic Name contains wildcard characters")
	ErrMessageExceedsMaxStringsLen   = errors.New("the length of the Message exceeds the maximum strings length")
)

// PUBLISHOptions represents options for a PUBLISH Packet.
type PUBLISHOptions struct {
	// DUP is the DUP flag of the fixed header.
	DUP bool
	// QoS is the QoS of the fixed header.
	QoS byte
	// Retain is the Retain of the fixed header.
	Retain bool
	// TopicName is the Topic Name of the varible header.
	TopicName []byte
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
	// Message is the Application Message of the payload.
	Message []byte
}

// validate validates the options.
func (opts *PUBLISHOptions) validate() error {
	// Check the QoS.
	if !mqtt.ValidQoS(opts.QoS) {
		return ErrInvalidQoS
	}

	// Check the length of the Topic Name.
	if len(opts.TopicName) > maxStringsLen {
		return ErrTopicNameExceedsMaxStringsLen
	}

	// Check if the Topic Name contains the wildcard characters.
	if bytes.IndexAny(opts.TopicName, wildcards) != -1 {
		return ErrTopicNameContainsWildcards
	}

	// Check the length of the Application Message.
	if len(opts.Message) > maxStringsLen {
		return ErrMessageExceedsMaxStringsLen
	}

	// End the validation if the QoS equals to QoS 0.
	if opts.QoS == mqtt.QoS0 {
		return nil
	}

	// Check the Packet Identifier.
	if opts.PacketID == 0 {
		return ErrInvalidPacketID
	}

	return nil
}
