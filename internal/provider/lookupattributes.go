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
	lookupAttributesLookupNameKey = "lookup_name"
	lookupAttributesFieldNamesKey = "field_names"
	lookupAttributesRowsKey       = "rows"
)

func resourceLookupAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get fields and rows for a specific lookup",
		ReadContext: resourceLookupAttributesRead,
		Schema: map[string]*schema.Schema{
			lookupAttributesLookupNameKey: {
				Description: "Name of the lookup",
				Type:        schema.TypeString,
				Required:    true,
			},
			lookupAttributesFieldNamesKey: {
				Description: "List of fields in the lookup",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			lookupAttributesRowsKey: {
				Description: "List of rows in the lookup",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func resourceLookupAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	lookupName := d.Get(lookupAttributesLookupNameKey).(string)

	d.SetId(lookupName)

	lookup, ok := suite.ExtrapolatedLookups().WithName(lookupName)
	if !ok {
		return diag.Errorf("unable to find lookup: %s", lookupName)
	}

	// fields
	fieldNames := make([]string, len(lookup.Fields))

	for i, field := range lookup.Fields {
		fieldNames[i] = field.Name
	}

	if err := d.Set(lookupAttributesFieldNamesKey, fieldNames); err != nil {
		return diag.FromErr(err)
	}

	// rows
	rows := make([]map[string]string, len(lookup.Rows))

	for i, row := range lookup.Rows {
		rows[i] = map[string]string(row.Values)
	}

	if err := d.Set(lookupAttributesRowsKey, rows); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
