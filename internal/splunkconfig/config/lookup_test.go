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

func TestLookup_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			// valid name and fields
			Lookup{Name: "valid", Fields: LookupFields{LookupField{Name: "field1"}}},
			false,
		},
		{
			// valid name, but no fields
			Lookup{Name: "valid"},
			true,
		},
		{
			// valid fields, but no name
			Lookup{Fields: LookupFields{LookupField{Name: "field1"}}},
			true,
		},
	}

	tests.test(t)
}

func TestLookup_TemplatedContent(t *testing.T) {
	tests := contentTemplaterTestCases{
		{
			Lookup{
				Fields: LookupFields{
					LookupField{Name: "fieldA"},
					LookupField{Name: "fieldB"},
					LookupField{Name: "fieldC"},
				},
				Rows: LookupRows{
					LookupRow{Values: LookupValues{"fieldA": "valueA"}},
					LookupRow{Values: LookupValues{"fieldB": `value"B"`}},
				},
			},
			// in CSV format, double quotes in values are escaped by doubling them, and the entire field is enclosed in double quotes
			// so `,"value""B""",` represents the literal: ,value"B",
			`fieldA,fieldB,fieldC
valueA,,
,"value""B""",
`,
		},
	}

	tests.test(t)
}

func TestLookup_defaultRows(t *testing.T) {
	tests := []struct {
		lookup         Lookup
		defaultDefiner LookupFields
		wantRows       LookupRows
	}{
		// defaultDefiner has no defaults defined
		{
			Lookup{
				Name:   "test_lookup",
				Fields: LookupFields{LookupField{Name: "fieldA"}},
			},
			LookupFields{},
			LookupRows{},
		},
		// defaultDefiner sets some, but not all, default row fields
		{
			Lookup{
				Name: "test_lookup",
				Fields: LookupFields{
					LookupField{Name: "fieldA", DefaultRowField: true},
					LookupField{Name: "fieldB", DefaultRowField: true},
				},
			},
			LookupFields{
				LookupField{Name: "fieldA", Default: "defaultA"},
			},
			LookupRows{},
		},
		// defaultDefiner sets all default row fields, gets a default row
		{
			Lookup{
				Name: "test_lookup",
				Fields: LookupFields{
					LookupField{Name: "fieldA", DefaultRowField: true},
					LookupField{Name: "fieldB", DefaultRowField: true},
				},
			},
			LookupFields{
				LookupField{Name: "fieldA", Default: "defaultA"},
				LookupField{Name: "fieldB", Default: "defaultB"},
			},
			LookupRows{
				LookupRow{
					LookupName: "test_lookup",
					Values:     LookupValues{"fieldA": "defaultA", "fieldB": "defaultB"},
				},
			},
		},
	}

	for _, test := range tests {
		gotRows := test.lookup.defaultRows(test.defaultDefiner)
		message := fmt.Sprintf("%T{%+v}.defaultRows(%T{%+v})", test.lookup, test.lookup, test.defaultDefiner, test.defaultDefiner)
		testEqual(gotRows, test.wantRows, message, t)
	}
}

func TestLookup_extrapolatedWithLookupRowsForLookupDefiners(t *testing.T) {
	tests := []struct {
		lookup     Lookup
		definers   []lookupRowsForLookupDefiner
		wantLookup Lookup
	}{
		{
			Lookup{
				Name: "indexes",
				Fields: LookupFields{
					LookupField{Name: "index", DefaultRowField: true},
					LookupField{Name: "contact"},
				},
			},
			[]lookupRowsForLookupDefiner{
				Indexes{
					// index_a has no custom row values, only default values
					Index{
						Name: "index_a",
					},
					// index_b sets the "contact" field for the "indexes" lookup
					Index{
						Name: "index_b",
						LookupRows: LookupRows{
							LookupRow{
								LookupName: "indexes",
								Values: LookupValues{
									"contact": "index_b_contact",
								},
							},
						},
					},
				},
			},
			Lookup{
				Name: "indexes",
				Fields: LookupFields{
					LookupField{Name: "index", DefaultRowField: true},
					LookupField{Name: "contact"},
				},
				Rows: LookupRows{
					LookupRow{LookupName: "indexes", Values: LookupValues{"index": "index_a"}},
					LookupRow{LookupName: "indexes", Values: LookupValues{"index": "index_b", "contact": "index_b_contact"}},
				},
			},
		},
	}

	for _, test := range tests {
		gotLookup := test.lookup.extrapolatedWithLookupRowsForLookupDefiners(test.definers...)
		message := fmt.Sprintf("%T{%+v}.extrapolatedWithLookupRowsForLookupDefiners(%T{%+v}...)", test.lookup, test.lookup, test.definers, test.definers)
		testEqual(gotLookup, test.wantLookup, message, t)
	}
}
