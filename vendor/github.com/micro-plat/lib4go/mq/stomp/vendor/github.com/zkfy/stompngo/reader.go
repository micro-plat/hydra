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
	"strconv"
	"strings"
	"time"
)

/*
	Logical network reader.

	Read STOMP frames from the connection, create MessageData
	structures from the received data, and push the MessageData to the client.
*/
func (c *Connection) reader() {
	//
	q := false // Shutdown indicator

	for {
		f, e := c.readFrame()
		if e != nil {
			f.Headers = append(f.Headers, "connection_read_error", e.Error())
			md := MessageData{Message(f), e}
			c.handleReadError(md)
			break
		}

		if f.Command == "" && q {
			break
		}

		m := Message(f)
		c.mets.tfr += 1 // Total frames read
		// Headers already decoded
		c.mets.tbr += m.Size(false) // Total bytes read
		d := MessageData{m, e}
		// TODO ? Maybe ? Rethink this logic.
		if sid, ok := f.Headers.Contains("subscription"); ok {
			// This is a read lock
			c.subsLock.RLock()
			c.subs[sid].md <- d
			c.subsLock.RUnlock()
		} else {
			c.input <- d
		}

		c.log("RECEIVE", m.Command, m.Headers)

		select {
		case q = <-c.rsd:
		default:
		}

		if q {
			break
		}

	}
	close(c.input)
	c.log("reader shutdown", time.Now())
}

/*
	Physical frame reader.

	This parses a single STOMP frame from data off of the wire, and
	returns a Frame, with a possible error.

	Note: this functionality could hang or exhibit other erroneous behavior
	if running against a non-compliant STOMP server.
*/
func (c *Connection) readFrame() (f Frame, e error) {
	f = Frame{"", Headers{}, NULLBUFF}
	// Read f.Command or line ends (maybe heartbeats)
	for {
		s, e := c.rdr.ReadString('\n')
		if s == "" {
			return f, e
		}
		if e != nil {
			return f, e
		}
		if c.hbd != nil {
			c.updateHBReads()
		}
		f.Command = s[0 : len(s)-1]
		if s != "\n" {
			break
		}
		// c.log("read slash n")
	}
	// Validate the command
	if _, ok := validCmds[f.Command]; !ok {
		return f, EINVBCMD
	}
	// Read f.Headers
	for {
		s, e := c.rdr.ReadString('\n')
		if e != nil {
			return f, e
		}
		if c.hbd != nil {
			c.updateHBReads()
		}
		if s == "\n" {
			break
		}
		s = s[0 : len(s)-1]
		p := strings.SplitN(s, ":", 2)
		if len(p) != 2 {
			return f, EUNKHDR
		}
		if c.Protocol() != SPL_10 {
			p[0] = decode(p[0])
			p[1] = decode(p[1])
		}
		f.Headers = append(f.Headers, p[0], p[1])
	}
	//
	e = checkHeaders(f.Headers, c.Protocol())
	if e != nil {
		return f, e
	}
	// Read f.Body
	if v, ok := f.Headers.Contains("content-length"); ok {
		l, e := strconv.Atoi(strings.TrimSpace(v))
		if e != nil {
			return f, e
		}
		if l == 0 {
			f.Body, e = readUntilNul(c.rdr)
		} else {
			f.Body, e = readBody(c.rdr, l)
		}
	} else {
		// content-length not present
		f.Body, e = readUntilNul(c.rdr)
	}
	if e != nil {
		return f, e
	}
	if c.hbd != nil {
		c.updateHBReads()
	}
	//
	return f, e
}

func (c *Connection) updateHBReads() {
	c.hbd.rdl.Lock()
	c.hbd.lr = time.Now().UnixNano() // Latest good read
	c.hbd.rdl.Unlock()
}
