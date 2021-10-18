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

	"gopkg.in/yaml.v2"
)

// yamlMarshalerTestCase defines a single test case for a yaml.Marshaler.
type yamlMarshalerTestCase struct {
	input         yaml.Marshaler
	wantInterface interface{}
	wantError     bool
}

// test performs the test defined in a yamlMarshalerTestCase.
func (testCase yamlMarshalerTestCase) test(t *testing.T) {
	gotInterface, err := testCase.input.MarshalYAML()
	gotError := err != nil

	messageError := fmt.Sprintf("%#v.MarshalYAML() returned error? %v", testCase.input, gotError)
	testEqual(gotError, testCase.wantError, messageError, t)

	messageContent := fmt.Sprintf("%#v.MarshalYAML() = %#v, want %#v", testCase.input, gotInterface, testCase.wantInterface)
	testEqual(gotInterface, testCase.wantInterface, messageContent, t)
}
