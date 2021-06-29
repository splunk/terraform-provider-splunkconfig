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

package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a minimal semantic version of the form Major.Minor.Patch.
type Version struct {
	Major int64
	Minor int64
	Patch int64
}

// NewVersionFromString returns a Version from a given string representation.
func NewVersionFromString(versionString string) (Version, error) {
	semVerRegex := regexp.MustCompile(`^[0-9]+(\.[0-9]+){0,2}$`)

	if !semVerRegex.MatchString(versionString) {
		return Version{}, fmt.Errorf("version string %q is invalid", versionString)
	}

	versionComponentStrings := strings.Split(versionString, ".")

	var err error
	var major int64
	var minor int64
	var patch int64

	if major, err = strconv.ParseInt(versionComponentStrings[0], 10, 64); err != nil {
		return Version{}, err
	}

	if len(versionComponentStrings) > 1 {
		if minor, err = strconv.ParseInt(versionComponentStrings[1], 10, 64); err != nil {
			return Version{}, err
		}
	}

	if len(versionComponentStrings) > 2 {
		if patch, err = strconv.ParseInt(versionComponentStrings[2], 10, 64); err != nil {
			return Version{}, err
		}
	}

	return Version{major, minor, patch}, nil
}

// deltaFrom returns a new Version that represents the delta between version and otherVersion.  It is
// effectively the result of version - otherVersion.
func (version Version) deltaFrom(otherVersion Version) Version {
	return Version{
		Major: version.Major - otherVersion.Major,
		Minor: version.Minor - otherVersion.Minor,
		Patch: version.Patch - otherVersion.Patch,
	}
}

// IsGreaterThan returns true if version is greater than otherVersion.
func (version Version) IsGreaterThan(otherVersion Version) bool {
	diffVersion := version.deltaFrom(otherVersion)

	// the first non-zero version level delta (from major -> patch) is sufficient to make a determination
	for _, versionLevel := range []int64{diffVersion.Major, diffVersion.Minor, diffVersion.Patch} {
		if versionLevel > 0 {
			return true
		}
		if versionLevel < 0 {
			return false
		}
	}

	return false
}

// AsString returns a string representation of the Version.
func (version Version) AsString() string {
	return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
}

// UnmarshalYAML implements custom unmarshalling for a Version.  It enables a Version to be unmarshalled from these
// types of content:
// * {major: X, minor: Y, patch: Z}   # explicitly define the Version components
// * X.Y.Z                            # define the Version as a string
func (version *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// realVersion only exists inside this function, and is used to allow attempting to unmarshal into what is
	// really just a Version directly.  Attempting to unmarshal(&Version) from inside this function will result in
	// infinite recursion back into this function, so we need another type to attempt that unmarshalling.
	type realVersion Version

	unmarshalledRealVersion := realVersion{}
	if err := unmarshal(&unmarshalledRealVersion); err == nil {
		*version = Version(unmarshalledRealVersion)
		return nil
	}

	unmarshalledVersionString := ""
	if err := unmarshal(&unmarshalledVersionString); err == nil {
		if *version, err = NewVersionFromString(unmarshalledVersionString); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("unable to unmarshall Version from YAML")
}

// PlusPatchCount returns a new Version that is patchCount higher than the original Version.
func (version Version) PlusPatchCount(patchCount int64) Version {
	newVersion := version

	newVersion.Patch += patchCount

	return newVersion
}
