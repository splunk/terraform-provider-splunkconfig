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
	samlGroupNamesKey     = "saml_group_names"
	samlGroupNamesIDValue = "splunkconfig_saml_group_names"
)

func resourceSAMLGroupNames() *schema.Resource {
	return &schema.Resource{
		Description: "Return SAML Group Names from the Splunk Configuration",
		ReadContext: resourceSAMLGroupNamesRead,
		Schema: map[string]*schema.Schema{
			samlGroupNamesKey: {
				Description: "List of SAML Group Names in the Splunk Configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSAMLGroupNamesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	d.SetId(samlGroupNamesIDValue)
	if err := d.Set(samlGroupNamesKey, suite.ExtrapolatedSAMLGroups().SAMLGroupNames()); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
