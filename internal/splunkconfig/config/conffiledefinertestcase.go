/*
 * Copyright 2021 Splunk, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"
	"testing"
)

// confFileDefinerTestCase describes a test case for a confFileDefiner.
type confFileDefinerTestCase struct {
	input confFileDefiner
	want  ConfFile
}

// test performs the test described by a confFileDefinerTestCase.
func (testCase confFileDefinerTestCase) test(t *testing.T) {
	got := testCase.input.confFile()
	message := fmt.Sprintf("%#v.confFile()", testCase.input)

	testEqual(got, testCase.want, message, t)
}
