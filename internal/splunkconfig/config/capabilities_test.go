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

func TestCapabilities_StanzaValues(t *testing.T) {
	tests := []struct {
		input Capabilities
		want  StanzaValues
	}{
		{
			Capabilities{"enabled_capability": true, "disabled_capability": false},
			StanzaValues{"enabled_capability": "enabled", "disabled_capability": "disabled"},
		},
	}

	for _, test := range tests {
		got := test.input.StanzaValues()
		message := fmt.Sprintf("%#v.StanzaValues()", test.input)
		testEqual(got, test.want, message, t)
	}
}

func TestCapabilities_CapabilityNamesByState(t *testing.T) {
	tests := []struct {
		input        Capabilities
		wantEnabled  CapabilityNames
		wantDisabled CapabilityNames
	}{
		{
			Capabilities{
				"enabled_capability":  true,
				"disabled_capability": false,
			},
			CapabilityNames{"enabled_capability"},
			CapabilityNames{"disabled_capability"},
		},
	}

	for _, test := range tests {
		gotEnabled, gotDisabled := test.input.CapabilityNamesByState()

		messageEnabled := fmt.Sprintf("%#v.CapabilityNamesByState() [enabled]", test.input)
		testEqual(gotEnabled, test.wantEnabled, messageEnabled, t)

		messageDisabled := fmt.Sprintf("#%v.CapabiltiyNamesByState() [disabled]", test.input)
		testEqual(gotDisabled, test.wantDisabled, messageDisabled, t)
	}
}
