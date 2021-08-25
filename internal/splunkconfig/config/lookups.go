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
	"reflect"
	"sort"
)

// Lookups is a list of Lookup objects.
type Lookups []Lookup

// extrapolatedWithLookupRowsForLookupDefiners returns a new Lookups object which includes rows for the given
// lookupRowsForLookupDefiners.
func (lookups Lookups) extrapolatedWithLookupRowsForLookupDefiners(definers ...lookupRowsForLookupDefiner) Lookups {
	extrapolatedLookups := make(Lookups, len(lookups))

	for i, lookup := range lookups {
		extrapolatedLookups[i] = lookup.extrapolatedWithLookupRowsForLookupDefiners(definers...)
	}

	return extrapolatedLookups
}

// validate returns an error if any of its member Lookup objects are invalid.
func (lookups Lookups) validate() error {
	for _, lookup := range lookups {
		if err := lookup.validate(); err != nil {
			return err
		}
	}

	return nil
}

// hasLookupName returns true if the given Lookup name is present in any of Lookups' items.
func (lookups Lookups) hasLookupName(lookupName string) bool {
	for _, lookup := range lookups {
		if lookup.Name == lookupName {
			return true
		}
	}

	return false
}

// WithName returns the Lookup object with the given Name. Returns ok=false if not found.
func (lookups Lookups) WithName(name string) (found Lookup, ok bool) {
	foundUIDer, ok := withUID(lookups, name)
	if !ok {
		return
	}

	foundValue := reflect.ValueOf(foundUIDer)
	found = foundValue.Interface().(Lookup)

	return
}

// WithNames returns a new Lookups with Lookup objects that match the names given.  Returns an error if any of the
// names were not found.
func (lookups Lookups) WithNames(names []string) (Lookups, error) {
	foundLookups := make(Lookups, len(names))

	for i, name := range names {
		foundLookup, ok := lookups.WithName(name)
		if !ok {
			return Lookups{}, fmt.Errorf("Lookup with name %s not found", name)
		}
		foundLookups[i] = foundLookup
	}

	return foundLookups, nil
}

// lookupNames returns a list of lookup names for Lookups.
func (lookups Lookups) lookupNames() []string {
	lookupNames := make([]string, len(lookups))

	for i, lookup := range lookups {
		lookupNames[i] = lookup.Name
	}

	sort.Strings(lookupNames)

	return lookupNames
}

// stanzas returns the Stanzas for Lookups.
func (lookups Lookups) stanzas() Stanzas {
	stanzas := make(Stanzas, len(lookups))

	// use lookups.lookupNames() to force sorting
	for i, lookupName := range lookups.lookupNames() {
		found, _ := lookups.WithName(lookupName)
		stanzas[i] = found.stanza()
	}

	return stanzas
}

// confFile returns the ConfFile for Lookups.
func (lookups Lookups) confFile() ConfFile {
	return ConfFile{
		Name:    "transforms",
		Stanzas: lookups.stanzas(),
	}
}

func (lookups Lookups) fileContenters() FileContenters {
	contenters := NewFileContentersFromList(lookups)
	contenters = append(contenters, lookups.confFile())

	return contenters
}
