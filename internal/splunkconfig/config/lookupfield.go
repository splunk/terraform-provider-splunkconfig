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

// LookupField is a single field of a lookup.
type LookupField struct {
	Name            string
	Required        bool
	DefaultRowField bool `yaml:"default_row_field"`
	Default         string
}

// validate returns an error if LookupField is invalid. It is invalid if it:
// * has an empty name
func (lookupField LookupField) validate() error {
	if lookupField.Name == "" {
		return fmt.Errorf("name of LookupField is empty")
	}

	return nil
}

// uid returns the unique identifier of LookupField.
func (lookupField LookupField) uid() string {
	return lookupField.Name
}
