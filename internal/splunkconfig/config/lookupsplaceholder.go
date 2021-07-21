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

import "fmt"

// LookupsPlaceholder represents a set of Lookups or an intent to import Lookups from elsewhere.
type LookupsPlaceholder struct {
	Lookups Lookups
	Import  []string
}

// validate returns an error if LookupsPlaceholder is invalid.  It is invalid if its Lookups are invalid.
func (lookupsPlaceholder LookupsPlaceholder) validate() error {
	if err := lookupsPlaceholder.Lookups.validate(); err != nil {
		return fmt.Errorf("invalid LookupsPlaceholder: invalid Lookups: %s", err)
	}

	return nil
}

// selectedLookups returns candidateLookups' matching Lookups if LookupsPlaceholder.Import has values.  Otherwise it
// returns LookupsPlaceholder.Lookups.
func (lookupsPlaceholder LookupsPlaceholder) selectedLookups(candidateLookups Lookups) (Lookups, error) {
	if len(lookupsPlaceholder.Import) > 0 {
		foundLookups, err := candidateLookups.WithNames(lookupsPlaceholder.Import)
		if err != nil {
			return Lookups{}, fmt.Errorf("unable to select lookups: %s", err)
		}

		return foundLookups, nil
	}

	return lookupsPlaceholder.Lookups, nil
}

// UnmarshalYAML implements custom unmarshalling for a LookupsPlaceholder.  It enables a LookupsPlaceholder to be
// unmarshalled from these types of content:
// * {lookups: [{name: my_lookup}]    # explicitly define its lookups within a LookupsPlaceholder structure
// * [{name: my_lookup}]              # list lookups directly
// * [my_lookup]                      # configure LookupsPlaceholder to import lookups with a list of names
func (lookupsPlaceholder *LookupsPlaceholder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// realLookupsPlaceholder only exists inside this function, and is used to allow attempting to unmarshal into what
	// is really just a LookupsPlaceholder directly.  Attempting to unmarshal(&LookupsPlaceholder) from inside this
	// function will result in infinite recursion back into this function, so we need another type to attempt that
	// unmarshalling.
	type realLookupsPlaceholder LookupsPlaceholder

	// first try to unmarshal into (effectively) an actual LookupsPlaceholder
	unmarshalledLookupsPlaceholder := realLookupsPlaceholder{}
	if err := unmarshal(&unmarshalledLookupsPlaceholder); err == nil {
		*lookupsPlaceholder = LookupsPlaceholder(unmarshalledLookupsPlaceholder)
		return nil
	}

	// then try to unmarshal into a Lookups object, to be embedded in the placeholder
	unmarshalledLookups := Lookups{}
	if err := unmarshal(&unmarshalledLookups); err == nil {
		lookupsPlaceholderFromLookups := LookupsPlaceholder{
			Lookups: unmarshalledLookups,
		}
		*lookupsPlaceholder = lookupsPlaceholderFromLookups

		return nil
	}

	// and finally try to unmarshal into a list of strings, to be embedded in the placeholder
	unmarshalledLookupsPlaceholderList := []string(nil)
	if err := unmarshal(&unmarshalledLookupsPlaceholderList); err == nil {
		lookupsPlaceHolderFromList := LookupsPlaceholder{
			Import: unmarshalledLookupsPlaceholderList,
		}
		*lookupsPlaceholder = lookupsPlaceHolderFromList

		return nil
	}

	// if none of the above unmarshal attempts succeed, return an error
	return fmt.Errorf("unable to unmarshall LookupsPlaceholder from YAML")
}
