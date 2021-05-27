package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceSAMLGroupNames(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSAMLGroupNamesConfig,
				Check: testCheckResourceAttrList("data.splunkconfig_saml_group_names.foo", "saml_group_names", []string{
					"saml_group_a",
					"saml_group_b",
				}),
			},
		},
	})
}

const testAccDataSourceSAMLGroupNamesConfig = `
provider "splunkconfig" {
	configuration = <<EOT
saml_groups:
  - name: saml_group_a
  - name: saml_group_b
EOT
}

data "splunkconfig_saml_group_names" "foo" {}
`
