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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceRoleNames(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRoleNamesConfig,
				Check: testCheckResourceAttrList("data.splunkconfig_role_names.foo", "role_names", []string{
					"role_a",
					"role_b",
				}),
			},
		},
	})
}

const testAccDataSourceRoleNamesConfig = `
provider "splunkconfig" {
	configuration = <<EOT
roles:
  - name: role_a
  - name: role_b
EOT
}

data "splunkconfig_role_names" "foo" {}
`
