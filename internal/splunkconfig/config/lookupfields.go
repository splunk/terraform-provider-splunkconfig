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

// LookupFields is a list of LookupField objects that make up the set of fields in a lookup.
type LookupFields []LookupField

// validate returns an error if LookupFields is invalid. It is invalid if it:
// * is empty
// * contains any duplicated field names
func (lookupFields LookupFields) validate() error {
	if len(lookupFields) == 0 {
		return fmt.Errorf("LookupFields is empty")
	}

	if err := allValidNoDuplicates(uniqueValidators(lookupFields)); err != nil {
		return err
	}

	return nil
}

// hasFieldName returns true if LookupFields contains a LookupField named fieldName.
func (lookupFields LookupFields) hasFieldName(fieldName string) bool {
	for _, lookupField := range lookupFields {
		if lookupField.Name == fieldName {
			return true
		}
	}

	return false
}

// headerValues returns a list of strings for the header of the CSV file.
func (lookupFields LookupFields) headerValues() []string {
	headerValues := make([]string, len(lookupFields))

	for i, lookupField := range lookupFields {
		headerValues[i] = lookupField.Name
	}

	return headerValues
}

// defaultLookupValues returns LookupValues for the default values of each LookupField.
func (lookupFields LookupFields) defaultLookupValues() LookupValues {
	defaults := LookupValues{}

	for _, lookupField := range lookupFields {
		if lookupField.Default != "" {
			defaults[lookupField.Name] = lookupField.Default
		}
	}

	return defaults
}

// defaultLookupValuesDefinerDefaultLookupValues returns the default LookupValues from a defaultLookupValuesDefiner
// specific to the given LookupFields.
func (lookupFields LookupFields) defaultLookupValuesDefinerValues(definer defaultLookupValuesDefiner) LookupValues {
	// defaults from lookupFields itself
	defaults := lookupFields.defaultLookupValues()

	// override values explicitly set by definer
	for definerFieldName, definerFieldValue := range definer.defaultLookupValues() {
		// but only if lookupFields has a field with that name
		if lookupFields.hasFieldName(definerFieldName) {
			defaults[definerFieldName] = definerFieldValue
		}
	}

	return defaults
}

// hasRequiredFields returns true if any of the fields in lookupFields are required.
func (lookupFields LookupFields) hasRequiredFields() bool {
	for _, lookupField := range lookupFields {
		if lookupField.Required {
			return true
		}
	}

	return false
}

// hasRequiredFields returns true if any of the fields in lookupFields are "default row fields".  This can be used to
// determine if a set of lookupFields should result in default rows being created.
func (lookupFields LookupFields) hasDefaultRowFields() bool {
	for _, lookupField := range lookupFields {
		if lookupField.DefaultRowField {
			return true
		}
	}

	return false
}

// valuesHaveAllDefaultRowFields returns true if the given LookupValues satisfies all members of this LookupFields
// that have DefaultRowField=true.
func (lookupFields LookupFields) valuesHaveAllDefaultRowFields(lookupValues LookupValues) bool {
	for _, lookupField := range lookupFields {
		if lookupField.DefaultRowField {
			if _, ok := lookupValues[lookupField.Name]; !ok {
				return false
			}
		}
	}

	return true
}
