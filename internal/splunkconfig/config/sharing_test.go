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

func TestSharing_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			Sharing(""),
			false,
		},
		{
			// "system" isn't a valid value for sharing, though it is what gets placed in .meta
			Sharing("system"),
			true,
		},
		{
			Sharing("app"),
			false,
		},
		{
			Sharing("user"),
			false,
		},
		{
			Sharing("global"),
			false,
		},
	}

	tests.test(t)
}

func TestSharing_metaValue(t *testing.T) {
	tests := []struct {
		input Sharing
		want  string
	}{
		{
			SHAREUNDEF,
			"",
		},
		{
			SHAREUSER,
			"",
		},
		{
			SHAREAPP,
			"",
		},
		{
			SHAREGLOBAL,
			"system",
		},
	}

	for _, test := range tests {
		got := test.input.metaValue()
		message := fmt.Sprintf("%T{%+v}.metaValue()", test.input, test.input)

		testEqual(got, test.want, message, t)
	}
}
