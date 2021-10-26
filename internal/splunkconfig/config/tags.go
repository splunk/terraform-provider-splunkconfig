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

// Tags is a list of Tag objects.
type Tags []Tag

// satisfiesTag returns true if a given Tag is satisfied by this Tags object.
func (t Tags) satisfiesTag(checkTag Tag) bool {
	for _, tag := range t {
		if tag.satisfiesTag(checkTag) {
			return true
		}
	}

	return false
}

// satisfiesTags returns true if all given Tags are satisfied by this Tags object.
func (t Tags) satisfiesTags(checkTags Tags) bool {
	for _, checkTag := range checkTags {
		if !t.satisfiesTag(checkTag) {
			return false
		}
	}

	return true
}

// excludesTags returns false if any of the given Tags are satisfied by this Tags object.
func (t Tags) excludesTags(checkTags Tags) bool {
	for _, checkTag := range checkTags {
		if t.satisfiesTag(checkTag) {
			return false
		}
	}

	return true
}
