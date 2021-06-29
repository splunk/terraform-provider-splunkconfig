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

func TestRoles_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			Roles{
				Role{Name: "role_a"},
				Role{Name: "role_b"},
			},
			false,
		},
		{
			Roles{
				// these roles are collisions even though they have different SearchIndexesAllowed values.
				Role{Name: "role_a", SearchIndexesAllowed: IndexNames{"index_a"}},
				Role{Name: "role_a", SearchIndexesAllowed: IndexNames{"index_b"}},
			},
			true,
		},
	}

	tests.test(t)
}

func TestRoles_roleNameExists(t *testing.T) {
	reusableRoles := Roles{Role{Name: "role_a"}, Role{Name: "role_b"}, Role{Name: "role_c"}}

	tests := []struct {
		roleName RoleName
		want     bool
	}{
		{
			roleName: "role_a",
			want:     true,
		},
		{
			roleName: "role_b",
			want:     true,
		},
		{
			roleName: "role_c",
			want:     true,
		},
		{
			roleName: "role_d",
			want:     false,
		},
	}

	for _, test := range tests {
		got := reusableRoles.roleNameExists(test.roleName)

		message := fmt.Sprintf("%T{%+v}.roleNameExists(%q)", reusableRoles, reusableRoles, test.roleName)
		testEqual(got, test.want, message, t)
	}
}

func TestRoles_extrapolateWithIndexes(t *testing.T) {
	roles := Roles{
		Role{Name: "role_a", SearchIndexesAllowed: IndexNames{"index_a"}},
		Role{Name: "role_b", SearchIndexesAllowed: IndexNames{"index_b"}},
		Role{Name: "role_c", SearchIndexesAllowed: IndexNames{"index_c"}},
	}

	indexes := Indexes{
		Index{Name: "index_ac", SearchRolesAllowed: RoleNames{"role_a", "role_c"}},
		Index{Name: "index_abc", SearchRolesAllowed: RoleNames{"role_a", "role_b", "role_c"}},
	}

	want := Roles{
		Role{Name: "role_a", SearchIndexesAllowed: IndexNames{"index_a", "index_abc", "index_ac"}},
		Role{Name: "role_b", SearchIndexesAllowed: IndexNames{"index_abc", "index_b"}},
		Role{Name: "role_c", SearchIndexesAllowed: IndexNames{"index_abc", "index_ac", "index_c"}},
	}

	got := roles.extrapolateWithIndexes(indexes)
	message := fmt.Sprintf("%T{%+v}.extrapolateWithIndexes(%T{%+v})", roles, roles, indexes, indexes)
	testEqual(got, want, message, t)
}

func TestRoles_lookupRowsForLookup(t *testing.T) {
	tests := lookupRowsForLookupDefinerTestCases{
		{
			Roles{
				Role{Name: "role_a", LookupRows: LookupRows{LookupRow{LookupName: "role_lookup", Values: LookupValues{"contact": "contact_a"}}}},
				Role{Name: "role_b"},
			},
			Lookup{
				Name: "role_lookup",
				Fields: LookupFields{
					LookupField{Name: "role", DefaultRowField: true},
					LookupField{Name: "contact"},
				},
			},
			LookupRows{
				LookupRow{
					LookupName: "role_lookup",
					Values:     LookupValues{"role": "role_a", "contact": "contact_a"},
				},
				LookupRow{
					LookupName: "role_lookup",
					Values:     LookupValues{"role": "role_b"},
				},
			},
		},
	}

	tests.test(t)
}

func TestRoles_stanzas(t *testing.T) {
	tests := stanzasDefinerTestCases{
		{
			Roles{
				Role{Name: "no_values_role"},
				Role{Name: "values_role", SearchIndexesAllowed: IndexNames{"indexA", "indexB"}, ImportRoles: RoleNames{"user", "admin"}},
			},
			Stanzas{
				Stanza{Name: "role_no_values_role", Values: StanzaValues{}},
				Stanza{Name: "role_values_role", Values: StanzaValues{"srchIndexesAllowed": "indexA;indexB", "importRoles": "admin;user"}},
			},
		},
	}

	tests.test(t)
}

func TestRoles_confFile(t *testing.T) {
	tests := confFileDefinerTestCases{
		{
			Roles{
				Role{Name: "testrole"},
			},
			ConfFile{
				Name: "authorize",
				Stanzas: Stanzas{
					Stanza{Name: "role_testrole", Values: StanzaValues{}},
				},
			},
		},
	}

	tests.test(t)
}
