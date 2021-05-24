package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccResourceRoleAttributes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRoleAttributesConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrList("data.splunkconfig_role_attributes.role_a", "search_indexes_allowed", []string{
						"index_a",
						"index_a_reverse",
					}),
					testCheckResourceAttrList("data.splunkconfig_role_attributes.role_a", "imported_roles", []string{
						"user",
					}),
					testCheckResourceAttrList("data.splunkconfig_role_attributes.role_a", "capabilities", []string{
						"admin_all_objects",
					}),
				),
			},
		},
	})
}

const testAccDataSourceRoleAttributesConfig = `
provider "splunkconfig" {
	configuration = <<EOT
roles:
  - name: role_a
    srchIndexesAllowed: ["index_a"]
    importRoles: ["user"]
    capabilities:
      admin_all_objects: true
      change_authentication: false

indexes:
  - name: index_a_reverse
    srchRolesAllowed: ["role_a"]
EOT
}

data "splunkconfig_role_attributes" "role_a" {
  role_name = "role_a"
}
`
