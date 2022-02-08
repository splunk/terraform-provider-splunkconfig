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

func TestAccResourceIndexNames(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIndexNamesConfig,
				Check: testCheckResourceAttrList("data.splunkconfig_index_names.foo", "index_names", []string{
					"index_a",
					"index_b",
				}),
			},
		},
	})
}

const testAccDataSourceIndexNamesConfig = `
provider "splunkconfig" {
	configuration = <<EOT
indexes:
  - name: index_a
  - name: index_b
EOT
}

data "splunkconfig_index_names" "foo" {}
`
