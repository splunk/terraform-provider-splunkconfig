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
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

// testEqual simply tests that two interfaces are equal and logs the appropriate message to the provided test state.
func testEqual(got interface{}, want interface{}, message string, t *testing.T) {
	gotJSON, _ := yaml.Marshal(got)
	wantJSON, _ := yaml.Marshal(want)
	prettyGotWantMessage := fmt.Sprintf("%s\ngot:\n%s\nwant:\n%s\n", message, gotJSON, wantJSON)

	if reflect.DeepEqual(got, want) {
		t.Logf("SUCCESS: %s", prettyGotWantMessage)
	} else {
		t.Errorf("FAIL: %s", prettyGotWantMessage)
	}
}
