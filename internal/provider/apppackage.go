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
	appPackagePathKey = "path"
	appPackageTGZKey  = "tarball_path"
	// the remainder of the fields used by this resource are defined in appautoversion.go,
	// as this resource is being deprecated in favor of it, and this resource makes use of
	// functions defined for its functionality in order to avoid code duplication.
)

func resourceAppPackage() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in favor of the splunkconfig_app_auto_version resource and splunkconfig_app_package data source.",
		Description:        "Create a tarball for an app. Generated app.conf's version will be automatically incremented when app content changes.",
		CustomizeDiff:      resourceAppAutoVersionCustomDiff,
		// create/read/update end up doing the same work, so they use the same function
		CreateContext: resourceAppPackageRead,
		ReadContext:   resourceAppPackageRead,
		UpdateContext: resourceAppPackageRead,
		DeleteContext: resourceAppPackageDelete,
		Schema: map[string]*schema.Schema{
			appAutoVersionAppIDKey: {
				Description: "ID of the app",
				Type:        schema.TypeString,
				Required:    true,
			},
			appPackagePathKey: {
				Description: "Path in which to create the app file",
				Type:        schema.TypeString,
				Required:    true,
			},
			appPackageTGZKey: {
				Description: "Full path of the generated tarball",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appAutoVersionBaseVersionKey: {
				Description: "Version of the app, directly from the provider",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appAutoVersionEffectiveVersionKey: {
				Description: "Version of the app, accounting for patch count",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appAutoVersionPatchCountKey: {
				Description: "Number of patches to the app since setting/changing its version",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			appAutoVersionFilesKey: {
				Description: "File content of the app",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						appAutoVersionFilePathKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						appAutoVersionFileContentKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// resourceAppPackageRead writes the package tarball at the specified location. It is also called for the "create"
// and "update" contexts, because we want to always write the tarball for the app to be used by other downstream
// resources.
func resourceAppPackageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	appID := d.Get(appAutoVersionAppIDKey).(string)
	d.SetId(appID)

	app, err := suite.ExtrapolatedAppWithId(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceAppPackageRead error: %s", err))
	}

	app = app.PlusPatchCount(int64(d.Get(appAutoVersionPatchCountKey).(int)))
	appPath := d.Get(appPackagePathKey).(string)

	tgzFile, err := app.WriteTar(appPath)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(appPackageTGZKey, tgzFile); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// resourceAppPackageDelete does nothing, as the is no deployed infrastructure or configuration to remove.
func resourceAppPackageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
