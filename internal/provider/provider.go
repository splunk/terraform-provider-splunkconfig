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
	suiteConfigYMLKey           = "configuration"
	suiteConfigFileKey          = "configuration_file"
	suiteConfigPathKey          = "configuration_path"
	roleNamesDataName           = "splunkconfig_role_names"
	roleAttributesdataName      = "splunkconfig_role_attributes"
	samlGroupNamesDataName      = "splunkconfig_saml_group_names"
	samlGroupAttributesDataName = "splunkconfig_saml_group_attributes"
	appPackageDataName          = "splunkconfig_app_package"
	appPackageResourceName      = "splunkconfig_app_package"
	appAutoVersionResourceName  = "splunkconfig_app_auto_version"
	appIdsDataName              = "splunkconfig_app_ids"
	appAttributesDataName       = "splunkconfig_app_attributes"
	userNamesDataName           = "splunkconfig_user_names"
	userAttributesdataName      = "splunkconfig_user_attributes"
	lookupAttributesDataName    = "splunkconfig_lookup_attributes"
	indexNamesDataName          = "splunkconfig_index_names"
	indexAttributesDataName     = "splunkconfig_index_attributes"
)

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		configContent := d.Get(suiteConfigYMLKey).(string)
		configFile := d.Get(suiteConfigFileKey).(string)
		configPath := d.Get(suiteConfigPathKey).(string)

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

		if configPath != "" {
			suite, err := config.NewSuiteFromYAMLPath(configPath)
			if err != nil {
				return config.Suite{}, diag.Errorf("unable to create NewSuiteFromYAMLPath: %s", err)
			}

			return suite, diag.Diagnostics{}
		}

		return config.Suite{}, diag.Errorf("must set %s, %s, or %s", suiteConfigYMLKey, suiteConfigFileKey, suiteConfigPathKey)
	}
}

// New returns a function that returns a pointer to a new schema.Provider for this provider.
func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			// provider schema
			Schema: map[string]*schema.Schema{
				suiteConfigYMLKey: {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{suiteConfigFileKey, suiteConfigPathKey},
					Description:   fmt.Sprintf("YAML content containing the abstracted configuration. Exactly one of `%s`, `%s`, or `%s` must be set.", suiteConfigYMLKey, suiteConfigFileKey, suiteConfigPathKey),
				},
				suiteConfigFileKey: {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{suiteConfigYMLKey, suiteConfigPathKey},
					Description:   fmt.Sprintf("Full path to YAML file containing the abstracted configuration. Exactly one of `%s`, `%s`, or `%s` must be set.", suiteConfigYMLKey, suiteConfigFileKey, suiteConfigPathKey),
				},
				suiteConfigPathKey: {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{suiteConfigYMLKey, suiteConfigFileKey},
					Description:   fmt.Sprintf("Full path to directory containing one or more YAML files containing the abstracted configuration. Exactly one of `%s`, `%s`, or `%s` must be set.", suiteConfigYMLKey, suiteConfigFileKey, suiteConfigPathKey),
				},
			},

			// data sources schema
			DataSourcesMap: map[string]*schema.Resource{
				roleNamesDataName:           dataRoleNames(),
				roleAttributesdataName:      dataRoleAttributes(),
				samlGroupNamesDataName:      dataSAMLGroupNames(),
				samlGroupAttributesDataName: dataSAMLGroupAttributes(),
				userNamesDataName:           dataUserNames(),
				userAttributesdataName:      dataUserAttributes(),
				lookupAttributesDataName:    dataLookupAttributes(),
				appIdsDataName:              dataAppIds(),
				appAttributesDataName:       dataAppAttributes(),
				appPackageDataName:          dataAppPackage(),
				indexNamesDataName:          dataIndexNames(),
				indexAttributesDataName:     dataIndexAttributes(),
			},

			// resources schema
			ResourcesMap: map[string]*schema.Resource{
				appPackageResourceName:     resourceAppPackage(),
				appAutoVersionResourceName: resourceAppAutoVersion(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}
