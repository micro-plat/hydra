package packet

// PINGREQ represents a PINGREQ Packet.
type PINGREQ struct {
	base
}

// NewPINGREQ creates and returns a PINGREQ Packet.
func NewPINGREQ() Packet {
	// Create a PINGREQ Packet.
	p := &PINGREQ{}

	// Set the fixed header to the Packet.
	p.fixedHeader = []byte{TypePINGREQ << 4, 0x00}

	// Return the Packet.
	return p
}
