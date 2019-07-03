package packet

// UNSUBSCRIBE represents an UNSUBSCRIBE Packet.
type UNSUBSCRIBE struct {
	base
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
	// TopicFilters represents a slice of the Topic Filters
	TopicFilters [][]byte
}

// setFixedHeader sets the fixed header to the Packet.
func (p *UNSUBSCRIBE) setFixedHeader() {
	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, TypeUNSUBSCRIBE<<4|0x02)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *UNSUBSCRIBE) setVariableHeader() {
	// Append the Packet Identifier to the variable header.
	p.variableHeader = append(p.variableHeader, encodeUint16(p.PacketID)...)
}

// setPayload sets the payload to the Packet.
func (p *UNSUBSCRIBE) setPayload() {
	// Append each Topic Filter to the payload.
	for _, tf := range p.TopicFilters {
		// Append the Topic Filter to the payload.
		p.payload = appendLenStr(p.payload, tf)
	}
}

// NewUNSUBSCRIBE creates and returns an UNSUBSCRIBE Packet.
func NewUNSUBSCRIBE(opts *UNSUBSCRIBEOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &UNSUBSCRIBEOptions{}
	}

	// Validate the options.
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Create a SUBSCRIBE Packet.
	p := &UNSUBSCRIBE{
		PacketID:     opts.PacketID,
		TopicFilters: opts.TopicFilters,
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
