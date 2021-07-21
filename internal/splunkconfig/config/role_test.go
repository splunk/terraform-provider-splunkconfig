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

// Role.validate() should return an error when expected.
func TestRole_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			Role{Name: "valid"},
			false,
		},
	}

	tests.test(t)
}

func TestRole_extrapolateFromIndexes(t *testing.T) {
	role := Role{
		Name:                 "role_a",
		SearchIndexesAllowed: IndexNames{"index_b", "index_c"},
	}

	indexes := Indexes{
		// this index will be added to role's SearchIndexesAllowed
		Index{
			Name:               "index_a",
			SearchRolesAllowed: RoleNames{"role_a"},
		},
		// this index already exists in role's SearchIndexesAllowed
		Index{
			Name:               "index_b",
			SearchRolesAllowed: RoleNames{"role_a"},
		},
	}

	want := IndexNames{"index_a", "index_b", "index_c"}
	got := role.extrapolateFromIndexes(indexes).SearchIndexesAllowed
	message := fmt.Sprintf("%T{%+v}.extrapolateFromIndexes(%T{%+v}).SearchIndexesAllowed", role, role, indexes, indexes)
	testEqual(got, want, message, t)
}

func TestRole_stanzas(t *testing.T) {
	tests := stanzaDefinerTestCases{
		{
			Role{Name: "testrole_a"},
			Stanza{Name: "role_testrole_a", Values: StanzaValues{}},
		},
		{
			Role{
				Name: "empty",
			},
			Stanza{
				Name:   "role_empty",
				Values: StanzaValues{},
			},
		},
		{
			Role{
				Name:                 "search_indexes",
				SearchIndexesAllowed: IndexNames{"index_a", "index_b"},
			},
			Stanza{
				Name: "role_search_indexes",
				Values: StanzaValues{
					"srchIndexesAllowed": "index_a;index_b",
				},
			},
		},
		{
			Role{
				Name:         "capabilities",
				Capabilities: Capabilities{"enabled_capability": true},
			},
			Stanza{
				Name: "role_capabilities",
				Values: StanzaValues{
					"enabled_capability": "enabled",
				},
			},
		},
		{
			Role{
				Name:         "search_filter",
				SearchFilter: "index::filtered_index",
			},
			Stanza{
				Name: "role_search_filter",
				Values: StanzaValues{
					"srchFilter": "index::filtered_index",
				},
			},
		},
		{
			Role{
				Name:          "search_time_window",
				SearchTimeWin: ExplicitlySetInt(0),
			},
			Stanza{
				Name: "role_search_time_window",
				Values: StanzaValues{
					"srchTimeWin": "0",
				},
			},
		},
		{
			Role{
				Name:            "search_disk_quota",
				SearchDiskQuota: ExplicitlySetInt(0),
			},
			Stanza{
				Name: "role_search_disk_quota",
				Values: StanzaValues{
					"srchDiskQuota": "0",
				},
			},
		},
		{
			Role{
				Name:            "search_jobs_quota",
				SearchJobsQuota: ExplicitlySetInt(0),
			},
			Stanza{
				Name: "role_search_jobs_quota",
				Values: StanzaValues{
					"srchJobsQuota": "0",
				},
			},
		},
		{
			Role{
				Name:              "rt_search_jobs_quota",
				RTSearchJobsQuota: ExplicitlySetInt(0),
			},
			Stanza{
				Name: "role_rt_search_jobs_quota",
				Values: StanzaValues{
					"rtSrchJobsQuota": "0",
				},
			},
		},
		{
			Role{
				Name:                      "cumulative_search_jobs_quota",
				CumulativeSearchJobsQuota: ExplicitlySetInt(0),
			},
			Stanza{
				Name: "role_cumulative_search_jobs_quota",
				Values: StanzaValues{
					"cumulativeSrchJobsQuota": "0",
				},
			},
		},
		{
			Role{
				Name:                        "cumulative_rt_search_jobs_quota",
				CumulativeRTSearchJobsQuota: ExplicitlySetInt(0),
			},
			Stanza{
				Name: "role_cumulative_rt_search_jobs_quota",
				Values: StanzaValues{
					"cumulativeRTSrchJobsQuota": "0",
				},
			},
		},
	}

	tests.test(t)
}

func TestRole_lookupRowsForLookup(t *testing.T) {
	tests := lookupRowsForLookupDefinerTestCases{
		{
			Role{Name: "role_a", LookupRows: LookupRows{LookupRow{LookupName: "role_lookup", Values: LookupValues{"contact": "contact_a"}}}},
			Lookup{
				Name: "role_lookup",
				Fields: LookupFields{
					LookupField{Name: "role", Required: true},
					LookupField{Name: "contact"},
				},
			},
			LookupRows{
				LookupRow{
					LookupName: "role_lookup",
					Values:     LookupValues{"role": "role_a", "contact": "contact_a"},
				},
			},
		},
	}

	tests.test(t)
}
