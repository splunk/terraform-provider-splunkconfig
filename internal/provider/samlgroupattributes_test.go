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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceSAMLGroupAttributes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSAMLGroupAttributesConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrList("data.splunkconfig_saml_group_attributes.saml_group_a", "roles", []string{
						"explicit_role",
						"implicit_role",
					}),
				),
			},
		},
	})
}

const testAccDataSourceSAMLGroupAttributesConfig = `
provider "splunkconfig" {
	configuration = <<EOT
saml_groups:
  - name: saml_group_a
    roles:
      - explicit_role

roles:
  - name: implicit_role
    saml_groups:
      - saml_group_a
EOT
}

data "splunkconfig_saml_group_attributes" "saml_group_a" {
  saml_group_name = "saml_group_a"
}
`
