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
	appAttributesIdKey              = "app_id"
	appAttributesNameKey            = "name"
	appAttributesDescriptionKey     = "description"
	appAttributesAuthorKey          = "author"
	appAttributesIsVisibleKey       = "is_visible"
	appAttributesCheckForUpdatesKey = "check_for_updates"
	appAttributesVersionKey         = "version"
	appAttributesAclReadKey         = "acl_read"
	appAttributesAclWriteKey        = "acl_write"
	appAttributesAclSharingKey      = "acl_sharing"
)

func resourceAppAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific app",
		ReadContext: resourceAppAttributesRead,
		Schema: map[string]*schema.Schema{
			appAttributesIdKey: {
				Description: "ID of the app",
				Type:        schema.TypeString,
				Required:    true,
			},
			appAttributesNameKey: {
				Description: "App name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appAttributesDescriptionKey: {
				Description: "App description",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appAttributesAuthorKey: {
				Description: "App author",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appAttributesIsVisibleKey: {
				Description: "App visibility",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			appAttributesCheckForUpdatesKey: {
				Description: "App updating checking",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			appAttributesVersionKey: {
				Description: "App version",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appAttributesAclReadKey: {
				Description: "App read roles",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			appAttributesAclWriteKey: {
				Description: "App write roles",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			appAttributesAclSharingKey: {
				Description: "App sharing",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceAppAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	appId := d.Get(appAttributesIdKey).(string)

	d.SetId(appId)

	apps, err := suite.ExtrapolatedApps()
	if err != nil {
		return diag.Errorf("unable to extrapolate apps: %s", err)
	}

	app, ok := apps.WithID(appId)
	if !ok {
		return diag.Errorf("Unable to find app with ID %s", appId)
	}

	c := conditionalConfigurations{
		{
			true,
			appAttributesNameKey,
			app.Name,
		},
		{
			app.Description != "",
			appAttributesDescriptionKey,
			app.Description,
		},
		{
			app.Author != "",
			appAttributesAuthorKey,
			app.Author,
		},
		{
			true,
			appAttributesIsVisibleKey,
			app.IsVisible,
		},
		{
			true,
			appAttributesCheckForUpdatesKey,
			app.CheckForUpdates,
		},
		{
			true,
			appAttributesVersionKey,
			app.Version.AsString(),
		},
		{
			app.ACL.Read != nil,
			appAttributesAclReadKey,
			app.ACL.Read,
		},
		{
			app.ACL.Write != nil,
			appAttributesAclWriteKey,
			app.ACL.Write,
		},
		{
			app.ACL.Sharing != "",
			appAttributesAclSharingKey,
			app.ACL.Sharing,
		},
	}

	if err := c.apply(d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
