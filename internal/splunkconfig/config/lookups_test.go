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

func TestLookups_validate(t *testing.T) {
	tests := validatorTestCases{
		// empty lookup, no reason to fail validation
		{
			Lookups{},
			false,
		},
		// lookup has a row with an invalid field
		{
			Lookups{
				Lookup{
					Name: "lookupA",
					Fields: LookupFields{
						LookupField{
							Name: "fieldA",
						},
					},
					Rows: LookupRows{
						LookupRow{
							Values: LookupValues{
								"fieldB": "valueB",
							},
						},
					},
				},
			},
			true,
		},
	}

	tests.test(t)
}

func TestLookups_extrapolatedWithLookupRowsForLookupDefiners(t *testing.T) {
	tests := []struct {
		lookups     Lookups
		definers    []lookupRowsForLookupDefiner
		wantLookups Lookups
	}{
		{
			// two lookups:
			// indexes: index,contact (contact is optional, thus each defined index will create an entry)
			// index_frequency: index,frequency (both default row fields, indexes that explicitly add a row will get an entry)
			Lookups{
				Lookup{
					Name: "indexes",
					Fields: LookupFields{
						LookupField{
							Name:            "index",
							DefaultRowField: true,
						},
						LookupField{
							Name: "contact",
						},
					},
				},
				Lookup{
					Name: "index_frequency",
					Fields: LookupFields{
						LookupField{
							Name:            "index",
							DefaultRowField: true,
						},
						LookupField{
							Name:            "frequency",
							DefaultRowField: true,
						},
					},
				},
			},
			[]lookupRowsForLookupDefiner{
				Indexes{
					Index{
						Name: "index_a",
					},
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
					Index{
						Name: "index_c",
						LookupRows: LookupRows{
							LookupRow{
								LookupName: "index_frequency",
								Values: LookupValues{
									"frequency": "86400",
								},
							},
						},
					},
				},
			},
			Lookups{
				Lookup{
					Name: "indexes",
					Fields: LookupFields{
						LookupField{
							Name:            "index",
							DefaultRowField: true,
						},
						LookupField{
							Name: "contact",
						},
					},
					Rows: LookupRows{
						LookupRow{
							LookupName: "indexes",
							Values: LookupValues{
								"index": "index_a",
							},
						},
						LookupRow{
							LookupName: "indexes",
							Values: LookupValues{
								"index":   "index_b",
								"contact": "index_b_contact",
							},
						},
						LookupRow{
							LookupName: "indexes",
							Values: LookupValues{
								"index": "index_c",
							},
						},
					},
				},
				Lookup{
					Name: "index_frequency",
					Fields: LookupFields{
						LookupField{
							Name:            "index",
							DefaultRowField: true,
						},
						LookupField{
							Name:            "frequency",
							DefaultRowField: true,
						},
					},
					Rows: LookupRows{
						LookupRow{
							LookupName: "index_frequency",
							Values: LookupValues{
								"index":     "index_c",
								"frequency": "86400",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		gotLookups := test.lookups.extrapolatedWithLookupRowsForLookupDefiners(test.definers...)
		message := fmt.Sprintf("%T{%+v}.extrapolatedWithLookupRowsForLookupDefiners(%T{%+v}...)", test.lookups, test.lookups, test.definers, test.definers)
		testEqual(gotLookups, test.wantLookups, message, t)
	}
}

func TestLookups_WithName(t *testing.T) {
	tests := []struct {
		lookups    Lookups
		findName   string
		wantLookup Lookup
		wantOk     bool
	}{
		{
			Lookups{},
			"missing_lookup",
			Lookup{},
			false,
		},
		{
			Lookups{
				Lookup{Name: "unmatched_lookup"},
				Lookup{Name: "found_lookup"},
			},
			"found_lookup",
			Lookup{Name: "found_lookup"},
			true,
		},
	}

	for _, test := range tests {
		gotLookup, gotOk := test.lookups.WithName(test.findName)

		messageOk := fmt.Sprintf("Lookups(%#v).WithName(%s) ok?", test.lookups, test.findName)
		testEqual(gotOk, test.wantOk, messageOk, t)

		messageLookup := fmt.Sprintf("Lookup(%#v).WithName(%s) returned Lookup", test.lookups, test.findName)
		testEqual(gotLookup, test.wantLookup, messageLookup, t)
	}
}

func TestLookups_WithNames(t *testing.T) {
	tests := []struct {
		lookups     Lookups
		findNames   []string
		wantLookups Lookups
		wantErr     bool
	}{
		{
			Lookups{},
			[]string{"missing_lookup"},
			Lookups{},
			true,
		},
		{
			Lookups{
				Lookup{Name: "unfetched_lookup_1"},
				Lookup{Name: "fetched_lookup_1"},
				Lookup{Name: "fetched_lookup_2"},
				Lookup{Name: "unfetched_lookup_2"},
			},
			[]string{"fetched_lookup_1", "fetched_lookup_2"},
			Lookups{
				Lookup{Name: "fetched_lookup_1"},
				Lookup{Name: "fetched_lookup_2"},
			},
			false,
		},
	}

	for _, test := range tests {
		gotLookups, err := test.lookups.WithNames(test.findNames)
		gotErr := err != nil

		messageErr := fmt.Sprintf("Lookups(%#v).WithNames(%#v) returned error?", test.lookups, test.findNames)
		testEqual(gotErr, test.wantErr, messageErr, t)

		messageLookups := fmt.Sprintf("Lookups(%#v).WithNames(%#v) returned Lookups", test.lookups, test.findNames)
		testEqual(gotLookups, test.wantLookups, messageLookups, t)
	}
}
