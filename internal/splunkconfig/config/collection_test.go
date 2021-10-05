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

import "testing"

func TestCollection_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			Collection{},
			true,
		},
		{
			Collection{
				Name: "validName",
			},
			false,
		},
		{
			Collection{
				Name: "validName",
				Fields: CollectionFields{
					"validField":   "string",
					"invalidField": "invalidType",
				},
			},
			true,
		},
	}

	tests.test(t)
}

func TestCollection_stanza(t *testing.T) {
	tests := stanzaDefinerTestCases{
		{
			Collection{
				Name: "test_collection",
			},
			Stanza{
				Name:   "test_collection",
				Values: StanzaValues{},
			},
		},
		{
			Collection{
				Name: "test_collection",
				Fields: CollectionFields{
					"bool_field":   "bool",
					"string_field": "string",
				},
			},
			Stanza{
				Name: "test_collection",
				Values: StanzaValues{
					"field.bool_field":   "bool",
					"field.string_field": "string",
				},
			},
		},
	}

	tests.test(t)
}
