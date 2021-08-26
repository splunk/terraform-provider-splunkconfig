// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import "fmt"

// Sharing is a string that defines the sharing value of an ACL.
type Sharing string

const (
	SHAREUNDEF  Sharing = ""
	SHAREUSER   Sharing = "user"
	SHAREAPP    Sharing = "app"
	SHAREGLOBAL Sharing = "global"
)

// validate returns an error if Sharing is invalid. It is invalid if:
// * it isn't one of the defined constants
func (sharing Sharing) validate() error {
	switch sharing {
	case SHAREUNDEF, SHAREUSER, SHAREAPP, SHAREGLOBAL:
		break
	default:
		return fmt.Errorf("invalid Sharing value: %s", sharing)
	}

	return nil
}

// metaValue returns a string for use in a .meta file's "export" configuration.
func (sharing Sharing) metaValue() string {
	if sharing == SHAREGLOBAL {
		return "system"
	}

	return ""
}
