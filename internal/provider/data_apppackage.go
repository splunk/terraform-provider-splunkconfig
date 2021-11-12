// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
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
	dataAppPackageAppIDKey            = "app_id"
	dataAppPackagePathKey             = "path"
	dataAppPackagePatchCountKey       = "patch_count"
	dataAppPackageEffectiveVersionKey = "effective_version"
	dataAppPackageTGZKey              = "tarball_path"
)

func dataAppPackage() *schema.Resource {
	return &schema.Resource{
		Description: "Create a tarball for an app.",
		ReadContext: dataAppPackageRead,
		Schema: map[string]*schema.Schema{
			dataAppPackageAppIDKey: {
				Description: "ID of the app",
				Type:        schema.TypeString,
				Required:    true,
			},
			dataAppPackagePathKey: {
				Description: "Path in which to create the app file",
				Type:        schema.TypeString,
				Required:    true,
			},
			dataAppPackagePatchCountKey: {
				Description: "Patch count to apply to the app's version",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			dataAppPackageEffectiveVersionKey: {
				Description: "Version of the app, accounting for patch count",
				Type:        schema.TypeString,
				Computed:    true,
			},
			dataAppPackageTGZKey: {
				Description: "Full path of the generated tarball",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataAppPackageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	suite := meta.(config.Suite)

	appID := d.Get(dataAppPackageAppIDKey).(string)
	d.SetId(appID)

	app, err := suite.ExtrapolatedAppWithId(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceAppPackageRead error: %s", err))
	}

	if patchCount, ok := d.GetOk(dataAppPackagePatchCountKey); ok {
		app = app.PlusPatchCount(int64(patchCount.(int)))
	}

	appPath := d.Get(dataAppPackagePathKey).(string)

	tgzFile, err := app.WriteTar(appPath)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(dataAppPackageTGZKey, tgzFile); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(dataAppPackageEffectiveVersionKey, app.Version.AsString()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
