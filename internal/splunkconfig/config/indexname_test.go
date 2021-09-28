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
	"testing"
)

// validateIndexName should return an error when one is expected.
func TestIndexName_validate(t *testing.T) {
	tests := validatorTestCases{
		{IndexName(""), true},
		{IndexName("*"), true},
		{IndexName("-blah"), true},
		{IndexName("_blah"), false},
		{IndexName("mykvstore"), true},
		{IndexName("mail@company"), true},
		{IndexName("mail.company"), true},
		{IndexName("summary7days"), false},
	}

	tests.test(t)
}

// validateIndexName should return an error when one is expected.
func TestIndexName_validatePattern(t *testing.T) {
	tests := []struct {
		input     IndexName
		wantError bool
	}{
		{IndexName(""), true},
		{IndexName("*"), false},
		{IndexName("_*"), false},
		{IndexName("-blah"), true},
		{IndexName("_blah"), false},
		// mykvstore isn't a valid index name, but it's a valid _pattern_
		{IndexName("mykvstore"), false},
		{IndexName("mail@company"), true},
		{IndexName("mail.company"), true},
		{IndexName("summary7days"), false},
	}

	for _, test := range tests {
		gotError := test.input.validatePattern() != nil
		message := fmt.Sprintf("%#v.validatePattern() returned error?", test.input)

		testEqual(gotError, test.wantError, message, t)
	}
}
