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

import "reflect"

// SAMLGroups is a list of SAMLGroup objects.
type SAMLGroups []SAMLGroup

// validate returns an error if SAMLGroups is invalid. It is invalid if any of its members are invalid.
func (samlGroups SAMLGroups) validate() error {
	if err := allValidNoDuplicates(uniqueValidators(samlGroups)); err != nil {
		return err
	}

	return nil
}

// extrapolateWithRoles returns a new SAMLGroups object with its members extrapolated from the passed Roles.
func (samlGroups SAMLGroups) extrapolateWithRoles(roles Roles) SAMLGroups {
	extrapolatedSAMLGroups := make(SAMLGroups, len(samlGroups))

	for i, samlGroup := range samlGroups {
		extrapolatedSAMLGroups[i] = samlGroup.extrapolateFromRoles(roles)
	}

	return extrapolatedSAMLGroups
}

// WithSAMLGroupName returns the SAMLGroup from SAMLGroups that has samlGroupName as its name. If none was found,
// returns ok=false.
func (samlGroups SAMLGroups) WithSAMLGroupName(samlGroupName string) (found SAMLGroup, ok bool) {
	foundUIDer, ok := withUID(samlGroups, samlGroupName)
	if !ok {
		return
	}

	foundValue := reflect.ValueOf(foundUIDer)
	found = foundValue.Interface().(SAMLGroup)

	return
}

// hasSAMLGroupName returns true if samlGroupName is present in samlGroups.
func (samlGroups SAMLGroups) hasSAMLGroupName(samlGroupName string) bool {
	_, found := samlGroups.WithSAMLGroupName(samlGroupName)

	return found
}

// SAMLGroupNames returns a list of names of SAMLGroups' members.
func (samlGroups SAMLGroups) SAMLGroupNames() []string {
	samlGroupNames := make([]string, len(samlGroups))

	for i, samlGroup := range samlGroups {
		samlGroupNames[i] = samlGroup.Name
	}

	return samlGroupNames
}
