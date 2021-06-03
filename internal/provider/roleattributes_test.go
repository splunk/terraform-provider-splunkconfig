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
					resource.TestCheckResourceAttr("data.splunkconfig_role_attributes.role_a", "cumulative_realtime_search_jobs_quota", "0"),
					resource.TestCheckResourceAttr("data.splunkconfig_role_attributes.role_a", "cumulative_search_jobs_quota", "0"),
					resource.TestCheckResourceAttr("data.splunkconfig_role_attributes.role_a", "realtime_search_jobs_quota", "0"),
					resource.TestCheckResourceAttr("data.splunkconfig_role_attributes.role_a", "search_disk_quota", "0"),
					resource.TestCheckResourceAttr("data.splunkconfig_role_attributes.role_a", "search_jobs_quota", "0"),
					resource.TestCheckResourceAttr("data.splunkconfig_role_attributes.role_a", "search_time_win", "0"),
				),
			},
			{
				Config: testAccDataSourceRoleAttributesConfigRemovedZeroed,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("data.splunkconfig_role_attributes.role_a", "cumulative_realtime_search_jobs_quota"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_role_attributes.role_a", "cumulative_search_jobs_quota"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_role_attributes.role_a", "realtime_search_jobs_quota"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_role_attributes.role_a", "search_disk_quota"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_role_attributes.role_a", "search_jobs_quota"),
					resource.TestCheckNoResourceAttr("data.splunkconfig_role_attributes.role_a", "search_time_win"),
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

    # explicitly settable as zero values
    cumulativeRTSrchJobsQuota: 0
    cumulativeSrchJobsQuota: 0
    rtSrchJobsQuota: 0
    srchDiskQuota: 0
    srchJobsQuota: 0
    srchTimeWin: 0

indexes:
  - name: index_a_reverse
    srchRolesAllowed: ["role_a"]
EOT
}

data "splunkconfig_role_attributes" "role_a" {
  role_name = "role_a"
}
`

const testAccDataSourceRoleAttributesConfigRemovedZeroed = `
provider "splunkconfig" {
	configuration = <<EOT
roles:
  - name: role_a
    srchIndexesAllowed: ["index_a"]
    importRoles: ["user"]
    capabilities:
      admin_all_objects: true
      change_authentication: false

    # explicitly settable zeroed values removed to ensure they're no longer set by the provider

indexes:
  - name: index_a_reverse
    srchRolesAllowed: ["role_a"]
EOT
}

data "splunkconfig_role_attributes" "role_a" {
  role_name = "role_a"
}
`
