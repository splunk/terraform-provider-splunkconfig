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
	"terraform-provider-splunkconfig/internal/splunkconfig/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tagNameKey   = "name"
	tagValuesKey = "values"
)

func tagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			tagNameKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the tag",
			},
			tagValuesKey: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Values of the tag",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func newTagsFromInterface(tagsInterface interface{}) config.Tags {
	// nested schemas are "fun".
	// we have to iterate through each piece as a list or map of interfaces, and type assert
	// to the true type at each layer that is nested.
	// this function does that, converting the []interface{} that the SDK gives us into a Tags object.

	// start with a list of interfaces, each of which is a tag block
	tagInterfaces := tagsInterface.([]interface{})
	tags := make(config.Tags, len(tagInterfaces))
	for tagNumber, tagInterface := range tagInterfaces {
		// type assert the item into a map (with keys of "name" and "values")
		tagMap := tagInterface.(map[string]interface{})
		tagName := tagMap[tagNameKey].(string)
		// when fetching the "values" key, we'll get a list of interfaces
		tagValueInterfaces := tagMap[tagValuesKey].([]interface{})
		tagValues := make([]string, len(tagValueInterfaces))
		for tagValueNumber, tagValueInterface := range tagValueInterfaces {
			// type assert the value interface into a string
			tagValues[tagValueNumber] = tagValueInterface.(string)
		}
		// add a tag with the determiend name/values to our real Tags object
		tags[tagNumber] = config.Tag{Name: tagName, Values: tagValues}
	}

	return tags
}
