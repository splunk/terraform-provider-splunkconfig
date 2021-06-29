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

func TestSAMLGroups_extrapolateWithRoles(t *testing.T) {
	tests := []struct {
		inputSAMLGroups SAMLGroups
		inputRoles      Roles
		wantSAMLGroups  SAMLGroups
	}{
		{
			SAMLGroups{
				SAMLGroup{
					Name:  "matched_group",
					Roles: RoleNames{"explicit_role"},
				},
			},
			Roles{
				Role{
					Name:       "implicit_role",
					SAMLGroups: []string{"matched_group"},
				},
				Role{
					Name:       "missing_role",
					SAMLGroups: []string{"unmatched_group"},
				},
			},
			SAMLGroups{
				SAMLGroup{
					Name:  "matched_group",
					Roles: RoleNames{"explicit_role", "implicit_role"},
				},
			},
		},
	}

	for _, test := range tests {
		gotSAMLGroups := test.inputSAMLGroups.extrapolateWithRoles(test.inputRoles)
		message := fmt.Sprintf("%#v.extrapolateWithRoles(%#v)", test.inputSAMLGroups, test.inputRoles)

		testEqual(gotSAMLGroups, test.wantSAMLGroups, message, t)
	}
}

func TestSAMLGroups_WithSAMLGroupName(t *testing.T) {
	tests := []struct {
		inputSAMLGroups    SAMLGroups
		inputSAMLGroupName string
		wantSAMLGroup      SAMLGroup
		wantOk             bool
	}{
		// find the right group
		{
			SAMLGroups{
				SAMLGroup{Name: "unmatched_saml_group_a"},
				SAMLGroup{Name: "matched_saml_group"},
				SAMLGroup{Name: "unmatched_saml_group_b"},
			},
			"matched_saml_group",
			SAMLGroup{Name: "matched_saml_group"},
			true,
		},
		// the searched-for group doesn't exist
		{
			SAMLGroups{
				SAMLGroup{Name: "unmatched_saml_group_a"},
				SAMLGroup{Name: "unmatched_saml_group_b"},
				SAMLGroup{Name: "unmatched_saml_group_c"},
			},
			"matched_saml_group",
			SAMLGroup{},
			false,
		},
	}

	for _, test := range tests {
		gotSAMLGroup, gotOk := test.inputSAMLGroups.WithSAMLGroupName(test.inputSAMLGroupName)

		messageSAMLGroup := fmt.Sprintf("%#v.WithSAMLGroupName(%#v)", test.inputSAMLGroups, test.inputSAMLGroupName)
		testEqual(gotSAMLGroup, test.wantSAMLGroup, messageSAMLGroup, t)

		messageOk := fmt.Sprintf("%#v.WithSAMLGroupName(%#v) returned error?", test.inputSAMLGroups, test.inputSAMLGroupName)
		testEqual(gotOk, test.wantOk, messageOk, t)
	}
}

func TestSAMLGroups_SAMLGroupNames(t *testing.T) {
	tests := []struct {
		input SAMLGroups
		want  []string
	}{
		{
			SAMLGroups{
				SAMLGroup{Name: "saml_group_1"},
				SAMLGroup{Name: "saml_group_2"},
				SAMLGroup{Name: "saml_group_3"},
			},
			[]string{
				"saml_group_1",
				"saml_group_2",
				"saml_group_3",
			},
		},
	}

	for _, test := range tests {
		got := test.input.SAMLGroupNames()
		message := fmt.Sprintf("%#v.SAMLGroupNames()", test.input)

		testEqual(got, test.want, message, t)
	}
}
