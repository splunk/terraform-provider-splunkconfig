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

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-splunkconfig/internal/splunkconfig/config"
)

const (
	roleNameKey                    = "role_name"
	searchIndexesAllowedKey        = "search_indexes_allowed"
	importRolesKey                 = "imported_roles"
	capabilitiesKey                = "capabilities"
	searchFilterKey                = "search_filter"
	cumulativeRTSearchJobsQuotaKey = "cumulative_realtime_search_jobs_quota"
	cumulativeSearchJobsQuotaKey   = "cumulative_search_jobs_quota"
	rtSearchJobsQuotaKey           = "realtime_search_jobs_quota"
	searchDiskQuotaKey             = "search_disk_quota"
	searchJobsQuotaKey             = "search_jobs_quota"
	searchTimeWinKey               = "search_time_win"
)

func resourceRoleAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific role",
		ReadContext: resourceRoleAttributesRead,
		Schema: map[string]*schema.Schema{
			roleNameKey: {
				Description: "Name of the role",
				Type:        schema.TypeString,
				Required:    true,
			},
			searchIndexesAllowedKey: {
				Description: "List of indexes searchable by the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			importRolesKey: {
				Description: "List of roles imported by the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			capabilitiesKey: {
				Description: "List of capabilities assigned to the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			searchFilterKey: {
				Description: "Search filter applied to the role",
				Type:        schema.TypeString,
				Computed:    true,
			},
			cumulativeRTSearchJobsQuotaKey: {
				Description: "Cumulative real-time search jobs quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			cumulativeSearchJobsQuotaKey: {
				Description: "Cumulative search jobs quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			rtSearchJobsQuotaKey: {
				Description: "Real-time search jobs quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			searchDiskQuotaKey: {
				Description: "Search disk quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			searchJobsQuotaKey: {
				Description: "Search jobs quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			searchTimeWinKey: {
				Description: "Search time window applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceRoleAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	roleName := d.Get(roleNameKey).(string)

	d.SetId(roleName)

	role, ok := suite.ExtrapolatedRoles().WithRoleName(config.RoleName(roleName))
	if !ok {
		return diag.Errorf("Unable to find role with name %q", roleName)
	}

	if len(role.SearchIndexesAllowed) > 0 {
		if err := d.Set(searchIndexesAllowedKey, role.SearchIndexesAllowed); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(role.ImportRoles) > 0 {
		if err := d.Set(importRolesKey, role.ImportRoles); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(role.EnabledCapabilityNames()) > 0 {
		if err := d.Set(capabilitiesKey, role.EnabledCapabilityNames()); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.SearchFilter != "" {
		if err := d.Set(searchFilterKey, role.SearchFilter); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.CumulativeRTSearchJobsQuota.Explicit {
		if err := d.Set(cumulativeRTSearchJobsQuotaKey, role.CumulativeRTSearchJobsQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.CumulativeSearchJobsQuota.Explicit {
		if err := d.Set(cumulativeSearchJobsQuotaKey, role.CumulativeSearchJobsQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.RTSearchJobsQuota.Explicit {
		if err := d.Set(rtSearchJobsQuotaKey, role.RTSearchJobsQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.SearchDiskQuota.Explicit {
		if err := d.Set(searchDiskQuotaKey, role.SearchDiskQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.SearchJobsQuota.Explicit {
		if err := d.Set(searchJobsQuotaKey, role.SearchJobsQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.SearchTimeWin.Explicit {
		if err := d.Set(searchTimeWinKey, role.SearchTimeWin.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}
