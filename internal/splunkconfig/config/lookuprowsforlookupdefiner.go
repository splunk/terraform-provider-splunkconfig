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

import "reflect"

// lookupRowsForLookupDefiner objects implement lookupRowsForLookup(lookup) to define its rows for a given lookup.
type lookupRowsForLookupDefiner interface {
	lookupRowsForLookup(Lookup) LookupRows
}

// lookupRowsForLookupDefiners returns a list of lookupRowsForLookupDefiner objects from a list of objects that adhere
// to the lookupRowsForLookupDefiner interface.
func lookupRowsForLookupDefiners(list interface{}) []lookupRowsForLookupDefiner {
	listValue := reflect.ValueOf(list)

	listLength := listValue.Len()
	lookupRowsForLookupDefiners := make([]lookupRowsForLookupDefiner, listLength)

	for i := range lookupRowsForLookupDefiners {
		iValue := listValue.Index(i)

		lookupRowsForLookupDefiners[i] = iValue.Interface().(lookupRowsForLookupDefiner)
	}

	return lookupRowsForLookupDefiners
}

// lookupRowsForLookup returns the LookupRows for the given Lookup from a list of lookupRowsForLookupDefiners.
func lookupRowsForLookup(lookup Lookup, definers []lookupRowsForLookupDefiner) LookupRows {
	rowsForLookup := LookupRows{}

	for _, definer := range definers {
		rowsForLookup = append(rowsForLookup, definer.lookupRowsForLookup(lookup)...)
	}

	return rowsForLookup
}
