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

// Stanzas is a list of Stanza objects.
type Stanzas []Stanza

// validate returns an error if Stanzas has any collisions within its Stanza members.
func (stanzas Stanzas) validate() error {
	// an empty or single-value Stanzas is valid, and needs to avoid recursion
	if len(stanzas) <= 1 {
		return nil
	}

	firstStanza := stanzas[0]
	remainingStanzas := stanzas[1:]

	if err := remainingStanzas.validateNoCollisionsWithStanza(firstStanza); err != nil {
		return err
	}

	return remainingStanzas.validate()
}

// validateNoCollisionsWithStanza returns an error if Stanzas has a collision with a single Stanza.
func (stanzas Stanzas) validateNoCollisionsWithStanza(otherStanza Stanza) error {
	for _, stanza := range stanzas {
		if err := stanza.validateNoCollisions(otherStanza); err != nil {
			return fmt.Errorf(
				"%T{%+v} has collision with %T{%+v}: %s",
				stanzas,
				stanzas,
				otherStanza,
				otherStanza,
				err)
		}
	}

	return nil
}

// validateNoCollisions returns an error if Stanzas has a collision with other Stanzas.
func (stanzas Stanzas) validateNoCollisions(otherStanzas Stanzas) error {
	for _, otherStanza := range otherStanzas {
		if err := stanzas.validateNoCollisionsWithStanza(otherStanza); err != nil {
			return fmt.Errorf(
				"%T{%+v} has collision with %T{%+v}: %s",
				stanzas,
				stanzas,
				otherStanzas,
				otherStanzas,
				err)
		}
	}

	return nil
}

// templateString returns a template string to use to template
func (stanzas Stanzas) templateString() string {
	return `{{ range . -}}
{{ .TemplatedContent }}
{{ end }}`
}

// TemplatedContent returns the templated Stanzas content.
func (stanzas Stanzas) TemplatedContent() string {
	return templateContent(stanzas)
}
