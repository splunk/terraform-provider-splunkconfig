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

import "fmt"

// LookupRows is a list of LookupRow items.
type LookupRows []LookupRow

// validateForLookupFields returns an error if any LookupRows items are invalid for a given LookupFields.
func (lookupRows LookupRows) validateForLookupFields(lookupFields LookupFields) error {
	for _, lookupRow := range lookupRows {
		if err := lookupRow.validateForLookupFields(lookupFields); err != nil {
			return fmt.Errorf("invalid LookupRows, has invalid invalid row: %s", err)
		}
	}

	return nil
}

// validateForLookups returns an error if any LookupRows items reference Lookup names that aren't part of the provided
// Lookup.
func (lookupRows LookupRows) validateForLookups(lookups Lookups) error {
	for _, lookupRow := range lookupRows {
		if err := lookupRow.validateWithLookups(lookups); err != nil {
			return err
		}
	}

	return nil
}

// forLookup returns a new LookupRows object containing the rows from this LookupRows that have the same name
// as the passed Lookup.  Applies defaults from the passed defaultLookupValuesDefiner.
func (lookupRows LookupRows) forLookup(lookup Lookup, definer defaultLookupValuesDefiner) LookupRows {
	rowsForLookup := LookupRows{}

	for _, lookupRow := range lookupRows {
		if lookupRow.LookupName == lookup.Name {
			rowsForLookup = append(rowsForLookup, lookupRow.withValuesForLookupFromDefaultLookupValuesDefiner(lookup, definer))
		}
	}

	return rowsForLookup
}
