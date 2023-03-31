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
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestApp_validate(t *testing.T) {
	tests := validatorTestCases{
		// App has empty Name and ID, because they are not set
		{
			App{},
			true,
		},
		// App has empty ID, because it was not set
		{
			App{
				Name: "My App",
			},
			true,
		},
		// App has empty Name, because it was not set
		{
			App{
				ID: "my_app",
			},
			true,
		},
		// App has invalid IndexesPlaceholder
		{
			App{
				Name: "My App",
				ID:   "my_app",
				IndexesPlaceholder: IndexesPlaceholder{
					Indexes: Indexes{
						// empty name is invalid
						Index{},
					},
				},
			},
			true,
		},
		// valid app (minimum valid Name and ID set)
		{
			App{
				Name: "My App",
				ID:   "my_app",
			},
			false,
		},
	}

	tests.test(t)
}

func TestApp_extrapolated(t *testing.T) {
	tests := []struct {
		app           App
		indexes       Indexes
		roles         Roles
		lookups       Lookups
		wantIndexes   Indexes
		wantConfFiles ConfFiles
		wantError     bool
	}{
		// empty app defines no indexes, doesn't set ImportIndexes, inherits no indexes
		{
			App{IndexesPlaceholder: IndexesPlaceholder{}},
			Indexes{Index{Name: "index_a"}},
			Roles{},
			Lookups{},
			Indexes(nil),
			ConfFiles{
				ConfFile{Name: "indexes", Stanzas: Stanzas{}},
				ConfFile{Name: "authorize", Stanzas: Stanzas{}},
			},
			false,
		},
		// app defines its own indexes
		{
			App{IndexesPlaceholder: IndexesPlaceholder{Indexes: Indexes{Index{Name: "index_a"}}}},
			Indexes{},
			Roles{},
			Lookups{},
			Indexes{Index{Name: "index_a"}},
			ConfFiles{
				ConfFile{
					Name: "indexes",
					Stanzas: Stanzas{
						Stanza{
							Name: "index_a",
							Values: StanzaValues{
								"homePath":   "$SPLUNK_DB/index_a/db",
								"coldPath":   "$SPLUNK_DB/index_a/colddb",
								"thawedPath": "$SPLUNK_DB/index_a/thaweddb",
							},
						},
					},
				},
				ConfFile{Name: "authorize", Stanzas: Stanzas{}},
			},
			false,
		},
		// app is configured to use external indexes (ImportIndexes = true)
		{
			App{IndexesPlaceholder: IndexesPlaceholder{Import: true}},
			Indexes{Index{Name: "index_a"}},
			Roles{},
			Lookups{},
			Indexes{Index{Name: "index_a"}},
			ConfFiles{
				ConfFile{
					Name: "indexes",
					Stanzas: Stanzas{
						Stanza{
							Name: "index_a",
							Values: StanzaValues{
								"homePath":   "$SPLUNK_DB/index_a/db",
								"coldPath":   "$SPLUNK_DB/index_a/colddb",
								"thawedPath": "$SPLUNK_DB/index_a/thaweddb",
							},
						},
					},
				},
				ConfFile{Name: "authorize", Stanzas: Stanzas{}},
			},
			false,
		},
	}

	for _, test := range tests {
		extrapolatedApp, err := test.app.extrapolated(test.indexes, test.roles, test.lookups)

		gotError := err != nil
		messageError := fmt.Sprintf(
			"%#v.extrapolated(%#v, %#v) returned error?",
			test.app,
			test.indexes,
			test.lookups)
		testEqual(gotError, test.wantError, messageError, t)

		gotIndexes := extrapolatedApp.IndexesPlaceholder.Indexes
		messageIndexes := fmt.Sprintf(
			"%#v.extrapolated(%#v, %#v).IncludeIndexes",
			test.app,
			test.indexes,
			test.lookups)
		testEqual(gotIndexes, test.wantIndexes, messageIndexes, t)

		gotConfFiles := extrapolatedApp.ConfFiles
		messageConfFiles := fmt.Sprintf(
			"%#v.extrapolated(%#v, %#v).ConfFiles",
			test.app,
			test.indexes,
			test.lookups)
		testEqual(gotConfFiles, test.wantConfFiles, messageConfFiles, t)
	}
}

func TestApp_PlusPatchCount(t *testing.T) {
	tests := []struct {
		inputApp        App
		inputPatchCount int64
		wantApp         App
	}{
		{
			App{},
			0,
			App{},
		},
		{
			App{},
			1,
			App{
				Version: Version{Patch: 1},
			},
		},
		{
			App{Version: Version{Patch: 1}},
			1,
			App{
				Version: Version{Patch: 2},
			},
		},
	}

	for _, test := range tests {
		gotApp := test.inputApp.PlusPatchCount(test.inputPatchCount)
		message := fmt.Sprintf("%#v.PlusPatchCount(%d)", test.inputPatchCount, test.inputPatchCount)

		testEqual(gotApp, test.wantApp, message, t)
	}
}

func TestApp_consistentTarball(t *testing.T) {
	app := App{Name: "Test App", ID: "test_app"}

	shaA := appSha1(app, t)
	shaB := appSha1(app, t)

	if shaA != shaB {
		t.Error("app created with differing SHA values")
	}
}

func appSha1(app App, t *testing.T) string {
	tempdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unable to create tempdir: %s", err)
	}
	defer os.RemoveAll(tempdir)

	tgzFile, err := app.WriteTar(tempdir)
	if err != nil {
		t.Fatalf("unable to create tgzFile: %s", err)
	}

	data, err := os.Open(tgzFile)
	if err != nil {
		t.Fatalf("unable to open tgzFile: %s", err)
	}
	defer data.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, data); err != nil {
		t.Fatalf("unable to calculate app SHA: %s", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
