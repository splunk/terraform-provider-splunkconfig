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

package main

import (
	"flag"
	"fmt"
	"os"
	"terraform-provider-splunkconfig/internal/splunkconfig/config"

	"gopkg.in/yaml.v2"
)

func main() {
	lookupName := flag.String("lookupName", "", "Name to give the new Lookup (required)")
	csvFilename := flag.String("csvFilename", "", "CSV file to import (required)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\template-lookup-csv: Print splunkconfig YAML from CSV file\n\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	argsErrors := []string{}
	if *lookupName == "" {
		argsErrors = append(argsErrors, "lookupName is required")
	}
	if *csvFilename == "" {
		argsErrors = append(argsErrors, "csvFilename is required")
	}
	if len(argsErrors) > 0 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "\nArgument errors:\n")
		for _, argsError := range argsErrors {
			fmt.Fprintf(os.Stderr, "  %s\n", argsError)
		}
		os.Exit(1)
	}

	csvFileReader, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open %s: %s\n", *csvFilename, err)
		os.Exit(1)
	}
	defer csvFileReader.Close()

	lookup, err := config.NewLookupFromIoReader(*lookupName, csvFileReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create new lookup from CSV: %s\n", err)
		os.Exit(1)
	}

	suite := config.Suite{
		Lookups: config.Lookups{lookup},
	}

	yamlBytes, err := yaml.Marshal(suite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to marshal lookup: %s\n", err)
		os.Exit(1)
	}

	fmt.Print(string(yamlBytes))
}
