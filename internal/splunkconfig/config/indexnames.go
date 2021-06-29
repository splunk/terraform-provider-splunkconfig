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
	"sort"
	"strings"
)

// IndexNames represents a list of IndexName objects.
type IndexNames []IndexName

// NewIndexNamesFromStrings creates and returns an IndexNames object from a list of strings.
func NewIndexNamesFromStrings(values []string) IndexNames {
	indexNames := make(IndexNames, len(values))

	for i, value := range values {
		indexNames[i] = IndexName(value)
	}

	return indexNames
}

// validate returns an error if IndexNames is invalid.  It is invalid if any of its members are invalid, or if there
// are any duplicated index names.
func (indexNames IndexNames) validate() error {
	return allValidNoDuplicates(uniqueValidators(indexNames))
}

// deduplicatedSorted returns a deduplicated and sorted IndexNames from one that potentiall has duplication.
func (indexNames IndexNames) deduplicatedSorted() IndexNames {
	deduplicatedNames := uniqueUIDsOfUIDers(indexNames)
	sort.Strings(deduplicatedNames)

	return NewIndexNamesFromStrings(deduplicatedNames)
}

// authorizeConfSrchIndexesAllowedValue returns a string suitable for use in authorize.conf for IndexNames.
func (indexNames IndexNames) authorizeConfSrchIndexesAllowedValue() string {
	return strings.Join(uidsOfUIDers(indexNames.deduplicatedSorted()), ";")
}
