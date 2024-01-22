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

// Index represents a single Splunk index.
type Index struct {
	Name                          IndexName
	FrozenTime                    TimePeriod `yaml:"frozenTimePeriod"`
	SearchRolesAllowed            RoleNames  `yaml:"srchRolesAllowed"`
	LookupRows                    LookupRows `yaml:"lookup_rows"`
	HomePath                      IndexPath  `yaml:"homePath"`
	ColdPath                      IndexPath  `yaml:"coldPath"`
	ThawedPath                    IndexPath  `yaml:"thawedPath"`
	DataType                      IndexDataType
	ColdStorageProvider           IndexArchiverProvider `yaml:"coldStorageProvider"`
	ColdStorageRetentionPeriod    TimePeriod            `yaml:"coldStorageRetentionPeriod"`
	EnableDataArchive             bool                  `yaml:"enableDataArchive "`
	MaxDataArchiveRetentionPeriod TimePeriod            `yaml:"maxDataArchiveRetentionPeriod"`
}

// validate returns an error if the Index is invalid.
func (index Index) validate() error {
	if err := index.Name.validate(); err != nil {
		return err
	}

	if err := index.SearchRolesAllowed.validate(); err != nil {
		return err
	}

	if err := index.DataType.validate(); err != nil {
		return err
	}

	return nil
}

// uid returns the name of the Index to determine uniqueness.
func (index Index) uid() string {
	return index.Name.uid()
}

// validateWithRoles returns an error if an Index references a RoleName not present in Roles.
func (index Index) validateWithRoles(roles Roles) error {
	for _, roleName := range index.SearchRolesAllowed {
		if !roles.roleNameExists(roleName) {
			return fmt.Errorf("SearchRolesAllowed RoleName %s not in provided Roles", roleName)
		}
	}

	return nil
}

// validateWithLookups returns an error if an Index references a Lookup name not present in Lookups.
func (index Index) validateWithLookups(lookups Lookups) error {
	if err := index.LookupRows.validateForLookups(lookups); err != nil {
		return fmt.Errorf("index %s has invalid LookupRows: %s", index.Name, err)
	}

	return nil
}

// searchableByRoleName returns true if this Index lists RoleName in SearchRolesAllowed.
func (index Index) searchableByRoleName(roleName RoleName) bool {
	for _, searchRoleAllowed := range index.SearchRolesAllowed {
		if searchRoleAllowed == roleName {
			return true
		}
	}

	return false
}

// stanzaName returns the Stanza's Name for an Index.
func (index Index) stanzaName() string {
	return string(index.Name)
}

// stanzaValues returns the StanzaValues for an Index.
func (index Index) stanzaValues() StanzaValues {
	stanzaValues := StanzaValues{}

	if index.FrozenTime.InSeconds() != 0 {
		stanzaValues["frozenTimePeriodInSecs"] = fmt.Sprintf("%d", index.FrozenTime.InSeconds())
	}

	if index.DataType != INDEXDATATYPEUNDEF {
		stanzaValues["datatype"] = string(index.DataType)
	}

	// defaultIndexPath will return a valid IndexPath, so we throw away the ok return value, expecting it will never be false
	stanzaValues["homePath"], _ = firstIndexPathString(index.HomePath, defaultIndexPath(index.Name, "db"))
	stanzaValues["coldPath"], _ = firstIndexPathString(index.HomePath, defaultIndexPath(index.Name, "colddb"))
	stanzaValues["thawedPath"], _ = firstIndexPathString(index.HomePath, defaultIndexPath(index.Name, "thaweddb"))

	// Dynamic Data Active Archive settings
	if index.ColdStorageProvider != ARCHIVERUNDEF {
		stanzaValues["archiver.coldStorageProvider"] = string(index.ColdStorageProvider)
	}
	//storage retention in days: https://docs.splunk.com/Documentation/Splunk/latest/Admin/indexesconf
	if index.ColdStorageRetentionPeriod.InDays() != 0 {
		stanzaValues["archiver.coldStorageRetentionPeriod"] = fmt.Sprintf("%d", index.ColdStorageRetentionPeriod.InDays())
	}
	if index.EnableDataArchive {
		stanzaValues["archiver.enableDataArchive"] = fmt.Sprintf("%v", index.EnableDataArchive)
	}
	if index.MaxDataArchiveRetentionPeriod.InSeconds() != 0 {
		stanzaValues["archiver.maxDataArchiveRetentionPeriod"] = fmt.Sprintf("%d", index.MaxDataArchiveRetentionPeriod.InSeconds())
	}

	return stanzaValues
}

// stanza returns the Stanza for an Index.
func (index Index) stanza() Stanza {
	return Stanza{
		index.stanzaName(),
		index.stanzaValues(),
	}
}

// defaultLookupValues returns a LookupValues object representing the defaults for an Index object.
// The defaults for an Index object are:
// * index - Name of the index
func (index Index) defaultLookupValues() LookupValues {
	return LookupValues{
		"index":          string(index.Name),
		"retention_days": fmt.Sprintf("%d", index.FrozenTime.InDays()),
	}
}

// lookupRowsForLookup returns this Index's LookupRows for the given Lookup.
func (index Index) lookupRowsForLookup(lookup Lookup) LookupRows {
	return rowsForLookupOrDefaultRows(index.LookupRows, lookup, index)
}
