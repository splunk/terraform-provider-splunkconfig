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

func TestLookupsPlaceholder_selectedLookups(t *testing.T) {
	tests := []struct {
		inputPlaceholder LookupsPlaceholder
		inputLookups     Lookups
		wantLookups      Lookups
		wantError        bool
	}{
		// unconfigured shouldn't fail
		{
			LookupsPlaceholder{},
			Lookups{},
			Lookups(nil),
			false,
		},
		// empty (unset) Import returns LookupPlaceholder's Lookups
		{
			LookupsPlaceholder{Lookups: Lookups{Lookup{Name: "my_lookup"}}},
			Lookups{},
			Lookups{Lookup{Name: "my_lookup"}},
			false,
		},
		// populated Import returns requested Lookups
		{
			LookupsPlaceholder{Import: []string{"requested_external_lookup"}},
			Lookups{
				Lookup{Name: "requested_external_lookup"},
				Lookup{Name: "unrequested_external_lookup"},
			},
			Lookups{Lookup{Name: "requested_external_lookup"}},
			false,
		},
	}

	for _, test := range tests {
		gotLookups, err := test.inputPlaceholder.selectedLookups(test.inputLookups)
		gotError := err != nil

		messageError := fmt.Sprintf("%#v.selectedLookups(%#v) returned error?", test.inputPlaceholder, test.inputLookups)
		testEqual(gotError, test.wantError, messageError, t)

		messageLookups := fmt.Sprintf("%#v.selectedLookups(%#v)", test.inputPlaceholder, test.inputLookups)
		testEqual(gotLookups, test.wantLookups, messageLookups, t)
	}
}

func TestLookupsPlaceholder_UnmarshalYAML(t *testing.T) {
	tests := yamlUnmarshallerTestCases{
		// empty definition isn't valid
		{
			&LookupsPlaceholder{},
			"",
			&LookupsPlaceholder{},
			true,
		},
		// explicit LookupsPlaceholder
		{
			&LookupsPlaceholder{},
			"{lookups: [{name: mylookup}]}",
			&LookupsPlaceholder{Lookups: Lookups{Lookup{Name: "mylookup"}}},
			false,
		},
		// list of lookups
		{
			&LookupsPlaceholder{},
			"[{name: mylookup}]",
			&LookupsPlaceholder{Lookups: Lookups{Lookup{Name: "mylookup"}}},
			false,
		},
		// list of lookup names
		{
			&LookupsPlaceholder{},
			"[mylookup]",
			&LookupsPlaceholder{Import: []string{"mylookup"}},
			false,
		},
	}

	tests.test(t)
}
