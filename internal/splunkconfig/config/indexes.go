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
	"reflect"
	"sort"
)

// Indexes is a list of Index objects.
type Indexes []Index

// validate returns an error if Indexes is invalid.
func (indexes Indexes) validate() error {
	return allValidNoDuplicates(uniqueValidators(indexes))
}

// validateWithRoles returns an error if any Index in Indexes references a RoleName not present in Roles.
func (indexes Indexes) validateWithRoles(roles Roles) error {
	for _, index := range indexes {
		if err := index.validateWithRoles(roles); err != nil {
			return err
		}
	}

	return nil
}

// validateWithLookups returns an error if any Index in Indexes references a Lookup name not present in Lookups.
func (indexes Indexes) validateWithLookups(lookups Lookups) error {
	for _, index := range indexes {
		if err := index.validateWithLookups(lookups); err != nil {
			return err
		}
	}

	return nil
}

// indexesSearchableByRoleName returns Indexes that are searchable by the provided RoleName.
func (indexes Indexes) indexesSearchableByRoleName(roleName RoleName) Indexes {
	var searchableIndexes Indexes

	for _, index := range indexes {
		if index.searchableByRoleName(roleName) {
			searchableIndexes = append(searchableIndexes, index)
		}
	}

	return searchableIndexes
}

// IndexNames returns IndexNames for Indexes.
func (indexes Indexes) IndexNames() IndexNames {
	uids := uidsOfUIDers(indexes)
	sort.Strings(uids)

	return NewIndexNamesFromStrings(uids)
}

// WithIndexName returns the Index object with the given IndexName. Returns ok=false if not found.
func (indexes Indexes) WithIndexName(indexName IndexName) (found Index, ok bool) {
	foundUIDer, ok := withUID(indexes, indexName.uid())
	if !ok {
		return
	}

	foundValue := reflect.ValueOf(foundUIDer)
	found = foundValue.Interface().(Index)

	return
}

// indexNamesSearchableByRoleName returns IndexNames that are searchable by the provided RoleName.
func (indexes Indexes) indexNamesSearchableByRoleName(roleName RoleName) IndexNames {
	return indexes.indexesSearchableByRoleName(roleName).IndexNames()
}

// indexNamesSearchableByRole returns IndexNames that are searchable by the provided Role.
func (indexes Indexes) indexNamesSearchableByRole(role Role) IndexNames {
	return indexes.indexNamesSearchableByRoleName(role.Name)
}

// lookupRowsForLookup returns this Indexes' LookupRows for the given Lookup.
func (indexes Indexes) lookupRowsForLookup(lookup Lookup) LookupRows {
	definers := lookupRowsForLookupDefiners(indexes)

	return lookupRowsForLookup(lookup, definers)
}

// stanzas returns the Stanzas for Indexes.
func (indexes Indexes) stanzas() Stanzas {
	stanzas := make(Stanzas, len(indexes))

	// use indexes.indexNames() to force sorting
	for i, indexName := range indexes.IndexNames() {
		found, _ := indexes.WithIndexName(indexName)
		stanzas[i] = found.stanza()
	}

	return stanzas
}

// confFile returns the ConfFile for Indexes.
func (indexes Indexes) confFile() ConfFile {
	return ConfFile{
		Name:    "indexes",
		Stanzas: indexes.stanzas(),
	}
}
