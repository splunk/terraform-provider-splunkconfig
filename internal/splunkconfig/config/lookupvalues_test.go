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
	"fmt"
	"testing"
)

func TestLookupValues_validateForLookupFields(t *testing.T) {
	tests := []struct {
		lookupValues LookupValues
		lookupFields LookupFields
		wantError    bool
	}{
		// empty values
		{
			LookupValues{},
			LookupFields{LookupField{Name: "fieldA"}},
			true,
		},
		// set only known fields
		{
			LookupValues{"fieldA": "valueA"},
			LookupFields{LookupField{Name: "fieldA"}},
			false,
		},
		// set unknown field
		{
			LookupValues{"fieldA": "valueA"},
			LookupFields{LookupField{Name: "fieldB"}},
			true,
		},
		// don't set a required field
		{
			LookupValues{"fieldA": "valueA"},
			LookupFields{LookupField{Name: "fieldA"}, LookupField{Name: "fieldB", Required: true}},
			true,
		},
		// explicitly set a required field to an empty string
		{
			LookupValues{"fieldA": ""},
			LookupFields{LookupField{Name: "fieldA", Required: true}},
			false,
		},
	}

	for _, test := range tests {
		gotError := test.lookupValues.validateForLookupFields(test.lookupFields) != nil
		message := fmt.Sprintf("%T{%+v}.validateForLookupFields(%T{%+v})", test.lookupValues, test.lookupValues, test.lookupFields, test.lookupFields)
		testEqual(gotError, test.wantError, message, t)
	}
}

func TestLookupValues_withDefaultLookupValues(t *testing.T) {
	tests := []struct {
		specifics LookupValues
		defaults  LookupValues
		want      LookupValues
	}{
		// no defaults
		{
			LookupValues{"fieldA": "specificValueA"},
			LookupValues{},
			LookupValues{"fieldA": "specificValueA"},
		},
		// default doesn't override specific
		{
			LookupValues{"fieldA": "specificValueA"},
			LookupValues{"fieldA": "defaultValueA"},
			LookupValues{"fieldA": "specificValueA"},
		},
		// default added if not in specific
		{
			LookupValues{"fieldA": "specificValueA"},
			LookupValues{"fieldB": "defaultValueB"},
			LookupValues{"fieldA": "specificValueA", "fieldB": "defaultValueB"},
		},
	}

	for _, test := range tests {
		got := test.specifics.withDefaultLookupValues(test.defaults)
		message := fmt.Sprintf("%T{%+v}.withDefaultLookupValues(%+v)", test.specifics, test.specifics, test.defaults)
		testEqual(got, test.want, message, t)
	}
}

func TestLookupValues_hasFieldNames(t *testing.T) {
	lookupValues := LookupValues{"fieldA": "valueA", "fieldB": "valueB"}
	tests := []struct {
		lookupValues LookupValues
		fieldNames   []string
		wantBool     bool
	}{
		// all fieldNames present
		{
			lookupValues,
			[]string{"fieldA", "fieldB"},
			true,
		},
		// no fieldNames present
		{
			lookupValues,
			[]string{"fieldC", "fieldD"},
			false,
		},
		// some fieldNames present
		{
			lookupValues,
			[]string{"fieldB", "fieldC"},
			false,
		},
	}

	for _, test := range tests {
		gotBool := test.lookupValues.hasFieldNames(test.fieldNames)
		message := fmt.Sprintf("%T{%+v}.hasFieldNames(%v)", test.lookupValues, test.lookupValues, test.fieldNames)
		testEqual(gotBool, test.wantBool, message, t)
	}
}
