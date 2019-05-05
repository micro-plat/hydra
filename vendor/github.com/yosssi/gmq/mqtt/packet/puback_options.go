package packet

// PUBACKOptions represents options for a PUBACK Packet.
type PUBACKOptions struct {
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
}

// validate validates the options.
func (opts *PUBACKOptions) validate() error {
	// Check the Packet Identifier.
	if opts.PacketID == 0 {
		return ErrInvalidPacketID
	}

	return nil
}
