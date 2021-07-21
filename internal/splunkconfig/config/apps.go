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
	"fmt"
	"reflect"
)

// Apps is a list of App objects.
type Apps []App

// validate returns an error if Apps is invalid.
func (apps Apps) validate() error {
	return allValidNoDuplicates(uniqueValidators(apps))
}

// extrapolated returns a new Apps object with each member App extrapolated with Indexes.
func (apps Apps) extrapolated(indexes Indexes, roles Roles, lookups Lookups) (Apps, error) {
	extrapolatedApps := make(Apps, len(apps))

	for i, app := range apps {
		extrapolatedApp, err := app.extrapolated(indexes, roles, lookups)
		if err != nil {
			return Apps{}, fmt.Errorf("unable to extrapolate app %s: %s", app.Name, err)
		}

		extrapolatedApps[i] = extrapolatedApp
	}

	return extrapolatedApps, nil
}

// WithID returns the App object with the given ID. Returns ok=false if not found.
func (apps Apps) WithID(name string) (found App, ok bool) {
	foundUIDer, ok := withUID(apps, name)
	if !ok {
		return
	}

	foundValue := reflect.ValueOf(foundUIDer)
	found = foundValue.Interface().(App)

	return
}
