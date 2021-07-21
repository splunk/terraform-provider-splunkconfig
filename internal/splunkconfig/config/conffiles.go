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

// ConfFiles represents a list of ConfFile object.
type ConfFiles []ConfFile

// validate returns an error if ConfFiles is invalid.  It is invalid if:
// * any of its ConfFile members are invalid
// * it has any Name collisions among its members
func (confFiles ConfFiles) validate() error {
	return allValidNoDuplicates(uniqueValidators(confFiles))
}

// WithConfFile returns a new ConfFiles object with additionalConfFile added to it or merged with an existing ConfFile
// if one exists with the same name.
func (confFiles ConfFiles) WithConfFile(additionalConfFile ConfFile) ConfFiles {
	newConfFiles := confFiles

	for i, confFile := range newConfFiles {
		if confFile.Name == additionalConfFile.Name {
			newConfFiles[i] = confFile.WithStanzas(additionalConfFile.Stanzas)
			// we've made our substitution, so we should just return now
			return newConfFiles
		}
	}

	// no substitution was made, so just add it to the list
	newConfFiles = append(newConfFiles, additionalConfFile)
	return newConfFiles
}
