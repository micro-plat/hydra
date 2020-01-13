package client

import "github.com/yosssi/gmq/mqtt/packet"

// session represents a Session which is a stateful interaction
// between a Client and a Server.
type session struct {
	// cleanSession is the Clean Session.
	cleanSession bool
	// clientID is the Client Identifier.
	clientID []byte
	// sendingPackets contains the pairs of the Packet Identifier
	// and the Packet.
	sendingPackets map[uint16]packet.Packet
	// receivingPackets contains the pairs of the Packet Identifier
	// and the Packet.
	receivingPackets map[uint16]packet.Packet
}

// newSession creates and returns a Session.
func newSession(cleanSession bool, clientID []byte) *session {
	return &session{
		cleanSession:     cleanSession,
		clientID:         clientID,
		sendingPackets:   make(map[uint16]packet.Packet),
		receivingPackets: make(map[uint16]packet.Packet),
	}
}
