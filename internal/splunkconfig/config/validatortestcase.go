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

// validatorTest represents a test case for a validator interface.
type validatorTestCase struct {
	validator validator
	wantError bool
}

// test performs the test defined by the validatorTest.
func (v validatorTestCase) test(t *testing.T) {
	gotError := v.validator.validate() != nil
	message := fmt.Sprintf("%T{%+v}.validate() returned error? %v, want %v", v.validator, v.validator, gotError, v.wantError)
	testEqual(gotError, v.wantError, message, t)
}
