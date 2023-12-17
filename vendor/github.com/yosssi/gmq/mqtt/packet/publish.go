package packet

import (
	"errors"

	"github.com/yosssi/gmq/mqtt"
)

// Maximum Remaining Length
const maxRemainingLength = 268435455

// Minimum length of the fixed header of the PUBLISH Packet
const minLenPUBLISHFixedHeader = 2

// Minimum length of the variable header of the PUBLISH Packet
const minLenPUBLISHVariableHeader = 2

// Error value
var ErrInvalidPacketID = errors.New("invalid Packet Identifier")

// PUBLISH represents a PUBLISH Packet.
type PUBLISH struct {
	base
	// dup is the DUP flag of the fixed header.
	DUP bool
	// qos is the QoS of the fixed header.
	QoS byte
	// retain is the Retain of the fixed header.
	retain bool
	// topicName is the Topic Name of the varible header.
	TopicName []byte
	// packetID is the Packet Identifier of the variable header.
	PacketID uint16
	// message is the Application Message of the payload.
	Message []byte
}

// setFixedHeader sets the fixed header to the Packet.
func (p *PUBLISH) setFixedHeader() {
	// Define the first byte of the fixed header.
	b := TypePUBLISH << 4

	// Set 1 to the Bit 3 if the DUP flag is true.
	if p.DUP {
		b |= 0x08
	}

	// Set the value of the Will QoS to the Bit 2 and 1.
	b |= p.QoS << 1

	// Set 1 to the Bit 0 if the Retain is true.
	if p.retain {
		b |= 0x01
	}

	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, b)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *PUBLISH) setVariableHeader() {
	// Append the Topic Name to the variable header.
	p.variableHeader = appendLenStr(p.variableHeader, p.TopicName)

	if p.QoS != mqtt.QoS0 {
		// Append the Packet Identifier to the variable header.
		p.variableHeader = append(p.variableHeader, encodeUint16(p.PacketID)...)
	}
}

// setPayload sets the payload to the Packet.
func (p *PUBLISH) setPayload() {
	p.payload = p.Message
}

// NewPUBLISH creates and returns a PUBLISH Packet.
func NewPUBLISH(opts *PUBLISHOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &PUBLISHOptions{}
	}

	// Validate the options.
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Create a PUBLISH Packet.
	p := &PUBLISH{
		DUP:       opts.DUP,
		QoS:       opts.QoS,
		retain:    opts.Retain,
		TopicName: opts.TopicName,
		PacketID:  opts.PacketID,
		Message:   opts.Message,
	}

	// Set the variable header to the Packet.
	p.setVariableHeader()

	// Set the payload to the Packet.
	p.setPayload()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p, nil
}

// NewPUBLISHFromBytes creates the PUBLISH Packet
// from the byte data and returns it.
func NewPUBLISHFromBytes(fixedHeader FixedHeader, remaining []byte) (Packet, error) {
	// Validate the byte data.
	if err := validatePUBLISHBytes(fixedHeader, remaining); err != nil {
		return nil, err
	}

	// Get the first byte from the fixedHeader.
	b := fixedHeader[0]

	// Create a PUBLISH Packet.
	p := &PUBLISH{
		DUP:    b&0x08 == 0x08,
		QoS:    b & 0x06 >> 1,
		retain: b&0x01 == 0x01,
	}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Extract the length of the Topic Name.
	lenTopicName, _ := decodeUint16(remaining[0:2])

	// Calculate the length of the variable header.
	var lenVariableHeader int

	if p.QoS == mqtt.QoS0 {
		lenVariableHeader = 2 + int(lenTopicName)
	} else {
		lenVariableHeader = 2 + int(lenTopicName) + 2
	}

	// Set the variable header to the Packet.
	p.variableHeader = remaining[:lenVariableHeader]

	// Set the payload to the Packet.
	p.payload = remaining[lenVariableHeader:]

	// Set the Topic Name to the Packet.
	p.TopicName = remaining[2 : 2+lenTopicName]

	// Extract the Packet Identifier.
	var packetID uint16

	if p.QoS != mqtt.QoS0 {
		packetID, _ = decodeUint16(remaining[2+lenTopicName : 2+lenTopicName+2])
	}

	// Set the Packet Identifier to the Packet.
	p.PacketID = packetID

	// Set the Application Message to the Packet.
	p.Message = p.payload

	// Return the Packet.
	return p, nil
}

// validatePUBLISHBytes validates the fixed header and the variable header.
func validatePUBLISHBytes(fixedHeader FixedHeader, remaining []byte) error {
	// Extract the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	// Check the length of the fixed header.
	if len(fixedHeader) < minLenPUBLISHFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	if ptype != TypePUBLISH {
		return ErrInvalidPacketType
	}

	// Get the QoS.
	qos := (fixedHeader[0] & 0x06) >> 1

	// Check the QoS.
	if !mqtt.ValidQoS(qos) {
		return ErrInvalidQoS
	}

	// Check the length of the remaining.
	if l := len(remaining); l < minLenPUBLISHVariableHeader || l > maxRemainingLength {
		return ErrInvalidRemainingLen
	}

	// Extract the length of the Topic Name.
	lenTopicName, _ := decodeUint16(remaining[0:2])

	// Calculate the length of the variable header.
	var lenVariableHeader int

	if qos == mqtt.QoS0 {
		lenVariableHeader = 2 + int(lenTopicName)
	} else {
		lenVariableHeader = 2 + int(lenTopicName) + 2
	}

	// Check the length of the remaining.
	if len(remaining) < lenVariableHeader {
		return ErrInvalidRemainingLength
	}

	// End the validation if the QoS equals to QoS 0.
	if qos == mqtt.QoS0 {
		return nil
	}

	// Extract the Packet Identifier.
	packetID, _ := decodeUint16(remaining[2+int(lenTopicName) : 2+int(lenTopicName)+2])

	// Check the Packet Identifier.
	if packetID == 0 {
		return ErrInvalidPacketID
	}

	return nil
}
