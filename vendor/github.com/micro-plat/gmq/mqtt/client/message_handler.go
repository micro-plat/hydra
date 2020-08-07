package client

// MessageHandler is the handler which handles
// the Application Message sent from the Server.
type MessageHandler func(topicName, message []byte)
