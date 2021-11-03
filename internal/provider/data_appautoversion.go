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
	appAutoVersionAppIDKey            = "app_id"
	appAutoVersionBaseVersionKey      = "base_version"
	appAutoVersionEffectiveVersionKey = "effective_version"
	appAutoVersionPatchCountKey       = "patch_count"
	appAutoVersionFilesKey            = "files"
	appAutoVersionFilePathKey         = "path"
	appAutoVersionFileContentKey      = "content"
)

func resourceAppAutoVersion() *schema.Resource {
	return &schema.Resource{
		Description:   "Generate an App's version based on content changes.",
		CustomizeDiff: resourceAppAutoVersionCustomDiff,
		// create/read/update end up doing the same work, so they use the same function
		CreateContext: resourceAppAutoVersionCreate,
		ReadContext:   resourceAppAutoVersionNoOp,
		UpdateContext: resourceAppAutoVersionNoOp,
		DeleteContext: resourceAppAutoVersionNoOp,
		Schema: map[string]*schema.Schema{
			appAutoVersionAppIDKey: {
				Description: "ID of the app",
				Type:        schema.TypeString,
				Required:    true,
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

// resourceAppAutoVersionCreate only sets the ID of the resource. All other logic is handled in resourceAppAutoVersionCustomDiff.
func resourceAppAutoVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appId := d.Get(appAutoVersionAppIDKey).(string)
	d.SetId(appId)

	return nil
}

// resourceAppAutoVersionNoOp performs no actions, but Read/Update/Delete must have functions defined for them.
func resourceAppAutoVersionNoOp(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// resourceAppAutoVersionFileContents returns a list of key/value pairs for an App that can be used to set the value
// of the "files" attribute.
func resourceAppAutoVersionFileContents(app config.App) []map[string]string {
	appFiles := app.FileContenters()
	fileContenters := make([]map[string]string, len(appFiles))

	for i, fileContenter := range app.FileContenters() {
		fileContenters[i] = map[string]string{
			appAutoVersionFilePathKey:    fileContenter.FilePath(),
			appAutoVersionFileContentKey: fileContenter.TemplatedContent(),
		}
	}

	return fileContenters
}

// resourceAppAutoVersionCustomDiff calculates and sets all attributes for the resource. This functionality is performed
// as a CustomizeDiff function to enable seeing the calculated views in the terraform plan diff *prior* to the apply.
func resourceAppAutoVersionCustomDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	suite := meta.(config.Suite)

	// CustomizeDiff is called before CreateContext, so we can't use d.Id() here
	appID := d.Get(appAutoVersionAppIDKey).(string)

	app, err := suite.ExtrapolatedAppWithId(appID)
	if err != nil {
		return fmt.Errorf("diff calculation error: %s", err)
	}

	// set the "version" from the fetched-from-suite app
	if err := d.SetNew(appAutoVersionBaseVersionKey, app.Version.AsString()); err != nil {
		return err
	}

	// reset (or set initially to 0) patch count
	if d.HasChange(appAutoVersionBaseVersionKey) {
		if err := d.SetNew(appAutoVersionPatchCountKey, 0); err != nil {
			return err
		}
	}

	// update app with patch count
	appPlusPatchCount := app.PlusPatchCount(int64(d.Get(appAutoVersionPatchCountKey).(int)))
	newVersion := appPlusPatchCount.Version

	if oldVersionInterface, oldVersionStringExists := d.GetOk(appAutoVersionEffectiveVersionKey); oldVersionStringExists && d.HasChange(appAutoVersionBaseVersionKey) {
		oldVersionString := oldVersionInterface.(string)
		oldVersion, err := config.NewVersionFromString(oldVersionString)
		if err != nil {
			return fmt.Errorf("unable to create NewVersionFromString %q: %s", oldVersionString, err)
		}

		if !newVersion.IsGreaterThan(oldVersion) {
			return fmt.Errorf("new effective version %q not greater than old effective version %q", newVersion.AsString(), oldVersionString)
		}
	}

	if err := d.SetNew(appAutoVersionEffectiveVersionKey, newVersion.AsString()); err != nil {
		return fmt.Errorf("unable to SetNew %q: %s", appAutoVersionEffectiveVersionKey, err)
	}

	// set "files" from app with patch count
	if err := d.SetNew(appAutoVersionFilesKey, resourceAppAutoVersionFileContents(appPlusPatchCount)); err != nil {
		return err
	}

	// but if "files" has changes (and "version" doesn't), bump patch count and re-calculate "files"
	// the exclusion of "version" changes is because a version change resets the patch count back to 0, and we don't
	// want to bump it immediately due to the resulting app.conf changes that would cause.
	if !d.HasChange(appAutoVersionBaseVersionKey) && d.HasChange(appAutoVersionFilesKey) {
		oldPatchCount := d.Get(appAutoVersionPatchCountKey).(int)
		newPatchCount := oldPatchCount + 1
		if err := d.SetNew(appAutoVersionPatchCountKey, newPatchCount); err != nil {
			return err
		}

		// re-create appPlusPatchCount from the *original* app to avoid adding patch count to a previously-bumped
		// version
		appPlusPatchCount = app.PlusPatchCount(int64(newPatchCount))
		if err := d.SetNew(appAutoVersionEffectiveVersionKey, appPlusPatchCount.Version.AsString()); err != nil {
			return err
		}

		// because we have a new patch count, we need to re-calculate the file contents to account for the new version
		if err := d.SetNew(appAutoVersionFilesKey, resourceAppAutoVersionFileContents(appPlusPatchCount)); err != nil {
			return nil
		}
	}

	return nil
}
