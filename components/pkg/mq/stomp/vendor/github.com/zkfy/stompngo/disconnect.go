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
	Disconnect from a STOMP broker.

	Shut down heart beats if necessary.
	Set connection status to false to disable further actions with this
	connection.


	Obtain a receipt unless the client specifically indicates a receipt request
	should be excluded.  If the client  actually asks for a receipt, use the
	supplied receipt id.  Otherwise generate a uniqueue receipt id and add that
	to the DISCONNECT headers.

	Example:
		h := stompngo.Headers{"receipt", "receipt-id1"} // Ask for a receipt
		e := c.Disconnect(h)
		if e != nil {
			// Do something sane ...
		}
		fmt.Printf("%q\n", c.DisconnectReceipt)
		// Or:
		h := stompngo.Headers{"noreceipt", "true"} // Ask for a receipt
		e := c.Disconnect(h)
		if e != nil {
			// Do something sane ...
		}
		fmt.Printf("%q\n", c.DisconnectReceipt)

*/
func (c *Connection) Disconnect(h Headers) error {
	c.discLock.Lock()
	defer c.discLock.Unlock()
	//
	if !c.connected {
		return EDISCPC
	}
	c.log(DISCONNECT, "start", h)
	if !c.connected {
		return ECONBAD
	}
	e := checkHeaders(h, c.Protocol())
	if e != nil {
		return e
	}
	ch := h.Clone()
	// If the caller dous not want a receipt do not ask for one.  Otherwise,
	// add a receipt request if caller did not specifically ask for one.  This is
	// in the spiritof the specification, and allows reasonable resource cleanup
	// in both the client and the message broker.
	_, cwr := ch.Contains("noreceipt")
	if !cwr {
		if _, ok := ch.Contains("receipt"); !ok {
			ch = append(ch, "receipt", Uuid())
		}
	}
	//
	f := Frame{DISCONNECT, ch, NULLBUFF}
	//
	r := make(chan error)
	c.output <- wiredata{f, r}
	e = <-r
	// Drive shutdown logic
	c.shutdown()
	// Only set DisconnectReceipt if we sucessfully received one.
	if !cwr && e == nil {
		// Receipt
		c.DisconnectReceipt = <-c.input
		c.log(DISCONNECT, "end", ch, c.DisconnectReceipt)
	}
	c.rsd <- true
	return e
}
