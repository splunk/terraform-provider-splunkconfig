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
	"path/filepath"
)

// FileContenter objects implement FilePath and TemplatedContent.
type FileContenter interface {
	FilePath() string
	TemplatedContent() string
}

// writeTarFileContents writes fileContenter contents to its filepath for a tar.Writer.
func writeTarFileContents(contenter FileContenter, tw *tar.Writer, basePath string) error {
	templatedContent := contenter.TemplatedContent()
	templatedLen := len(templatedContent)

	hdr := &tar.Header{
		Name: filepath.Join(basePath, contenter.FilePath()),
		Mode: 0644,
		Size: int64(templatedLen),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("unable to write tar header for file %s: %s", contenter.FilePath(), err)
	}

	if _, err := tw.Write([]byte(templatedContent)); err != nil {
		return fmt.Errorf("unable to write tar content for file %s: %s", contenter.FilePath(), err)
	}

	return nil
}
