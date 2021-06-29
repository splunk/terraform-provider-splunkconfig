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
)

// Role represents a Splunk role
type Role struct {
	Name                        RoleName
	SAMLGroups                  []string   `yaml:"saml_groups"`
	SearchIndexesAllowed        IndexNames `yaml:"srchIndexesAllowed"`
	ImportRoles                 RoleNames  `yaml:"importRoles"`
	Capabilities                Capabilities
	LookupRows                  LookupRows  `yaml:"lookup_rows"`
	SearchFilter                string      `yaml:"srchFilter"`
	SearchTimeWin               ExplicitInt `yaml:"srchTimeWin"`
	SearchDiskQuota             ExplicitInt `yaml:"srchDiskQuota"`
	SearchJobsQuota             ExplicitInt `yaml:"srchJobsQuota"`
	RTSearchJobsQuota           ExplicitInt `yaml:"rtSrchJobsQuota"`
	CumulativeSearchJobsQuota   ExplicitInt `yaml:"cumulativeSrchJobsQuota"`
	CumulativeRTSearchJobsQuota ExplicitInt `yaml:"cumulativeRTSrchJobsQuota"`
}

// validate returns an error if the Role configuration is not valid.
func (r Role) validate() error {
	if err := r.Name.validate(); err != nil {
		return err
	}

	if err := r.SearchIndexesAllowed.validate(); err != nil {
		return err
	}

	if err := r.ImportRoles.validate(); err != nil {
		return err
	}

	if err := r.Capabilities.validate(); err != nil {
		return err
	}

	return nil
}

// validateForLookups returns an error if the Role's LookupRows reference a Lookup name that doesn't exist in Lookups.
func (r Role) validateForLookups(lookups Lookups) error {
	if err := r.LookupRows.validateForLookups(lookups); err != nil {
		return fmt.Errorf("role %s has invalid LookupRows: %s", r.Name, err)
	}

	return nil
}

// validateForSAMLGroups returns an error if the Role's SAMLGroups reference a SAMLGroup that doesn't exist in
// SAMLGroups.
func (r Role) validateForSAMLGroups(samlGroups SAMLGroups) error {
	for _, samlGroupName := range r.SAMLGroups {
		if !samlGroups.hasSAMLGroupName(samlGroupName) {
			return fmt.Errorf("role %s is has invalid SAMLGroup name: %s", r.Name, samlGroupName)
		}
	}

	return nil
}

// uid returns the Role's Name as a string to determine uniqueness.
func (r Role) uid() string {
	return r.Name.uid()
}

// extrapolateFromIndexes returns a Role that incorporates SearchRolesAllowed from Indexes.
func (r Role) extrapolateFromIndexes(indexes Indexes) Role {
	searchIndexesAllowed := append(r.SearchIndexesAllowed, indexes.indexNamesSearchableByRole(r)...)
	r.SearchIndexesAllowed = searchIndexesAllowed.deduplicatedSorted()

	return r
}

// stanzaName returns the Stanza's Name value for a Role.
func (r Role) stanzaName() string {
	return fmt.Sprintf("role_%s", r.Name)
}

// stanzaValues returns the StanzaValues for a Role.
func (r Role) stanzaValues() StanzaValues {
	stanzaValues := StanzaValues{}

	if len(r.SearchIndexesAllowed) > 0 {
		stanzaValues["srchIndexesAllowed"] = r.SearchIndexesAllowed.authorizeConfSrchIndexesAllowedValue()
	}

	if len(r.ImportRoles) > 0 {
		stanzaValues["importRoles"] = r.ImportRoles.authorizeConfImportRolesValue()
	}

	for capabilityName, capabilityValue := range r.Capabilities.StanzaValues() {
		stanzaValues[capabilityName] = capabilityValue
	}

	if r.SearchFilter != "" {
		stanzaValues["srchFilter"] = r.SearchFilter
	}

	if r.SearchTimeWin.Explicit {
		stanzaValues["srchTimeWin"] = fmt.Sprintf("%d", r.SearchTimeWin.Value)
	}

	if r.SearchDiskQuota.Explicit {
		stanzaValues["srchDiskQuota"] = fmt.Sprintf("%d", r.SearchDiskQuota.Value)
	}

	if r.SearchJobsQuota.Explicit {
		stanzaValues["srchJobsQuota"] = fmt.Sprintf("%d", r.SearchJobsQuota.Value)
	}

	if r.RTSearchJobsQuota.Explicit {
		stanzaValues["rtSrchJobsQuota"] = fmt.Sprintf("%d", r.RTSearchJobsQuota.Value)
	}

	if r.CumulativeSearchJobsQuota.Explicit {
		stanzaValues["cumulativeSrchJobsQuota"] = fmt.Sprintf("%d", r.CumulativeSearchJobsQuota.Value)
	}

	if r.CumulativeRTSearchJobsQuota.Explicit {
		stanzaValues["cumulativeRTSrchJobsQuota"] = fmt.Sprintf("%d", r.CumulativeRTSearchJobsQuota.Value)
	}

	return stanzaValues
}

// stanza returns the Stanza for a role.
func (r Role) stanza() Stanza {
	return Stanza{
		r.stanzaName(),
		r.stanzaValues(),
	}
}

// defaultLookupValues returns a LookupValues object representing the defaults for a Role object.
// The defaults for a Role object are:
// * role - Name of the role
func (r Role) defaultLookupValues() LookupValues {
	return LookupValues{
		"role": string(r.Name),
	}
}

// lookupRowsForLookup returns this Role's LookupRows for the given Lookup.
func (r Role) lookupRowsForLookup(lookup Lookup) LookupRows {
	return rowsForLookupOrDefaultRows(r.LookupRows, lookup, r)
}

// EnabledCapabilityNames returns the CapabilityNames that are enabled for this Role.
func (r Role) EnabledCapabilityNames() CapabilityNames {
	return r.Capabilities.EnabledCapabilityNames()
}

// hasSAMLGroupName returns true if samlGroupName is in the role's SAMLGroups.
func (r Role) hasSAMLGroupName(samlGroupName string) bool {
	for _, foundSAMLGroupName := range r.SAMLGroups {
		if foundSAMLGroupName == samlGroupName {
			return true
		}
	}

	return false
}

// hasSAMLGroup returns true if samlGroup's Name is in the role's SAMLGroups.
func (r Role) hasSAMLGroup(samlGroup SAMLGroup) bool {
	return r.hasSAMLGroupName(samlGroup.Name)
}
