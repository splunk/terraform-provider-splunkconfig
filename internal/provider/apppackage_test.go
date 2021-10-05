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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAppPackage(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// initial creation
			{
				Config: testAccDataSourceAppFileConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "base_version", "1.0.0"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "effective_version", "1.0.0"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "patch_count", "0"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "tarball_path", "/tmp/indexes_app-1.0.0.tgz"),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.0.content", regexp.MustCompile("version = 1.0.0")),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "files.0.path", "default/app.conf"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "files.1.path", "default/indexes.conf"),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.1.content", regexp.MustCompile(`\[original_index]`)),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "files.2.path", "default/collections.conf"),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.2.content", regexp.MustCompile(`\[collection_a]`)),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.2.content", regexp.MustCompile("field.field_a = string")),
				),
			},

			// perform updates that result in a bumped patch count
			{
				Config: testAccDataSourceAppFileConfigPatchIncrease,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "base_version", "1.0.0"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "effective_version", "1.0.1"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "patch_count", "1"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "tarball_path", "/tmp/indexes_app-1.0.1.tgz"),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.0.content", regexp.MustCompile("version = 1.0.1")),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.1.content", regexp.MustCompile("[patch_increase_index]")),
				),
			},

			// perform another update that result in a bumped patch count, to ensure the templated version matches the
			// expected patch count (ie, patch counts aren't cumulative additions)
			{
				Config: testAccDataSourceAppFileConfigPatchIncreaseAgain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "base_version", "1.0.0"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "effective_version", "1.0.2"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "patch_count", "2"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "tarball_path", "/tmp/indexes_app-1.0.2.tgz"),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.0.content", regexp.MustCompile("version = 1.0.2")),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.1.content", regexp.MustCompile("[patch_increase_again_index]")),
				),
			},

			// perform updates that result in a reset patch count
			{
				Config: testAccDataSourceAppFileConfigPatchReset,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "base_version", "1.1.0"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "effective_version", "1.1.0"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "patch_count", "0"),
					resource.TestCheckResourceAttr("splunkconfig_app_package.indexes", "tarball_path", "/tmp/indexes_app-1.1.0.tgz"),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.0.content", regexp.MustCompile("version = 1.1.0")),
					resource.TestMatchResourceAttr("splunkconfig_app_package.indexes", "files.1.content", regexp.MustCompile("[new_version_index]")),
				),
			},

			// perform updates that result in a lowered effective version, which is disallowed
			{
				Config:      testAccDataSourceAppFileConfigPatchResetInvalid,
				ExpectError: regexp.MustCompile("not greater than old effective version"),
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

resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes_app"
  path = "/tmp"
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

resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes_app"
  path = "/tmp"
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

resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes_app"
  path = "/tmp"
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

resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes_app"
  path = "/tmp"
}
`

const testAccDataSourceAppFileConfigPatchResetInvalid = `
provider "splunkconfig" {
	configuration = <<EOT
apps:
  - name: Indexes App
    id: indexes_app
    # *** START CHANGE ***
    version: 1.0.0
    # *** END CHANGE ***
    indexes:
      - name: new_version_invalid_index
EOT
}

resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes_app"
  path = "/tmp"
}
`
