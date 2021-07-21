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
)

// uider is an interface for objects that return a unique identifier as a string.
type uider interface {
	uid() string
}

// uidersFromUIDersList returns a list of uniqueIdentifiers from a list of objects that implement the uider interface.
func uidersFromUIDersList(list interface{}) []uider {
	listValue := reflect.ValueOf(list)
	listLength := listValue.Len()
	uniqueIdentifiers := make([]uider, listLength)

	for i := 0; i < listLength; i++ {
		iValue := listValue.Index(i)
		iUniqueIdentifier := iValue.Interface().(uider)
		uniqueIdentifiers[i] = iUniqueIdentifier
	}

	return uniqueIdentifiers
}

// withUID returns a uider from the provided list matching the uid given. If not found, returns ok=false.
func withUID(list interface{}, uid string) (found uider, ok bool) {
	for _, identifierFromList := range uidersFromUIDersList(list) {
		if identifierFromList.uid() == uid {
			return identifierFromList, true
		}
	}

	return
}

// hasUID returns true if identifier exists in list.
func hasUID(list interface{}, identifier uider) bool {
	if _, ok := withUID(list, identifier.uid()); ok {
		return true
	}

	return false
}

// uidsOfUIDers returns a list of unique identifier strings from a list of uniqueIdentifiers.
func uidsOfUIDers(list interface{}) []string {
	identifiersFromList := uidersFromUIDersList(list)
	identifierStrings := make([]string, len(identifiersFromList))

	for i, identifierFromList := range identifiersFromList {
		identifierStrings[i] = identifierFromList.uid()
	}

	return identifierStrings
}

// uniqueUIDsOfUIDers returns a deduplicated list of unique identifier strings from a list of uniqueIdentifiers.
func uniqueUIDsOfUIDers(list interface{}) []string {
	seen := make(map[string]bool)

	for _, seenIdentifier := range uidsOfUIDers(list) {
		seen[seenIdentifier] = true
	}

	uniqueIdentifierStrings := make([]string, len(seen))

	i := 0
	for seenIdentifier := range seen {
		uniqueIdentifierStrings[i] = seenIdentifier
		i++
	}

	return uniqueIdentifierStrings
}
