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
	userNamesKey     = "user_names"
	userNamesIDValue = "splunkconfig_user_names"
)

func resourceUserNames() *schema.Resource {
	return &schema.Resource{
		Description: "Return User Names from the Splunk Configuration",
		ReadContext: resourceUserNamesRead,
		Schema: map[string]*schema.Schema{
			userNamesKey: {
				Description: "List of User Names in the Splunk Configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserNamesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	d.SetId(userNamesIDValue)
	if err := d.Set(userNamesKey, suite.Users.Names()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
