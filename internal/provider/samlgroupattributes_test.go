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
  role_name = "saml_group_a"
}
`
