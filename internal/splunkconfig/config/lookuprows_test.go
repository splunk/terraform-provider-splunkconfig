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

func TestLookupRows_validateForLookupFields(t *testing.T) {
	tests := []struct {
		lookupRows   LookupRows
		lookupFields LookupFields
		wantError    bool
	}{
		{
			LookupRows{LookupRow{Values: LookupValues{"fieldA": "valueA"}}},
			LookupFields{LookupField{Name: "fieldA"}},
			false,
		},
		{
			LookupRows{LookupRow{Values: LookupValues{"fieldB": "valueB"}}},
			LookupFields{LookupField{Name: "fieldA"}},
			true,
		},
	}

	for _, test := range tests {
		gotError := test.lookupRows.validateForLookupFields(test.lookupFields) != nil
		message := fmt.Sprintf("%T{%+v}.validateForLookupFields(%T{%+v}) returned error?", test.lookupRows, test.lookupRows, test.lookupFields, test.lookupFields)
		testEqual(gotError, test.wantError, message, t)
	}
}

func TestLookupRows_forLookup(t *testing.T) {
	tests := []struct {
		lookupRows LookupRows
		lookup     Lookup
		definer    defaultLookupValuesDefiner
		wantRows   LookupRows
	}{
		{
			LookupRows{
				LookupRow{
					Values: LookupValues{"fieldA": "unnamedLookupRowA"},
				},
				LookupRow{
					LookupName: "lookupA",
					Values:     LookupValues{"fieldA": "lookupARowA"},
				},
				LookupRow{
					LookupName: "lookupB",
					Values:     LookupValues{"fieldA": "lookupBRowA"},
				},
			},
			Lookup{Name: "lookupA"},
			// no object default rows
			LookupFields{},
			LookupRows{
				LookupRow{
					LookupName: "lookupA",
					Values:     LookupValues{"fieldA": "lookupARowA"},
				},
			},
		},
	}

	for _, test := range tests {
		gotRows := test.lookupRows.forLookup(test.lookup, test.definer)
		message := fmt.Sprintf("%T{%+v}.forLookup(%T{%+v})v", test.lookupRows, test.lookupRows, test.lookup, test.lookup)
		testEqual(gotRows, test.wantRows, message, t)
	}
}
