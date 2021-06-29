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

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-provider-splunkconfig/internal/splunkconfig/config"
)

const (
	samlGroupNameKey = "saml_group_name"
	rolesKey         = "roles"
)

func resourceSAMLGroupAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific SAML group",
		ReadContext: resourceSAMLGroupAttributesRead,
		Schema: map[string]*schema.Schema{
			samlGroupNameKey: {
				Description: "Name of the SAML group",
				Type:        schema.TypeString,
				Required:    true,
			},
			rolesKey: {
				Description: "List of roles associated with the SAML group",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSAMLGroupAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	samlGroupName := d.Get(samlGroupNameKey).(string)

	d.SetId(samlGroupName)

	samlGroup, ok := suite.ExtrapolatedSAMLGroups().WithSAMLGroupName(samlGroupName)
	if !ok {
		return diag.Errorf("Unable to find SAML group with name %q", samlGroupName)
	}

	if len(samlGroup.Roles) > 0 {
		d.Set(rolesKey, samlGroup.Roles)
	}

	return diag.Diagnostics{}
}
