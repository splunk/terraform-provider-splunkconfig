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

func TestStanzaValues_validateNoCollisions(t *testing.T) {
	tests := []struct {
		stanzaValues      StanzaValues
		otherStanzaValues StanzaValues
		wantError         bool
	}{
		{
			StanzaValues{"keyA": "valueA"},
			StanzaValues{"keyB": "valueB"},
			false,
		},
		{
			StanzaValues{"keyA": "valueA"},
			StanzaValues{"keyA": "newValueA"},
			true,
		},
		{
			StanzaValues{"keyA": "valueA", "keyB": "valueB", "keyC": "valueC"},
			StanzaValues{"keyB": "newValueB"},
			true,
		},
	}

	for _, test := range tests {
		gotError := test.stanzaValues.validateNoCollisions(test.otherStanzaValues) != nil
		message := fmt.Sprintf(
			"%T{%+v}.validateNoCollisions(%T{%+v}) returned error?",
			test.stanzaValues,
			test.stanzaValues,
			test.otherStanzaValues,
			test.otherStanzaValues)
		testEqual(gotError, test.wantError, message, t)
	}
}

func TestStanzaValues_TemplatedContent(t *testing.T) {
	tests := contentTemplaterTestCases{
		{
			StanzaValues{
				"frozenTimePeriodInSecs": "86400",
			},
			"frozenTimePeriodInSecs = 86400\n",
		},
		{
			StanzaValues{
				"frozenTimePeriodInSecs": "86400",
				"maxTotalDataSizeMB":     "500000",
			},
			"frozenTimePeriodInSecs = 86400\nmaxTotalDataSizeMB = 500000\n",
		},
	}

	tests.test(t)
}
