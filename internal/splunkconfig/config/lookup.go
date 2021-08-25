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
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
)

// Lookup is a combination of a LookupDefinition and a list of LookupRows to make a complete representation of a lookup.
type Lookup struct {
	Name            string
	Fields          LookupFields
	ExternalCommand string `yaml:"external_cmd"`
	ExternalType    string `yaml:"external_type"`
	Collection      string
	Rows            LookupRows
}

// validate returns an error if the Lookup is invalid. It is invalid if its Definition is invalid, or if its Rows
// are invalid in the context of the Definition's Fields.
func (lookup Lookup) validate() error {
	if lookup.Name == "" {
		return fmt.Errorf("invalid Lookup, has an invalid Name: %s", lookup.Name)
	}

	if err := lookup.Fields.validate(); err != nil {
		return fmt.Errorf("invalid Lookup, has invalid Fields: %s", err)
	}

	if err := lookup.Rows.validateForLookupFields(lookup.Fields); err != nil {
		return fmt.Errorf("invalid Lookup, has invalid Rows: %s", err)
	}

	return nil
}

// uid returns the Name of the lookup to determine uniqueness.
func (lookup Lookup) uid() string {
	return lookup.Name
}

// writeCSV writes a Lookup's header and rows to an io.Writer.
func (lookup Lookup) writeCSV(writer io.Writer) error {
	w := csv.NewWriter(writer)

	if err := w.Write(lookup.Fields.headerValues()); err != nil {
		return fmt.Errorf("unable to write csv header: %s", err)
	}

	for _, row := range lookup.Rows {
		if err := w.Write(row.valuesForLookupFields(lookup.Fields)); err != nil {
			return fmt.Errorf("unable to write csv row: %s", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return fmt.Errorf("unable to write csv: %s", err)
	}

	return nil
}

// filename returns the filename to use when writing Lookup to a file.
func (lookup Lookup) filename() string {
	return fmt.Sprintf("%s.csv", lookup.Name)
}

// FilePath returns the path for the lookup file, relative to the app directory.
func (lookup Lookup) FilePath() string {
	return fmt.Sprintf("lookups/%s", lookup.filename())
}

// TemplatedContent returns the templated CSV content, or an empty string if no content should be created.
func (lookup Lookup) TemplatedContent() string {
	if lookup.ExternalType != "" {
		return ""
	}

	buf := new(bytes.Buffer)

	if err := lookup.writeCSV(buf); err != nil {
		log.Fatalf("unable to template CSV content: %s", err)
	}

	return buf.String()
}

// stanzaValues returns the StanzaValues for the Lookup.
func (lookup Lookup) stanzaValues() StanzaValues {
	stanzaValues := StanzaValues{
		"filename": lookup.filename(),
	}

	if lookup.ExternalType != "" {
		stanzaValues["external_type"] = lookup.ExternalType
		// fields_list only comes into play when external_type is set
		stanzaValues["fields_list"] = strings.Join(lookup.Fields.FieldNames(), ", ")
	}

	if lookup.ExternalCommand != "" {
		stanzaValues["external_cmd"] = lookup.ExternalCommand
	}

	if lookup.Collection != "" {
		stanzaValues["collection"] = lookup.Collection
	}

	return stanzaValues
}

// stanza returns the Stanza for the Lookup.
func (lookup Lookup) stanza() Stanza {
	return Stanza{
		Name:   lookup.Name,
		Values: lookup.stanzaValues(),
	}
}

// defaultRow returns the default row for a defaultLookupValuesDefiner if it would be valid for this Lookup's Fields.
// if the defaultRow would not be valid, an empty row is returned with ok=false.
func (lookup Lookup) defaultRow(definer defaultLookupValuesDefiner) (defaultRow LookupRow, ok bool) {
	defaultValues := lookup.Fields.defaultLookupValuesDefinerValues(definer)

	// our default row will be a LookupRow with this Lookup's name, and the defaultValues for the definer/fields
	defaultRow = LookupRow{
		LookupName: lookup.Name,
		Values:     defaultValues,
	}

	// if the lookup's LookupFields don't have any default row fields, don't create a default row.
	if !lookup.Fields.hasDefaultRowFields() {
		return LookupRow{}, false
	}

	// if the defaultValues is missing any default row fields, don't create a default row.
	if !lookup.Fields.valuesHaveAllDefaultRowFields(defaultValues) {
		return LookupRow{}, false
	}

	// if we made it this far, return our default row.
	return defaultRow, true
}

// defaultRows returns a LookupRows object for the default LookupRows for this Lookup based on a
// defaultLookupValuesDefiner's default values.  if no valid default row for this combination exists, LookupRows
// will be empty.
func (lookup Lookup) defaultRows(definer defaultLookupValuesDefiner) LookupRows {
	defaultRow, ok := lookup.defaultRow(definer)
	if !ok {
		return LookupRows{}
	}

	return LookupRows{defaultRow}
}

// extrapolatedWithLookupRowsForLookupDefiners returns a new Lookup which includes rows for the given
// lookupRowsForLookupDefiners.
func (lookup Lookup) extrapolatedWithLookupRowsForLookupDefiners(definers ...lookupRowsForLookupDefiner) Lookup {
	extrapolatedLookup := lookup
	extrapolatedRows := extrapolatedLookup.Rows

	for _, definer := range definers {
		extrapolatedRows = append(extrapolatedRows, definer.lookupRowsForLookup(lookup)...)
	}

	extrapolatedLookup.Rows = extrapolatedRows

	return extrapolatedLookup
}
