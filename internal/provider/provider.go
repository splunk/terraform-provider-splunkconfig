// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-splunkconfig/internal/splunkconfig/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	suiteConfigYMLKey               = "configuration"
	suiteConfigFileKey              = "configuration_file"
	roleNamesResourceName           = "splunkconfig_role_names"
	roleAttributesResourceName      = "splunkconfig_role_attributes"
	samlGroupNamesResourceName      = "splunkconfig_saml_group_names"
	samlGroupAttributesResourceName = "splunkconfig_saml_group_attributes"
	appPackageResourceName          = "splunkconfig_app_package"
	appIdsResourceName              = "splunkconfig_app_ids"
	appAttributesResourceName       = "splunkconfig_app_attributes"
	userNamesResourceName           = "splunkconfig_user_names"
	userAttributesResourceName      = "splunkconfig_user_attributes"
	lookupAttributesResourceName    = "splunkconfig_lookup_attributes"
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
					Description: fmt.Sprintf("YAML content containing the abstracted configruation. Either this or `%s` must be set.", suiteConfigFileKey),
				},
				suiteConfigFileKey: &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: fmt.Sprintf("Full path to YAML file containing the abstracted configuration. Either this or `%s` must be set.", suiteConfigYMLKey),
				},
			},

			// data sources schema
			DataSourcesMap: map[string]*schema.Resource{
				roleNamesResourceName:           resourceRoleNames(),
				roleAttributesResourceName:      resourceRoleAttributes(),
				samlGroupNamesResourceName:      resourceSAMLGroupNames(),
				samlGroupAttributesResourceName: resourceSAMLGroupAttributes(),
				userNamesResourceName:           resourceUserNames(),
				userAttributesResourceName:      resourceUserAttributes(),
				lookupAttributesResourceName:    resourceLookupAttributes(),
				appIdsResourceName:              resourceAppIds(),
				appAttributesResourceName:       resourceAppAttributes(),
			},

			// resources schema
			ResourcesMap: map[string]*schema.Resource{
				appPackageResourceName: resourceAppFile(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}
