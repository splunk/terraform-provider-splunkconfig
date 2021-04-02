package provider

import (
	"cd.splunkdev.com/sas/libraries/go/splunk/config/pkg/config"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	roleNameKey             = "role_name"
	searchIndexesAllowedKey = "search_indexes_allowed"
)

func resourceRoleAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific role",
		ReadContext: resourceRoleAttributesRead,
		Schema: map[string]*schema.Schema{
			roleNameKey: {
				Description: "Name of the role",
				Type:        schema.TypeString,
				Required:    true,
			},
			searchIndexesAllowedKey: {
				Description: "List of indexes searchable by the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceRoleAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	roleName := d.Get(roleNameKey).(string)

	d.SetId(roleName)

	role, ok := suite.ExtrapolatedRoles().WithRoleName(config.RoleName(roleName))
	if !ok {
		return diag.Errorf("Unable to find role with name %q", roleName)
	}

	d.Set(searchIndexesAllowedKey, role.SearchIndexesAllowed)

	return diag.Diagnostics{}
}
