// Implementation of MQTT V3.1 encoding and decoding.
//
// See http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html
// for the MQTT protocol specification. This package does not implement the
// semantics of MQTT, but purely the encoding and decoding of its messages.
//
// Decoding Messages:
//
// Use the DecodeOneMessage function to read a Message from an io.Reader, it
// will return a Message value. The function can be implemented using the public
// API of this package if more control is required. For example:
//
//   for {
//     msg, err := mqtt.DecodeOneMessage(conn, nil)
//     if err != nil {
//       // handle err
//     }
//     switch msg := msg.(type) {
//     case *Connect:
//       // ...
//     case *Publish:
//       // ...
//       // etc.
//     }
//   }
//
// Encoding Messages:
//
// Create a message value, and use its Encode method to write it to an
// io.Writer. For example:
//
//   someData := []byte{1, 2, 3}
//   msg := &Publish{
//     Header: {
//       DupFlag: false,
//       QosLevel: QosAtLeastOnce,
//       Retain: false,
//     },
//     TopicName: "a/b",
//     MessageId: 10,
//     Payload: BytesPayload(someData),
//   }
//   if err := msg.Encode(conn); err != nil {
//     // handle err
//   }
//
// Advanced PUBLISH payload handling:
//
// The default behaviour for decoding PUBLISH payloads, and most common way to
// supply payloads for encoding, is the BytesPayload, which is a []byte
// derivative.
//
// More complex handling is possible by implementing the Payload interface,
// which can be injected into DecodeOneMessage via the `config` parameter, or
// into an outgoing Publish message via its Payload field.  Potential benefits
// of this include:
//
// * Data can be (un)marshalled directly on a connection, without an unecessary
// round-trip via bytes.Buffer.
//
// * Data can be streamed directly on readers/writers (e.g files, other
// connections, pipes) without the requirement to buffer an entire message
// payload in memory at once.
//
// The limitations of these streaming features are:
//
// * When encoding a payload, the encoded size of the payload must be known and
// declared upfront.
//
// * The payload size (and PUBLISH variable header) can be no more than 256MiB
// minus 1 byte. This is a specified limitation of MQTT v3.1 itself.
package mqtt

import (
	"errors"
	"io"
)

var (
	badMsgTypeError        = errors.New("mqtt: message type is invalid")
	badQosError            = errors.New("mqtt: QoS is invalid")
	badWillQosError        = errors.New("mqtt: will QoS is invalid")
	badLengthEncodingError = errors.New("mqtt: remaining length field exceeded maximum of 4 bytes")
	badReturnCodeError     = errors.New("mqtt: is invalid")
	dataExceedsPacketError = errors.New("mqtt: data exceeds packet length")
	msgTooLongError        = errors.New("mqtt: message is too long")
)

const (
	QosAtMostOnce = QosLevel(iota)
	QosAtLeastOnce
	QosExactlyOnce

	qosFirstInvalid
)

type QosLevel uint8

func (qos QosLevel) IsValid() bool {
	return qos < qosFirstInvalid
}

func (qos QosLevel) HasId() bool {
	return qos == QosAtLeastOnce || qos == QosExactlyOnce
}

const (
	RetCodeAccepted = ReturnCode(iota)
	RetCodeUnacceptableProtocolVersion
	RetCodeIdentifierRejected
	RetCodeServerUnavailable
	RetCodeBadUsernameOrPassword
	RetCodeNotAuthorized

	retCodeFirstInvalid
)

type ReturnCode uint8

func (rc ReturnCode) IsValid() bool {
	return rc >= RetCodeAccepted && rc < retCodeFirstInvalid
}

// DecoderConfig provides configuration for decoding messages.
type DecoderConfig interface {
	// MakePayload returns a Payload for the given Publish message. r is a Reader
	// that will read the payload data, and n is the number of bytes in the
	// payload. The Payload.ReadPayload method is called on the returned payload
	// by the decoding process.
	MakePayload(msg *Publish, r io.Reader, n int) (Payload, error)
}

type DefaultDecoderConfig struct{}

func (c DefaultDecoderConfig) MakePayload(msg *Publish, r io.Reader, n int) (Payload, error) {
	return make(BytesPayload, n), nil
}

// ValueConfig always returns the given Payload when MakePayload is called.
type ValueConfig struct {
	Payload Payload
}

func (c *ValueConfig) MakePayload(msg *Publish, r io.Reader, n int) (Payload, error) {
	return c.Payload, nil
}

// DecodeOneMessage decodes one message from r. config provides specifics on
// how to decode messages, nil indicates that the DefaultDecoderConfig should
// be used.
func DecodeOneMessage(r io.Reader, config DecoderConfig) (msg Message, err error) {
	var hdr Header
	var msgType MessageType
	var packetRemaining int32
	msgType, packetRemaining, err = hdr.Decode(r)
	if err != nil {
		return
	}

	msg, err = NewMessage(msgType)
	if err != nil {
		return
	}

	if config == nil {
		config = DefaultDecoderConfig{}
	}

	return msg, msg.Decode(r, hdr, packetRemaining, config)
}

// NewMessage creates an instance of a Message value for the given message
// type. An error is returned if msgType is invalid.
func NewMessage(msgType MessageType) (msg Message, err error) {
	switch msgType {
	case MsgConnect:
		msg = new(Connect)
	case MsgConnAck:
		msg = new(ConnAck)
	case MsgPublish:
		msg = new(Publish)
	case MsgPubAck:
		msg = new(PubAck)
	case MsgPubRec:
		msg = new(PubRec)
	case MsgPubRel:
		msg = new(PubRel)
	case MsgPubComp:
		msg = new(PubComp)
	case MsgSubscribe:
		msg = new(Subscribe)
	case MsgUnsubAck:
		msg = new(UnsubAck)
	case MsgSubAck:
		msg = new(SubAck)
	case MsgUnsubscribe:
		msg = new(Unsubscribe)
	case MsgPingReq:
		msg = new(PingReq)
	case MsgPingResp:
		msg = new(PingResp)
	case MsgDisconnect:
		msg = new(Disconnect)
	default:
		return nil, badMsgTypeError
	}

	return
}

// panicErr wraps an error that caused a problem that needs to bail out of the
// API, such that errors can be recovered and returned as errors from the
// public API.
type panicErr struct {
	err error
}

func (p panicErr) Error() string {
	return p.err.Error()
}

func raiseError(err error) {
	panic(panicErr{err})
}

// recoverError recovers any panic in flight and, iff it's an error from
// raiseError, will return the error. Otherwise re-raises the panic value.
// If no panic is in flight, it returns existingErr.
//
// This must be used in combination with a defer in all public API entry
// points where raiseError could be called.
func recoverError(existingErr error, recovered interface{}) error {
	if recovered != nil {
		if pErr, ok := recovered.(panicErr); ok {
			return pErr.err
		} else {
			panic(recovered)
		}
	}
	return existingErr
}
