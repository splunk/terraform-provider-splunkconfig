package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceUserNames(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserNamesConfig,
				Check: testCheckResourceAttrList("data.splunkconfig_user_names.foo", "user_names", []string{
					"user_a",
					"user_b",
				}),
			},
		},
	})
}

const testAccDataSourceUserNamesConfig = `
provider "splunkconfig" {
	configuration = <<EOT
users:
  - name: user_a
  - name: user_b
EOT
}

data "splunkconfig_user_names" "foo" {}
`
