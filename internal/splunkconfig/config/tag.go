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

// Tag contains the list of values for a given tag name.
type Tag struct {
	Name   string
	Values []string
}

// hasValue returns true if the given value exists in the Tag.
func (t Tag) hasValue(checkValue string) bool {
	for _, value := range t.Values {
		if value == checkValue {
			return true
		}
	}

	return false
}

// hasValues returns true if all of the given values exist in the Tag.
func (t Tag) hasValues(checkValues []string) bool {
	for _, checkValue := range checkValues {
		if !t.hasValue(checkValue) {
			return false
		}
	}

	return true
}

// satisfiesTag returns true if Tag has the same name and at least all of the Values
// as another Tag.
func (t Tag) satisfiesTag(checkTag Tag) bool {
	if t.Name != checkTag.Name {
		return false
	}

	if !t.hasValues(checkTag.Values) {
		return false
	}

	return true
}
