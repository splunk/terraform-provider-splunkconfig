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

import (
	"reflect"
	"sort"
)

// Roles is a list of Role objects.
type Roles []Role

// validate returns an error if any Roles' members are invalid or duplicates.
func (roles Roles) validate() error {
	return allValidNoDuplicates(uniqueValidators(roles))
}

// validateForLookups returns an error if any of Roles' members reference a Lookup name not present in Lookups.
func (roles Roles) validateForLookups(lookups Lookups) error {
	for _, role := range roles {
		if err := role.validateForLookups(lookups); err != nil {
			return err
		}
	}

	return nil
}

// validateForSAMLGroups returns an error if any of its members reference a SAMLGroup not present in SAMLGroups.
func (roles Roles) validateForSAMLGroups(samlGroups SAMLGroups) error {
	for _, role := range roles {
		if err := role.validateForSAMLGroups(samlGroups); err != nil {
			return err
		}
	}

	return nil
}

// roleNameExists returns true if the given RoleName is present in Roles.
func (roles Roles) roleNameExists(roleName RoleName) bool {
	return hasUID(roles, roleName)
}

// extrapolateWithIndexes returns a new Roles object that incorporates appropriate changes from Indexes.
func (roles Roles) extrapolateWithIndexes(indexes Indexes) Roles {
	extrapolatedRoles := make(Roles, len(roles))

	for i, role := range roles {
		extrapolatedRoles[i] = role.extrapolateFromIndexes(indexes)
	}

	return extrapolatedRoles
}

// rolesWithSAMLGroup returns a new Roles object that includes members of the original roles that have the given
// samlGroup.
func (roles Roles) rolesWithSAMLGroup(samlGroup SAMLGroup) Roles {
	foundRoles := Roles{}

	for _, role := range roles {
		if role.hasSAMLGroup(samlGroup) {
			foundRoles = append(foundRoles, role)
		}
	}

	return foundRoles
}

// roleNamesWithSAMLGroup returns a RoleNames object with the names of the member roles that have the given samlGroup.
func (roles Roles) roleNamesWithSAMLGroup(samlGroup SAMLGroup) RoleNames {
	return roles.rolesWithSAMLGroup(samlGroup).RoleNames()
}

// RoleNames returns a RoleNames object for the names of each Role in Roles, sorted by RoleName.
func (roles Roles) RoleNames() RoleNames {
	uids := uidsOfUIDers(roles)
	sort.Strings(uids)
	return NewRoleNamesFromStrings(uids)
}

// WithRoleName returns the Role object with the given RoleName. Returns ok=false if not found.
func (roles Roles) WithRoleName(roleName RoleName) (found Role, ok bool) {
	foundUIDer, ok := withUID(roles, roleName.uid())
	if !ok {
		return
	}

	foundValue := reflect.ValueOf(foundUIDer)
	found = foundValue.Interface().(Role)

	return
}

// lookupRowsForLookup returns this Roles' LookupRows for the given Lookup.
func (roles Roles) lookupRowsForLookup(lookup Lookup) LookupRows {
	definers := lookupRowsForLookupDefiners(roles)

	return lookupRowsForLookup(lookup, definers)
}

// stanzas returns the Stanzas for Roles.
func (roles Roles) stanzas() Stanzas {
	stanzas := make(Stanzas, len(roles))
	for i, role := range roles {
		stanzas[i] = role.stanza()
	}

	return stanzas
}

// confFile returns the ConfFile for Roles.
func (roles Roles) confFile() ConfFile {
	return ConfFile{
		Name:    "authorize",
		Stanzas: roles.stanzas(),
	}
}
