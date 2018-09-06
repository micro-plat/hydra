//
// Copyright © 2016 Guy M. Allard
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
	Provide package version information.  A nod to the concept of semver.

	Example:
		fmt.Println("current stompngo version", stompngo.Version())

*/

import (
	"fmt"
)

var (
	pref  = "v"       // Prefix
	major = "1"       // Major
	minor = "0"       // Minor
	patch = "2"       // Patch
)

func Version() string {
	return fmt.Sprintf("%s%s.%s.%s", pref, major, minor, patch)
}
