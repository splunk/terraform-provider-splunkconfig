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

// defaultLookupValuesDefiner objects implement defaultLookupValues which returns LookupValues
// for its default values.
type defaultLookupValuesDefiner interface {
	defaultLookupValues() LookupValues
}

// rowsForLookupOrDefaultRows returns a LookupRows object made up of the LookupRow items in LookupRows that belong
// to the given Lookup.  If the resulting LookupRows is empty, it instead returns the default LookupRows for the
// given Lookup and defaultLookupValuesDefiner.
func rowsForLookupOrDefaultRows(rows LookupRows, lookup Lookup, definer defaultLookupValuesDefiner) LookupRows {
	rowsForLookup := rows.forLookup(lookup, definer)

	if len(rowsForLookup) > 0 {
		return rowsForLookup
	}

	return lookup.defaultRows(definer)
}

// defaultLookupValuesForLookupFields return a LookupValues object for the proper default values defined by a
// defaultLookupValuesDefiner object for a specific LookupFields.
func defaultLookupValuesForLookupFields(lookupFields LookupFields, definer defaultLookupValuesDefiner) LookupValues {
	lookupValues := LookupValues{}

	for defaultLookupKey, defaultLookupValue := range definer.defaultLookupValues() {
		if lookupFields.hasFieldName(defaultLookupKey) {
			lookupValues[defaultLookupKey] = defaultLookupValue
		}
	}

	return lookupValues
}

// defaultLookupValuesForLookup returns a LookupValues object for the proper default values defined by a
// defaultLookupValuesDefiner object for a specific Lookup.
func defaultLookupValuesForLookup(lookup Lookup, definer defaultLookupValuesDefiner) LookupValues {
	return defaultLookupValuesForLookupFields(lookup.Fields, definer)
}
