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
	Abort a STOMP transaction.

	Headers MUST contain a "transaction" header key
	with a value that is not an empty string.

	Example:
		h := stompngo.Headers{"transaction", "transaction-id1"}
		e := c.Abort(h)
		if e != nil {
			// Do something sane ...
		}
*/
func (c *Connection) Abort(h Headers) error {
	c.log(ABORT, "start", h)
	if !c.connected {
		return ECONBAD
	}
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return e
	}
	if _, ok := h.Contains("transaction"); !ok {
		return EREQTIDABT
	}
	if h.Value("transaction") == "" {
		return EREQTIDABT
	}
	e = c.transmitCommon(ABORT, h) // transmitCommon Clones() the headers
	c.log(ABORT, "end", h)
	return e
}
