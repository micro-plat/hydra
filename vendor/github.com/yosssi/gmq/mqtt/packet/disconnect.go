package packet

// DISCONNECT represents a DISCONNECT Packet.
type DISCONNECT struct {
	base
}

// NewDISCONNECT creates and returns a DISCONNECT Packet.
func NewDISCONNECT() Packet {
	// Create a DISCONNECT Packet.
	p := &DISCONNECT{}

	// Set the fixed header to the Packet.
	p.fixedHeader = []byte{TypeDISCONNECT << 4, 0x00}

	// Return the Packet.
	return p
}
