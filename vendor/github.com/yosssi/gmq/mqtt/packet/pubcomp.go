package packet

// Length of the fixed header of the PUBCOMP Packet
const lenPUBCOMPFixedHeader = 2

// Length of the variable header of the PUBCOMP Packet
const lenPUBCOMPVariableHeader = 2

// PUBCOMP represents a PUBCOMP Packet.
type PUBCOMP struct {
	base
	// packetID is the Packet Identifier of the variable header.
	PacketID uint16
}

// setFixedHeader sets the fixed header to the Packet.
func (p *PUBCOMP) setFixedHeader() {
	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, TypePUBCOMP<<4)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *PUBCOMP) setVariableHeader() {
	// Append the Packet Identifier to the variable header.
	p.variableHeader = append(p.variableHeader, encodeUint16(p.PacketID)...)
}

// NewPUBCOMP creates and returns a PUBCOMP Packet.
func NewPUBCOMP(opts *PUBCOMPOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &PUBCOMPOptions{}
	}

	// Validate the options.
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Create a PUBCOMP Packet.
	p := &PUBCOMP{
		PacketID: opts.PacketID,
	}

	// Set the variable header to the Packet.
	p.setVariableHeader()

	// Set the Fixed header to the Packet.
	p.setFixedHeader()

	// Return the Packet.
	return p, nil
}

// NewPUBCOMPFromBytes creates a PUBCOMP Packet
// from the byte data and returns it.
func NewPUBCOMPFromBytes(fixedHeader FixedHeader, variableHeader []byte) (Packet, error) {
	// Validate the byte data.
	if err := validatePUBCOMPBytes(fixedHeader, variableHeader); err != nil {
		return nil, err
	}

	// Decode the Packet Identifier.
	// No error occur because of the precedent validation and
	// the returned error is not be taken care of.
	packetID, _ := decodeUint16(variableHeader)

	// Create a PUBCOMP Packet.
	p := &PUBCOMP{
		PacketID: packetID,
	}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Set the variable header to the Packet.
	p.variableHeader = variableHeader

	// Return the Packet.
	return p, nil
}

// validatePUBCOMPBytes validates the fixed header and the variable header.
func validatePUBCOMPBytes(fixedHeader FixedHeader, variableHeader []byte) error {
	// Extract the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	// Check the length of the fixed header.
	if len(fixedHeader) != lenPUBCOMPFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	if ptype != TypePUBCOMP {
		return ErrInvalidPacketType
	}

	// Check the reserved bits of the fixed header.
	if fixedHeader[0]<<4 != 0x00 {
		return ErrInvalidFixedHeader
	}

	// Check the Remaining Length of the fixed header.
	if fixedHeader[1] != lenPUBCOMPVariableHeader {
		return ErrInvalidRemainingLength
	}

	// Check the length of the variable header.
	if len(variableHeader) != lenPUBCOMPVariableHeader {
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
