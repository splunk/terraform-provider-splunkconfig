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

// SAMLGroup represents a SAML group.
type SAMLGroup struct {
	Name  string
	Roles RoleNames
}

// validate returns an error if the SAMLGroup is invalid.  It is invalid if:
// * its Roles object is invalid
func (samlGroup SAMLGroup) validate() error {
	if err := samlGroup.Roles.validate(); err != nil {
		return err
	}

	return nil
}

// uid returns the name of samlGroup to be used as a unique identifier.
func (samlGroup SAMLGroup) uid() string {
	return samlGroup.Name
}

// extrapolateFromRoles returns a new SAMLGroup that incorporates the appropriate Role members from the passed Roles
// object.
func (samlGroup SAMLGroup) extrapolateFromRoles(roles Roles) SAMLGroup {
	roleNames := append(samlGroup.Roles, roles.roleNamesWithSAMLGroup(samlGroup)...)
	samlGroup.Roles = roleNames.deduplicatedSorted()

	return samlGroup
}
