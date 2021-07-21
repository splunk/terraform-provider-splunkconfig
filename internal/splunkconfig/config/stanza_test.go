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

func TestStanza_validateNoCollisions(t *testing.T) {
	tests := []struct {
		stanza      Stanza
		otherStanza Stanza
		wantError   bool
	}{
		// different stanza names aren't collisions
		{
			Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
			Stanza{Name: "stanzaB", Values: StanzaValues{"keyA": "valueA"}},
			false,
		},
		// same stanza names, same key names are collisions
		{
			Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
			Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
			true,
		},
	}

	for _, test := range tests {
		gotError := test.stanza.validateNoCollisions(test.otherStanza) != nil
		message := fmt.Sprintf(
			"%T{%+v}.validateNoCollisions(%T{%+v}) returned error?",
			test.stanza,
			test.stanza,
			test.otherStanza,
			test.otherStanza)
		testEqual(gotError, test.wantError, message, t)
	}
}

func TestStanza_TemplatedContent(t *testing.T) {
	tests := contentTemplaterTestCases{
		{
			Stanza{
				Name: "myindex",
				Values: StanzaValues{
					"frozenTimePeriodInSecs": "86400",
					"maxTotalDataSizeMB":     "500000",
				},
			},
			"[myindex]\nfrozenTimePeriodInSecs = 86400\nmaxTotalDataSizeMB = 500000\n",
		},
	}

	tests.test(t)
}
