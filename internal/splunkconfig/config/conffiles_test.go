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

func TestConfFiles_validate(t *testing.T) {
	tests := validatorTestCases{
		// empty set of ConfFiles
		{
			ConfFiles{},
			false,
		},
		// only valid ConfFiles - names don't collide
		{
			ConfFiles{
				ConfFile{Name: "confFileA"},
				ConfFile{Name: "confFileB"},
			},
			false,
		},
		// invalid ConfFiles - names collide
		{
			ConfFiles{
				ConfFile{Name: "confFileA"},
				ConfFile{Name: "confFileA"},
			},
			true,
		},
		// invalid ConfFiles - has invalid stanzas
		{
			ConfFiles{
				ConfFile{
					Name: "confFileA",
					Stanzas: Stanzas{
						Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
						Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "newValueA"}},
					},
				},
			},
			true,
		},
	}

	tests.test(t)
}

func TestConfFiles_WithConfFile(t *testing.T) {
	tests := []struct {
		confFiles        ConfFiles
		newConfFile      ConfFile
		wantNewConfFiles ConfFiles
	}{
		// start with empty ConfFiles
		{
			ConfFiles{},
			ConfFile{Name: "confFileA"},
			ConfFiles{
				ConfFile{Name: "confFileA"},
			},
		},
		// add different file
		{
			ConfFiles{
				ConfFile{Name: "confFileA", Stanzas: Stanzas{Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}}}},
			},
			ConfFile{Name: "confFileB", Stanzas: Stanzas{Stanza{Name: "stanzaB", Values: StanzaValues{"keyB": "valueB"}}}},
			ConfFiles{
				ConfFile{Name: "confFileA", Stanzas: Stanzas{Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}}}},
				ConfFile{Name: "confFileB", Stanzas: Stanzas{Stanza{Name: "stanzaB", Values: StanzaValues{"keyB": "valueB"}}}},
			},
		},
		// merge file with different stanzas
		{
			ConfFiles{
				ConfFile{Name: "confFileA", Stanzas: Stanzas{Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}}}},
			},
			ConfFile{Name: "confFileA", Stanzas: Stanzas{Stanza{Name: "stanzaB", Values: StanzaValues{"keyB": "valueB"}}}},
			ConfFiles{
				ConfFile{
					Name: "confFileA",
					Stanzas: Stanzas{
						Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
						Stanza{Name: "stanzaB", Values: StanzaValues{"keyB": "valueB"}},
					},
				},
			},
		},
		// merge file stanza
		{
			ConfFiles{
				ConfFile{Name: "confFileA", Stanzas: Stanzas{Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}}}},
			},
			ConfFile{Name: "confFileA", Stanzas: Stanzas{Stanza{Name: "stanzaA", Values: StanzaValues{"keyB": "valueB"}}}},
			ConfFiles{
				ConfFile{
					Name: "confFileA",
					Stanzas: Stanzas{
						Stanza{Name: "stanzaA", Values: StanzaValues{"keyA": "valueA"}},
						Stanza{Name: "stanzaA", Values: StanzaValues{"keyB": "valueB"}},
					},
				},
			},
		},
	}

	for _, test := range tests {
		gotConfFiles := test.confFiles.WithConfFile(test.newConfFile)
		message := fmt.Sprintf(
			"%T{%+v}.WithConfFile(%T{%+v}) = %+v, want %+v",
			test.confFiles,
			test.confFiles,
			test.newConfFile,
			test.newConfFile,
			gotConfFiles,
			test.wantNewConfFiles)
		testEqual(gotConfFiles, test.wantNewConfFiles, message, t)
	}
}
