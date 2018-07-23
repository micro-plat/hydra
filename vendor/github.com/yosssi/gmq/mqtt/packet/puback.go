package packet

// Length of the fixed header of the PUBACK Packet
const lenPUBACKFixedHeader = 2

// Length of the variable header of the PUBACK Packet
const lenPUBACKVariableHeader = 2

// PUBACK represents a PUBACK Packet.
type PUBACK struct {
	base
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// setFixedHeader sets the fixed header to the Packet.
func (p *PUBACK) setFixedHeader() {
	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, TypePUBACK<<4)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *PUBACK) setVariableHeader() {
	// Append the Packet Identifier to the variable header.
	p.variableHeader = append(p.variableHeader, encodeUint16(p.PacketID)...)
}

// NewPUBACK creates and returns a PUBACK Packet.
func NewPUBACK(opts *PUBACKOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &PUBACKOptions{}
	}

	// Validate the options.
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Create a PUBACK Packet.
	p := &PUBACK{
		PacketID: opts.PacketID,
	}

	// Set the variable header to the Packet.
	p.setVariableHeader()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p, nil
}

// NewPUBACKFromBytes creates a PUBACK Packet
// from the byte data and returns it.
func NewPUBACKFromBytes(fixedHeader FixedHeader, variableHeader []byte) (Packet, error) {
	// Validate the byte data.
	if err := validatePUBACKBytes(fixedHeader, variableHeader); err != nil {
		return nil, err
	}

	// Decode the Packet Identifier.
	// No error occur because of the precedent validation and
	// the returned error is not be taken care of.
	packetID, _ := decodeUint16(variableHeader)

	// Create a PUBACK Packet.
	p := &PUBACK{
		PacketID: packetID,
	}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Set the variable header to the Packet.
	p.variableHeader = variableHeader

	// Return the Packet.
	return p, nil
}

// validatePUBACKBytes validates the fixed header and the variable header.
func validatePUBACKBytes(fixedHeader FixedHeader, variableHeader []byte) error {
	// Extract the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	// Check the length of the fixed header.
	if len(fixedHeader) != lenPUBACKFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	if ptype != TypePUBACK {
		return ErrInvalidPacketType
	}

	// Check the reserved bits of the fixed header.
	if fixedHeader[0]<<4 != 0x00 {
		return ErrInvalidFixedHeader
	}

	// Check the Remaining Length of the fixed header.
	if fixedHeader[1] != lenPUBACKVariableHeader {
		return ErrInvalidRemainingLength
	}

	// Check the length of the variable header.
	if len(variableHeader) != lenPUBACKVariableHeader {
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
