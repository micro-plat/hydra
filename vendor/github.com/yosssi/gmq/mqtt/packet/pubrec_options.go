package packet

// PUBRECOptions represents options for a PUBREC Packet.
type PUBRECOptions struct {
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// validate validates the options.
func (opts *PUBRECOptions) validate() error {
	// Check the Packet Identifier.
	if opts.PacketID == 0 {
		return ErrInvalidPacketID
	}

	return nil
}
