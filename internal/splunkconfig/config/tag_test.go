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

func TestTag_hasValue(t *testing.T) {
	tests := []struct {
		inputTag   Tag
		inputValue string
		want       bool
	}{
		{
			Tag{Values: []string{"present"}},
			"present",
			true,
		},
		{
			Tag{Values: []string{"present"}},
			"missing",
			false,
		},
	}

	for _, test := range tests {
		got := test.inputTag.hasValue(test.inputValue)
		message := fmt.Sprintf("%#v.hasValue(%q)", test.inputTag, test.inputValue)

		testEqual(got, test.want, message, t)
	}
}

func TestTag_hasValues(t *testing.T) {
	tests := []struct {
		inputTag    Tag
		inputValues []string
		want        bool
	}{
		{
			Tag{Values: []string{"present", "also present"}},
			[]string{"present", "also present"},
			true,
		},
		{
			Tag{Values: []string{"present", "also present"}},
			[]string{"present", "missing"},
			false,
		},
	}

	for _, test := range tests {
		got := test.inputTag.hasValues(test.inputValues)
		message := fmt.Sprintf("%#v.hasValues(%#v)", test.inputTag, test.inputValues)

		testEqual(got, test.want, message, t)
	}
}

func TestTag_satisfiesTag(t *testing.T) {
	tests := []struct {
		inputTag Tag
		checkTag Tag
		want     bool
	}{
		{
			Tag{Name: "matched name"},
			Tag{Name: "matched name"},
			true,
		},
		{
			Tag{Name: "matched name"},
			Tag{Name: "unmatched name"},
			false,
		},
		{
			Tag{Name: "matched name", Values: []string{"matched value", "extra value"}},
			Tag{Name: "matched name", Values: []string{"matched value"}},
			true,
		},
		{
			Tag{Name: "matched name", Values: []string{"matched value", "extra value"}},
			Tag{Name: "matched name", Values: []string{"unmatched value"}},
			false,
		},
	}

	for _, test := range tests {
		got := test.inputTag.satisfiesTag(test.checkTag)
		message := fmt.Sprintf("%#v.satisfiesTag(%#v)", test.inputTag, test.checkTag)

		testEqual(got, test.want, message, t)
	}
}
