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
	"sort"
	"strings"
)

// RoleNames represents a list of RoleName objects.
type RoleNames []RoleName

// NewRoleNamesFromStrings creates and returns an RoleNames object from a list of strings.
func NewRoleNamesFromStrings(values []string) RoleNames {
	roleNames := make(RoleNames, len(values))

	for i, value := range values {
		roleNames[i] = RoleName(value)
	}

	return roleNames
}

// validate() returns an error if RoleNames is invalid.  It is invalid if any of its members are invalid, or if there
// are any duplicated role names.
func (roleNames RoleNames) validate() error {
	return allValidNoDuplicates(uniqueValidators(roleNames))
}

// deduplicatedSorted returns a deduplicated and sorted RoleNames from one that potentially has duplication.
func (roleNames RoleNames) deduplicatedSorted() RoleNames {
	deduplicatedNames := uniqueUIDsOfUIDers(roleNames)
	sort.Strings(deduplicatedNames)

	return NewRoleNamesFromStrings(deduplicatedNames)
}

// authorizeConfImportRolesValue returns a string suitable for use in authorize.conf for importRoles.
func (roleNames RoleNames) authorizeConfImportRolesValue() string {
	return strings.Join(uidsOfUIDers(roleNames.deduplicatedSorted()), ";")
}
