package provider

import (
	"cd.splunkdev.com/sas/libraries/go/splunk/config/pkg/config"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	suiteConfigYMLKey          = "configuration"
	suiteConfigFileKey         = "configuration_file"
	roleNamesResourceName      = "splunkconfig_role_names"
	roleAttributesResourceName = "splunkconfig_role_attributes"
)

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		configContent := d.Get(suiteConfigYMLKey).(string)
		configFile := d.Get(suiteConfigFileKey).(string)

		if configContent != "" && configFile != "" {
			return config.Suite{}, diag.Errorf("Both %s and %s specified, only one allowed", suiteConfigYMLKey, suiteConfigFileKey)
		}

		if configContent != "" {
			suite, err := config.NewSuiteFromYAML([]byte(configContent))
			if err != nil {
				return config.Suite{}, diag.Errorf("Unable to create NewSuiteFromYAML: %s", err)
			}

			return suite, diag.Diagnostics{}
		}

		if configFile != "" {
			suite, err := config.NewSuiteFromYAMLFile(configFile)
			if err != nil {
				return config.Suite{}, diag.Errorf("Unable to create NewSuiteFromYAMLFile: %s", err)
			}

			return suite, diag.Diagnostics{}
		}

		return config.Suite{}, diag.Errorf("must set %s or %s", suiteConfigYMLKey, suiteConfigFileKey)
	}
}

// New returns a function that returns a pointer to a new schema.Provider for this provider.
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			// provider schema
			Schema: map[string]*schema.Schema{
				suiteConfigYMLKey: &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: fmt.Sprintf("YAML content containing the abstracted configruation. Either this or %s must be set.", suiteConfigFileKey),
				},
				suiteConfigFileKey: &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: fmt.Sprintf("Full path to YAML file containing the abstracted configuration. Either this or %s must be set.", suiteConfigYMLKey),
				},
			},

			// data sources schema
			DataSourcesMap: map[string]*schema.Resource{
				roleNamesResourceName:      resourceRoleNames(),
				roleAttributesResourceName: resourceRoleAttributes(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}
