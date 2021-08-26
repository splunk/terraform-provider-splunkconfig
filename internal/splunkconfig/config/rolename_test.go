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
	"testing"
)

// validateRoleName should return an error when one is expected
func TestRoleName_validate(t *testing.T) {
	tests := validatorTestCases{
		{RoleName(""), true},
		{RoleName("12345"), false},
		{RoleName("abcde"), false},
		{RoleName("ab&de"), true},
		{RoleName("*"), false},
	}

	tests.test(t)
}
