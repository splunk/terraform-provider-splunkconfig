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

// RoleNames.validate() should return an error when expected.
func TestRoleNames_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			RoleNames{"ok", "anotherok"},
			false,
		},
		{
			RoleNames{"duplicate", "duplicate"},
			true,
		},
	}

	tests.test(t)
}

func TestRoleNames_metaAccessValue(t *testing.T) {
	tests := []struct {
		input RoleNames
		want  string
	}{
		{
			RoleNames{"one", "two"},
			"[ one, two ]",
		},
		{
			RoleNames{},
			"[ ]",
		},
		{
			nil,
			"",
		},
	}

	for _, test := range tests {
		got := test.input.metaAccessValue()
		message := fmt.Sprintf("%T{%+v}.metaAccessValue()", test.input, test.input)

		testEqual(got, test.want, message, t)
	}
}
