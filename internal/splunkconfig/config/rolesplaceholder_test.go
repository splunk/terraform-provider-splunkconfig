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

func TestRolesPlaceholder_validate(t *testing.T) {
	tests := validatorTestCases{
		// empty placeholder needs to be valid
		{
			RolesPlaceholder{},
			false,
		},
		// has only valid roles
		{
			RolesPlaceholder{
				Roles: Roles{
					Role{Name: "valid_role"},
				},
			},
			false,
		},
		// has invalid role
		{
			RolesPlaceholder{
				Roles: Roles{
					Role{Name: "invalid role"},
				},
			},
			true,
		},
	}

	tests.test(t)
}

func TestRolesPlaceholder_selectedRoles(t *testing.T) {
	tests := []struct {
		inputPlaceholder RolesPlaceholder
		inputRoles       Roles
		wantRoles        Roles
	}{
		// placeholder configured to import roles
		{
			RolesPlaceholder{Import: true},
			Roles{Role{Name: "imported_role"}},
			Roles{Role{Name: "imported_role"}},
		},
		// placeholder doesn't import roles, also doesn't define any
		{
			RolesPlaceholder{},
			Roles{Role{Name: "imported_role"}},
			Roles(nil),
		},
		// placeholder doesn't import roles, uses its own
		{
			RolesPlaceholder{Roles: Roles{Role{Name: "specified_role"}}},
			Roles{Role{Name: "imported_role"}},
			Roles{Role{Name: "specified_role"}},
		},
	}

	for _, test := range tests {
		gotRoles := test.inputPlaceholder.selectedRoles(test.inputRoles)
		message := fmt.Sprintf("%#v.selectedRoles(%#v)", test.inputPlaceholder, test.inputRoles)
		testEqual(gotRoles, test.wantRoles, message, t)
	}
}

func TestRolesPlaceholder_UnmarshalYAML(t *testing.T) {
	tests := yamlUnmarshallerTestCases{
		// empty definition isn't valid
		{
			&RolesPlaceholder{},
			"",
			&RolesPlaceholder{},
			true,
		},
		// explicit RolesPlaceholder
		{
			&RolesPlaceholder{},
			"{roles: [{name: myrole}]}",
			&RolesPlaceholder{Roles: Roles{Role{Name: "myrole"}}},
			false,
		},
		// list of roles
		{
			&RolesPlaceholder{},
			"[{name: myrole}]",
			&RolesPlaceholder{Roles: Roles{Role{Name: "myrole"}}},
			false,
		},
		// boolean
		{
			&RolesPlaceholder{},
			"true",
			&RolesPlaceholder{Import: true},
			false,
		},
	}

	tests.test(t)
}
