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

// CapabilityName is the name for an authorization capability.
type CapabilityName string

// validate returns an error if the CapabilityName is invalid. Splunk's documentation defines valid capability names:
//
// Only alphanumeric characters, "_" (underscore), and "-" (dash) are allowed in capability names.
func (capabilityName CapabilityName) validate() error {
	validRegex := regexp.MustCompile("^[a-zA-Z0-9_-]+$")

	if !validRegex.MatchString(string(capabilityName)) {
		return fmt.Errorf("invalid CapabilityName %s, may only consist of alphanumeric characters, underscores, and dashes", capabilityName)
	}

	return nil
}

// uid returns the string value of a CapabilityName to uniquely identify it.
func (capabilityName CapabilityName) uid() string {
	return string(capabilityName)
}
