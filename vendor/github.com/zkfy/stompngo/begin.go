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

/*
	Begin a STOMP transaction.

	Headers MUST contain a "transaction" header key
	with a value that is not an empty string.

	Example:
		h := stompngo.Headers{"transaction", "transaction-id1",
			"destination", "/queue/mymessages"}
		e := c.Begin(h)
		if e != nil {
			// Do something sane ...
		}
*/
func (c *Connection) Begin(h Headers) error {
	c.log(BEGIN, "start", h)
	if !c.connected {
		return ECONBAD
	}
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return e
	}
	if _, ok := h.Contains("transaction"); !ok {
		return EREQTIDBEG
	}
	if h.Value("transaction") == "" {
		return EREQTIDBEG
	}
	e = c.transmitCommon(BEGIN, h) // transmitCommon Clones() the headers
	c.log(BEGIN, "end", h)
	return e
}
