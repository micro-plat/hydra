package packet

// UNSUBSCRIBEOptions represents options for an UNSUBSCRIBE Packet.
type UNSUBSCRIBEOptions struct {
	// PacketID is the Packet Identifier of the variable header.
	PacketID uint16
	// TopicFilters represents a slice of the Topic Filters
	TopicFilters [][]byte
}

// validate validates the options.
func (opts *UNSUBSCRIBEOptions) validate() error {
	// Check the Packet Identifier.
	if opts.PacketID == 0 {
		return ErrInvalidPacketID
	}

	// Check the existence of the subscription requests.
	if len(opts.TopicFilters) == 0 {
		return ErrNoTopicFilter
	}

	// Check the Topic Filters.
	for _, topicFilter := range opts.TopicFilters {
		// Check the length of the Topic Filter.
		l := len(topicFilter)

		if l == 0 {
			return ErrNoTopicFilter
		}

		if l > maxStringsLen {
			return ErrTopicFilterExceedsMaxStringsLen
		}
	}

	return nil
}
