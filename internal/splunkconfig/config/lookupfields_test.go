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

func TestLookupFields_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			LookupFields{
				LookupField{Name: "fieldA"},
				LookupField{Name: "fieldB"},
			},
			false,
		},
		{
			LookupFields{},
			true,
		},
		{
			LookupFields{
				LookupField{Name: "fieldA"},
				LookupField{Name: ""},
			},
			true,
		},
	}

	tests.test(t)
}

func TestLookupFields_defaultLookupValues(t *testing.T) {
	tests := defaultLookupValuesDefinerTestCases{
		// no defaults
		{
			LookupFields{
				LookupField{Name: "fieldA"},
			},
			LookupValues{},
		},
		// all fields have defaults
		{
			LookupFields{
				LookupField{Name: "fieldA", Default: "defaultA"},
			},
			LookupValues{
				"fieldA": "defaultA",
			},
		},
		// some fields have defaults
		{
			LookupFields{
				LookupField{Name: "fieldA", Default: "defaultA"},
				LookupField{Name: "fieldB"},
			},
			LookupValues{
				"fieldA": "defaultA",
			},
		},
	}

	tests.test(t)
}

func TestLookupFields_defaultLookupValuesDefinerValues(t *testing.T) {
	tests := []struct {
		fields LookupFields
		// LookupFields itself is a defaultLookupValuesDefiner, so we can use it to test here
		defaultsDefiner LookupFields
		wantValues      LookupValues
	}{
		// defaultsDefiner overrides fields defaults
		{
			LookupFields{LookupField{Name: "fieldA", Default: "fieldsDefaultA"}},
			LookupFields{LookupField{Name: "fieldA", Default: "defaultsDefinerDefaultA"}},
			LookupValues{"fieldA": "defaultsDefinerDefaultA"},
		},
		// defaultsDefiner has fields not present in fields
		{
			LookupFields{LookupField{Name: "fieldA"}},
			LookupFields{LookupField{Name: "fieldB"}},
			LookupValues{},
		},
	}

	for _, test := range tests {
		gotValues := test.fields.defaultLookupValuesDefinerValues(test.defaultsDefiner)
		message := fmt.Sprintf("%T{%+v}.defaultLookupValuesDefinerValues(%T{%+v})", test.fields, test.fields, test.defaultsDefiner, test.defaultsDefiner)
		testEqual(gotValues, test.wantValues, message, t)
	}
}

func TestLookupFields_stanzas(t *testing.T) {
	tests := stanzasDefinerTestCases{
		{
			Lookups{
				Lookup{
					Name:   "first_lookup",
					Fields: LookupFields{{Name: "any_field"}},
				},
				Lookup{
					// this Lookup is out of order, and should be reordered in the resulting Stanzas
					Name:   "third_lookup",
					Fields: LookupFields{{Name: "third_field"}},
				},
				Lookup{
					Name:   "second_lookup",
					Fields: LookupFields{{Name: "another_field"}},
				},
			},
			Stanzas{
				{
					Name: "first_lookup",
					Values: StanzaValues{
						// "fields_list": "any_field",
						"filename": "first_lookup.csv",
					},
				},
				{
					Name: "second_lookup",
					Values: StanzaValues{
						// "fields_list": "another_field",
						"filename": "second_lookup.csv",
					},
				},
				{
					Name: "third_lookup",
					Values: StanzaValues{
						// "fields_list": "third_field",
						"filename": "third_lookup.csv",
					},
				},
			},
		},
	}

	tests.test(t)
}

func TestLookups_confFile(t *testing.T) {
	tests := confFileDefinerTestCases{
		{
			Lookups{
				{
					Name:   "test_lookup",
					Fields: LookupFields{{Name: "test_field"}},
				},
			},
			ConfFile{
				Name: "transforms",
				Stanzas: Stanzas{
					{
						Name: "test_lookup",
						Values: StanzaValues{
							"filename": "test_lookup.csv",
						},
					},
				},
			},
		},
	}

	tests.test(t)
}

func TestLookupFields_FieldNames(t *testing.T) {
	tests := []struct {
		input LookupFields
		want  []string
	}{
		{
			LookupFields{
				{Name: "unsorted_field"},
				{Name: "field_1"},
				{Name: "field_2"},
				{Name: "field_3"},
			},
			// fields are not sorted
			[]string{"unsorted_field", "field_1", "field_2", "field_3"},
		},
	}

	for _, test := range tests {
		got := test.input.FieldNames()
		message := fmt.Sprintf("%T{%+v}.FieldNames()", test.input, test.input)
		testEqual(got, test.want, message, t)
	}
}
