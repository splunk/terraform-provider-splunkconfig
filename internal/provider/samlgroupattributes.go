package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"cd.splunkdev.com/sas/go/libraries/splunkconfig/pkg/config"
)

const (
	samlGroupNameKey = "role_name"
	rolesKey         = "roles"
)

func resourceSAMLGroupAttributes() *schema.Resource {
	return &schema.Resource{
		Description: "Get attributes for a specific SAML group",
		ReadContext: resourceSAMLGroupAttributesRead,
		Schema: map[string]*schema.Schema{
			samlGroupNameKey: {
				Description: "Name of the SAML group",
				Type:        schema.TypeString,
				Required:    true,
			},
			rolesKey: {
				Description: "List of roles associated with the SAML group",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSAMLGroupAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	samlGroupName := d.Get(roleNameKey).(string)

	d.SetId(samlGroupName)

	samlGroup, ok := suite.ExtrapolatedSAMLGroups().WithSAMLGroupName(samlGroupName)
	if !ok {
		return diag.Errorf("Unable to find SAML group with name %q", samlGroupName)
	}

	if len(samlGroup.Roles) > 0 {
		d.Set(rolesKey, samlGroup.Roles)
	}

	return diag.Diagnostics{}
}
