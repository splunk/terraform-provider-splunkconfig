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
	"regexp"
)

// AppID represents the "id" value under the "package" stanza in app.conf.
type AppID string

// validate returns an error if AppID is invalid.  As per Splunk's app.conf .spec:
// * id must adhere to these cross-platform folder name restrictions:
//   - must contain only letters, numbers, "." (dot), and "_" (underscore)
//     characters.
// "-" (dash) characters are also permitted as many official Splunk apps have them in their IDs
func (appID AppID) validate() error {
	validRegex := regexp.MustCompile("^[A-Za-z0-9_.-]+$")

	if !validRegex.MatchString(string(appID)) {
		return fmt.Errorf("AppID is invalid, may only contain letters, numbers, dots, and underscores")
	}

	return nil
}
