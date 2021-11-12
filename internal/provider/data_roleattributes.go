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
	"terraform-provider-splunkconfig/internal/splunkconfig/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	roleAttributesRoleNameKey                    = "role_name"
	roleAttributesSearchIndexesAllowedKey        = "search_indexes_allowed"
	roleAttributesImportRolesKey                 = "imported_roles"
	roleAttributesCapabilitiesKey                = "capabilities"
	roleAttributesSearchFilterKey                = "search_filter"
	roleAttributesCumulativeRTSearchJobsQuotaKey = "cumulative_realtime_search_jobs_quota"
	roleAttributesCumulativeSearchJobsQuotaKey   = "cumulative_search_jobs_quota"
	roleAttributesRtSearchJobsQuotaKey           = "realtime_search_jobs_quota"
	roleAttributesSearchDiskQuotaKey             = "search_disk_quota"
	roleAttributesSearchJobsQuotaKey             = "search_jobs_quota"
	roleAttributesSearchTimeWinKey               = "search_time_win"
)

func dataRoleAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific role",
		ReadContext: resourceRoleAttributesRead,
		Schema: map[string]*schema.Schema{
			roleAttributesRoleNameKey: {
				Description: "Name of the role",
				Type:        schema.TypeString,
				Required:    true,
			},
			roleAttributesSearchIndexesAllowedKey: {
				Description: "List of indexes searchable by the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			roleAttributesImportRolesKey: {
				Description: "List of roles imported by the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			roleAttributesCapabilitiesKey: {
				Description: "List of capabilities assigned to the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			roleAttributesSearchFilterKey: {
				Description: "Search filter applied to the role",
				Type:        schema.TypeString,
				Computed:    true,
			},
			roleAttributesCumulativeRTSearchJobsQuotaKey: {
				Description: "Cumulative real-time search jobs quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			roleAttributesCumulativeSearchJobsQuotaKey: {
				Description: "Cumulative search jobs quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			roleAttributesRtSearchJobsQuotaKey: {
				Description: "Real-time search jobs quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			roleAttributesSearchDiskQuotaKey: {
				Description: "Search disk quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			roleAttributesSearchJobsQuotaKey: {
				Description: "Search jobs quota applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			roleAttributesSearchTimeWinKey: {
				Description: "Search time window applied to the role",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceRoleAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	roleName := d.Get(roleAttributesRoleNameKey).(string)

	d.SetId(roleName)

	role, ok := suite.ExtrapolatedRoles().WithRoleName(config.RoleName(roleName))
	if !ok {
		return diag.Errorf("Unable to find role with name %q", roleName)
	}

	if len(role.SearchIndexesAllowed) > 0 {
		if err := d.Set(roleAttributesSearchIndexesAllowedKey, role.SearchIndexesAllowed); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(role.ImportRoles) > 0 {
		if err := d.Set(roleAttributesImportRolesKey, role.ImportRoles); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(role.EnabledCapabilityNames()) > 0 {
		if err := d.Set(roleAttributesCapabilitiesKey, role.EnabledCapabilityNames()); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.SearchFilter != "" {
		if err := d.Set(roleAttributesSearchFilterKey, role.SearchFilter); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.CumulativeRTSearchJobsQuota.Explicit {
		if err := d.Set(roleAttributesCumulativeRTSearchJobsQuotaKey, role.CumulativeRTSearchJobsQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.CumulativeSearchJobsQuota.Explicit {
		if err := d.Set(roleAttributesCumulativeSearchJobsQuotaKey, role.CumulativeSearchJobsQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.RTSearchJobsQuota.Explicit {
		if err := d.Set(roleAttributesRtSearchJobsQuotaKey, role.RTSearchJobsQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.SearchDiskQuota.Explicit {
		if err := d.Set(roleAttributesSearchDiskQuotaKey, role.SearchDiskQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.SearchJobsQuota.Explicit {
		if err := d.Set(roleAttributesSearchJobsQuotaKey, role.SearchJobsQuota.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	if role.SearchTimeWin.Explicit {
		if err := d.Set(roleAttributesSearchTimeWinKey, role.SearchTimeWin.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}
