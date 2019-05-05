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
	"bufio"
	"strconv"
	"time"
)

/*
	Logical network writer.  Read wiredata structures from the communication
	channel, and put them on the wire.
*/
func (c *Connection) writer() {
	q := false
	for {

		select {
		case d := <-c.output:
			c.wireWrite(d)
		case q = <-c.wsd:
			break
		}

		if q {
			break
		}

	}
	c.log("writer shutdown", time.Now())
}

/*
	Connection logical write.
*/
func (c *Connection) wireWrite(d wiredata) {
	f := &d.frame
	switch f.Command {
	case "\n": // HeartBeat frame
		if _, e := c.wtr.WriteString(f.Command); e != nil {
			d.errchan <- e
			return
		}
	default: // Other frames
		if e := f.writeFrame(c.wtr, c.Protocol()); e != nil {
			d.errchan <- e
			return
		}
		if e := c.wtr.Flush(); e != nil {
			d.errchan <- e
			return
		}
		if e := c.wtr.WriteByte('\x00'); e != nil {
			d.errchan <- e
			return
		}
	}
	if e := c.wtr.Flush(); e != nil {
		d.errchan <- e
		return
	}
	//
	if c.hbd != nil {
		c.hbd.sdl.Lock()
		c.hbd.ls = time.Now().UnixNano() // Latest good send
		c.hbd.sdl.Unlock()
	}
	c.mets.tfw += 1             // Frame written count
	c.mets.tbw += f.Size(false) // Bytes written count
	//
	d.errchan <- nil
	return
}

/*
	Frame physical write.
*/
func (f *Frame) writeFrame(w *bufio.Writer, l string) error {
	// Write the frame Command
	if _, e := w.WriteString(f.Command + "\n"); e != nil {
		return e
	}
	// Content length - Always add it if client does not suppress it and
	// does not supply it.
	if _, ok := f.Headers.Contains("suppress-content-length"); !ok {
		if _, clok := f.Headers.Contains("content-length"); !clok {
			f.Headers = append(f.Headers, "content-length", strconv.Itoa(len(f.Body)))
		}
	}
	// Write the frame Headers
	for i := 0; i < len(f.Headers); i += 2 {
		if l > SPL_10 && f.Command != CONNECT {
			f.Headers[i] = encode(f.Headers[i])
			f.Headers[i+1] = encode(f.Headers[i+1])
		}
		_, e := w.WriteString(f.Headers[i] + ":" + f.Headers[i+1] + "\n")
		if e != nil {
			return e
		}
	}
	// Write the last Header LF
	if e := w.WriteByte('\n'); e != nil {
		return e
	}
	// Write the body
	if len(f.Body) != 0 { // Foolish to write 0 length data
		if _, e := w.Write(f.Body); e != nil {
			return e
		}
	}
	return nil
}
