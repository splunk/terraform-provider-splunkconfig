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
	"path"
)

// ConfFile represents a configuration file.
type ConfFile struct {
	Name      string
	Extension string
	Location  string
	Stanzas   Stanzas
}

// validate returns an error if ConfFile is invalid.  It is invalid if it:
// * has an empty Name
// * has Stanza collisions
func (confFile ConfFile) validate() error {
	if confFile.Name == "" {
		return fmt.Errorf("invalid %T{%+v}, has empty name", confFile, confFile)
	}

	if err := confFile.Stanzas.validate(); err != nil {
		return fmt.Errorf("invalid %T{%+v}, has invalid Stanzas: %s", confFile, confFile, err)
	}

	return nil
}

// uid returns the name of the ConfFile to determine uniqueness.
func (confFile ConfFile) uid() string {
	return confFile.Name
}

// filename returns the filename for a configuration file, which is just its Name with a suffix of .conf.
func (confFile ConfFile) filename() string {
	extension := "conf"
	if confFile.Extension != "" {
		extension = confFile.Extension
	}

	return fmt.Sprintf("%s.%s", confFile.Name, extension)
}

// FilePath returns the path for the configuration file, relative to the app directory.
func (confFile ConfFile) FilePath() string {
	location := "default"
	if confFile.Location != "" {
		location = confFile.Location
	}

	return path.Join(location, confFile.filename())
}

// WithStanzas returns a new ConfFile that is a copy of this ConfFile with additional Stanzas added.  Stanzas are not
// merged by Name, and instead are simply appended.
func (confFile ConfFile) WithStanzas(stanzas Stanzas) ConfFile {
	newConfFile := confFile
	newStanzas := append(confFile.Stanzas, stanzas...)
	newConfFile.Stanzas = newStanzas

	return newConfFile
}

// TemplatedContent returns the content for a ConfFile.
func (confFile ConfFile) TemplatedContent() string {
	return confFile.Stanzas.TemplatedContent()
}
