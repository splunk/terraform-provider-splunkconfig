package provider

import (
	"cd.splunkdev.com/sas/libraries/go/splunk/config/pkg/config"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	roleNamesKey     = "role_names"
	roleNamesIDValue = "splunkconfig_role_names"
)

func resourceRoleNames() *schema.Resource {
	return &schema.Resource{
		Description: "Return Role Names from the Splunk Configuration",
		ReadContext: resourceRoleNamesRead,
		Schema: map[string]*schema.Schema{
			roleNamesKey: {
				Description: "List of Role Names in the Splunk Configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceRoleNamesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	d.SetId(roleNamesIDValue)
	d.Set(roleNamesKey, suite.ExtrapolatedRoles().RoleNames())

	return diag.Diagnostics{}
}
