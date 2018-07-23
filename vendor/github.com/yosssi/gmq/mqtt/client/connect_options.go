package client

import (
	"crypto/tls"
	"time"
)

// ConnectOptions represents options for the Connect method
// of the Client.
type ConnectOptions struct {
	// Network is the network on which the Client connects to.
	Network string
	// Address is the address which the Client connects to.
	Address string
	// TLSConfig is the configuration for the TLS connection.
	TLSConfig *tls.Config
	// CONNACKTimeout is timeout in seconds for the Client
	// to wait for receiving the CONNACK Packet after sending
	// the CONNECT Packet.
	CONNACKTimeout time.Duration
	// PINGRESPTimeout is timeout in seconds for the Client
	// to wait for receiving the PINGRESP Packet after sending
	// the PINGREQ Packet.
	PINGRESPTimeout time.Duration
	// ClientID is the Client Identifier of the payload.
	ClientID []byte
	// UserName is the User Name of the payload.
	UserName []byte
	// Password is the Password of the payload.
	Password []byte
	// CleanSession is the Clean Session of the variable header.
	CleanSession bool
	// KeepAlive is the Keep Alive of the variable header.
	KeepAlive uint16
	// WillTopic is the Will Topic of the payload.
	WillTopic []byte
	// WillMessage is the Will Message of the payload.
	WillMessage []byte
	// WillQoS is the Will QoS of the variable header.
	WillQoS byte
	// WillRetain is the Will Retain of the variable header.
	WillRetain bool
}
