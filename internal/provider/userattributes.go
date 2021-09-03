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
	userAttributesUserNameKey            = "user_name"
	userAttributesUserEmailKey           = "email"
	userAttributesUserForceChangePassKey = "force_change_pass"
	userAttributesUserPasswordKey        = "password"
	userAttributesUserRealNameKey        = "realname"
	userAttributesUserRolesKey           = "roles"
)

func resourceUserAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific user",
		ReadContext: resourceUserAttributesRead,
		Schema: map[string]*schema.Schema{
			userAttributesUserNameKey: {
				Description: "Name of the user",
				Type:        schema.TypeString,
				Required:    true,
			},
			userAttributesUserEmailKey: {
				Description: "Email address of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},
			userAttributesUserForceChangePassKey: {
				Description: "Force password change status of the user",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			userAttributesUserPasswordKey: {
				Description: "Password of the user",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			userAttributesUserRealNameKey: {
				Description: "Real name of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},
			userAttributesUserRolesKey: {
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

	userName := d.Get(userAttributesUserNameKey).(string)

	d.SetId(userName)

	user, ok := suite.Users.WithName(userName)
	if !ok {
		return diag.Errorf("Unable to find user with name %q", userName)
	}

	if user.Email != "" {
		if err := d.Set(userAttributesUserEmailKey, user.Email); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set(userAttributesUserForceChangePassKey, user.ForceChangePass); err != nil {
		return diag.FromErr(err)
	}

	if user.Password != "" {
		if err := d.Set(userAttributesUserPasswordKey, user.Password); err != nil {
			return diag.FromErr(err)
		}
	}

	if user.RealName != "" {
		if err := d.Set(userAttributesUserRealNameKey, user.RealName); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(user.Roles) > 0 {
		if err := d.Set(userAttributesUserRolesKey, user.Roles); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}
