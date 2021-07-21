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

func TestFirstIndexPath(t *testing.T) {
	tests := []struct {
		input         []IndexPath
		wantIndexPath IndexPath
		wantOk        bool
	}{
		{
			[]IndexPath{""},
			IndexPath(""),
			false,
		},
		{
			[]IndexPath{"first", "", "third"},
			IndexPath("first"),
			true,
		},
		{
			[]IndexPath{"", "second", "third"},
			IndexPath("second"),
			true,
		},
		{
			[]IndexPath{"", "", "third"},
			IndexPath("third"),
			true,
		},
	}

	for _, test := range tests {
		gotIndexPath, gotOk := firstIndexPath(test.input...)
		messageIndexPath := fmt.Sprintf("firstIndexPath(%#v)", test.input)
		messageOk := fmt.Sprintf("firstIndexPath(%#v) returned error?", test.input)

		testEqual(gotIndexPath, test.wantIndexPath, messageIndexPath, t)
		testEqual(gotOk, test.wantOk, messageOk, t)
	}
}
