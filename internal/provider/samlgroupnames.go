package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"cd.splunkdev.com/sas/go/libraries/splunkconfig/pkg/config"
)

const (
	samlGroupNamesKey     = "saml_group_names"
	samlGroupNamesIDValue = "splunkconfig_saml_group_names"
)

func resourceSAMLGroupNames() *schema.Resource {
	return &schema.Resource{
		Description: "Return SAML Group Names from the Splunk Configuration",
		ReadContext: resourceSAMLGroupNamesRead,
		Schema: map[string]*schema.Schema{
			samlGroupNamesKey: {
				Description: "List of SAML Group Names in the Splunk Configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceSAMLGroupNamesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	d.SetId(samlGroupNamesIDValue)
	d.Set(samlGroupNamesKey, suite.ExtrapolatedSAMLGroups().SAMLGroupNames())

	return diag.Diagnostics{}
}
