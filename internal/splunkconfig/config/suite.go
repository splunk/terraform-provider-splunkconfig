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
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"
)

// Suite is a type that represents a collection of component configurations.
type Suite struct {
	Indexes    Indexes    `yaml:"indexes,omitempty"`
	Roles      Roles      `yaml:"roles,omitempty"`
	SAMLGroups SAMLGroups `yaml:"saml_groups,omitempty"`
	Lookups    Lookups    `yaml:"lookups,omitempty"`
	Apps       Apps       `yaml:"apps,omitempty"`
	Users      Users      `yaml:"users,omitempty"`
	// Anchors isn't actually part of the configuration, it just gives you somewhere to define
	// YAML anchors while still disallowing unknown keys.
	Anchors interface{} `yaml:"anchors,omitempty"`
}

// validate returns an error if any of a Suite's configurations are invalid.
func (suite Suite) validate() error {
	if err := suite.Indexes.validate(); err != nil {
		return err
	}

	if err := suite.Roles.validate(); err != nil {
		return err
	}

	if err := suite.SAMLGroups.validate(); err != nil {
		return err
	}

	// if an Index references a Role that doesn't exist, fail validation
	if err := suite.Indexes.validateWithRoles(suite.Roles); err != nil {
		return err
	}

	// if an Index references a Lookup that doesn't exist, fail validation
	if err := suite.Indexes.validateWithLookups(suite.Lookups); err != nil {
		return err
	}

	// if a Role references a Lookup that doesn't exist, fail validation
	if err := suite.Roles.validateForLookups(suite.Lookups); err != nil {
		return err
	}

	// if a Role references a SAMLGroup that doesn't exist, fail validation
	if err := suite.Roles.validateForSAMLGroups(suite.SAMLGroups); err != nil {
		return err
	}

	// validate extrapolated lookups
	// there's no reason to validate lookups prior to extrapolation, because extrapolation can't fix them
	if err := suite.ExtrapolatedLookups().validate(); err != nil {
		return err
	}

	if err := suite.Apps.validate(); err != nil {
		return err
	}

	if err := suite.Users.validate(); err != nil {
		return err
	}

	return nil
}

// NewSuiteFromYAML returns a new Suite object from the YAML contents passed in. It returns an error if any errors
// were encountered while attempting to unmarshal the content or if the resulting Suite is invalid.
func NewSuiteFromYAML(yamlContent []byte) (suite Suite, err error) {
	decoder := yaml.NewDecoder(bytes.NewReader(yamlContent))
	decoder.SetStrict(true)

	if err = decoder.Decode(&suite); err != nil {
		return
	}

	if err = suite.validate(); err != nil {
		// return empty Suite object if invalid
		suite = Suite{}
		return
	}

	return
}

// NewSuiteFromYAMLFile returns a new Suite object from the YAML contents passed in. It returns an error if any errors
// were encountered while attempting to unmarshal the content or if the resulting Suite is invalid.
func NewSuiteFromYAMLFile(path string) (suite Suite, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	suite, err = NewSuiteFromYAML(content)
	return
}

// mergeSuite returns a new Suite by merging the contents of additionalSuite.  This method assumes that Suite
// only ever contains members that are slices (other than Anchors, which is not merged).
func (suite Suite) mergeSuite(additionalSuite Suite) (mergedSuite Suite) {
	suiteV := reflect.ValueOf(&suite).Elem()
	additionalSuiteV := reflect.ValueOf(&additionalSuite).Elem()
	mergedSuiteV := reflect.ValueOf(&mergedSuite).Elem()

	// perform merge for each field (except Anchors)
	for i := 0; i < suiteV.NumField(); i++ {
		suiteField := suiteV.Type().Field(i)

		// no merging occurs for Anchors, they're only used within a single YAML file
		if suiteField.Name == "Anchors" {
			continue
		}

		// get this field's reflect.Value for existing, additional, merged Suites
		suiteFieldValue := suiteV.Field(i)
		additionalSuiteFieldValue := additionalSuiteV.Field(i)
		mergedSuiteFieldValue := mergedSuiteV.Field(i)

		// set merged Suite's value to the result of appending existing and additional values
		mergedSuiteFieldValue.Set(reflect.AppendSlice(suiteFieldValue, additionalSuiteFieldValue))
	}

	return
}

// mergeSuites returns a new Suite by merging the contents of additionalSuites.
func (suite Suite) mergeSuites(additionalSuites ...Suite) (mergedSuite Suite) {
	// starting state is identical to suite
	mergedSuite = suite

	// add on each additionalSuite
	for _, additionalSuite := range additionalSuites {
		mergedSuite = mergedSuite.mergeSuite(additionalSuite)
	}

	return
}

// ExtrapolatedRoles returns the Suite's Roles extrapolated against its Indexes.
func (suite Suite) ExtrapolatedRoles() Roles {
	return suite.Roles.extrapolateWithIndexes(suite.Indexes)
}

// ExtrapolatedSAMLGroups returns the Suite's SAMLGroups extrapolated against its Roles.
func (suite Suite) ExtrapolatedSAMLGroups() SAMLGroups {
	return suite.SAMLGroups.extrapolateWithRoles(suite.Roles)
}

// ExtrapolatedLookups returns the Suite's Lookups extrapolated against its Indexes and Roles.
func (suite Suite) ExtrapolatedLookups() Lookups {
	return suite.Lookups.extrapolatedWithLookupRowsForLookupDefiners(suite.Indexes, suite.Roles)
}

// ExtrapolatedApps returns the Suite's Apps extrapolated against its Indexes.
func (suite Suite) ExtrapolatedApps() (Apps, error) {
	extrapolatedApps, err := suite.Apps.extrapolated(suite.Indexes, suite.ExtrapolatedRoles(), suite.ExtrapolatedLookups())
	if err != nil {
		return Apps{}, fmt.Errorf("ExtrapolatedApps error: %s", err)
	}

	return extrapolatedApps, nil
}
