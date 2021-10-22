// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
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
	appIdsRequireTagKey = "require_tag"
	appIdsExcludeTagKey = "exclude_tag"
	appIdsAppIdsKey     = "app_ids"
	appIdsIdValue       = "splunkconfig_app_ids"
)

func resourceAppIds() *schema.Resource {
	return &schema.Resource{
		Description: "Return App IDs from the Splunk Configuration",
		ReadContext: resourceAppIdsRead,
		Schema: map[string]*schema.Schema{
			appIdsRequireTagKey: {
				Description: "Tags to require for returned App IDs",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        tagsSchema(),
			},
			appIdsExcludeTagKey: {
				Description: "Tags to exclude for returned App IDs",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        tagsSchema(),
			},
			appIdsAppIdsKey: {
				Description: "List of App IDs in the Splunk Configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAppIdsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	d.SetId(appIdsIdValue)

	requireTags := newTagsFromInterface(d.Get(appIdsRequireTagKey))
	excludeTags := newTagsFromInterface(d.Get(appIdsExcludeTagKey))

	apps, err := suite.ExtrapolatedApps()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(appIdsAppIdsKey, apps.SatisfyingTags(requireTags, excludeTags).AppIDs()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
