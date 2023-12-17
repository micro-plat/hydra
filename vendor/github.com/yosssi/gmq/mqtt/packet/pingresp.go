package packet

// Length of the fixed header of the PINGRESP Packet
const lenPINGRESPFixedHeader = 2

// PINGRESP represents a PINGRESP Packet.
type PINGRESP struct {
	base
}

// NewPINGRESPFromBytes creates a PINGRESP Packet from
// the byte data and returns it.
func NewPINGRESPFromBytes(fixedHeader FixedHeader, remaining []byte) (Packet, error) {
	// Validate the byte data.
	if err := validatePINGRESPBytes(fixedHeader, remaining); err != nil {
		return nil, err
	}

	// Create a PINGRESP Packet.
	p := &PINGRESP{}

	// Set the fixed header to the Packet.
	p.fixedHeader = fixedHeader

	// Return the Packet.
	return p, nil
}

// validatePINGRESPBytes validates the fixed header and the remaining.
func validatePINGRESPBytes(fixedHeader FixedHeader, remaining []byte) error {
	// Extract the MQTT Control Packet type.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return err
	}

	// Check the length of the fixed header.
	if len(fixedHeader) != lenPINGRESPFixedHeader {
		return ErrInvalidFixedHeaderLen
	}

	// Check the MQTT Control Packet type.
	if ptype != TypePINGRESP {
		return ErrInvalidPacketType
	}

	// Check the reserved bits of the fixed header.
	if fixedHeader[0]<<4 != 0x00 {
		return ErrInvalidFixedHeader
	}

	// Check the Remaining Length of the fixed header.
	if fixedHeader[1] != 0x00 {
		return ErrInvalidRemainingLength
	}

	// Check the length of the remaining.
	if len(remaining) != 0 {
		return ErrInvalidRemainingLen
	}

	return nil
}
