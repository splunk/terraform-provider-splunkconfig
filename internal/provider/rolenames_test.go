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
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrList("data.splunkconfig_role_names.foo", "role_names", []string{
						"role_a",
						"role_b",
					})...,
				),
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
