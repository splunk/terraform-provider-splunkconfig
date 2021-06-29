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

func TestIndexes_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			Indexes{
				Index{Name: "index_a"},
				Index{Name: "index_b"},
			},
			false,
		},
		{
			Indexes{
				// even though these objects have different FrozenTimes, they are considered collisions by name
				Index{Name: "index_a", FrozenTime: TimePeriod{Seconds: 1000}},
				Index{Name: "index_a", FrozenTime: TimePeriod{Seconds: 2000}},
			},
			true,
		},
	}

	tests.test(t)
}

func TestIndexes_validateWithRoles(t *testing.T) {
	indexes := Indexes{
		Index{Name: "index_a", SearchRolesAllowed: RoleNames{"role_a"}},
		Index{Name: "index_b", SearchRolesAllowed: RoleNames{"role_b"}},
		Index{Name: "index_c", SearchRolesAllowed: RoleNames{"role_c"}},
	}
	validAgainstRoles := Roles{Role{Name: "role_a"}, Role{Name: "role_b"}, Role{Name: "role_c"}}
	invalidAgainstRoles := Roles{Role{Name: "role_d"}, Role{Name: "role_e"}, Role{Name: "role_f"}}

	tests := []struct {
		roles     Roles
		wantError bool
	}{
		{validAgainstRoles, false},
		{invalidAgainstRoles, true},
	}

	for _, test := range tests {
		gotError := indexes.validateWithRoles(test.roles) != nil
		message := fmt.Sprintf("%T{%+v}.validateWithRoles(%T{%+v}) returned error?", indexes, indexes, test.roles, test.roles)

		testEqual(gotError, test.wantError, message, t)
	}
}

func TestIndexes_indexNamesSearchableByRoleName(t *testing.T) {
	indexA := Index{Name: "index_a", SearchRolesAllowed: RoleNames{"role_a"}}
	indexB := Index{Name: "index_b", SearchRolesAllowed: RoleNames{"role_b"}}
	indexC := Index{Name: "index_c", SearchRolesAllowed: RoleNames{"role_c"}}
	indexAB := Index{Name: "index_a_b", SearchRolesAllowed: RoleNames{"role_a", "role_b"}}

	indexes := Indexes{indexA, indexB, indexC, indexAB}

	tests := []struct {
		roleName RoleName
		want     Indexes
	}{
		{roleName: "role_a", want: Indexes{indexA, indexAB}},
		{roleName: "role_b", want: Indexes{indexB, indexAB}},
		{roleName: "role_c", want: Indexes{indexC}},
	}

	for _, test := range tests {
		got := indexes.indexesSearchableByRoleName(test.roleName)
		message := fmt.Sprintf("%T{%+v}.indexNamesSearchableByRoleName(%q)", indexes, indexes, test.roleName)
		testEqual(got, test.want, message, t)
	}
}

func TestIndexes_lookupRowsForLookup(t *testing.T) {
	tests := lookupRowsForLookupDefinerTestCases{
		{
			Indexes{
				Index{Name: "index_a", LookupRows: LookupRows{LookupRow{LookupName: "index_lookup", Values: LookupValues{"contact": "contact_a"}}}},
				Index{Name: "index_b"},
			},
			Lookup{
				Name: "index_lookup",
				Fields: LookupFields{
					LookupField{Name: "index", DefaultRowField: true},
					LookupField{Name: "contact"},
				},
			},
			LookupRows{
				LookupRow{
					LookupName: "index_lookup",
					Values:     LookupValues{"index": "index_a", "contact": "contact_a"},
				},
				LookupRow{
					LookupName: "index_lookup",
					Values:     LookupValues{"index": "index_b"},
				},
			},
		},
	}

	tests.test(t)
}

func TestIndexes_stanzas(t *testing.T) {
	tests := stanzasDefinerTestCases{
		{
			Indexes{
				Index{Name: "index_no_values"},
				Index{Name: "index_with_values", FrozenTime: TimePeriod{Seconds: 86400}},
			},
			Stanzas{
				Stanza{Name: "index_no_values", Values: StanzaValues{}},
				Stanza{Name: "index_with_values", Values: StanzaValues{"frozenTimePeriodInSecs": "86400"}},
			},
		},
	}

	tests.test(t)
}

func TestIndexes_confFile(t *testing.T) {
	tests := confFileDefinerTestCases{
		{
			Indexes{
				Index{Name: "testindex"},
			},
			ConfFile{
				Name: "indexes",
				Stanzas: Stanzas{
					Stanza{Name: "testindex", Values: StanzaValues{}},
				},
			},
		},
	}

	tests.test(t)
}
