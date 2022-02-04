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
	indexAttributesNameKey       = "name"
	indexAttributesFrozenTimeKey = "frozen_time_period_in_secs"
	indexAttributesDatatypeKey   = "datatype"
)

func dataIndexAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific index",
		ReadContext: resourceIndexAttributesRead,
		Schema: map[string]*schema.Schema{
			indexAttributesNameKey: {
				Description: "Index name",
				Type:        schema.TypeString,
				Required:    true,
			},
			indexAttributesFrozenTimeKey: {
				Description: "Retention period of the index, in seconds",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			indexAttributesDatatypeKey: {
				Description: "Data type of the index",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceIndexAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	indexName := d.Get(indexAttributesNameKey).(string)

	d.SetId(indexName)

	index, ok := suite.Indexes.WithIndexName(config.IndexName(indexName))
	if !ok {
		return diag.Errorf("index not found: %s", indexName)
	}

	c := conditionalConfigurations{
		{
			index.FrozenTime.InSeconds() != 0,
			indexAttributesFrozenTimeKey,
			index.FrozenTime.InSeconds(),
		},
		{
			index.DataType != "",
			indexAttributesDatatypeKey,
			index.DataType,
		},
	}

	if err := c.apply(d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
