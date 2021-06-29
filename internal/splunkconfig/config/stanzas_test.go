/*
 * Copyright 2021 Splunk, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"
	"testing"
)

func TestStanzas_validate(t *testing.T) {
	tests := validatorTestCases{
		// no duplicates, because only one stanza/value present
		{
			Stanzas{
				Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA", "keyB": "valueB"}},
			},
			false,
		},
		// duplicates, because keyB exists in two stanzas with the name stanzaA
		{
			Stanzas{
				Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA", "keyB": "valueB"}},
				Stanza{Name: "stanzaA", Values: StanzaValues{"keyC": "valueC", "keyD": "valueD"}},
				Stanza{Name: "stanzaA", Values: StanzaValues{"keyE": "valueE", "keyB": "newValueB"}},
			},
			true,
		},
		// not duplicates, because though keyA exists in both stanzas, the stanzas have different names
		{
			Stanzas{
				Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
				Stanza{Name: "stanzaB", Values: StanzaValues{"keyA": "valueA"}},
			},
			false,
		},
	}

	tests.test(t)
}

func TestStanzas_validateNoCollisions(t *testing.T) {
	tests := []struct {
		stanzas      Stanzas
		otherStanzas Stanzas
		wantError    bool
	}{
		{
			Stanzas{
				Stanza{Values: StanzaValues{"keyA": "valueA"}},
			},
			Stanzas{
				Stanza{Values: StanzaValues{"keyB": "valueA"}},
			},
			false,
		},
		{
			Stanzas{
				Stanza{Values: StanzaValues{"keyA": "valueA"}},
				Stanza{Values: StanzaValues{"keyB": "valueB"}},
				Stanza{Values: StanzaValues{"keyC": "valueC"}},
			},
			Stanzas{
				Stanza{Values: StanzaValues{"keyD": "valueD"}},
				Stanza{Values: StanzaValues{"keyC": "valueC"}},
			},
			true,
		},
	}

	for _, test := range tests {
		gotError := test.stanzas.validateNoCollisions(test.otherStanzas) != nil
		message := fmt.Sprintf(
			"%T{%+v}.validateNoCollisions(%T{%+v}) returned error?",
			test.stanzas,
			test.stanzas,
			test.otherStanzas,
			test.otherStanzas)
		testEqual(gotError, test.wantError, message, t)
	}
}

func TestStanzas_TemplatedContent(t *testing.T) {
	tests := contentTemplaterTestCases{
		{
			Stanzas{
				Stanza{
					"index_a",
					StanzaValues{
						"frozenTimePeriodInSecs": "86400",
						"maxTotalDataSizeMB":     "500000",
					},
				},
				Stanza{
					"index_b",
					StanzaValues{
						"frozenTimePeriodInSecs": "86400",
						"maxTotalDataSizeMB":     "500000",
					},
				},
			},
			`[index_a]
frozenTimePeriodInSecs = 86400
maxTotalDataSizeMB = 500000

[index_b]
frozenTimePeriodInSecs = 86400
maxTotalDataSizeMB = 500000

`,
		},
	}

	tests.test(t)
}
