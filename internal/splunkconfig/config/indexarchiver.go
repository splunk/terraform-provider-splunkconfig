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

package config

import "fmt"

// Archiver represents the value of an Index's "archiver.coldStorageProvider" fields.
type Archiver string

const (
	ARCHIVERUNDEF  Archiver = ""
	ARCHIVERAWS    Archiver = "Glacier"
	ARCHIVERGCP    Archiver = "GCSArchive"
)

// validate returns an error if Archiver is invalid. It is invalid if:
// * it isn't one of the defined constants
func (archiver Archiver) validate() error {
	switch archiver {
	case ARCHIVERUNDEF, ARCHIVERAWS, ARCHIVERGCP:
		break
	default:
		return fmt.Errorf("invalid Archiver value: %s", archiver)
	}

	return nil
}
