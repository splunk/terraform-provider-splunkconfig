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

// LookupRow is a map of field names to field values.
type LookupRow struct {
	// LookupName is used to indirectly associate this Row with a LookupDefinition
	LookupName string `yaml:"lookup_name"`
	Values     LookupValues
}

// validateForLookupFields returns an error if LookupRow is not valid in the context of LookupFields.
// A LookupRow has no way to validate itself without such a context.
func (lookupRow LookupRow) validateForLookupFields(lookupFields LookupFields) error {
	if err := lookupRow.Values.validateForLookupFields(lookupFields); err != nil {
		return fmt.Errorf("invalid LookupRow, has invalid Values: %s", err)
	}

	return nil
}

// validateWithLookups returns an error if a LookupRow references a Lookup name not present in Lookups.
func (lookupRow LookupRow) validateWithLookups(lookups Lookups) error {
	if !lookups.hasLookupName(lookupRow.LookupName) {
		return fmt.Errorf("%T{%+v} references unknown lookup %q", lookupRow, lookupRow, lookupRow.LookupName)
	}

	return nil
}

// valuesForLookupFields returns a list of strings representing a LookupRow in the context of a LookupFields object.
func (lookupRow LookupRow) valuesForLookupFields(lookupFields LookupFields) []string {
	return lookupRow.Values.valuesForLookupFields(lookupFields)
}

// withDefaultLookupValues returns a new LookupRow with default values applied from the given LookupValues.
func (lookupRow LookupRow) withDefaultLookupValues(lookupValues LookupValues) LookupRow {
	rowWithDefaults := lookupRow

	rowWithDefaults.Values = rowWithDefaults.Values.withDefaultLookupValues(lookupValues)

	return rowWithDefaults
}

// withValuesFromDefaultLookupValuesDefiner returns a new LookupRow with default values applied from the given
// defaultLookupValuesDefiner.
func (lookupRow LookupRow) withValuesForLookupFromDefaultLookupValuesDefiner(lookup Lookup, definer defaultLookupValuesDefiner) LookupRow {
	return lookupRow.withDefaultLookupValues(defaultLookupValuesForLookup(lookup, definer))
}
