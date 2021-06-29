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
	"reflect"
)

// FileContenters is a list of FileContenter objects.
type FileContenters []FileContenter

// WithContent returns a new FileContenters containing the members of the original FileContenters that have non-empty
// content.
func (contenters FileContenters) WithContent() FileContenters {
	contentersWithContent := FileContenters{}

	for _, contenter := range contenters {
		if len(contenter.TemplatedContent()) > 0 {
			contentersWithContent = append(contentersWithContent, contenter)
		}
	}

	return contentersWithContent
}

// NewFileContentersFromList returns FileContenters from a list of objects that implement the FileContenter methods.
func NewFileContentersFromList(list interface{}) FileContenters {
	listValue := reflect.ValueOf(list)
	listLength := listValue.Len()
	newContenters := make(FileContenters, listLength)

	for i := 0; i < listLength; i++ {
		iValue := listValue.Index(i)
		contenter := iValue.Interface().(FileContenter)
		newContenters[i] = contenter
	}

	return newContenters
}
