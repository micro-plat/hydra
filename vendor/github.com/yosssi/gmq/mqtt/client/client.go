package client

import (
	"errors"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/packet"
)

// Minimum and maximum Packet Identifiers
const (
	minPacketID uint16 = 1
	maxPacketID uint16 = 65535
)

// Error values
var (
	ErrAlreadyConnected = errors.New("the Client has already connected to the Server")
	ErrNotYetConnected  = errors.New("the Client has not yet connected to the Server")
	ErrCONNACKTimeout   = errors.New("the CONNACK Packet was not received within a reasonalbe amount of time")
	ErrPINGRESPTimeout  = errors.New("the PINGRESP Packet was not received within a reasonalbe amount of time")
	ErrPacketIDExhaused = errors.New("Packet Identifiers are exhausted")
	ErrInvalidPINGRESP  = errors.New("invalid PINGRESP Packet")
	ErrInvalidSUBACK    = errors.New("invalid SUBACK Packet")
)

// Client represents a Client.
type Client struct {
	// muConn is the Mutex for the Network Connection.
	muConn sync.RWMutex
	// conn is the Network Connection.
	conn *connection

	// muSess is the Mutex for the Session.
	muSess sync.RWMutex
	// sess is the Session.
	sess *session

	// wg is the Wait Group for the goroutines
	// which are launched by the New method.
	wg sync.WaitGroup
	// disconnc is the channel which handles the signal
	// to disconnect the Network Connection.
	disconnc chan struct{}
	// disconnEndc is the channel which ends the goroutine
	// which disconnects the Network Connection.
	disconnEndc chan struct{}

	// errorHandler is the error handler.
	errorHandler ErrorHandler
}

// Connect establishes a Network Connection to the Server and
// sends a CONNECT Packet to the Server.
func (cli *Client) Connect(opts *ConnectOptions) error {
	// Lock for the connection.
	cli.muConn.Lock()

	// Unlock.
	defer cli.muConn.Unlock()

	// Return an error if the Client has already connected to the Server.
	if cli.conn != nil {
		return ErrAlreadyConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &ConnectOptions{}
	}

	// Establish a Network Connection.
	conn, err := newConnection(opts.Network, opts.Address, opts.TLSConfig)
	if err != nil {
		return err
	}

	// Set the Network Connection to the Client.
	cli.conn = conn

	// Lock for reading and updating the Session.
	cli.muSess.Lock()

	// Create a Session or reuse the current Session.
	if opts.CleanSession || cli.sess == nil {
		// Create a Session and set it to the Client.
		cli.sess = newSession(opts.CleanSession, opts.ClientID)
	} else {
		// Reuse the Session and set its Client Identifier to the options.
		opts.ClientID = cli.sess.clientID
	}

	// Unlock.
	cli.muSess.Unlock()

	// Send a CONNECT Packet to the Server.
	err = cli.sendCONNECT(&packet.CONNECTOptions{
		ClientID:     opts.ClientID,
		UserName:     opts.UserName,
		Password:     opts.Password,
		CleanSession: opts.CleanSession,
		KeepAlive:    opts.KeepAlive,
		WillTopic:    opts.WillTopic,
		WillMessage:  opts.WillMessage,
		WillQoS:      opts.WillQoS,
		WillRetain:   opts.WillRetain,
	})

	if err != nil {
		// Close the Network Connection.
		cli.conn.Close()

		// Clean the Network Connection and the Session if necessary.
		cli.clean()

		return err
	}

	// Launch a goroutine which waits for receiving the CONNACK Packet.
	cli.conn.wg.Add(1)
	go cli.waitPacket(cli.conn.connack, opts.CONNACKTimeout, ErrCONNACKTimeout)

	// Launch a goroutine which receives a Packet from the Server.
	cli.conn.wg.Add(1)
	go cli.receivePackets()

	// Launch a goroutine which sends a Packet to the Server.
	cli.conn.wg.Add(1)
	go cli.sendPackets(time.Duration(opts.KeepAlive), opts.PINGRESPTimeout)

	// Resend the unacknowledged PUBLISH and PUBREL Packets to the Server
	// if the Clean Session is false.
	if !opts.CleanSession {
		// Lock for reading and updating the Session.
		cli.muSess.Lock()

		// Unlock.
		defer cli.muSess.Unlock()

		for id, p := range cli.sess.sendingPackets {
			// Extract the MQTT Control MQTT Control Packet type.
			ptype, err := p.Type()
			if err != nil {
				return err
			}

			switch ptype {
			case packet.TypePUBLISH:
				// Set the DUP flag of the PUBLISH Packet to true.
				p.(*packet.PUBLISH).DUP = true
				// Resend the PUBLISH Packet to the Server.
				cli.conn.send <- p
			case packet.TypePUBREL:
				// Resend the PUBREL Packet to the Server.
				cli.conn.send <- p
			default:
				// Delete the Packet from the Session.
				delete(cli.sess.sendingPackets, id)
			}
		}
	}

	return nil
}

// Disconnect sends a DISCONNECT Packet to the Server and
// closes the Network Connection.
func (cli *Client) Disconnect() error {
	// Lock for the disconnection.
	cli.muConn.Lock()

	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		// Unlock.
		cli.muConn.Unlock()

		return ErrNotYetConnected
	}

	// Send a DISCONNECT Packet to the Server.
	// Ignore the error returned by the send method because
	// we proceed to the subsequent disconnecting processing
	// even if the send method returns the error.
	cli.send(packet.NewDISCONNECT())

	// Close the Network Connection.
	if err := cli.conn.Close(); err != nil {
		// Unlock.
		cli.muConn.Unlock()

		return err
	}

	// Change the state of the Network Connection to disconnected.
	cli.conn.disconnected = true

	// Send the end signal to the goroutine via the channels.
	select {
	case cli.conn.sendEnd <- struct{}{}:
	default:
	}

	// Unlock.
	cli.muConn.Unlock()

	// Wait until all goroutines end.
	cli.conn.wg.Wait()

	// Lock for cleaning the Network Connection.
	cli.muConn.Lock()

	// Lock for cleaning the Session.
	cli.muSess.Lock()

	// Clean the Network Connection and the Session.
	cli.clean()

	// Unlock.
	cli.muSess.Unlock()

	// Unlock.
	cli.muConn.Unlock()

	return nil
}

// Publish sends a PUBLISH Packet to the Server.
func (cli *Client) Publish(opts *PublishOptions) error {
	// Lock for reading.
	cli.muConn.RLock()

	// Unlock.
	defer cli.muConn.RUnlock()

	// Check the Network Connection.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Initialize the options.
	if opts == nil {
		opts = &PublishOptions{}
	}

	// Create a PUBLISH Packet.
	p, err := cli.newPUBLISHPacket(opts)
	if err != nil {
		return err
	}

	// Send the Packet to the Server.
	cli.conn.send <- p

	return nil
}

// Subscribe sends a SUBSCRIBE Packet to the Server.
func (cli *Client) Subscribe(opts *SubscribeOptions) error {
	// Lock for reading and updating.
	cli.muConn.Lock()

	// Unlock.
	defer cli.muConn.Unlock()

	// Check the Network Connection.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Check the existence of the options.
	if opts == nil || len(opts.SubReqs) == 0 {
		return packet.ErrInvalidNoSubReq
	}

	// Define a Packet Identifier.
	var packetID uint16

	// Define an error.
	var err error

	// Lock for updating the Session.
	cli.muSess.Lock()

	defer cli.muSess.Unlock()

	// Generate a Packet Identifer.
	if packetID, err = cli.generatePacketID(); err != nil {
		return err
	}

	// Create subscription requests for the SUBSCRIBE Packet.
	var subReqs []*packet.SubReq

	for _, s := range opts.SubReqs {
		subReqs = append(subReqs, &packet.SubReq{
			TopicFilter: s.TopicFilter,
			QoS:         s.QoS,
		})
	}

	// Create a SUBSCRIBE Packet.
	p, err := packet.NewSUBSCRIBE(&packet.SUBSCRIBEOptions{
		PacketID: packetID,
		SubReqs:  subReqs,
	})
	if err != nil {
		return err
	}

	// Set the Packet to the Session.
	cli.sess.sendingPackets[packetID] = p

	// Set the subscription information to
	// the Network Connection.
	for _, s := range opts.SubReqs {
		cli.conn.unackSubs[string(s.TopicFilter)] = s.Handler
	}

	// Send the Packet to the Server.
	cli.conn.send <- p

	return nil
}

// Unsubscribe sends an UNSUBSCRIBE Packet to the Server.
func (cli *Client) Unsubscribe(opts *UnsubscribeOptions) error {
	// Lock for reading and updating.
	cli.muConn.Lock()

	// Unlock.
	defer cli.muConn.Unlock()

	// Check the Network Connection.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Check the existence of the options.
	if opts == nil || len(opts.TopicFilters) == 0 {
		return packet.ErrNoTopicFilter
	}

	// Define a Packet Identifier.
	var packetID uint16

	// Define an error.
	var err error

	// Lock for updating the Session.
	cli.muSess.Lock()

	defer cli.muSess.Unlock()

	// Generate a Packet Identifer.
	if packetID, err = cli.generatePacketID(); err != nil {
		return err
	}

	// Create an UNSUBSCRIBE Packet.
	p, err := packet.NewUNSUBSCRIBE(&packet.UNSUBSCRIBEOptions{
		PacketID:     packetID,
		TopicFilters: opts.TopicFilters,
	})
	if err != nil {
		return err
	}

	// Set the Packet to the Session.
	cli.sess.sendingPackets[packetID] = p

	// Send the Packet to the Server.
	cli.conn.send <- p

	return nil
}

// Terminate ternimates the Client.
func (cli *Client) Terminate() {
	// Send the end signal to the disconnecting goroutine.
	cli.disconnEndc <- struct{}{}

	// Wait until all goroutines end.
	cli.wg.Wait()
}

// send sends an MQTT Control Packet to the Server.
func (cli *Client) send(p packet.Packet) error {
	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		return ErrNotYetConnected
	}

	// Write the Packet to the buffered writer.
	if _, err := p.WriteTo(cli.conn.w); err != nil {
		return err
	}

	// Flush the buffered writer.
	return cli.conn.w.Flush()
}

// sendCONNECT creates a CONNECT Packet and sends it to the Server.
func (cli *Client) sendCONNECT(opts *packet.CONNECTOptions) error {
	// Initialize the options.
	if opts == nil {
		opts = &packet.CONNECTOptions{}
	}

	// Create a CONNECT Packet.
	p, err := packet.NewCONNECT(opts)
	if err != nil {
		return err
	}

	// Send a CONNECT Packet to the Server.
	return cli.send(p)
}

// receive receives an MQTT Control Packet from the Server.
func (cli *Client) receive() (packet.Packet, error) {
	// Return an error if the Client has not yet connected to the Server.
	if cli.conn == nil {
		return nil, ErrNotYetConnected
	}

	// Get the first byte of the Packet.
	b, err := cli.conn.r.ReadByte()
	if err != nil {
		return nil, err
	}

	// Create the Fixed header.
	fixedHeader := packet.FixedHeader([]byte{b})

	// Get and decode the Remaining Length.
	var mp uint32 = 1 // multiplier
	var rl uint32     // the Remaining Length
	for {
		// Get the next byte of the Packet.
		b, err = cli.conn.r.ReadByte()
		if err != nil {
			return nil, err
		}

		fixedHeader = append(fixedHeader, b)

		rl += uint32(b&0x7F) * mp

		if b&0x80 == 0 {
			break
		}

		mp *= 128
	}

	// Create the Remaining (the Variable header and the Payload).
	remaining := make([]byte, rl)

	if rl > 0 {
		// Get the remaining of the Packet.
		if _, err = io.ReadFull(cli.conn.r, remaining); err != nil {
			return nil, err
		}
	}

	// Create and return a Packet.
	return packet.NewFromBytes(fixedHeader, remaining)
}

// clean cleans the Network Connection and the Session if necessary.
func (cli *Client) clean() {
	// Clean the Network Connection.
	cli.conn = nil

	// Clean the Session if the Clean Session is true.
	if cli.sess != nil && cli.sess.cleanSession {
		cli.sess = nil
	}
}

// waitPacket waits for receiving the Packet.
func (cli *Client) waitPacket(packetc <-chan struct{}, timeout time.Duration, errTimeout error) {
	defer cli.conn.wg.Done()

	var timeoutc <-chan time.Time

	if timeout > 0 {
		timeoutc = time.After(timeout * time.Second)
	}

	select {
	case <-packetc:
	case <-timeoutc:
		// Handle the timeout error.
		cli.handleErrorAndDisconn(errTimeout)
	}
}

// receivePackets receives Packets from the Server.
func (cli *Client) receivePackets() {
	defer func() {
		// Close the channel which handles a signal which
		// notifies the arrival of the CONNACK Packet.
		close(cli.conn.connack)

		cli.conn.wg.Done()
	}()

	for {
		// Receive a Packet from the Server.
		p, err := cli.receive()
		if err != nil {
			// Handle the error and disconnect
			// the Network Connection.
			cli.handleErrorAndDisconn(err)

			// End the goroutine.
			return
		}

		// Handle the Packet.
		if err := cli.handlePacket(p); err != nil {
			// Handle the error and disconnect
			// the Network Connection.
			cli.handleErrorAndDisconn(err)

			// End the goroutine.
			return
		}
	}
}

// handlePacket handles the Packet.
func (cli *Client) handlePacket(p packet.Packet) error {
	// Get the MQTT Control Packet type.
	ptype, err := p.Type()
	if err != nil {
		return err
	}

	switch ptype {
	case packet.TypeCONNACK:
		cli.handleCONNACK()
		return nil
	case packet.TypePUBLISH:
		return cli.handlePUBLISH(p)
	case packet.TypePUBACK:
		return cli.handlePUBACK(p)
	case packet.TypePUBREC:
		return cli.handlePUBREC(p)
	case packet.TypePUBREL:
		return cli.handlePUBREL(p)
	case packet.TypePUBCOMP:
		return cli.handlePUBCOMP(p)
	case packet.TypeSUBACK:
		return cli.handleSUBACK(p)
	case packet.TypeUNSUBACK:
		return cli.handleUNSUBACK(p)
	case packet.TypePINGRESP:
		return cli.handlePINGRESP()
	default:
		return packet.ErrInvalidPacketType
	}
}

// handleCONNACK handles the CONNACK Packet.
func (cli *Client) handleCONNACK() {
	// Notify the arrival of the CONNACK Packet if possible.
	select {
	case cli.conn.connack <- struct{}{}:
	default:
	}
}

// handlePUBLISH handles the PUBLISH Packet.
func (cli *Client) handlePUBLISH(p packet.Packet) error {
	// Get the PUBLISH Packet.
	publish := p.(*packet.PUBLISH)

	switch publish.QoS {
	case mqtt.QoS0:
		// Lock for reading.
		cli.muConn.RLock()

		// Unlock.
		defer cli.muConn.RUnlock()

		// Handle the Application Message.
		cli.handleMessage(publish.TopicName, publish.Message)

		return nil
	case mqtt.QoS1:
		// Lock for reading.
		cli.muConn.RLock()

		// Unlock.
		defer cli.muConn.RUnlock()

		// Handle the Application Message.
		cli.handleMessage(publish.TopicName, publish.Message)

		// Create a PUBACK Packet.
		puback, err := packet.NewPUBACK(&packet.PUBACKOptions{
			PacketID: publish.PacketID,
		})
		if err != nil {
			return err
		}

		// Send the Packet to the Server.
		cli.conn.send <- puback

		return nil
	default:
		// Lock for update.
		cli.muSess.Lock()

		// Unlock.
		defer cli.muSess.Unlock()

		// Validate the Packet Identifier.
		if _, exist := cli.sess.receivingPackets[publish.PacketID]; exist {
			return packet.ErrInvalidPacketID
		}

		// Set the Packet to the Session.
		cli.sess.receivingPackets[publish.PacketID] = p

		// Create a PUBREC Packet.
		pubrec, err := packet.NewPUBREC(&packet.PUBRECOptions{
			PacketID: publish.PacketID,
		})
		if err != nil {
			return err
		}

		// Send the Packet to the Server.
		cli.conn.send <- pubrec

		return nil
	}
}

// handlePUBACK handles the PUBACK Packet.
func (cli *Client) handlePUBACK(p packet.Packet) error {
	// Lock for update.
	cli.muSess.Lock()

	// Unlock.
	defer cli.muSess.Unlock()

	// Extract the Packet Identifier of the Packet.
	id := p.(*packet.PUBACK).PacketID

	// Validate the Packet Identifier.
	if err := cli.validatePacketID(cli.sess.sendingPackets, id, packet.TypePUBLISH); err != nil {
		return err
	}

	// Delete the PUBLISH Packet from the Session.
	delete(cli.sess.sendingPackets, id)

	return nil
}

// handlePUBREC handles the PUBREC Packet.
func (cli *Client) handlePUBREC(p packet.Packet) error {
	// Lock for update.
	cli.muSess.Lock()

	// Unlock.
	defer cli.muSess.Unlock()

	// Extract the Packet Identifier of the Packet.
	id := p.(*packet.PUBREC).PacketID

	// Validate the Packet Identifier.
	if err := cli.validatePacketID(cli.sess.sendingPackets, id, packet.TypePUBLISH); err != nil {
		return err
	}

	// Create a PUBREL Packet.
	pubrel, err := packet.NewPUBREL(&packet.PUBRELOptions{
		PacketID: id,
	})
	if err != nil {
		return err
	}

	// Set the PUBREL Packet to the Session.
	cli.sess.sendingPackets[id] = pubrel

	// Send the Packet to the Server.
	cli.conn.send <- pubrel

	return nil
}

// handlePUBREL handles the PUBREL Packet.
func (cli *Client) handlePUBREL(p packet.Packet) error {
	// Lock for update.
	cli.muSess.Lock()

	// Unlock.
	defer cli.muSess.Unlock()

	// Extract the Packet Identifier of the Packet.
	id := p.(*packet.PUBREL).PacketID

	// Validate the Packet Identifier.
	if err := cli.validatePacketID(cli.sess.receivingPackets, id, packet.TypePUBLISH); err != nil {
		return err
	}

	// Get the Packet from the Session.
	publish := cli.sess.receivingPackets[id].(*packet.PUBLISH)

	// Lock for reading.
	cli.muConn.RLock()

	// Handle the Application Message.
	cli.handleMessage(publish.TopicName, publish.Message)

	// Unlock.
	cli.muConn.RUnlock()

	// Delete the Packet from the Session
	delete(cli.sess.receivingPackets, id)

	// Create a PUBCOMP Packet.
	pubcomp, err := packet.NewPUBCOMP(&packet.PUBCOMPOptions{
		PacketID: id,
	})
	if err != nil {
		return err
	}

	// Send the Packet to the Server.
	cli.conn.send <- pubcomp

	return nil
}

// handlePUBCOMP handles the PUBCOMP Packet.
func (cli *Client) handlePUBCOMP(p packet.Packet) error {
	// Lock for update.
	cli.muSess.Lock()

	// Unlock.
	defer cli.muSess.Unlock()

	// Extract the Packet Identifier of the Packet.
	id := p.(*packet.PUBCOMP).PacketID

	// Validate the Packet Identifier.
	if err := cli.validatePacketID(cli.sess.sendingPackets, id, packet.TypePUBREL); err != nil {
		return err
	}

	// Delete the PUBREL Packet from the Session.
	delete(cli.sess.sendingPackets, id)

	return nil
}

// handleSUBACK handles the SUBACK Packet.
func (cli *Client) handleSUBACK(p packet.Packet) error {
	// Lock for update.
	cli.muConn.Lock()
	cli.muSess.Lock()

	// Unlock.
	defer cli.muConn.Unlock()
	defer cli.muSess.Unlock()

	// Extract the Packet Identifier of the Packet.
	id := p.(*packet.SUBACK).PacketID

	// Validate the Packet Identifier.
	if err := cli.validatePacketID(cli.sess.sendingPackets, id, packet.TypeSUBSCRIBE); err != nil {
		return err
	}

	// Get the subscription requests of the SUBSCRIBE Packet.
	subreqs := cli.sess.sendingPackets[id].(*packet.SUBSCRIBE).SubReqs

	// Delete the SUBSCRIBE Packet from the Session.
	delete(cli.sess.sendingPackets, id)

	// Get the Return Codes of the SUBACK Packet.
	returnCodes := p.(*packet.SUBACK).ReturnCodes

	// Check the lengths of the Return Codes.
	if len(returnCodes) != len(subreqs) {
		return ErrInvalidSUBACK
	}

	// Set the subscriptions to the Network Connection.
	for i, code := range returnCodes {
		// Skip if the Return Code is failure.
		if code == packet.SUBACKRetFailure {
			continue
		}

		// Get the Topic Filter.
		topicFilter := string(subreqs[i].TopicFilter)

		// Move the subscription information from
		// unackSubs to ackedSubs.
		cli.conn.ackedSubs[topicFilter] = cli.conn.unackSubs[topicFilter]
		delete(cli.conn.unackSubs, topicFilter)
	}

	return nil
}

// handleUNSUBACK handles the UNSUBACK Packet.
func (cli *Client) handleUNSUBACK(p packet.Packet) error {
	// Lock for update.
	cli.muConn.Lock()
	cli.muSess.Lock()

	// Unlock.
	defer cli.muConn.Unlock()
	defer cli.muSess.Unlock()

	// Extract the Packet Identifier of the Packet.
	id := p.(*packet.UNSUBACK).PacketID

	// Validate the Packet Identifier.
	if err := cli.validatePacketID(cli.sess.sendingPackets, id, packet.TypeUNSUBSCRIBE); err != nil {
		return err
	}

	// Get the Topic Filters of the UNSUBSCRIBE Packet.
	topicFilters := cli.sess.sendingPackets[id].(*packet.UNSUBSCRIBE).TopicFilters

	// Delete the UNSUBSCRIBE Packet from the Session.
	delete(cli.sess.sendingPackets, id)

	// Delete the Topic Filters from the Network Connection.
	for _, topicFilter := range topicFilters {
		delete(cli.conn.ackedSubs, string(topicFilter))
	}

	return nil
}

// handlePINGRESP handles the PINGRESP Packet.
func (cli *Client) handlePINGRESP() error {
	// Lock for reading and updating pingrespcs.
	cli.conn.muPINGRESPs.Lock()

	// Check the length of pingrespcs.
	if len(cli.conn.pingresps) == 0 {
		// Unlock.
		cli.conn.muPINGRESPs.Unlock()

		// Return an error if there is no channel in pingrespcs.
		return ErrInvalidPINGRESP
	}

	// Get the first channel in pingrespcs.
	pingrespc := cli.conn.pingresps[0]

	// Remove the first channel from pingrespcs.
	cli.conn.pingresps = cli.conn.pingresps[1:]

	// Unlock.
	cli.conn.muPINGRESPs.Unlock()

	// Notify the arrival of the PINGRESP Packet if possible.
	select {
	case pingrespc <- struct{}{}:
	default:
	}

	return nil
}

// handleError handles the error and disconnects
// the Network Connection.
func (cli *Client) handleErrorAndDisconn(err error) {
	// Lock for reading.
	cli.muConn.RLock()

	// Ignore the error and end the process
	// if the Network Connection has already
	// been disconnected.
	if cli.conn == nil || cli.conn.disconnected {
		// Unlock.
		cli.muConn.RUnlock()

		return
	}

	// Unlock.
	cli.muConn.RUnlock()

	// Handle the error.
	if cli.errorHandler != nil {
		cli.errorHandler(err)
	}

	// Send a disconnect signal to the goroutine
	// via the channel if possible.
	select {
	case cli.disconnc <- struct{}{}:
	default:
	}
}

// sendPackets sends Packets to the Server.
func (cli *Client) sendPackets(keepAlive time.Duration, pingrespTimeout time.Duration) {
	defer func() {
		// Lock for reading and updating pingrespcs.
		cli.conn.muPINGRESPs.Lock()

		// Close the channels which handle a signal which
		// notifies the arrival of the PINGREQ Packet.
		for _, pingresp := range cli.conn.pingresps {
			close(pingresp)
		}

		// Initialize pingrespcs
		cli.conn.pingresps = make([]chan struct{}, 0)

		// Unlock.
		cli.conn.muPINGRESPs.Unlock()

		cli.conn.wg.Done()
	}()

	for {
		var keepAlivec <-chan time.Time

		if keepAlive > 0 {
			keepAlivec = time.After(keepAlive * time.Second)
		}

		select {
		case p := <-cli.conn.send:
			// Lock for sending the Packet.
			cli.muConn.RLock()

			// Send the Packet to the Server.
			err := cli.send(p)

			// Unlock.
			cli.muConn.RUnlock()

			if err != nil {
				// Handle the error and disconnect the Network Connection.
				cli.handleErrorAndDisconn(err)

				// End this function.
				return
			}
		case <-keepAlivec:
			// Lock for appending the channel to pingrespcs.
			cli.conn.muPINGRESPs.Lock()

			// Create a channel which handles the signal to notify the arrival of
			// the PINGRESP Packet.
			pingresp := make(chan struct{}, 1)

			// Append the channel to pingrespcs.
			cli.conn.pingresps = append(cli.conn.pingresps, pingresp)

			// Launch a goroutine which waits for receiving the PINGRESP Packet.
			cli.conn.wg.Add(1)
			go cli.waitPacket(pingresp, pingrespTimeout, ErrPINGRESPTimeout)

			// Unlock.
			cli.conn.muPINGRESPs.Unlock()

			// Lock for sending the Packet.
			cli.muConn.RLock()

			// Send a PINGREQ Packet to the Server.
			err := cli.send(packet.NewPINGREQ())

			// Unlock.
			cli.muConn.RUnlock()

			if err != nil {
				// Handle the error and disconnect the Network Connection.
				cli.handleErrorAndDisconn(err)

				// End this function.
				return
			}
		case <-cli.conn.sendEnd:
			// End this function.
			return
		}
	}
}

// generatePacketID generates and returns a Packet Identifier.
func (cli *Client) generatePacketID() (uint16, error) {
	// Define a Packet Identifier.
	id := minPacketID

	for {
		// Find a Packet Identifier which does not used.
		if _, exist := cli.sess.sendingPackets[id]; !exist {
			// Return the Packet Identifier.
			return id, nil
		}

		if id == maxPacketID {
			break
		}

		id++
	}

	// Return an error if available ids are not found.
	return 0, ErrPacketIDExhaused
}

// newPUBLISHPacket creates and returns a PUBLISH Packet.
func (cli *Client) newPUBLISHPacket(opts *PublishOptions) (packet.Packet, error) {
	// Define a Packet Identifier.
	var packetID uint16

	if opts.QoS != mqtt.QoS0 {
		// Lock for reading and updating the Session.
		cli.muSess.Lock()

		defer cli.muSess.Unlock()

		// Define an error.
		var err error

		// Generate a Packet Identifer.
		if packetID, err = cli.generatePacketID(); err != nil {
			return nil, err
		}
	}

	// Create a PUBLISH Packet.
	p, err := packet.NewPUBLISH(&packet.PUBLISHOptions{
		QoS:       opts.QoS,
		Retain:    opts.Retain,
		TopicName: opts.TopicName,
		PacketID:  packetID,
		Message:   opts.Message,
	})
	if err != nil {
		return nil, err
	}

	if opts.QoS != mqtt.QoS0 {
		// Set the Packet to the Session.
		cli.sess.sendingPackets[packetID] = p
	}

	// Return the Packet.
	return p, nil
}

// validateSendingPacketID checks if the Packet which has
// the Packet Identifier and the MQTT Control Packet type
// specified by the parameters exists in the Session's
// sendingPackets.
func (cli *Client) validatePacketID(packets map[uint16]packet.Packet, id uint16, ptype byte) error {
	// Extract the Packet.
	p, exist := packets[id]

	if !exist {
		// Return an error if there is no Packet which has the Packet Identifier
		// specified by the parameter.
		return packet.ErrInvalidPacketID
	}

	// Extract the MQTT Control Packet type of the Packet.
	t, err := p.Type()
	if err != nil {
		return err
	}

	if t != ptype {
		// Return an error if the Packet's MQTT Control Packet type does not
		// equal to one specified by the parameter.
		return packet.ErrInvalidPacketID
	}

	return nil
}

// handleMessage handles the Application Message.
func (cli *Client) handleMessage(topicName, message []byte) {
	// Get the string of the Topic Name.
	topicNameStr := string(topicName)

	for topicFilter, handler := range cli.conn.ackedSubs {
		if handler == nil || !match(topicNameStr, topicFilter) {
			continue
		}

		// Execute the handler.
		go handler(topicName, message)
	}
}

// New creates and returns a Client.
func New(opts *Options) *Client {
	// Initialize the options.
	if opts == nil {
		opts = &Options{}
	}
	// Create a Client.
	cli := &Client{
		disconnc:     make(chan struct{}, 1),
		disconnEndc:  make(chan struct{}),
		errorHandler: opts.ErrorHandler,
	}

	// Launch a goroutine which disconnects the Network Connection.
	cli.wg.Add(1)
	go func() {
		defer func() {
			cli.wg.Done()
		}()

		for {
			select {
			case <-cli.disconnc:
				if err := cli.Disconnect(); err != nil {
					if cli.errorHandler != nil {
						cli.errorHandler(err)
					}
				}
			case <-cli.disconnEndc:
				// End the goroutine.
				return
			}
		}

	}()

	// Return the Client.
	return cli
}

// match checks if the Topic Name matches the Topic Filter.
func match(topicName, topicFilter string) bool {
	// Tokenize the Topic Name.
	nameTokens := strings.Split(topicName, "/")
	nameTokensLen := len(nameTokens)

	// Tolenize the Topic Filter.
	filterTokens := strings.Split(topicFilter, "/")

	for i, t := range filterTokens {
		switch t {
		case "#":
			return i != 0 || !strings.HasPrefix(nameTokens[0], "$")
		case "+":
			if i == 0 && strings.HasPrefix(nameTokens[0], "$") {
				return false
			}

			if nameTokensLen <= i {
				return false
			}
		default:
			if nameTokensLen <= i || t != nameTokens[i] {
				return false
			}
		}
	}

	return len(filterTokens) == nameTokensLen
}
