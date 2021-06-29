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
	userNameKey            = "user_name"
	userEmailKey           = "email"
	userForceChangePassKey = "force_change_pass"
	userPasswordKey        = "password"
	userRealNameKey        = "realname"
	userRolesKey           = "roles"
)

func resourceUserAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific user",
		ReadContext: resourceUserAttributesRead,
		Schema: map[string]*schema.Schema{
			userNameKey: {
				Description: "Name of the user",
				Type:        schema.TypeString,
				Required:    true,
			},
			userEmailKey: {
				Description: "Email address of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},
			userForceChangePassKey: {
				Description: "Force password change status of the user",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			userPasswordKey: {
				Description: "Password of the user",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			userRealNameKey: {
				Description: "Real name of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},
			userRolesKey: {
				Description: "Cumulative real-time search jobs quota applied to the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	userName := d.Get(userNameKey).(string)

	d.SetId(userName)

	user, ok := suite.Users.WithName(userName)
	if !ok {
		return diag.Errorf("Unable to find user with name %q", userName)
	}

	if user.Email != "" {
		d.Set(userEmailKey, user.Email)
	}

	d.Set(userForceChangePassKey, user.ForceChangePass)

	if user.Password != "" {
		d.Set(userPasswordKey, user.Password)
	}

	if user.RealName != "" {
		d.Set(userRealNameKey, user.RealName)
	}

	if len(user.Roles) > 0 {
		d.Set(userRolesKey, user.Roles)
	}

	return diag.Diagnostics{}
}
