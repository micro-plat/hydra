package packet

import "errors"

// Error value
var ErrInvalidNoSubReq = errors.New("subscription request must be specified")

// SUBSCRIBEOptions represents options for a SUBSCRIBE Packet.
type SUBSCRIBEOptions struct {
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
	// SubReqs is a slice of the subscription requests.
	SubReqs []*SubReq
}

// validate validates the options.
func (opts *SUBSCRIBEOptions) validate() error {
	// Check the Packet Identifier.
	if opts.PacketID == 0 {
		return ErrInvalidPacketID
	}

	// Check the existence of the subscription requests.
	if len(opts.SubReqs) == 0 {
		return ErrInvalidNoSubReq
	}

	// Validate each subscription request.
	for _, s := range opts.SubReqs {
		if err := s.validate(); err != nil {
			return err
		}
	}

	return nil
}
