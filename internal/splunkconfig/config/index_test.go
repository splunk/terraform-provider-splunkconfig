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
	"fmt"
	"testing"
)

// Index.validate should return an error when one is expected.
func TestIndex_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			Index{},
			true,
		},
		{
			Index{Name: "main"},
			false,
		},
		{
			Index{
				Name:               "main",
				SearchRolesAllowed: RoleNames{"admin"},
			},
			false,
		},
		{
			Index{
				Name:               "main",
				SearchRolesAllowed: RoleNames{"duplicate", "duplicate"},
			},
			true,
		},
	}

	tests.test(t)
}

func TestIndex_validateWithRoles(t *testing.T) {
	// validateWithRoles only validates SearchRolesAllowed, so we don't have to set an index name here
	index := Index{SearchRolesAllowed: RoleNames{"role_a", "role_b", "role_c"}}
	validAgainstRoles := Roles{Role{Name: "role_a"}, Role{Name: "role_b"}, Role{Name: "role_c"}}
	invalidAgainstRoles := Roles{Role{Name: "role_d"}, Role{Name: "role_e"}, Role{Name: "role_f"}}

	tests := []struct {
		roles     Roles
		wantError bool
	}{
		{validAgainstRoles, false},
		{invalidAgainstRoles, true},
	}

	for _, test := range tests {
		gotError := index.validateWithRoles(test.roles) != nil
		message := fmt.Sprintf("%T{%+v}.validateWithRoles(%T{%+v}) returned error?", index, index, test.roles, test.roles)

		testEqual(gotError, test.wantError, message, t)
	}
}

func TestIndex_searchableByRoleName(t *testing.T) {
	index := Index{SearchRolesAllowed: RoleNames{"role_a", "role_b"}}

	tests := []struct {
		roleName RoleName
		want     bool
	}{
		{roleName: "role_a", want: true},
		{roleName: "role_b", want: true},
		{roleName: "role_c", want: false},
	}

	for _, test := range tests {
		got := index.searchableByRoleName(test.roleName)
		message := fmt.Sprintf("%T{%+v}.searchableByRoleName(%q)", index, index, test.roleName)
		testEqual(got, test.want, message, t)
	}
}

func TestIndex_stanza(t *testing.T) {
	tests := stanzaDefinerTestCases{
		{
			Index{Name: "index_a"},
			Stanza{
				Name: "index_a",
				Values: StanzaValues{
					"homePath":   "$SPLUNK_DB/index_a/db",
					"coldPath":   "$SPLUNK_DB/index_a/colddb",
					"thawedPath": "$SPLUNK_DB/index_a/thaweddb",
				},
			},
		},
		{
			Index{Name: "index_a", FrozenTime: TimePeriod{Seconds: 86400}},
			Stanza{
				Name: "index_a",
				Values: StanzaValues{
					"homePath":               "$SPLUNK_DB/index_a/db",
					"coldPath":               "$SPLUNK_DB/index_a/colddb",
					"thawedPath":             "$SPLUNK_DB/index_a/thaweddb",
					"frozenTimePeriodInSecs": "86400",
				},
			},
		},
		{
			Index{Name: "index_a", DataType: "event"},
			Stanza{
				Name: "index_a",
				Values: StanzaValues{
					"homePath":   "$SPLUNK_DB/index_a/db",
					"coldPath":   "$SPLUNK_DB/index_a/colddb",
					"datatype":   "event",
					"thawedPath": "$SPLUNK_DB/index_a/thaweddb",
				},
			},
		},
		{
			Index{Name: "index_a", EnableDataArchive: true, ColdStorageRetentionPeriod: TimePeriod{Days: 5}, MaxDataArchiveRetentionPeriod: TimePeriod{Seconds: 800}, ColdStorageProvider: "Glacier"},
			Stanza{
				Name: "index_a",
				Values: StanzaValues{
					"homePath":                               "$SPLUNK_DB/index_a/db",
					"coldPath":                               "$SPLUNK_DB/index_a/colddb",
					"thawedPath":                             "$SPLUNK_DB/index_a/thaweddb",
					"archiver.enableDataArchive":             "true",
					"archiver.coldStorageRetentionPeriod":    "5",
					"archiver.maxDataArchiveRetentionPeriod": "800",
					"archiver.coldStorageProvider":           "Glacier",
				},
			},
		},
	}

	tests.test(t)
}

func TestIndex_lookupRowsForLookup(t *testing.T) {
	tests := lookupRowsForLookupDefinerTestCases{
		{
			Index{Name: "index_a", LookupRows: LookupRows{LookupRow{LookupName: "index_lookup", Values: LookupValues{"contact": "contact_a"}}}},
			Lookup{
				Name: "index_lookup",
				Fields: LookupFields{
					LookupField{Name: "index", Required: true},
					LookupField{Name: "contact"},
				},
			},
			LookupRows{
				LookupRow{
					LookupName: "index_lookup",
					Values:     LookupValues{"index": "index_a", "contact": "contact_a"},
				},
			},
		},
	}

	tests.test(t)
}
