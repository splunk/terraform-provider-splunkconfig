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

func TestAccResourceUserAttributes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserAttributesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.splunkconfig_user_attributes.user_a", "email", "user_a@example.com"),
					resource.TestCheckResourceAttr("data.splunkconfig_user_attributes.user_a", "force_change_pass", "true"),
					resource.TestCheckResourceAttr("data.splunkconfig_user_attributes.user_a", "password", "user_a_password"),
					resource.TestCheckResourceAttr("data.splunkconfig_user_attributes.user_a", "realname", "User A"),
					testCheckResourceAttrList("data.splunkconfig_user_attributes.user_a", "roles", []string{
						"role_a",
					}),
				),
			},
		},
	})
}

const testAccDataSourceUserAttributesConfig = `
provider "splunkconfig" {
	configuration = <<EOT
users:
  - name: user_a
    email: user_a@example.com
    force_change_pass: true
    password: user_a_password
    realname: User A
    roles: ["role_a"]
EOT
}

data "splunkconfig_user_attributes" "user_a" {
  user_name = "user_a"
}
`
