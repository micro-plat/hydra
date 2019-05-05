package packet

// PUBCOMPOptions represents options for a PUBCOMP Packet.
type PUBCOMPOptions struct {
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// validate validates the options.
func (opts *PUBCOMPOptions) validate() error {
	// Check the Packet Identifier.
	if opts.PacketID == 0 {
		return ErrInvalidPacketID
	}

	return nil
}
