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

// Stanza represents a single stanza in a configuration file.
type Stanza struct {
	Name   string
	Values StanzaValues
}

// validateNoCollisions returns an error if a Stanza has collisions with another Stanza.
func (stanza Stanza) validateNoCollisions(otherStanza Stanza) error {
	// no need to compare for differing stanza names
	if stanza.Name != otherStanza.Name {
		return nil
	}

	return stanza.Values.validateNoCollisions(otherStanza.Values)
}

// templateString returns a template string to use to template
func (stanza Stanza) templateString() string {
	return `[{{ .Name }}]
{{ .Values.TemplatedContent }}`
}

// TemplatedContent returns the templated Stanza content.
func (stanza Stanza) TemplatedContent() string {
	return templateContent(stanza)
}
