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

package config

import (
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// App represents a Splunk app that can be packaged into a tarball and installed in $SPLUNK_HOME/etc/apps.
type App struct {
	Name               string
	Description        string
	ID                 AppID
	Author             string
	IsVisible          bool
	CheckForUpdates    bool
	Version            Version
	ConfFiles          ConfFiles
	IndexesPlaceholder IndexesPlaceholder `yaml:"indexes"`
	RolesPlaceholder   RolesPlaceholder   `yaml:"roles"`
	LookupsPlaceholder LookupsPlaceholder `yaml:"lookups"`
}

// validate returns an error if App is invalid.  It is invalid if:
// * it has an empty Name
// * has an invalid ID
// * has invalid ConfFiles
// * has invalid IndexesPlaceholder
// * has invalid RolesPlaceholder
// * has invalid LookupsPlaceholder
func (app App) validate() error {
	if app.Name == "" {
		return fmt.Errorf("invalid App (%v), has an empty Name", app)
	}

	validators := map[string]validator{
		"ID":                 app.ID,
		"ConfFiles":          app.ConfFiles,
		"IndexesPlaceholder": app.IndexesPlaceholder,
		"RolesPlaceholder":   app.RolesPlaceholder,
		"LookupsPlaceholder": app.LookupsPlaceholder,
	}

	for vName, v := range validators {
		if err := v.validate(); err != nil {
			return fmt.Errorf("invalid App, has invalid %s: %s", vName, err)
		}
	}

	return nil
}

// uid returns the ID of the App to determine uniqueness.
func (app App) uid() string {
	return string(app.ID)
}

// appStanzas returns the Stanzas for an App's app.conf.
func (app App) appStanzas() Stanzas {
	return Stanzas{
		Stanza{
			Name: "ui",
			Values: StanzaValues{
				"label":      app.Name,
				"is_visible": strconv.FormatBool(app.IsVisible),
			},
		},
		Stanza{
			Name: "launcher",
			Values: StanzaValues{
				"author":      app.Author,
				"description": app.Description,
				"version":     app.Version.AsString(),
			},
		},
		Stanza{
			Name: "package",
			Values: StanzaValues{
				"id":                string(app.ID),
				"check_for_updates": strconv.FormatBool(app.CheckForUpdates),
			},
		},
	}
}

// appConfFile returns a ConfFile for an App's app.conf.
func (app App) appConfFile() ConfFile {
	return ConfFile{
		Name:    "app",
		Stanzas: app.appStanzas(),
	}
}

// extrapolated returns a new copy of App that has external components (Indexes, Lookups) substituted for any true
// placeholders.
func (app App) extrapolated(indexes Indexes, roles Roles, lookups Lookups) (App, error) {
	newApp := app

	extrapolatedIndexes := app.IndexesPlaceholder.selectedIndexes(indexes)
	newApp.IndexesPlaceholder = IndexesPlaceholder{Indexes: extrapolatedIndexes}

	extrapolatedRoles := app.RolesPlaceholder.selectedRoles(roles)
	newApp.RolesPlaceholder = RolesPlaceholder{Roles: extrapolatedRoles}

	newApp.ConfFiles = newApp.ConfFiles.WithConfFile(extrapolatedIndexes.confFile())
	newApp.ConfFiles = newApp.ConfFiles.WithConfFile(extrapolatedRoles.confFile())

	extrapolatedLookups, err := app.LookupsPlaceholder.selectedLookups(lookups)
	if err != nil {
		return App{}, fmt.Errorf("unable to extrapolate App %s: %s", app.Name, err)
	}

	newApp.LookupsPlaceholder = LookupsPlaceholder{Lookups: extrapolatedLookups}

	return newApp, nil
}

// tarFilename returns the filename to use when when creating a tarfile for an App.  It is:
// <app name>-<version>.tgz
func (app App) tarFilename() string {
	return fmt.Sprintf("%s-%s.tgz", app.ID, app.Version.AsString())
}

// FileContenters returns FileContenters for the App.
func (app App) FileContenters() FileContenters {
	contenters := FileContenters{app.appConfFile()}
	contenters = append(contenters, NewFileContentersFromList(app.ConfFiles)...)
	contenters = append(contenters, app.LookupsPlaceholder.Lookups.fileContenters()...)

	return contenters.WithContent()
}

// WriteTar creates a tarfile for this app at the given path.  It returns the absolute path of the created tarball, or
// an error if one was encountered.
func (app App) WriteTar(path string) (tgzPath string, err error) {
	tarFilePath := filepath.Join(path, app.tarFilename())
	tarFileAbsPath, err := filepath.Abs(tarFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to create absolute path from %q: %s", tarFilePath, err)
	}

	tf, err := os.Create(tarFileAbsPath)
	if err != nil {
		return "", fmt.Errorf("unable to create %s: %s", tarFilePath, err)
	}
	defer tf.Close()

	tw := tar.NewWriter(tf)
	defer tw.Close()

	for _, contenter := range app.FileContenters() {
		if err := writeTarFileContents(contenter, tw, string(app.ID)); err != nil {
			return "", fmt.Errorf("App unable to write tar contents: %s", err)
		}
	}

	return tarFileAbsPath, nil
}

// PlusPatchCount returns a new App with a Version adjusted for changes.
func (app App) PlusPatchCount(patchCount int64) App {
	newApp := app
	newApp.Version = newApp.Version.PlusPatchCount(patchCount)

	return newApp
}
