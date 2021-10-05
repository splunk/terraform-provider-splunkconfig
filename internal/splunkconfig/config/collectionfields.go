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

// CollectionFields is a mapping of field names to CollectionFieldType.
type CollectionFields map[string]CollectionFieldType

// validate returns an error if CollectionFields is invalid. It is invalid if
// any member has invalid:
//   * CollectionFieldType
func (collectionFields CollectionFields) validate() error {
	for _, fieldType := range collectionFields {
		if err := fieldType.validate(); err != nil {
			return err
		}
	}

	return nil
}
