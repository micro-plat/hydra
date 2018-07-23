package packet

// SUBSCRIBE represents a SUBSCRIBE Packet.
type SUBSCRIBE struct {
	base
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
	// SubReqs is a slice of the subscription requests.
	SubReqs []*SubReq
}

// setFixedHeader sets the fixed header to the Packet.
func (p *SUBSCRIBE) setFixedHeader() {
	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, TypeSUBSCRIBE<<4|0x02)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *SUBSCRIBE) setVariableHeader() {
	// Append the Packet Identifier to the variable header.
	p.variableHeader = append(p.variableHeader, encodeUint16(p.PacketID)...)
}

// setPayload sets the payload to the Packet.
func (p *SUBSCRIBE) setPayload() {
	// Append each subscription request to the payload.
	for _, s := range p.SubReqs {
		// Append the Topic Filter to the payload.
		p.payload = appendLenStr(p.payload, s.TopicFilter)

		// Append the QoS to the payload.
		p.payload = append(p.payload, s.QoS)
	}
}

// NewSUBSCRIBE creates and returns a SUBSCRIBE Packet.
func NewSUBSCRIBE(opts *SUBSCRIBEOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &SUBSCRIBEOptions{}
	}

	// Validate the options.
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Create a SUBSCRIBE Packet.
	p := &SUBSCRIBE{
		PacketID: opts.PacketID,
		SubReqs:  opts.SubReqs,
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
