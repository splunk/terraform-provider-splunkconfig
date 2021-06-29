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
	"regexp"
)

// IndexName represents the name of a Splunk index.
type IndexName string

// validate() returns an error if IndexName is invalid.
func (indexName IndexName) validate() error {
	// from the documentation:
	//   Index names must consist of only numbers, lowercase letters, underscores,
	//   and hyphens. They cannot begin with an underscore or hyphen, or contain
	//   the word "kvstore".

	kvstoreRegex := regexp.MustCompile("kvstore")

	// NOTE: we allow leading underscores, because not all Index objects are new indexes, some describe defaults
	generalRegex := regexp.MustCompile("^[a-z0-9_][a-z0-9_-]+$")

	if kvstoreRegex.MatchString(string(indexName)) {
		return fmt.Errorf("index name (%s) can not contain the word kvstore", indexName)
	}

	if !generalRegex.MatchString(string(indexName)) {
		return fmt.Errorf("index name (%s) must consist of only numbers, lowercase letters, underscores, and hyphens, and must not beigin with a hyphen", indexName)
	}

	return nil
}

// uid returns the string value of an IndexName to determine uniqueness.
func (indexName IndexName) uid() string {
	return string(indexName)
}
