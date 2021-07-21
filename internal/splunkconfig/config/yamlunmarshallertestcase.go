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
	"bytes"
	"gopkg.in/yaml.v2"
	"testing"
)

// yamlUnmarshallerTestCase defines a single test case for a yaml.Unmarshaller.
type yamlUnmarshallerTestCase struct {
	inputUnmarshallerPtr yaml.Unmarshaler
	inputContent         string
	wantUnmarshalledPtr  yaml.Unmarshaler
	wantError            bool
}

// test performs the test defined in a yamlUnmarshallerTestCase.
func (testCase yamlUnmarshallerTestCase) test(t *testing.T) {
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(testCase.inputContent)))
	decoder.SetStrict(true)

	gotError := decoder.Decode(testCase.inputUnmarshallerPtr) != nil

	messageError := "decoder.Decode returned error?"
	testEqual(gotError, testCase.wantError, messageError, t)

	messageUnmarshalled := "decoder.Decode unmarshalled"
	testEqual(testCase.inputUnmarshallerPtr, testCase.wantUnmarshalledPtr, messageUnmarshalled, t)
}
