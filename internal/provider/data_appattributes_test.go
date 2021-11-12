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

func TestAccResourceAppAttributes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppAttributesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "name", "Bare Minimum"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "description"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "author"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "is_visible", "false"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "check_for_updates", "false"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "version", "0.0.0"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "acl_read"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "acl_write"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_app_attributes.bare_minimum", "acl_sharing"),

					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "name", "Full Details"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "description", "App with all attributes set"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "author", "App Author"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "is_visible", "true"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "check_for_updates", "true"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "version", "1.0.0"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "acl_read.0", "*"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "acl_write.0", "admin"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_attributes.full_details", "acl_sharing", "global"),
				),
			},
		},
	})
}

const testAccDataSourceAppAttributesConfig = `
provider "splunkconfig" {
    configuration = <<EOT
apps:
  - id: bare_minimum
    name: Bare Minimum

  - id: full_details
    name: Full Details
    description: App with all attributes set
    author: App Author
    is_visible: true
    check_for_updates: true
    version: 1.0.0
    acl:
      read: ["*"]
      write: [admin]
      sharing: global
EOT
}

data "splunkconfig_app_attributes" "bare_minimum" {
    app_id = "bare_minimum"
}

data "splunkconfig_app_attributes" "full_details" {
    app_id = "full_details"
}
`
