//
// Copyright © 2011-2016 Guy M. Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package stompngo

import (
	"log"
	"time"
)

// Exported Connection methods

/*
	Connected returns the current connection status.
*/
func (c *Connection) Connected() bool {
	return c.connected
}

/*
	Session returns the broker assigned session id.
*/
func (c *Connection) Session() string {
	return c.session
}

/*
	Protocol returns the current connection protocol level.
*/
func (c *Connection) Protocol() string {
	return c.protocol
}

/*
	SetLogger enables a client defined logger for this connection.

	Set to "nil" to disable logging.

	Example:
		// Start logging
		l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
		c.SetLogger(l)
*/
func (c *Connection) SetLogger(l *log.Logger) {
	logLock.Lock()
	c.logger = l
	logLock.Unlock()
}

/*
	SendTickerInterval returns any heartbeat send ticker interval in ms.  A return
	value of zero means	no heartbeats are being sent.
*/
func (c *Connection) SendTickerInterval() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.sti / 1000000
}

/*
	ReceiveTickerInterval returns any heartbeat receive ticker interval in ms.
	A return value of zero means no heartbeats are being received.
*/
func (c *Connection) ReceiveTickerInterval() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.rti / 1000000
}

/*
	SendTickerCount returns any heartbeat send ticker count.  A return value of
	zero usually indicates no send heartbeats are enabled.
*/
func (c *Connection) SendTickerCount() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.sc
}

/*
	ReceiveTickerCount returns any heartbeat receive ticker count. A return
	value of zero usually indicates no read heartbeats are enabled.
*/
func (c *Connection) ReceiveTickerCount() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.rc
}

// Package exported functions

/*
	Supported checks if a particular STOMP version is supported in the current
	implementation.
*/
func Supported(v string) bool {
	return hasValue(supported, v)
}

/*
	Protocols returns a slice of client supported protocol levels.
*/
func Protocols() []string {
	return supported
}

/*
	FramesRead returns a count of the number of frames read on the connection.
*/
func (c *Connection) FramesRead() int64 {
	return c.mets.tfr
}

/*
	BytesRead returns a count of the number of bytes read on the connection.
*/
func (c *Connection) BytesRead() int64 {
	return c.mets.tbr
}

/*
	FramesWritten returns a count of the number of frames written on the connection.
*/
func (c *Connection) FramesWritten() int64 {
	return c.mets.tfw
}

/*
	BytesWritten returns a count of the number of bytes written on the connection.
*/
func (c *Connection) BytesWritten() int64 {
	return c.mets.tbw
}

/*
	Running returns a time duration since connection start.
*/
func (c *Connection) Running() time.Duration {
	return time.Since(c.mets.st)
}

/*
	SubChanCap returns the current subscribe channel capacity.
*/
func (c *Connection) SubChanCap() int {
	return c.scc
}

/*
	SetSubChanCap sets a new subscribe channel capacity, to be used during future
	SUBSCRIBE operations.
*/
func (c *Connection) SetSubChanCap(nc int) {
	c.scc = nc
	return
}

// Unexported Connection methods

/*
	Log data if possible.
*/
func (c *Connection) log(v ...interface{}) {
	logLock.Lock()
	if c.logger != nil {
		c.logger.Print(c.session, v)
	}
	logLock.Unlock()
	return
}

/*
	Shutdown logic.
*/
func (c *Connection) shutdown() {
	c.log("SHUTDOWN", "starts")
	// Shutdown heartbeats if necessary
	if c.hbd != nil {
		if c.hbd.hbs {
			c.hbd.ssd <- true
		}
		if c.hbd.hbr {
			c.hbd.rsd <- true
		}
	}
	// Stop writer go routine
	c.wsd <- true
	// Close all individual subscribe channels
	// This is a write lock
	c.subsLock.Lock()
	for key := range c.subs {
		close(c.subs[key].md)
	}
	c.connected = false
	c.subsLock.Unlock()
	c.log("SHUTDOWN", "ends")
	return
}

/*
	Read error handler.
*/
func (c *Connection) handleReadError(md MessageData) {
	// Notify any general subscriber of error
	c.input <- md
	// Notify all individual subscribers of error
	// This is a read lock
	c.subsLock.RLock()
	if c.connected {
		for key := range c.subs {
			c.subs[key].md <- md
		}
	}
	c.subsLock.RUnlock()
	// Let further shutdown logic proceed normally.
	return
}
