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

func TestVersion_newFromString(t *testing.T) {
	tests := []struct {
		input       string
		wantVersion Version
		wantError   bool
	}{
		{
			"",
			Version{},
			true,
		},
		{
			"1",
			Version{1, 0, 0},
			false,
		},
		{
			"1.2",
			Version{1, 2, 0},
			false,
		},
		{
			"1.2.3",
			Version{1, 2, 3},
			false,
		},
		{
			"1.2.3.4",
			Version{},
			true,
		},
	}

	for _, test := range tests {
		gotVersion, err := NewVersionFromString(test.input)
		gotError := err != nil

		versionMessage := fmt.Sprintf("newVersionFromString(%q) = %v, want %v", test.input, gotVersion, test.wantVersion)
		testEqual(gotVersion, test.wantVersion, versionMessage, t)

		errorMessage := fmt.Sprintf("newVersionFromString(%q) gave error? %v, want %v", test.input, gotError, test.wantError)
		testEqual(gotError, test.wantError, errorMessage, t)
	}
}

func TestVersion_IsGreaterThan(t *testing.T) {
	tests := []struct {
		versionA Version
		versionB Version
		wantAgtB bool
	}{
		{
			Version{1, 2, 3},
			Version{1, 2, 2},
			true,
		},
		{
			Version{1, 2, 3},
			Version{1, 2, 4},
			false,
		},
		{
			Version{1, 2, 3},
			Version{1, 2, 3},
			false,
		},
		{
			Version{1, 2, 0},
			Version{2, 1, 0},
			false,
		},
	}

	for _, test := range tests {
		gotAgtB := test.versionA.IsGreaterThan(test.versionB)
		message := fmt.Sprintf("%T{%+v}.IsGreaterThan(%T{%+v} = %v, want %v", test.versionA, test.versionA, test.versionB, test.versionB, gotAgtB, test.wantAgtB)
		testEqual(gotAgtB, test.wantAgtB, message, t)
	}
}

func TestVersion_PlusPatchCount(t *testing.T) {
	tests := []struct {
		inputVersion    Version
		inputPatchCount int64
		wantVersion     Version
	}{
		{
			Version{},
			1,
			Version{0, 0, 1},
		},
		{
			Version{0, 0, 1},
			1,
			Version{0, 0, 2},
		},
	}

	for _, test := range tests {
		gotVersion := test.inputVersion.PlusPatchCount(test.inputPatchCount)
		message := fmt.Sprintf("%#v.PlusPatchCount(%d)", test.inputVersion, test.inputPatchCount)

		testEqual(gotVersion, test.wantVersion, message, t)
	}
}

func TestVersion_UnmarshalYAML(t *testing.T) {
	tests := yamlUnmarshallerTestCases{
		{
			&Version{},
			"",
			&Version{},
			true,
		},
		{
			&Version{},
			"{major: 1, minor: 2, patch: 3}",
			&Version{1, 2, 3},
			false,
		},
		{
			&Version{},
			"1.2.3",
			&Version{1, 2, 3},
			false,
		},
		{
			&Version{},
			"1.2.3.4",
			&Version{},
			true,
		},
	}

	tests.test(t)
}
