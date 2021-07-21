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

func TestUsers_Names(t *testing.T) {
	tests := []struct {
		input Users
		want  []string
	}{
		{
			Users{
				User{Name: "user_c"},
				User{Name: "user_a"},
				User{Name: "user_b"},
			},
			[]string{
				"user_a",
				"user_b",
				"user_c",
			},
		},
	}

	for _, test := range tests {
		got := test.input.Names()
		message := fmt.Sprintf("%#v.Names()", test.input)

		testEqual(got, test.want, message, t)
	}
}

func TestUsers_WithName(t *testing.T) {
	tests := []struct {
		inputUsers Users
		inputName  string
		wantUser   User
		wantOk     bool
	}{
		{
			Users{
				User{Name: "matched_user"},
			},
			"matched_user",
			User{Name: "matched_user"},
			true,
		},
		{
			Users{
				User{Name: "matched_user"},
			},
			"unmatched_user",
			User{},
			false,
		},
	}

	for _, test := range tests {
		gotUser, gotOk := test.inputUsers.WithName(test.inputName)
		messageUser := fmt.Sprintf("%#v.WithName(%#v)", test.inputUsers, test.inputName)
		messageOk := fmt.Sprintf("%#v.WithName(%#v) returned error?", test.inputUsers, test.inputName)

		testEqual(gotUser, test.wantUser, messageUser, t)
		testEqual(gotOk, test.wantOk, messageOk, t)
	}
}
