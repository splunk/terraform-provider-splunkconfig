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

func TestAccDataAppPackage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// initial creation
			{
				Config: testAccDataSourceAppFileConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "effective_version", "1.0.0"),
					// patch_count not set, so it shouldn't be be found
					resource.TestCheckNoResourceAttr("data.splunkconfig_app_package.indexes", "patch_count"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "tarball_path", "/tmp/indexes_app-1.0.0.tgz"),
				),
			},

			// perform updates that result in a bumped patch count
			{
				Config: testAccDataSourceAppFileConfigPatchIncrease,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "effective_version", "1.0.1"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "patch_count", "1"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "tarball_path", "/tmp/indexes_app-1.0.1.tgz"),
				),
			},

			// perform another update that result in a bumped patch count, to ensure the templated version matches the
			// expected patch count (ie, patch counts aren't cumulative additions)
			{
				Config: testAccDataSourceAppFileConfigPatchIncreaseAgain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "effective_version", "1.0.2"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "patch_count", "2"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "tarball_path", "/tmp/indexes_app-1.0.2.tgz"),
				),
			},

			// perform updates that result in a reset patch count
			{
				Config: testAccDataSourceAppFileConfigPatchReset,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "effective_version", "1.1.0"),
					// patch_count not set, so it shouldn't be be found
					resource.TestCheckNoResourceAttr("data.splunkconfig_app_package.indexes", "patch_count"),
					resource.TestCheckResourceAttr("data.splunkconfig_app_package.indexes", "tarball_path", "/tmp/indexes_app-1.1.0.tgz"),
				),
			},
		},
	})
}

const testAccDataSourceAppFileConfig = `
provider "splunkconfig" {
	configuration = <<EOT
apps:
  - name: Indexes App
    id: indexes_app
    version: 1.0.0
    indexes:
      - name: original_index
    collections:
      - name: collection_a
        fields:
          field_a: string
EOT
}

data "splunkconfig_app_package" "indexes" {
  app_id = "indexes_app"
  path   = "/tmp"
}
`

const testAccDataSourceAppFileConfigPatchIncrease = `
provider "splunkconfig" {
	configuration = <<EOT
apps:
  - name: Indexes App
    id: indexes_app
    version: 1.0.0
    indexes:
      # *** START CHANGE ***
      - name: patch_increase_index
      # *** END CHANGE ***
EOT
}

data "splunkconfig_app_package" "indexes" {
  app_id      = "indexes_app"
  patch_count = 1
  path        = "/tmp"
}
`

const testAccDataSourceAppFileConfigPatchIncreaseAgain = `
provider "splunkconfig" {
	configuration = <<EOT
apps:
  - name: Indexes App
    id: indexes_app
    version: 1.0.0
    indexes:
      # *** START CHANGE ***
      - name: patch_increase_again_index
      # *** END CHANGE ***
EOT
}

data "splunkconfig_app_package" "indexes" {
  app_id      = "indexes_app"
  patch_count = 2
  path        = "/tmp"
}
`

const testAccDataSourceAppFileConfigPatchReset = `
provider "splunkconfig" {
	configuration = <<EOT
apps:
  - name: Indexes App
    id: indexes_app
    # *** START CHANGE ***
    version: 1.1.0
    # *** END CHANGE ***
    indexes:
      - name: new_version_index
EOT
}

data "splunkconfig_app_package" "indexes" {
  app_id = "indexes_app"
  # patch_count not set
  path   = "/tmp"
}
`
