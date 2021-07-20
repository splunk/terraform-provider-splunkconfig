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
	roleNamesKey     = "role_names"
	roleNamesIDValue = "splunkconfig_role_names"
)

func resourceRoleNames() *schema.Resource {
	return &schema.Resource{
		Description: "Return Role Names from the Splunk Configuration",
		ReadContext: resourceRoleNamesRead,
		Schema: map[string]*schema.Schema{
			roleNamesKey: {
				Description: "List of Role Names in the Splunk Configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceRoleNamesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	d.SetId(roleNamesIDValue)
	if err := d.Set(roleNamesKey, suite.ExtrapolatedRoles().RoleNames()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
