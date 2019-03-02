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
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
)

/*
	Encode a string per STOMP 1.1+ specifications.
*/
func encode(s string) string {
	r := s
	for _, tr := range codecValues {
		if strings.Index(r, tr.decoded) >= 0 {
			r = strings.Replace(r, tr.decoded, tr.encoded, -1)
		}
	}
	return r
}

/*
	Decode a string per STOMP 1.1+ specifications.
*/
func decode(s string) string {
	r := s
	for _, tr := range codecValues {
		if strings.Index(r, tr.encoded) >= 0 {
			r = strings.Replace(r, tr.encoded, tr.decoded, -1)
		}
	}
	return r
}

/*
	A network helper.  Read from the wire until a 0x00 byte is encountered.
*/
func readUntilNul(r *bufio.Reader) ([]uint8, error) {
	b, e := r.ReadBytes(0)
	if e != nil {
		return b, e
	}
	if len(b) == 1 {
		b = NULLBUFF
	} else {
		b = b[0 : len(b)-1]
	}
	return b, e
}

/*
	A network helper.  Read a full message body with a known length that is
	> 0.  Then read the trailing 'null' byte expected for STOMP frames.
*/
func readBody(r *bufio.Reader, l int) ([]uint8, error) {
	b := make([]byte, l)
	n, e := io.ReadFull(r, b)
	if n < l { // Short read, e is ErrUnexpectedEOF
		return b[0 : n-1], e
	}
	if e != nil { // Other erors
		return b, e
	}
	_, _ = r.ReadByte() // trailing NUL
	return b, e
}

/*
	Handle data from the wire after CONNECT is sent. Attempt to create a Frame
	from the wire data.

	Called one time per connection at connection start.
*/
func connectResponse(s string) (*Frame, error) {
	//
	f := new(Frame)
	f.Headers = Headers{}
	f.Body = make([]uint8, 0)

	// Get f.Command
	c := strings.SplitN(s, "\n", 2)
	if len(c) < 2 {
		return nil, EBADFRM
	}
	f.Command = c[0]
	if f.Command != CONNECTED && f.Command != ERROR {
		return f, EUNKFRM
	}

	switch c[1] {
	case "\x00", "\n": // No headers, malformed bodies
		f.Body = []uint8(c[1])
		return f, EBADFRM
	case "\n\x00": // No headers, no body is OK
		return f, nil
	default: // Otherwise continue
	}

	b := strings.SplitN(c[1], "\n\n", 2)
	if len(b) == 1 { // No Headers, b[0] == body
		w := []uint8(b[0])
		f.Body = w[0 : len(w)-1]
		if f.Command == CONNECTED && len(f.Body) > 0 {
			return f, EBDYDATA
		}
		return f, nil
	}

	// Here:
	// b[0] - the headers
	// b[1] - the body

	// Get f.Headers
	for _, l := range strings.Split(b[0], "\n") {
		p := strings.SplitN(l, ":", 2)
		if len(p) < 2 {
			f.Body = []uint8(p[0]) // Bad feedback
			return f, EUNKHDR
		}
		f.Headers = append(f.Headers, p[0], p[1])
	}
	// get f.Body
	w := []uint8(b[1])
	f.Body = w[0 : len(w)-1]
	if f.Command == CONNECTED && len(f.Body) > 0 {
		return f, EBDYDATA
	}

	return f, nil
}

/*
	Sha1 returns a SHA1 hash for a specified string.
*/
func Sha1(q string) string {
	g := sha1.New()
	g.Write([]byte(q))
	return fmt.Sprintf("%x", g.Sum(nil))
}

/*
	Uuid returns a type 4 UUID.
*/
func Uuid() string {
	b := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, b)
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

/*
	Common Header Validation.
*/
func checkHeaders(h Headers, p string) error {
	if h == nil {
		return EHDRNIL
	}
	// Length check
	if e := h.Validate(); e != nil {
		return e
	}
	// Empty key / value check
	for i := 0; i < len(h); i += 2 {
		if h[i] == "" {
			return EHDRMTK
		}
		if p == SPL_10 && h[i+1] == "" {
			return EHDRMTV
		}
	}
	// UTF8 check
	if p != SPL_10 {
		_, e := h.ValidateUTF8()
		if e != nil {
			return e
		}
	}
	return nil
}

/*
	Internal function used by heartbeat initialization.
*/
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

/*
	Internal function, used only during CONNECT processing.
*/
func hasValue(a []string, w string) bool {
	for _, v := range a {
		if v == w {
			return true
		}
	}
	return false
}
