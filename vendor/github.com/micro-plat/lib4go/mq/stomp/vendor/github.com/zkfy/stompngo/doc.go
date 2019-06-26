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

/*
	Package stompngo implements a STOMP 1.1+ compatible client library.
	For more STOMP information, see the specification at:
	http://stomp.github.com/.


	Preparation

	Network Connect:

	You are responsible for first establishing a network connection.

	This network connection will be used when you create a stompngo.Connection to
	interact with the STOMP broker.

		h := "localhost"
		p := "61613"
		n, err := net.Dial("tcp", net.JoinHostPort(h, p))
		if err != nil {
			// Do something sane ...
		}


	Shutdown

	Network Disconnect:

	When processing is complete, you MUST close the network
	connection.  If you fail to do this, you will leak goroutines!

		err = n.Close() // Could be defered above, think about it!
		if err != nil {
			// Do something sane ...
		}


	STOMP Frames

	The STOMP specification defines these physical frames that can be sent from a client to a STOMP broker:
		CONNECT		connect to a STOMP broker, any version.
		STOMP		connect to a STOMP broker, specification version 1.1+ only.
		DISCONNECT	disconnect from a STOMP broker.
		SEND		Send a message to a named queue or topic.
		SUBSCRIBE	Prepare to read messages from a named queue or topic.
		UNSUBSCRIBE	Complete reading messages from a named queue or topic.
		ACK		Acknowledge that a message has been received and processed.
		NACK		Deny that a message has been received and processed, specification version 1.1+ only.
		BEGIN		Begin a transaction.
		COMMIT		Commit a transaction.
		ABORT		Abort a transaction.

	The STOMP specification defines these physical frames that a client can receive from a STOMP broker:
		CONNECTED	Broker response upon connection success.
		ERROR		Broker emitted upon any error at any time during an active STOMP connection.
		MESSAGE		A STOMP message frame, possibly with headers and a data payload.
		RECEIPT		A receipt from the broker for a previous frame sent by the client.


	Subscribe and MessageData Channels

	The Subscribe method returns a channel from which you receive MessageData values.

	The channel returned has different characteristics depending on the Stomp Version, and the Headers you pass to Subscribe.

	For details on Subscribe requirements and behavior, see: https://github.com/gmallard/stompngo/wiki/subscribe-and-messagedata


	RECEIPTs

	Receipts are never received on a subscription unique MessageData channel.

	They are always queued to the shared connection level
	stompgo.Connection.MessageData channel.

	The reason for this behavior is because RECEIPT frames do not contain a subscription Header
	(per the STOMP specifications).  See the:

	https://github.com/gmallard/stompngo_examples

	package for several examples.

*/
package stompngo
