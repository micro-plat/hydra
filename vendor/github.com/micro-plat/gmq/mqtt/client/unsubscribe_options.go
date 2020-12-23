package client

// UnsubscribeOptions represents options for
// the Unsubscribe method of the Client.
type UnsubscribeOptions struct {
	// TopicFilters represents a slice of the Topic Filters.
	TopicFilters [][]byte
}
