// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
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

func TestACL_stanzaValues(t *testing.T) {
	tests := []struct {
		input ACL
		want  StanzaValues
	}{
		{
			ACL{},
			StanzaValues{},
		},
		{
			ACL{Read: RoleNames{}, Write: RoleNames{}},
			StanzaValues{"access": "read : [ ], write : [ ]"},
		},
		{
			ACL{
				Read:  RoleNames{"read_role_1", "read_role_2"},
				Write: RoleNames{"write_role_1", "write_role_2"},
			},
			StanzaValues{"access": "read : [ read_role_1, read_role_2 ], write : [ write_role_1, write_role_2 ]"},
		},
		{
			ACL{
				Sharing: "global",
			},
			StanzaValues{"export": "system"},
		},
	}

	for _, test := range tests {
		got := test.input.stanzaValues()
		message := fmt.Sprintf("%T{%+v}.stanzaValues()", test.input, test.input)

		testEqual(got, test.want, message, t)
	}
}
