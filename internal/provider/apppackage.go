/*
 * Copyright 2021 Splunk, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-splunkconfig/internal/splunkconfig/config"
)

const (
	appPackageIDKey               = "app_id"
	appPackagePathKey             = "path"
	appPackageTGZKey              = "tarball_path"
	appPackageBaseVersionKey      = "base_version"
	appPackageEffectiveVersionKey = "effective_version"
	appPackagePatchCountKey       = "patch_count"
	appPackageFilesKey            = "files"
	appPackageFilePathKey         = "path"
	appPackageFileContentKey      = "content"
)

func resourceAppFile() *schema.Resource {
	return &schema.Resource{
		Description:   "Create a tarball for an app",
		CustomizeDiff: resourceAppPackageCustomDiff,
		CreateContext: resourceAppPackageCreate,
		// create/read/update end up doing the same work, so they use the same function
		ReadContext:   resourceAppPackageCreate,
		UpdateContext: resourceAppPackageCreate,
		DeleteContext: resourceAppPackageDelete,
		Schema: map[string]*schema.Schema{
			appPackageIDKey: {
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
			appPackageBaseVersionKey: {
				Description: "Version of the app, directly from the provider",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appPackageEffectiveVersionKey: {
				Description: "Version of the app, accounting for patch count",
				Type:        schema.TypeString,
				Computed:    true,
			},
			appPackagePatchCountKey: {
				Description: "Number of patches to the app since setting/changing its version",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			appPackageFilesKey: {
				Description: "File content of the app",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						appPackageFilePathKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						appPackageFileContentKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// resourceAppPackageGetApp returns the given App from an appID. It returns an error if one was encountered.
func resourceAppPackageGetApp(appID string, meta interface{}) (config.App, error) {
	suite := meta.(config.Suite)

	apps, err := suite.ExtrapolatedApps()
	if err != nil {
		return config.App{}, fmt.Errorf("resourceAppPackageGetApp unable to extrapolate apps: %s", err)
	}

	app, ok := apps.WithID(appID)
	if !ok {
		return config.App{}, fmt.Errorf("resourceAppPackageGetApp unable to find app with ID %q", appID)
	}

	return app, nil
}

// resourceAppPackageCreate writes the package tarball at the specified location. It is also called for the "read"
// and "update" contexts, because we want to always write the tarball for the app to be used by other downstream
// resources.
func resourceAppPackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appID := d.Get(appPackageIDKey).(string)
	d.SetId(appID)

	app, err := resourceAppPackageGetApp(d.Id(), meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("resourceAppPackageCreate error: %s", err))
	}

	app = app.PlusPatchCount(int64(d.Get(appPackagePatchCountKey).(int)))
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

// resourceAppPackageFileContents returns a list of key/value pairs for an App that can be used to set the value
// of the "files" attribute.
func resourceAppPackageFileContents(app config.App) []map[string]string {
	appFiles := app.FileContenters()
	fileContenters := make([]map[string]string, len(appFiles))

	for i, fileContenter := range app.FileContenters() {
		fileContenters[i] = map[string]string{
			appPackageFilePathKey:    fileContenter.FilePath(),
			appPackageFileContentKey: fileContenter.TemplatedContent(),
		}
	}

	return fileContenters
}

// resourceAppPackageCustomDiff calculates and sets all attributes for the resource. This functionality is performed
// as a CustomizeDiff function to enable seeing the calculated views in the terraform plan diff *prior* to the apply.
func resourceAppPackageCustomDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// CustomizeDiff is called before CreateContext, so we can't use d.Id() here
	appID := d.Get(appPackageIDKey).(string)

	app, err := resourceAppPackageGetApp(appID, meta)
	if err != nil {
		return fmt.Errorf("diff calculation error: %s", err)
	}

	// set the "version" from the fetched-from-suite app
	if err := d.SetNew(appPackageBaseVersionKey, app.Version.AsString()); err != nil {
		return err
	}

	// reset (or set initially to 0) patch count
	if d.HasChange(appPackageBaseVersionKey) {
		if err := d.SetNew(appPackagePatchCountKey, 0); err != nil {
			return err
		}
	}

	// update app with patch count
	appPlusPatchCount := app.PlusPatchCount(int64(d.Get(appPackagePatchCountKey).(int)))
	newVersion := appPlusPatchCount.Version

	if oldVersionInterface, oldVersionStringExists := d.GetOk(appPackageEffectiveVersionKey); oldVersionStringExists && d.HasChange(appPackageBaseVersionKey) {
		oldVersionString := oldVersionInterface.(string)
		oldVersion, err := config.NewVersionFromString(oldVersionString)
		if err != nil {
			return fmt.Errorf("unable to create NewVersionFromString %q: %s", oldVersionString, err)
		}

		if !newVersion.IsGreaterThan(oldVersion) {
			return fmt.Errorf("new effective version %q not greater than old effective version %q", newVersion.AsString(), oldVersionString)
		}
	}

	if err := d.SetNew(appPackageEffectiveVersionKey, newVersion.AsString()); err != nil {
		return fmt.Errorf("unable to SetNew %q: %s", appPackageEffectiveVersionKey, err)
	}

	// set "files" from app with patch count
	if err := d.SetNew(appPackageFilesKey, resourceAppPackageFileContents(appPlusPatchCount)); err != nil {
		return err
	}

	// but if "files" has changes (and "version" doesn't), bump patch count and re-calculate "files"
	// the exclusion of "version" changes is because a version change resets the patch count back to 0, and we don't
	// want to bump it immediately due to the resulting app.conf changes that would cause.
	if !d.HasChange(appPackageBaseVersionKey) && d.HasChange(appPackageFilesKey) {
		oldPatchCount := d.Get(appPackagePatchCountKey).(int)
		newPatchCount := oldPatchCount + 1
		if err := d.SetNew(appPackagePatchCountKey, newPatchCount); err != nil {
			return err
		}

		// re-create appPlusPatchCount from the *original* app to avoid adding patch count to a previously-bumped
		// version
		appPlusPatchCount = app.PlusPatchCount(int64(newPatchCount))
		if err := d.SetNew(appPackageEffectiveVersionKey, appPlusPatchCount.Version.AsString()); err != nil {
			return err
		}

		// because we have a new patch count, we need to re-calculate the file contents to account for the new version
		if err := d.SetNew(appPackageFilesKey, resourceAppPackageFileContents(appPlusPatchCount)); err != nil {
			return nil
		}
	}

	return nil
}
