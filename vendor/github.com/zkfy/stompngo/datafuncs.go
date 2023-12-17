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
	"unicode/utf8"
)

/*
	Error returns a string for a particular Error.
*/
func (e Error) Error() string {
	return string(e)
}

/*
	BodyString returns a Message body as a string.
*/
func (m *Message) BodyString() string {
	return string(m.Body)
}

/*
	Size returns the size of Message on the wire, in bytes.
*/
func (m *Message) Size(e bool) int64 {
	var r int64 = 0
	r += int64(len(m.Command)) + 1 + m.Headers.Size(e) + 1 + int64(len(m.Body)) + 1
	return r
}

/*
	Size returns the size of Frame on the wire, in bytes.
*/
func (f *Frame) Size(e bool) int64 {
	var r int64 = 0
	r += int64(len(f.Command)) + 1 + f.Headers.Size(e) + 1 + int64(len(f.Body)) + 1
	return r
}

// Headers

/*
	Add appends a key and value pair as a header to a set of Headers.
*/
func (h Headers) Add(k, v string) Headers {
	r := append(h, k, v)
	return r
}

/*
	AddHeaders appends one set of Headers to another.
*/
func (h Headers) AddHeaders(o Headers) Headers {
	r := append(h, o...)
	return r
}

/*
	Compare compares one set of Headers with another.
*/
func (h Headers) Compare(other Headers) bool {
	if len(h) != len(other) {
		return false
	}
	for i, v := range h {
		if v != other[i] {
			return false
		}
	}
	return true
}

/*
	Contains returns true if a set of Headers contains a key.
*/
func (h Headers) Contains(k string) (string, bool) {
	for i := 0; i < len(h); i += 2 {
		if h[i] == k {
			return h[i+1], true
		}
	}
	return "", false
}

/*
	ContainsKV returns true if a set of Headers contains a key and value pair.
*/
func (h Headers) ContainsKV(k string, v string) bool {
	for i := 0; i < len(h); i += 2 {
		if h[i] == k && h[i+1] == v {
			return true
		}
	}
	return false
}

/*
	Value returns a header value for a specified key.  If the key is not present
	an empty string is returned.
*/
func (h Headers) Value(k string) string {
	for i := 0; i < len(h); i += 2 {
		if h[i] == k {
			return h[i+1]
		}
	}
	return ""
}

/*
	Index returns the index of a keader key in Headers.  Return -1 if the
	key is not present.
*/
func (h Headers) Index(k string) int {
	r := -1
	for i := 0; i < len(h); i += 2 {
		if h[i] == k {
			r = i
			break
		}
	}
	return r
}

/*
	Validate performs bacic validation of a set of Headers.
*/
func (h Headers) Validate() error {
	if len(h)%2 != 0 {
		return EHDRLEN
	}
	return nil
}

/*
	ValidateUTF8 validates that header strings are UTF8.
*/
func (h Headers) ValidateUTF8() (string, error) {
	for i := range h {
		if !utf8.ValidString(h[i]) {
			return h[i], EHDRUTF8
		}
	}
	return "", nil
}

/*
	Clone copies a set of Headers.
*/
func (h Headers) Clone() Headers {
	r := make(Headers, len(h))
	copy(r, h)
	return r
}

/*
	Delete removes a key and value pair from a set of Headers.
*/
func (h Headers) Delete(k string) Headers {
	r := h.Clone()
	i := r.Index(k)
	if i >= 0 {
		r = append(r[:i], r[i+2:]...)
	}
	return r
}

/*
	Size returns the size of Headers on the wire, in bytes.
*/
func (h Headers) Size(e bool) int64 {
	l := 0
	for i := 0; i < len(h); i += 2 {
		if e {
			l += len(encode(h[i])) + 1 + len(encode(h[i+1])) + 1
		} else {
			l += len(h[i]) + 1 + len(h[i+1]) + 1
		}
	}
	return int64(l)
}
