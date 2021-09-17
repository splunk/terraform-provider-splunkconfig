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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAppIds(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppIdsConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrList("data.splunkconfig_app_ids.all", "app_ids", []string{
						"app_a",
						"app_b",
					}),
					testCheckResourceAttrList("data.splunkconfig_app_ids.filtered", "app_ids", []string{
						"app_a",
					}),
				),
			},
		},
	})
}

const testAccDataSourceAppIdsConfig = `
provider "splunkconfig" {
	configuration = <<EOT
apps:
  - id: app_a
    name: App A
    tags:
      - name: env
        values: [prod]
  - id: app_b
    name: App B
EOT
}

data "splunkconfig_app_ids" "all" {}

data "splunkconfig_app_ids" "filtered" {
	require_tag {
		name   = "env"
		values = ["prod"]
	}
}
`
