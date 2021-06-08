package provider

import (
	"cd.splunkdev.com/sas/go/libraries/splunkconfig/pkg/config"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	userNamesKey     = "user_names"
	userNamesIDValue = "splunkconfig_user_names"
)

func resourceUserNames() *schema.Resource {
	return &schema.Resource{
		Description: "Return User Names from the Splunk Configuration",
		ReadContext: resourceUserNamesRead,
		Schema: map[string]*schema.Schema{
			userNamesKey: {
				Description: "List of User Names in the Splunk Configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserNamesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	d.SetId(userNamesIDValue)
	d.Set(userNamesKey, suite.Users.Names())

	return diag.Diagnostics{}
}
