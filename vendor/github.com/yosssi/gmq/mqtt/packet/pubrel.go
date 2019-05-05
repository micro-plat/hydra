package packet

// Length of the fixed header of the PUBREL Packet
const lenPUBRELFixedHeader = 2

// Length of the variable header of the PUBREL Packet
const lenPUBRELVariableHeader = 2

// PUBREL represents a PUBREL Packet.
type PUBREL struct {
	base
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// setFixedHeader sets the fixed header to the Packet.
func (p *PUBREL) setFixedHeader() {
	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, TypePUBREL<<4|0x02)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *PUBREL) setVariableHeader() {
	// Append the Packet Identifier to the variable header.
	p.variableHeader = append(p.variableHeader, encodeUint16(p.PacketID)...)
}

// NewPUBREL creates and returns a PUBREL Packet.
func NewPUBREL(opts *PUBRELOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &PUBRELOptions{}
	}

	// Validate the options.
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Create a PUBREL Packet.
	p := &PUBREL{
		PacketID: opts.PacketID,
	}

	// Set the variable header to the Packet.
	p.setVariableHeader()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p, nil
}

// NewPUBRELFromBytes creates a PUBREL Packet
// from the byte data and returns it.
func NewPUBRELFromBytes(fixedHeader FixedHeader, variableHeader []byte) (Packet, error) {
	// Validate the byte data.
	if err := validatePUBRELBytes(fixedHeader, variableHeader); err != nil {
		return nil, err
	}

	// Decode the Packet Identifier.
	// No error occur because of the precedent validation and
	// the returned error is not be taken care of.
	packetID, _ := decodeUint16(variableHeader)

	// Create a PUBREL Packet.
	p := &PUBREL{
		PacketID: packetID,
	}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Set the variable header to the Packet.
	p.variableHeader = variableHeader

	// Return the Packet.
	return p, nil
}

// validatePUBRELBytes validates the fixed header and the variable header.
func validatePUBRELBytes(fixedHeader FixedHeader, variableHeader []byte) error {
	// Extract the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	// Check the length of the fixed header.
	if len(fixedHeader) != lenPUBRELFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	if ptype != TypePUBREL {
		return ErrInvalidPacketType
	}

	// Check the reserved bits of the fixed header.
	if fixedHeader[0]&0x02 != 0x02 {
		return ErrInvalidFixedHeader
	}

	// Check the Remaining Length of the fixed header.
	if fixedHeader[1] != lenPUBRELVariableHeader {
		return ErrInvalidRemainingLength
	}

	// Check the length of the variable header.
	if len(variableHeader) != lenPUBRELVariableHeader {
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
