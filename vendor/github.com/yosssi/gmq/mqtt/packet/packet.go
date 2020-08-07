package packet

import (
	"errors"
	"io"
)

// Error value
var ErrInvalidPacketType = errors.New("invalid MQTT Control Packet type")

// Packet represents an MQTT Control Packet.
type Packet interface {
	io.WriterTo
	// Type returns the MQTT Control Packet type of the Packet.
	Type() (byte, error)
}

// NewFromBytes creates a Packet from the byte data and returns it.
func NewFromBytes(fixedHeader FixedHeader, remaining []byte) (Packet, error) {
	// Extract the MQTT Control Packet type from the fixed header.
	ptype, err := fixedHeader.ptype()
	if err != nil {
		return nil, err
	}

	// Create and return a Packet.
	switch ptype {
	case TypeCONNACK:
		return NewCONNACKFromBytes(fixedHeader, remaining)
	case TypePUBLISH:
		return NewPUBLISHFromBytes(fixedHeader, remaining)
	case TypePUBACK:
		return NewPUBACKFromBytes(fixedHeader, remaining)
	case TypePUBREC:
		return NewPUBRECFromBytes(fixedHeader, remaining)
	case TypePUBREL:
		return NewPUBRELFromBytes(fixedHeader, remaining)
	case TypePUBCOMP:
		return NewPUBCOMPFromBytes(fixedHeader, remaining)
	case TypeSUBACK:
		return NewSUBACKFromBytes(fixedHeader, remaining)
	case TypeUNSUBACK:
		return NewUNSUBACKFromBytes(fixedHeader, remaining)
	case TypePINGRESP:
		return NewPINGRESPFromBytes(fixedHeader, remaining)
	default:
		return nil, ErrInvalidPacketType
	}
}
