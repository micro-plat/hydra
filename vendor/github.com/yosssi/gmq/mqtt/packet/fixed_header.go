package packet

import "errors"

// Error value
var ErrInvalidFixedHeaderLen = errors.New("the length of the fixed header is invalid")

// FixedHeader represents the fixed header of the Packet.
type FixedHeader []byte

// ptype extracts the MQTT Control Packet type from
// the fixed header and returns it.
func (fixedHeader FixedHeader) ptype() (byte, error) {
	// Check the length of the fixed header.
	if len(fixedHeader) < 1 {
		return 0x00, ErrInvalidFixedHeaderLen
	}

	// Extract the MQTT Control Packet type from
	// the fixed header and return it.
	return fixedHeader[0] >> 4, nil
}
