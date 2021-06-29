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

// Test that the value returned by InSeconds is correct.
func TestTimePeriod_InSeconds(t *testing.T) {
	// These tests are currently all just different ways to represent one full day (86400 seconds), but "want"
	// is included per-test to avoid refactoring if we need to add non-day tests later.
	tests := []struct {
		input TimePeriod
		want  int64
	}{
		{
			TimePeriod{Seconds: 86400},
			86400,
		},
		{
			TimePeriod{Minutes: 1440},
			86400,
		},
		{
			TimePeriod{Hours: 24},
			86400,
		},
		{
			TimePeriod{Days: 1},
			86400,
		},
		{
			TimePeriod{Hours: 23, Minutes: 59, Seconds: 60},
			86400,
		},
	}

	for _, test := range tests {
		got := test.input.InSeconds()
		message := fmt.Sprintf("TimePeriod%+v.InSeconds()", test.input)

		testEqual(got, test.want, message, t)
	}
}
