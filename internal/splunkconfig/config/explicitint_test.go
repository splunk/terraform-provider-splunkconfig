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

import "testing"

func TestExplicitInt_UnmarshalYAML(t *testing.T) {
	tests := yamlUnmarshallerTestCases{
		{
			// empty definition isn't valid
			&ExplicitInt{},
			"",
			&ExplicitInt{},
			true,
		},
		{
			// valid ExplicitInt content
			&ExplicitInt{},
			"{value: 0, explicit: true}",
			&ExplicitInt{Value: 0, Explicit: true},
			false,
		},
		{
			// implicit ExplicitInt via integer value
			&ExplicitInt{},
			"0",
			&ExplicitInt{Value: 0, Explicit: true},
			false,
		},
	}

	for _, test := range tests {
		test.test(t)
	}
}
