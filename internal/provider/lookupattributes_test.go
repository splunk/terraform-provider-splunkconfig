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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceLookupAttributes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLookupAttributesConfig,
				Check: resource.ComposeTestCheckFunc(
					// unnumbered
					testCheckResourceAttrList("data.splunkconfig_lookup_attributes.unnumbered", "field_names", []string{
						"field_a",
						"field_b",
					}),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.unnumbered", "rows.0.field_a", "row_1_value_a"),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.unnumbered", "rows.0.field_b", "row_1_value_b"),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.unnumbered", "rows.1.field_a", "row_2_value_a"),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.unnumbered", "rows.1.field_b", "row_2_value_b"),
					// numbered
					testCheckResourceAttrList("data.splunkconfig_lookup_attributes.numbered", "field_names", []string{
						"row_number",
						"field_a",
						"field_b",
					}),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.numbered", "rows.0.row_number", "1"),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.numbered", "rows.0.field_a", "row_1_value_a"),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.numbered", "rows.0.field_b", "row_1_value_b"),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.numbered", "rows.1.row_number", "2"),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.numbered", "rows.1.field_a", "row_2_value_a"),
					resource.TestCheckResourceAttr("data.splunkconfig_lookup_attributes.numbered", "rows.1.field_b", "row_2_value_b"),
				),
			},
		},
	})
}

const testAccDataSourceLookupAttributesConfig = `
provider "splunkconfig" {
    configuration = <<EOT
lookups:
  - name: test_lookup
    fields:
      - name: field_a
      - name: field_b
    rows:
      - values:
          field_a: row_1_value_a
          field_b: row_1_value_b
      - values:
          field_a: row_2_value_a
          field_b: row_2_value_b
EOT
}

data "splunkconfig_lookup_attributes" "unnumbered" {
    lookup_name = "test_lookup"
}

data "splunkconfig_lookup_attributes" "numbered" {
	lookup_name = "test_lookup"
    row_number_field = "row_number"
}
`
