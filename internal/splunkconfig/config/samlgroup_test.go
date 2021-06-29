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

func TestSAMLGroup_extrapolatedFromRoles(t *testing.T) {
	tests := []struct {
		inputSAMLGroup SAMLGroup
		inputRoles     Roles
		wantSAMLGroup  SAMLGroup
	}{
		{
			SAMLGroup{
				Name:  "matched_group",
				Roles: RoleNames{"explicit_role"},
			},
			Roles{
				Role{Name: "implicit_role", SAMLGroups: []string{"matched_group"}},
			},
			SAMLGroup{
				Name:  "matched_group",
				Roles: RoleNames{"explicit_role", "implicit_role"},
			},
		},
		{
			SAMLGroup{
				Name:  "matched_group",
				Roles: RoleNames{"explicit_role"},
			},
			Roles{
				Role{Name: "implicit_role", SAMLGroups: []string{"unmatched_group"}},
			},
			SAMLGroup{
				Name:  "matched_group",
				Roles: RoleNames{"explicit_role"},
			},
		},
	}

	for _, test := range tests {
		gotSAMLGroup := test.inputSAMLGroup.extrapolateFromRoles(test.inputRoles)
		message := fmt.Sprintf("%#v.extrapolatedFromRoles(%#v)", test.inputSAMLGroup, test.inputRoles)

		testEqual(gotSAMLGroup, test.wantSAMLGroup, message, t)
	}
}
