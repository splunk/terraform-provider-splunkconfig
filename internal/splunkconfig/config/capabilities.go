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

// Capabilities is a map of CapabilityNames to a boolean state indicating if it is enabled or not.
type Capabilities map[CapabilityName]bool

func (capabilities Capabilities) validate() error {
	for capabilityName := range capabilities {
		if err := capabilityName.validate(); err != nil {
			return err
		}
	}

	return nil
}

// StanzaValues returns the StanzaValues for Capabilities.
func (capabilities Capabilities) StanzaValues() StanzaValues {
	stanzaValues := StanzaValues{}

	for name, state := range capabilities {
		if state {
			stanzaValues[string(name)] = "enabled"
		} else {
			stanzaValues[string(name)] = "disabled"
		}
	}

	return stanzaValues
}

// CapabilityNamesByState returns CapabilityNames, separately, for enabled and disabled capabilities.
func (capabilities Capabilities) CapabilityNamesByState() (enabled CapabilityNames, disabled CapabilityNames) {
	for name, state := range capabilities {
		if state {
			enabled = append(enabled, name)
		} else {
			disabled = append(disabled, name)
		}
	}

	return enabled, disabled
}

// EnabledCapabilityNames returns CapabilityNames that are enabled for Capabilities.
func (capabilities Capabilities) EnabledCapabilityNames() CapabilityNames {
	enabled, _ := capabilities.CapabilityNamesByState()

	return enabled
}

// DisabledCapabilityNames returns CapabilityNames that are disabled for Capabilities.
func (capabilities Capabilities) DisabledCapabilityNames() CapabilityNames {
	_, disabled := capabilities.CapabilityNamesByState()

	return disabled
}
