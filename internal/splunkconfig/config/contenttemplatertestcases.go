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

// contentTemplaterTestCases is a list of contentTemplaterTestCase objects.
type contentTemplaterTestCases []contentTemplaterTestCase

// test performs the test for each item in contentTemplaterTestCases.
func (c contentTemplaterTestCases) test(t *testing.T) {
	for _, testCase := range c {
		testCase.test(t)
	}
}
