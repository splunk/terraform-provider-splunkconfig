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

// LookupValues is a map of field names to field values for an individual row of a lookup.
type LookupValues map[string]string

// valuesForLookupFields returns a list of strings representing a LookupRow in the context of a LookupFields object.
func (lookupValues LookupValues) valuesForLookupFields(lookupFields LookupFields) []string {
	fieldValues := make([]string, len(lookupFields))

	for i, lookupField := range lookupFields {
		lookupValue, ok := lookupValues[lookupField.Name]
		if !ok {
			lookupValue = ""
		}

		fieldValues[i] = lookupValue
	}

	return fieldValues
}

// validateForLookupFields returns an error if LookupValues is not valid in the context of LookupFields.
// A LookupValues object has no way to validate itself without such a context.
func (lookupValues LookupValues) validateForLookupFields(lookupFields LookupFields) error {
	// check that at least one field is set
	if len(lookupValues) == 0 {
		return fmt.Errorf("LookupValues is empty")
	}

	// check that all set fields are available
	for fieldName := range lookupValues {
		if !lookupFields.hasFieldName(fieldName) {
			return fmt.Errorf("field %q not in LookupFields %v", fieldName, lookupFields)
		}
	}

	// check that all required fields are set
	for _, lookupField := range lookupFields {
		_, lookupFieldSet := lookupValues[lookupField.Name]
		if lookupField.Required && !lookupFieldSet {
			return fmt.Errorf("field %q is required, but not set in LookupValues %v", lookupField.Name, lookupValues)
		}
	}

	return nil
}

// withDefaultLookupValues returns a new LookupValues object with defaults applied.
func (lookupValues LookupValues) withDefaultLookupValues(defaults LookupValues) LookupValues {
	withDefaults := defaults

	for fieldName, fieldValue := range lookupValues {
		withDefaults[fieldName] = fieldValue
	}

	return withDefaults
}

// hasFieldName returns a boolean indicating if fieldName is present.
func (lookupValues LookupValues) hasFieldName(fieldName string) bool {
	if _, ok := lookupValues[fieldName]; ok {
		return true
	}

	return false
}

// hasFieldNames returns a boolean indicating if all fieldNames are present.
func (lookupValues LookupValues) hasFieldNames(fieldNames []string) bool {
	for _, fieldName := range fieldNames {
		if !lookupValues.hasFieldName(fieldName) {
			return false
		}
	}

	return true
}
