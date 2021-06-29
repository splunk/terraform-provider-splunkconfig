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

// contentTemplaterTestCase defines a test for a ContentTemplater object.
type contentTemplaterTestCase struct {
	input ContentTemplater
	want  string
}

// test performs the test defined in the contentTemplaterTestCase.
func (c contentTemplaterTestCase) test(t *testing.T) {
	got := c.input.TemplatedContent()
	message := fmt.Sprintf("%T{%+v}.TemplatedContent()\ngot:\n%s\nwant:\n%s\n", c.input, c.input, got, c.want)
	testEqual(got, c.want, message, t)
}
