package mqtt

import (
	"bytes"
	"io"
)

const (
	// Maximum payload size in bytes (256MiB - 1B).
	MaxPayloadSize = (1 << (4 * 7)) - 1
)

// Header contains the common attributes of all messages. Some attributes are
// not applicable to some message types.
type Header struct {
	DupFlag, Retain bool
	QosLevel        QosLevel
}

func (hdr *Header) Encode(w io.Writer, msgType MessageType, remainingLength int32) error {
	buf := new(bytes.Buffer)
	err := hdr.encodeInto(buf, msgType, remainingLength)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func (hdr *Header) encodeInto(buf *bytes.Buffer, msgType MessageType, remainingLength int32) error {
	if !hdr.QosLevel.IsValid() {
		return badQosError
	}
	if !msgType.IsValid() {
		return badMsgTypeError
	}

	val := byte(msgType) << 4
	val |= (boolToByte(hdr.DupFlag) << 3)
	val |= byte(hdr.QosLevel) << 1
	val |= boolToByte(hdr.Retain)
	buf.WriteByte(val)
	encodeLength(remainingLength, buf)
	return nil
}

func (hdr *Header) Decode(r io.Reader) (msgType MessageType, remainingLength int32, err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	var buf [1]byte

	if _, err = io.ReadFull(r, buf[:]); err != nil {
		return
	}

	byte1 := buf[0]
	msgType = MessageType(byte1 & 0xF0 >> 4)

	*hdr = Header{
		DupFlag:  byte1&0x08 > 0,
		QosLevel: QosLevel(byte1 & 0x06 >> 1),
		Retain:   byte1&0x01 > 0,
	}

	remainingLength = decodeLength(r)

	return
}

// Message is the interface that all MQTT messages implement.
type Message interface {
	// Encode writes the message to w.
	Encode(w io.Writer) error

	// Decode reads the message extended headers and payload from
	// r. Typically the values for hdr and packetRemaining will
	// be returned from Header.Decode.
	Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) error
}

// MessageType constants.
const (
	MsgConnect = MessageType(iota + 1)
	MsgConnAck
	MsgPublish
	MsgPubAck
	MsgPubRec
	MsgPubRel
	MsgPubComp
	MsgSubscribe
	MsgSubAck
	MsgUnsubscribe
	MsgUnsubAck
	MsgPingReq
	MsgPingResp
	MsgDisconnect

	msgTypeFirstInvalid
)

type MessageType uint8

// IsValid returns true if the MessageType value is valid.
func (mt MessageType) IsValid() bool {
	return mt >= MsgConnect && mt < msgTypeFirstInvalid
}

func writeMessage(w io.Writer, msgType MessageType, hdr *Header, payloadBuf *bytes.Buffer, extraLength int32) error {
	totalPayloadLength := int64(len(payloadBuf.Bytes())) + int64(extraLength)
	if totalPayloadLength > MaxPayloadSize {
		return msgTooLongError
	}

	buf := new(bytes.Buffer)
	err := hdr.encodeInto(buf, msgType, int32(totalPayloadLength))
	if err != nil {
		return err
	}

	buf.Write(payloadBuf.Bytes())
	_, err = w.Write(buf.Bytes())

	return err
}

// Connect represents an MQTT CONNECT message.
type Connect struct {
	Header
	ProtocolName               string
	ProtocolVersion            uint8
	WillRetain                 bool
	WillFlag                   bool
	CleanSession               bool
	WillQos                    QosLevel
	KeepAliveTimer             uint16
	ClientId                   string
	WillTopic, WillMessage     string
	UsernameFlag, PasswordFlag bool
	Username, Password         string
}

func (msg *Connect) Encode(w io.Writer) (err error) {
	if !msg.WillQos.IsValid() {
		return badWillQosError
	}

	buf := new(bytes.Buffer)

	flags := boolToByte(msg.UsernameFlag) << 7
	flags |= boolToByte(msg.PasswordFlag) << 6
	flags |= boolToByte(msg.WillRetain) << 5
	flags |= byte(msg.WillQos) << 3
	flags |= boolToByte(msg.WillFlag) << 2
	flags |= boolToByte(msg.CleanSession) << 1

	setString(msg.ProtocolName, buf)
	setUint8(msg.ProtocolVersion, buf)
	buf.WriteByte(flags)
	setUint16(msg.KeepAliveTimer, buf)
	setString(msg.ClientId, buf)
	if msg.WillFlag {
		setString(msg.WillTopic, buf)
		setString(msg.WillMessage, buf)
	}
	if msg.UsernameFlag {
		setString(msg.Username, buf)
	}
	if msg.PasswordFlag {
		setString(msg.Password, buf)
	}

	return writeMessage(w, MsgConnect, &msg.Header, buf, 0)
}

func (msg *Connect) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	msg.Header = hdr

	protocolName := getString(r, &packetRemaining)
	protocolVersion := getUint8(r, &packetRemaining)
	flags := getUint8(r, &packetRemaining)
	keepAliveTimer := getUint16(r, &packetRemaining)
	clientId := getString(r, &packetRemaining)

	*msg = Connect{
		ProtocolName:    protocolName,
		ProtocolVersion: protocolVersion,
		UsernameFlag:    flags&0x80 > 0,
		PasswordFlag:    flags&0x40 > 0,
		WillRetain:      flags&0x20 > 0,
		WillQos:         QosLevel(flags & 0x18 >> 3),
		WillFlag:        flags&0x04 > 0,
		CleanSession:    flags&0x02 > 0,
		KeepAliveTimer:  keepAliveTimer,
		ClientId:        clientId,
	}

	if msg.WillFlag {
		msg.WillTopic = getString(r, &packetRemaining)
		msg.WillMessage = getString(r, &packetRemaining)
	}
	if msg.UsernameFlag {
		msg.Username = getString(r, &packetRemaining)
	}
	if msg.PasswordFlag {
		msg.Password = getString(r, &packetRemaining)
	}

	if packetRemaining != 0 {
		return msgTooLongError
	}

	return nil
}

// ConnAck represents an MQTT CONNACK message.
type ConnAck struct {
	Header
	ReturnCode ReturnCode
}

func (msg *ConnAck) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)

	buf.WriteByte(byte(0)) // Reserved byte.
	setUint8(uint8(msg.ReturnCode), buf)

	return writeMessage(w, MsgConnAck, &msg.Header, buf, 0)
}

func (msg *ConnAck) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	msg.Header = hdr

	getUint8(r, &packetRemaining) // Skip reserved byte.
	msg.ReturnCode = ReturnCode(getUint8(r, &packetRemaining))
	if !msg.ReturnCode.IsValid() {
		return badReturnCodeError
	}

	if packetRemaining != 0 {
		return msgTooLongError
	}

	return nil
}

// Publish represents an MQTT PUBLISH message.
type Publish struct {
	Header
	TopicName string
	MessageId uint16
	Payload   Payload
}

func (msg *Publish) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)

	setString(msg.TopicName, buf)
	if msg.Header.QosLevel.HasId() {
		setUint16(msg.MessageId, buf)
	}

	if err = writeMessage(w, MsgPublish, &msg.Header, buf, int32(msg.Payload.Size())); err != nil {
		return
	}

	return msg.Payload.WritePayload(w)
}

func (msg *Publish) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	msg.Header = hdr

	msg.TopicName = getString(r, &packetRemaining)
	if msg.Header.QosLevel.HasId() {
		msg.MessageId = getUint16(r, &packetRemaining)
	}

	payloadReader := &io.LimitedReader{r, int64(packetRemaining)}

	if msg.Payload, err = config.MakePayload(msg, payloadReader, int(packetRemaining)); err != nil {
		return
	}

	return msg.Payload.ReadPayload(payloadReader)
}

// PubAck represents an MQTT PUBACK message.
type PubAck struct {
	Header
	MessageId uint16
}

func (msg *PubAck) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgPubAck)
}

func (msg *PubAck) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId, config)
}

// PubRec represents an MQTT PUBREC message.
type PubRec struct {
	Header
	MessageId uint16
}

func (msg *PubRec) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgPubRec)
}

func (msg *PubRec) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId, config)
}

// PubRel represents an MQTT PUBREL message.
type PubRel struct {
	Header
	MessageId uint16
}

func (msg *PubRel) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgPubRel)
}

func (msg *PubRel) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId, config)
}

// PubComp represents an MQTT PUBCOMP message.
type PubComp struct {
	Header
	MessageId uint16
}

func (msg *PubComp) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgPubComp)
}

func (msg *PubComp) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId, config)
}

// Subscribe represents an MQTT SUBSCRIBE message.
type Subscribe struct {
	Header
	MessageId uint16
	Topics    []TopicQos
}

type TopicQos struct {
	Topic string
	Qos   QosLevel
}

func (msg *Subscribe) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)
	if msg.Header.QosLevel.HasId() {
		setUint16(msg.MessageId, buf)
	}
	for _, topicSub := range msg.Topics {
		setString(topicSub.Topic, buf)
		setUint8(uint8(topicSub.Qos), buf)
	}

	return writeMessage(w, MsgSubscribe, &msg.Header, buf, 0)
}

func (msg *Subscribe) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	msg.Header = hdr

	if msg.Header.QosLevel.HasId() {
		msg.MessageId = getUint16(r, &packetRemaining)
	}
	var topics []TopicQos
	for packetRemaining > 0 {
		topics = append(topics, TopicQos{
			Topic: getString(r, &packetRemaining),
			Qos:   QosLevel(getUint8(r, &packetRemaining)),
		})
	}
	msg.Topics = topics

	return nil
}

// SubAck represents an MQTT SUBACK message.
type SubAck struct {
	Header
	MessageId uint16
	TopicsQos []QosLevel
}

func (msg *SubAck) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)
	setUint16(msg.MessageId, buf)
	for i := 0; i < len(msg.TopicsQos); i += 1 {
		setUint8(uint8(msg.TopicsQos[i]), buf)
	}

	return writeMessage(w, MsgSubAck, &msg.Header, buf, 0)
}

func (msg *SubAck) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	msg.Header = hdr

	msg.MessageId = getUint16(r, &packetRemaining)
	topicsQos := make([]QosLevel, 0)
	for packetRemaining > 0 {
		grantedQos := QosLevel(getUint8(r, &packetRemaining) & 0x03)
		topicsQos = append(topicsQos, grantedQos)
	}
	msg.TopicsQos = topicsQos

	return nil
}

// Unsubscribe represents an MQTT UNSUBSCRIBE message.
type Unsubscribe struct {
	Header
	MessageId uint16
	Topics    []string
}

func (msg *Unsubscribe) Encode(w io.Writer) (err error) {
	buf := new(bytes.Buffer)
	if msg.Header.QosLevel.HasId() {
		setUint16(msg.MessageId, buf)
	}
	for _, topic := range msg.Topics {
		setString(topic, buf)
	}

	return writeMessage(w, MsgUnsubscribe, &msg.Header, buf, 0)
}

func (msg *Unsubscribe) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	msg.Header = hdr

	if qos := msg.Header.QosLevel; qos == 1 || qos == 2 {
		msg.MessageId = getUint16(r, &packetRemaining)
	}
	topics := make([]string, 0)
	for packetRemaining > 0 {
		topics = append(topics, getString(r, &packetRemaining))
	}
	msg.Topics = topics

	return nil
}

// UnsubAck represents an MQTT UNSUBACK message.
type UnsubAck struct {
	Header
	MessageId uint16
}

func (msg *UnsubAck) Encode(w io.Writer) error {
	return encodeAckCommon(w, &msg.Header, msg.MessageId, MsgUnsubAck)
}

func (msg *UnsubAck) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) (err error) {
	msg.Header = hdr
	return decodeAckCommon(r, packetRemaining, &msg.MessageId, config)
}

// PingReq represents an MQTT PINGREQ message.
type PingReq struct {
	Header
}

func (msg *PingReq) Encode(w io.Writer) error {
	return msg.Header.Encode(w, MsgPingReq, 0)
}

func (msg *PingReq) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) error {
	if packetRemaining != 0 {
		return msgTooLongError
	}
	return nil
}

// PingResp represents an MQTT PINGRESP message.
type PingResp struct {
	Header
}

func (msg *PingResp) Encode(w io.Writer) error {
	return msg.Header.Encode(w, MsgPingResp, 0)
}

func (msg *PingResp) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) error {
	if packetRemaining != 0 {
		return msgTooLongError
	}
	return nil
}

// Disconnect represents an MQTT DISCONNECT message.
type Disconnect struct {
	Header
}

func (msg *Disconnect) Encode(w io.Writer) error {
	return msg.Header.Encode(w, MsgDisconnect, 0)
}

func (msg *Disconnect) Decode(r io.Reader, hdr Header, packetRemaining int32, config DecoderConfig) error {
	if packetRemaining != 0 {
		return msgTooLongError
	}
	return nil
}

func encodeAckCommon(w io.Writer, hdr *Header, messageId uint16, msgType MessageType) error {
	buf := new(bytes.Buffer)
	setUint16(messageId, buf)
	return writeMessage(w, msgType, hdr, buf, 0)
}

func decodeAckCommon(r io.Reader, packetRemaining int32, messageId *uint16, config DecoderConfig) (err error) {
	defer func() {
		err = recoverError(err, recover())
	}()

	*messageId = getUint16(r, &packetRemaining)

	if packetRemaining != 0 {
		return msgTooLongError
	}

	return nil
}
