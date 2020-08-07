package client

// PublishOptions represents options for
// the Publish method of the Client.
type PublishOptions struct {
	// QoS is the QoS of the fixed header.
	QoS byte
	// Retain is the Retain of the fixed header.
	Retain bool
	// TopicName is the Topic Name of the varible header.
	TopicName []byte
	// Message is the Application Message of the payload.
	Message []byte
}
