package client

// SubscribeOptions represents options for
// the Subscribe method of the Client.
type SubscribeOptions struct {
	// SubReqs is a slice of the subscription requests.
	SubReqs []*SubReq
}
