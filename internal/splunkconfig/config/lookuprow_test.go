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

func TestLookupRow_validateForLookupFields(t *testing.T) {
	tests := []struct {
		lookupRow    LookupRow
		lookupFields LookupFields
		wantError    bool
	}{
		{
			LookupRow{Values: LookupValues{"fieldA": "valueA"}},
			LookupFields{LookupField{Name: "fieldA"}},
			false,
		},
		{
			LookupRow{Values: LookupValues{"fieldB": "valueB"}},
			LookupFields{LookupField{Name: "fieldA"}},
			true,
		},
	}

	for _, test := range tests {
		gotError := test.lookupRow.validateForLookupFields(test.lookupFields) != nil
		message := fmt.Sprintf("%T{%+v}.validateForLookupFields(%T{%+v}) returned error?", test.lookupRow, test.lookupRow, test.lookupFields, test.lookupFields)
		testEqual(gotError, test.wantError, message, t)
	}
}

func TestLookupRow_valuesForLookupFields(t *testing.T) {
	tests := []struct {
		lookupRow    LookupRow
		lookupFields LookupFields
		wantValues   []string
	}{
		{
			LookupRow{Values: LookupValues{"fieldA": "valueA", "fieldB": "valueB"}},
			LookupFields{
				LookupField{Name: "fieldA"},
				LookupField{Name: "fieldB"},
			},
			[]string{"valueA", "valueB"},
		},
		{
			LookupRow{Values: LookupValues{"fieldA": "valueA", "fieldB": "valueB"}},
			LookupFields{
				LookupField{Name: "fieldB"},
				LookupField{Name: "fieldA"},
			},
			[]string{"valueB", "valueA"},
		},
	}

	for _, test := range tests {
		gotValues := test.lookupRow.valuesForLookupFields(test.lookupFields)
		message := fmt.Sprintf("%T{%+v}.valuesForLookupFields(%T{%+v})", test.lookupRow, test.lookupRow, test.lookupFields, test.lookupFields)
		testEqual(gotValues, test.wantValues, message, t)
	}
}
