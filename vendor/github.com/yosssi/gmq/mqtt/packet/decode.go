package packet

import "errors"

// Error value
var ErrInvalidByteLen = errors.New("invalid byte length")

// decodeUint16 converts the slice of bytes in big-endian order
// into an unsigned 16-bit integer.
func decodeUint16(b []byte) (uint16, error) {
	// Check the length of the slice of bytes.
	if len(b) != 2 {
		return 0, ErrInvalidByteLen
	}

	return uint16(b[0])<<8 | uint16(b[1]), nil
}
