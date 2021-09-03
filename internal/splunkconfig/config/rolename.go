// Copyright 2021 Splunk, Inc.
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

package config

import (
	"fmt"
	"regexp"
)

// RoleName represents the name of a Splunk role.
type RoleName string

// validate returns an error if RoleName is not valid.
func (roleName RoleName) validate() error {
	// from the Splunk documentation:
	//   * Role names cannot have uppercase characters. Usernames, however, are
	//     case-insensitive.
	//   * Role names cannot contain spaces, colons, semicolons, or forward slashes.

	// the above documentation doesn't have a full list of disallowed characters, so we'll define a valid name
	// as containing only lowercase letters, numbers, underscores, and dashes
	// EXCEPTION: a role name is permitted to be exactly "*" to match "any role"
	generalRegex := regexp.MustCompile("^[a-z0-9_-]+$")

	if roleName != "*" && !generalRegex.MatchString(string(roleName)) {
		return fmt.Errorf("role name (%s) is invalid, should only contain lowercase letters, numbers, underscores, and dashes", roleName)
	}

	return nil
}

// uid returns the string value of a RoleName to determine uniqueness.
func (roleName RoleName) uid() string {
	return string(roleName)
}
