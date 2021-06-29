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

// IndexesPlaceholder represents a set of Indexes or an intent to import Indexes from elsewhere.
type IndexesPlaceholder struct {
	Indexes Indexes
	Import  bool
}

// validate returns an error if IndexesPlaceholder is invalid.  It is invalid if its Indexes are invalid.
func (indexesPlaceholder IndexesPlaceholder) validate() error {
	if err := indexesPlaceholder.Indexes.validate(); err != nil {
		return fmt.Errorf("invalid IndexesPlaceholder: invalid Indexes: %s", err)
	}

	return nil
}

// selectedIndexes returns the candidateIndexes if IndexesPlaceholder.Import is true.  Otherwise it returns
// IndexesPlacholder.Indexes.
func (indexesPlaceholder IndexesPlaceholder) selectedIndexes(candidateIndexes Indexes) Indexes {
	if indexesPlaceholder.Import {
		return candidateIndexes
	}

	return indexesPlaceholder.Indexes
}

// UnmarshalYAML implements custom unmarshalling for an IndexesPlaceholder.  It enables an IndexesPlaceholder to be
// unmarshalled from these types of content:
// * {indexes: [{name: my_index}]     # explicitly define its indexes within an IndexesPlaceholder structure
// * [{name: my_index}]               # list indexes directly
// * true                             # configure IndexesPlaceholder to import indexes instead
func (indexesPlaceholder *IndexesPlaceholder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// realIndexesPlaceholder only exists inside this function, and is used to allow attempting to unmarshal into what
	// is really just an IndexesPlaceholder directly.  Attempting to unmarshal(&IndexesPlaceholder) from inside this
	// function will result in infinite recursion back into this function, so we need another type to attempt that
	// unmarshalling.
	type realIndexesPlaceholder IndexesPlaceholder

	// first try to unmarshal into (effectively) an actual IndexesPlaceholder
	unmarshalledIndexesPlaceholder := realIndexesPlaceholder{}
	if err := unmarshal(&unmarshalledIndexesPlaceholder); err == nil {
		*indexesPlaceholder = IndexesPlaceholder(unmarshalledIndexesPlaceholder)
		return nil
	}

	// then try to unmarshal into an Indexes object, to be embedded in the placeholder
	unmarshalledIndexes := Indexes{}
	if err := unmarshal(&unmarshalledIndexes); err == nil {
		indexesPlaceholderFromIndexes := IndexesPlaceholder{
			Indexes: unmarshalledIndexes,
		}
		*indexesPlaceholder = indexesPlaceholderFromIndexes

		return nil
	}

	// and finally try to unmarshal into a boolean, to be embedded in the placeholder
	unmarshalledIndexesPlaceholderBool := false
	if err := unmarshal(&unmarshalledIndexesPlaceholderBool); err == nil {
		indexesPlaceHolderFromBool := IndexesPlaceholder{
			Import: unmarshalledIndexesPlaceholderBool,
		}
		*indexesPlaceholder = indexesPlaceHolderFromBool

		return nil
	}

	// if none of the above unmarshal attempts succeed, return an error
	return fmt.Errorf("unable to unmarshall IndexesPlaceholder from YAML")
}
