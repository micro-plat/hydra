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
	"fmt"
)

var _ = fmt.Println

/*
	Nack a STOMP 1.1+ message.

	For Stomp 1.1 Headers must contain a "message-id" key and a "subscription"
	header key.


	For Stomp 1.2 Headers must contain a unique "id" header key.


	See the specifications at http://stomp.github.com/ for details.


	Disallowed for an established STOMP 1.0 connection, and EBADVERNAK is returned.

	Example:
		h := stompngo.Headers{"message-id", "message-id1",
			"subscription", "d2cbe608b70a54c8e69d951b246999fbc20df694"}
		e := c.Nack(h)
		if e != nil {
			// Do something sane ...
		}

*/
func (c *Connection) Nack(h Headers) error {
	c.log(NACK, "start", h, c.Protocol())
	if !c.connected {
		return ECONBAD
	}
	if c.Protocol() == SPL_10 {
		return EBADVERNAK
	}
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return e
	}

	switch c.Protocol() {
	case SPL_12:
		if _, ok := h.Contains("id"); !ok {
			return EREQIDNAK
		}
	default: // SPL_11
		if _, ok := h.Contains("subscription"); !ok {
			return EREQSUBNAK
		}
		if _, ok := h.Contains("message-id"); !ok {
			return EREQMIDNAK
		}
	}

	e = c.transmitCommon(NACK, h) // transmitCommon Clones() the headers
	c.log(NACK, "end", h, c.Protocol())
	return e
}
