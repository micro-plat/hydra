package packet

// PUBRELOptions represents options for a PUBREL Packet.
type PUBRELOptions struct {
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// validate validates the options.
func (opts *PUBRELOptions) validate() error {
	// Check the Packet Identifier.
	if opts.PacketID == 0 {
		return ErrInvalidPacketID
	}

	return nil
}
