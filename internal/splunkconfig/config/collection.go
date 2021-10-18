// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import "fmt"

// Collection represents a KVStore Collection.
type Collection struct {
	Name         string
	EnforceTypes bool `yaml:"enforceTypes"`
	Fields       CollectionFields
	Replicate    bool
}

// validate returns an error if Collection is invalid. It is invalid if it
// has invalid:
//   * Name
//   * Fields
func (collection Collection) validate() error {
	if len(collection.Name) == 0 {
		return fmt.Errorf("Collection name can not be empty")
	}

	if err := collection.Fields.validate(); err != nil {
		return err
	}

	return nil
}

// stanzaName returns the stanza's name for a Collection.
func (collection Collection) stanzaName() string {
	return string(collection.Name)
}

// stanzaValues returns the StanzaValues for a Collection.
func (collection Collection) stanzaValues() StanzaValues {
	stanzaValues := StanzaValues{}

	if collection.EnforceTypes {
		stanzaValues["enforceTypes"] = fmt.Sprintf("%v", collection.EnforceTypes)
	}

	if collection.Replicate {
		stanzaValues["replicate"] = fmt.Sprintf("%v", collection.Replicate)
	}

	for fieldName, fieldType := range collection.Fields {
		fieldKey := fmt.Sprintf("field.%s", fieldName)
		stanzaValues[fieldKey] = string(fieldType)
	}

	return stanzaValues
}

// stanza returns the Stanza for a Collection.
func (collection Collection) stanza() Stanza {
	return Stanza{
		Name:   collection.stanzaName(),
		Values: collection.stanzaValues(),
	}
}
