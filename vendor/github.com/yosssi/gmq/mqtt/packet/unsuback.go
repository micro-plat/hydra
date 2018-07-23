package packet

// Length of the fixed header of the UNSUBACK Packet
const lenUNSUBACKFixedHeader = 2

// Length of the variable header of the UNSUBACK Packet
const lenUNSUBACKVariableHeader = 2

// UNSUBACK represents an UNSUBACK Packet.
type UNSUBACK struct {
	base
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// NewUNSUBACKFromBytes creates an UNSUBACK Packet
// from the byte data and returns it.
func NewUNSUBACKFromBytes(fixedHeader FixedHeader, variableHeader []byte) (Packet, error) {
	// Validate the byte data.
	if err := validateUNSUBACKBytes(fixedHeader, variableHeader); err != nil {
		return nil, err
	}

	// Decode the Packet Identifier.
	// No error occur because of the precedent validation and
	// the returned error is not be taken care of.
	packetID, _ := decodeUint16(variableHeader)

	// Create a PUBACK Packet.
	p := &UNSUBACK{
		PacketID: packetID,
	}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Set the variable header to the Packet.
	p.variableHeader = variableHeader

	// Return the Packet.
	return p, nil
}

// validateUNSUBACKBytes validates the fixed header and the variable header.
func validateUNSUBACKBytes(fixedHeader FixedHeader, variableHeader []byte) error {
	// Extract the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	// Check the length of the fixed header.
	if len(fixedHeader) != lenUNSUBACKFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	if ptype != TypeUNSUBACK {
		return ErrInvalidPacketType
	}

	// Check the reserved bits of the fixed header.
	if fixedHeader[0]<<4 != 0x00 {
		return ErrInvalidFixedHeader
	}

	// Check the Remaining Length of the fixed header.
	if fixedHeader[1] != lenUNSUBACKVariableHeader {
		return ErrInvalidRemainingLength
	}

	// Check the length of the variable header.
	if len(variableHeader) != lenUNSUBACKVariableHeader {
		return ErrInvalidVariableHeaderLen
	}

	// Extract the Packet Identifier.
	packetID, _ := decodeUint16(variableHeader)

	// Check the Packet Identifier.
	if packetID == 0 {
		return ErrInvalidPacketID
	}

	return nil
}
