package packet

import (
	"errors"

	"github.com/yosssi/gmq/mqtt"
)

// Minimum length of the fixed header of the SUBACK Packet
const minLenSUBACKFixedHeader = 2

// Length of the variable header of the SUBACK Packet
const lenSUBACKVariableHeader = 2

// Return Code Failure
const SUBACKRetFailure byte = 0x80

// Error value
var ErrInvalidSUBACKReturnCode = errors.New("invalid SUBACK Return Code")

// SUBACK represents a SUBACK Packet.
type SUBACK struct {
	base
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
	// ReturnCodes is the Return Codes of the payload.
	ReturnCodes []byte
}

// NewSUBACKFromBytes creates a SUBACK Packet
// from the byte data and returns it.
func NewSUBACKFromBytes(fixedHeader FixedHeader, remaining []byte) (Packet, error) {
	// Validate the byte data.
	if err := validateSUBACKBytes(fixedHeader, remaining); err != nil {
		return nil, err
	}

	// Extract the variable header.
	variableHeader := remaining[0:lenSUBACKVariableHeader]

	// Extract the payload.
	payload := remaining[lenSUBACKVariableHeader:]

	// Decode the Packet Identifier.
	// No error occur because of the precedent validation and
	// the returned error is not be taken care of.
	packetID, _ := decodeUint16(variableHeader)

	// Create a PUBACK Packet.
	p := &SUBACK{
		PacketID:    packetID,
		ReturnCodes: payload,
	}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Set the variable header to the Packet.
	p.variableHeader = variableHeader

	// Set the payload to the Packet.
	p.payload = payload

	// Return the Packet.
	return p, nil
}

// validateSUBACKBytes validates the fixed header and the remaining.
func validateSUBACKBytes(fixedHeader FixedHeader, remaining []byte) error {
	// Extract the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	// Check the length of the fixed header.
	if len(fixedHeader) < minLenSUBACKFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	if ptype != TypeSUBACK {
		return ErrInvalidPacketType
	}

	// Check the reserved bits of the fixed header.
	if fixedHeader[0]<<4 != 0x00 {
		return ErrInvalidFixedHeader
	}

	// Check the length of the remaining.
	if len(remaining) < lenSUBACKVariableHeader+1 {
		return ErrInvalidRemainingLen
	}

	// Extract the Packet Identifier.
	packetID, _ := decodeUint16(remaining[0:lenSUBACKVariableHeader])

	// Check the Packet Identifier.
	if packetID == 0 {
		return ErrInvalidPacketID
	}

	// Check each Return Code.
	for _, b := range remaining[lenSUBACKVariableHeader:] {
		if !mqtt.ValidQoS(b) && b != SUBACKRetFailure {
			return ErrInvalidSUBACKReturnCode
		}
	}

	return nil
}
