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

func TestRowsForLookupOrDefaultRows(t *testing.T) {
	tests := []struct {
		rows     LookupRows
		lookup   Lookup
		definer  defaultLookupValuesDefiner
		wantRows LookupRows
	}{
		// non-empty rows returned directly
		{
			LookupRows{LookupRow{LookupName: "test_lookup", Values: LookupValues{"lookup_field": "row_defined_value"}}},
			Lookup{Name: "test_lookup", Fields: LookupFields{LookupField{Name: "lookup_field", Required: true}}},
			LookupFields{LookupField{Name: "lookup_field", Default: "default_value"}},
			LookupRows{LookupRow{LookupName: "test_lookup", Values: LookupValues{"lookup_field": "row_defined_value"}}},
		},
		// empty rows (none for specified lookup) yields default rows for definer
		{
			LookupRows{LookupRow{LookupName: "unmatched_lookup", Values: LookupValues{"lookup_field": "row_defined_value"}}},
			Lookup{Name: "test_lookup", Fields: LookupFields{LookupField{Name: "lookup_field", DefaultRowField: true}}},
			LookupFields{LookupField{Name: "lookup_field", Default: "default_value"}},
			LookupRows{LookupRow{LookupName: "test_lookup", Values: LookupValues{"lookup_field": "default_value"}}},
		},
		// empty rows (none for specified lookup), defaults don't meet lookup's required fields
		{
			LookupRows{LookupRow{LookupName: "unmatched_lookup", Values: LookupValues{"lookup_field": "row_defined_value"}}},
			Lookup{Name: "test_lookup", Fields: LookupFields{LookupField{Name: "lookup_field", Required: true}}},
			LookupFields{LookupField{Name: "unmatched_field", Default: "default_value"}},
			LookupRows{},
		},
		// non-empty rows have definer defaults applied
		{
			LookupRows{LookupRow{LookupName: "test_lookup", Values: LookupValues{"row_field": "row_defined_value"}}},
			Lookup{Name: "test_lookup", Fields: LookupFields{LookupField{Name: "lookup_field", DefaultRowField: true}}},
			LookupFields{LookupField{Name: "lookup_field", Default: "default_value"}},
			LookupRows{LookupRow{LookupName: "test_lookup", Values: LookupValues{"row_field": "row_defined_value", "lookup_field": "default_value"}}},
		},
		// lookup defines no required fields, so no default rows created
		{
			LookupRows{},
			Lookup{Name: "test_lookup", Fields: LookupFields{LookupField{Name: "lookup_field"}}},
			LookupFields{LookupField{Name: "lookup_field", Default: "default_value"}},
			LookupRows{},
		},
		// definer has default values not defined in the lookup, using default rows
		{
			LookupRows{},
			Lookup{Name: "test_lookup", Fields: LookupFields{LookupField{Name: "lookup_field", DefaultRowField: true}}},
			LookupFields{
				LookupField{Name: "lookup_field", Default: "default_value"},
				LookupField{Name: "definer_field", Default: "default_value"},
			},
			LookupRows{LookupRow{LookupName: "test_lookup", Values: LookupValues{"lookup_field": "default_value"}}},
		},
		// definer has default values not defined in the lookup, applied to explicitly defined rows
		{
			LookupRows{LookupRow{LookupName: "test_lookup"}},
			Lookup{Name: "test_lookup", Fields: LookupFields{LookupField{Name: "lookup_field", DefaultRowField: true}}},
			LookupFields{
				LookupField{Name: "lookup_field", Default: "default_value"},
				LookupField{Name: "definer_field", Default: "default_value"},
			},
			LookupRows{LookupRow{LookupName: "test_lookup", Values: LookupValues{"lookup_field": "default_value"}}},
		},
	}

	for _, test := range tests {
		gotRows := rowsForLookupOrDefaultRows(test.rows, test.lookup, test.definer)
		message := fmt.Sprintf("rowsForLookupOrDefaultRows(%T{%+v}, %T{%+v}, %T{%+v})", test.rows, test.rows, test.lookup, test.lookup, test.definer, test.definer)
		testEqual(gotRows, test.wantRows, message, t)
	}
}
