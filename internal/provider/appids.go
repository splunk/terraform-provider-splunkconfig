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
	appIdsTagKey       = "require_tag"
	appIdsTagNameKey   = "name"
	appIdsTagsValueKey = "values"
	appIdsAppIdsKey    = "app_ids"
	appIdsIdValue      = "splunkconfig_app_ids"
)

func resourceAppIds() *schema.Resource {
	return &schema.Resource{
		Description: "Return App IDs from the Splunk Configuration",
		ReadContext: resourceAppIdsRead,
		Schema: map[string]*schema.Schema{
			appIdsTagKey: {
				Description: "Tags to require for returned App IDs",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						appIdsTagNameKey: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the tag to require",
						},
						appIdsTagsValueKey: {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Values of the tag to require",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
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

	// nested schemas are "fun".
	// we have to iterate through each piece as a list or map of interfaces, and type assert
	// to the true type at each layer that is nested.
	// this block does that, converting the []interface{} that the SDK gives us into a Tags object.
	// start with a list of interfaces, each of which is a tag block
	tagInterfaces := d.Get(appIdsTagKey).([]interface{})
	tags := make(config.Tags, len(tagInterfaces))
	for tagNumber, tagInterface := range tagInterfaces {
		// type assert the item into a map (with keys of "name" and "values")
		tagMap := tagInterface.(map[string]interface{})
		tagName := tagMap[appIdsTagNameKey].(string)
		// when fetching the "values" key, we'll get a list of interfaces
		tagValueInterfaces := tagMap[appIdsTagsValueKey].([]interface{})
		tagValues := make([]string, len(tagValueInterfaces))
		for tagValueNumber, tagValueInterface := range tagValueInterfaces {
			// type assert the value interface into a string
			tagValues[tagValueNumber] = tagValueInterface.(string)
		}
		// add a tag with the determiend name/values to our real Tags object
		tags[tagNumber] = config.Tag{Name: tagName, Values: tagValues}
	}

	apps, err := suite.ExtrapolatedApps()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(appIdsAppIdsKey, apps.AppIDsSatisfyingTags(tags)); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
