package packet

// CONNECT represents a CONNECT Packet.
type CONNECT struct {
	base
	// clientID is the Client Identifier of the payload.
	clientID []byte
	// userName is the User Name of the payload.
	userName []byte
	// password is the Password of the payload.
	password []byte
	// cleanSession is the Clean Session of the variable header.
	cleanSession bool
	// keepAlive is the Keep Alive of the variable header.
	keepAlive uint16
	// willTopic is the Will Topic of the payload.
	willTopic []byte
	// willMessage is the Will Message of the payload.
	willMessage []byte
	// willQoS is the Will QoS of the variable header.
	willQoS byte
	// willRetain is the Will Retain of the variable header.
	willRetain bool
}

// setFixedHeader sets the fixed header to the Packet.
func (p *CONNECT) setFixedHeader() {
	// Append the first byte to the fixed header.
	p.fixedHeader = append(p.fixedHeader, TypeCONNECT<<4)

	// Append the Remaining Length to the fixed header.
	p.appendRemainingLength()
}

// setVariableHeader sets the variable header to the Packet.
func (p *CONNECT) setVariableHeader() {
	// Convert the Keep Alive to the slice.
	keepAlive := encodeUint16(p.keepAlive)

	// Create a variable header and set it to the Packet.
	p.variableHeader = []byte{
		0x00,             // Length MSB (0)
		0x04,             // Length LSB (4)
		0x4D,             // 'M'
		0x51,             // 'Q'
		0x54,             // 'T'
		0x54,             // 'T'
		0x04,             // Level(4)
		p.connectFlags(), // Connect Flags
		keepAlive[0],     // Keep Alive MSB
		keepAlive[1],     // Keep Alive LSB
	}
}

// setPayload sets the payload to the Packet.
func (p *CONNECT) setPayload() {
	// Append the Client Identifier to the payload.
	p.payload = appendLenStr(p.payload, p.clientID)

	// Append the Will Topic and the Will Message to the payload
	// if the Packet has them.
	if p.will() {
		p.payload = appendLenStr(p.payload, p.willTopic)
		p.payload = appendLenStr(p.payload, p.willMessage)
	}

	// Append the User Name to the payload if the Packet has it.
	if len(p.userName) > 0 {
		p.payload = appendLenStr(p.payload, p.userName)
	}

	// Append the Password to the payload if the Packet has it.
	if len(p.password) > 0 {
		p.payload = appendLenStr(p.payload, p.password)
	}
}

// connectFlags creates and returns a byte which represents the Connect Flags.
func (p *CONNECT) connectFlags() byte {
	// Create a byte which represents the Connect Flags.
	var b byte

	// Set 1 to the Bit 7 if the Packet has the User Name.
	if len(p.userName) > 0 {
		b |= 0x80
	}

	// Set 1 to the Bit 6 if the Packet has the Password.
	if len(p.password) > 0 {
		b |= 0x40
	}

	// Set 1 to the Bit 5 if the Will Retain is true.
	if p.willRetain {
		b |= 0x20
	}

	// Set the value of the Will QoS to the Bit 4 and 3.
	b |= p.willQoS << 3

	// Set 1 to the Bit 2 if the Packet has the Will Topic and the Will Message.
	if p.will() {
		b |= 0x04
	}

	// Set 1 to the Bit 1 if the Clean Session is true.
	if p.cleanSession {
		b |= 0x02
	}

	// Return the byte.
	return b
}

// will return true if both the Will Topic and the Will Message are not zero-byte.
func (p *CONNECT) will() bool {
	return len(p.willTopic) > 0 && len(p.willMessage) > 0
}

// NewCONNECT creates and returns a CONNECT Packet.
func NewCONNECT(opts *CONNECTOptions) (Packet, error) {
	// Initialize the options.
	if opts == nil {
		opts = &CONNECTOptions{}
	}

	// Validate the options.
	if err := opts.validate(); err != nil {
		return nil, err
	}

	// Create a CONNECT Packet.
	p := &CONNECT{
		clientID:     opts.ClientID,
		userName:     opts.UserName,
		password:     opts.Password,
		cleanSession: opts.CleanSession,
		keepAlive:    opts.KeepAlive,
		willTopic:    opts.WillTopic,
		willMessage:  opts.WillMessage,
		willQoS:      opts.WillQoS,
		willRetain:   opts.WillRetain,
	}

	// Set the variable header to the Packet.
	p.setVariableHeader()

	// Set the payload to the Packet.
	p.setPayload()

	// Set the fixed header to the packet.
	p.setFixedHeader()

	// Return the Packet.
	return p, nil
}
