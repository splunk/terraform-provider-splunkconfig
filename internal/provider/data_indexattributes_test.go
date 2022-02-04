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

func TestAccResourceIndexAttributes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIndexAttributesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.splunkconfig_index_attributes.empty", "name", "empty"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_index_attributes.empty", "frozen_time_period_in_secs"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_index_attributes.empty", "datatype"),

					resource.TestCheckResourceAttr("data.splunkconfig_index_attributes.frozen_time", "name", "frozen_time"),
					resource.TestCheckResourceAttr("data.splunkconfig_index_attributes.frozen_time", "frozen_time_period_in_secs", "86400"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_index_attributes.frozen_time", "datatype"),

					resource.TestCheckResourceAttr("data.splunkconfig_index_attributes.datatype", "name", "datatype"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_index_attributes.datatype", "frozen_time_period_in_secs"),
					resource.TestCheckResourceAttr("data.splunkconfig_index_attributes.datatype", "datatype", "event"),
				),
			},
		},
	})
}

const testAccDataSourceIndexAttributesConfig = `
provider "splunkconfig" {
    configuration = <<EOT
indexes:
  - name: empty

  - name: frozen_time
    frozenTimePeriod: {days: 1}

  - name: datatype
    datatype: event
EOT
}

data "splunkconfig_index_attributes" "empty" {
    name = "empty"
}

data "splunkconfig_index_attributes" "frozen_time" {
    name = "frozen_time"
}

data "splunkconfig_index_attributes" "datatype" {
    name = "datatype"
}
`
