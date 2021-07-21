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

func TestConfFile_validate(t *testing.T) {
	tests := validatorTestCases{
		// no name defined
		{
			ConfFile{},
			true,
		},
		// name defined, no stanzas defined
		{
			ConfFile{Name: "confFileA"},
			false,
		},
		// name defined, valid stanzas
		{
			ConfFile{
				Name: "confFileA",
				Stanzas: Stanzas{
					Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
				},
			},
			false,
		},
		// name defined, invalid stanzas
		{
			ConfFile{
				Name: "confFileA",
				Stanzas: Stanzas{
					Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
					Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "newValueA"}},
				},
			},
			true,
		},
	}

	tests.test(t)
}

func TestConfFile_WithStanzas(t *testing.T) {
	tests := []struct {
		confFile ConfFile
		stanzas  Stanzas
		want     ConfFile
	}{
		// empty to start, empty to add
		{
			ConfFile{},
			Stanzas{},
			ConfFile{},
		},
		// empty to start, stanzas to add
		{
			ConfFile{},
			Stanzas{
				Stanza{Name: "added_stanza_1", Values: StanzaValues{"added_stanza_1_key": "added_stanza_1_value"}},
			},
			ConfFile{
				Stanzas: Stanzas{
					Stanza{Name: "added_stanza_1", Values: StanzaValues{"added_stanza_1_key": "added_stanza_1_value"}},
				},
			},
		},
		// stanzas to start, empty to add
		{
			ConfFile{
				Stanzas: Stanzas{
					Stanza{Name: "existing_stanza_1", Values: StanzaValues{"existing_stanza_1_key": "existing_stanza_1_value"}},
				},
			},
			Stanzas{},
			ConfFile{
				Stanzas: Stanzas{
					Stanza{Name: "existing_stanza_1", Values: StanzaValues{"existing_stanza_1_key": "existing_stanza_1_value"}},
				},
			},
		},
		// stanzas to start, stanzas to add
		{
			ConfFile{
				Stanzas: Stanzas{
					Stanza{Name: "existing_stanza_1", Values: StanzaValues{"existing_stanza_1_key": "existing_stanza_1_value"}},
				},
			},
			Stanzas{
				Stanza{Name: "added_stanza_1", Values: StanzaValues{"added_stanza_1_key": "added_stanza_1_value"}},
			},
			ConfFile{
				Stanzas: Stanzas{
					Stanza{Name: "existing_stanza_1", Values: StanzaValues{"existing_stanza_1_key": "existing_stanza_1_value"}},
					Stanza{Name: "added_stanza_1", Values: StanzaValues{"added_stanza_1_key": "added_stanza_1_value"}},
				},
			},
		},
	}

	for _, test := range tests {
		got := test.confFile.WithStanzas(test.stanzas)
		message := fmt.Sprintf("%T{%+v}.WithStanzas(%T{%+v})", test.confFile, test.confFile, test.stanzas, test.stanzas)
		testEqual(got, test.want, message, t)
	}
}
