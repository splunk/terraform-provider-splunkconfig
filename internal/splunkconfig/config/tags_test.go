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

func TestTags_satisfiesTag(t *testing.T) {
	tests := []struct {
		inputTags Tags
		checkTag  Tag
		want      bool
	}{
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tag{Name: "matched name", Values: []string{"matched value"}},
			true,
		},
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tag{Name: "unmatched name", Values: []string{"matched value"}},
			false,
		},
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tag{Name: "matched name", Values: []string{"unmatched value"}},
			false,
		},
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tag{Name: "unmatched name", Values: []string{"unmatched value"}},
			false,
		},
	}

	for _, test := range tests {
		got := test.inputTags.satisfiesTag(test.checkTag)
		message := fmt.Sprintf("%#v.satisfiesTag(%#v)", test.inputTags, test.checkTag)

		testEqual(got, test.want, message, t)
	}
}

func TestTags_satisfiesTags(t *testing.T) {
	tests := []struct {
		inputTags Tags
		checkTags Tags
		want      bool
	}{
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			true,
		},
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tags{
				{Name: "matched name", Values: []string{"matched value"}},
				{Name: "unmatched name", Values: []string{"matched value"}},
			},
			false,
		},
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tags{
				{Name: "matched name", Values: []string{"matched value"}},
				{Name: "matched name", Values: []string{"unmatched value"}},
			},
			false,
		},
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tags{
				{Name: "matched name", Values: []string{"matched value"}},
				{Name: "unmatched name", Values: []string{"unmatched value"}},
			},
			false,
		},
	}

	for _, test := range tests {
		got := test.inputTags.satisfiesTags(test.checkTags)
		message := fmt.Sprintf("%#v.satisfiesTags(%#v)", test.inputTags, test.checkTags)

		testEqual(got, test.want, message, t)
	}
}

func TestTags_excludesTags(t *testing.T) {
	tests := []struct {
		inputTags Tags
		checkTags Tags
		want      bool
	}{
		// empty checkTags against empty Tags
		{
			Tags{},
			Tags{},
			true,
		},
		// empty checkTags against populated Tags
		{
			Tags{{Name: "any name", Values: []string{"any value"}}},
			Tags{},
			true,
		},
		// one of checkTags exists in Tags
		{
			Tags{{Name: "matched name", Values: []string{"matched value"}}},
			Tags{
				{Name: "matched name", Values: []string{"matched value"}},
				{Name: "unmatched name", Values: []string{"unmatched value"}},
			},
			false,
		},
		// all of checkTags exist in Tags
		{
			Tags{
				{Name: "matched name", Values: []string{"matched value"}},
				{Name: "another matched name", Values: []string{"anothermatched value"}},
			},
			Tags{
				{Name: "matched name", Values: []string{"matched value"}},
				{Name: "another matched name", Values: []string{"anothermatched value"}},
			},
			false,
		},
		// not all values in checkTags exist in Tags
		{
			Tags{
				{Name: "matched name", Values: []string{"matched value"}},
			},
			Tags{
				{Name: "matched name", Values: []string{"matched value", "unmatched value"}},
			},
			true,
		},
		// all values in checkTags exist in Tags (exact match)
		{
			Tags{
				{Name: "matched name", Values: []string{"matched value", "another matched value"}},
			},
			Tags{
				{Name: "matched name", Values: []string{"matched value", "another matched value"}},
			},
			false,
		},
		// all values in checkTags exist in Tags (unordered match)
		{
			Tags{
				{Name: "matched name", Values: []string{"matched value", "another matched value"}},
			},
			Tags{
				{Name: "matched name", Values: []string{"another matched value", "matched value"}},
			},
			false,
		},
		// all values in checkTags exist in Tags (but Tags has more)
		{
			Tags{
				{Name: "matched name", Values: []string{"matched value", "unmatched value"}},
			},
			Tags{
				{Name: "matched name", Values: []string{"matched value"}},
			},
			false,
		},
	}

	for _, test := range tests {
		got := test.inputTags.excludesTags(test.checkTags)
		message := fmt.Sprintf("%#v.excludesTags(%#v)", test.inputTags, test.checkTags)

		testEqual(got, test.want, message, t)
	}
}
