package packet

import (
	"bytes"
	"io"
)

// base holds the fields and methods which are common
// among the MQTT Control Packets.
type base struct {
	// fixedHeader represents the fixed header of the Packet.
	fixedHeader FixedHeader
	// VariableHeader represents the variable header of the Packet.
	variableHeader []byte
	// Payload represents the payload of the Packet.
	payload []byte
}

// WriteTo writes the Packet data to the writer.
func (b *base) WriteTo(w io.Writer) (int64, error) {
	// Create a byte buffer.
	var bf bytes.Buffer

	// Write the Packet data to the buffer.
	bf.Write(b.fixedHeader)
	bf.Write(b.variableHeader)
	bf.Write(b.payload)

	// Write the buffered data to the writer.
	n, err := w.Write(bf.Bytes())

	// Return the result.
	return int64(n), err
}

// Type extracts the MQTT Control Packet type from
// the fixed header and returns it.
func (b *base) Type() (byte, error) {
	return b.fixedHeader.ptype()
}

// appendRemainingLength appends the Remaining Length
// to the fixed header.
func (b *base) appendRemainingLength() {
	// Calculate and encode the Remaining Length.
	rl := encodeLength(uint32(len(b.variableHeader) + len(b.payload)))

	// Append the Remaining Length to the fixed header.
	b.fixedHeader = appendRemainingLength(b.fixedHeader, rl)
}

// appendRemainingLength appends the Remaining Length
// to the slice and returns it.
func appendRemainingLength(b []byte, rl uint32) []byte {
	// Append the Remaining Length to the slice.
	switch {
	case rl&0xFF000000 > 0:
		b = append(b, byte((rl&0xFF000000)>>24))
		fallthrough
	case rl&0x00FF0000 > 0:
		b = append(b, byte((rl&0x00FF0000)>>16))
		fallthrough
	case rl&0x0000FF00 > 0:
		b = append(b, byte((rl&0x0000FF00)>>8))
		fallthrough
	default:
		b = append(b, byte(rl&0x000000FF))
	}

	// Return the slice.
	return b
}
