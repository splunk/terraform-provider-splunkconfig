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

import "fmt"

// StanzaValues is a map of key and value pairs under a stanza.
type StanzaValues map[string]string

// hasKey returns true if StanzaValues has an entry for a given key.
func (stanzaValues StanzaValues) hasKey(key string) bool {
	_, hasKey := stanzaValues[key]

	return hasKey
}

// validateNoCollisions returns an error if stanzaValues and otherStanzaValues have entries for any of the same keys.
func (stanzaValues StanzaValues) validateNoCollisions(otherStanzaValues StanzaValues) error {
	for key := range stanzaValues {
		if otherStanzaValues.hasKey(key) {
			return fmt.Errorf("StanzaValues conflict: key %q present in both %+v and %+v", key, stanzaValues, otherStanzaValues)
		}
	}

	return nil
}

// templateString returns a template string to use to template
func (stanzaValues StanzaValues) templateString() string {
	return `{{ range $key, $value := . -}}
{{ $key }} = {{ $value }}
{{ end }}`
}

// TemplatedContent returns the templated StanzaValues content.
func (stanzaValues StanzaValues) TemplatedContent() string {
	return templateContent(stanzaValues)
}
