package client

// SubReq represents subscription request.
type SubReq struct {
	// TopicFilter is the Topic Filter of the Subscription.
	TopicFilter []byte
	// QoS is the requsting QoS.
	QoS byte
	// Handler is the handler which handles the Application Message
	// sent from the Server.
	Handler MessageHandler
}
