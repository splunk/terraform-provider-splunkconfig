package provider

import (
	"cd.splunkdev.com/sas/go/libraries/splunkconfig/pkg/config"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	roleNameKey             = "role_name"
	searchIndexesAllowedKey = "search_indexes_allowed"
	importRolesKey          = "imported_roles"
	capabilitiesKey         = "capabilities"
	searchFilterKey         = "search_filter"
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
			importRolesKey: {
				Description: "List of roles imported by the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			capabilitiesKey: {
				Description: "List of capabilities assigned to the role",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			searchFilterKey: {
				Description: "Search filter applied to the role",
				Type:        schema.TypeString,
				Computed:    true,
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

	if len(role.SearchIndexesAllowed) > 0 {
		d.Set(searchIndexesAllowedKey, role.SearchIndexesAllowed)
	}

	if len(role.ImportRoles) > 0 {
		d.Set(importRolesKey, role.ImportRoles)
	}

	if len(role.EnabledCapabilityNames()) > 0 {
		d.Set(capabilitiesKey, role.EnabledCapabilityNames())
	}

	if role.SearchFilter != "" {
		d.Set(searchFilterKey, role.SearchFilter)
	}

	return diag.Diagnostics{}
}
