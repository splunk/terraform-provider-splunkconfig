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
	"path"
)

// IndexPath represents a path to be used for homePath, coldPath, etc.
type IndexPath string

// defaultIndexPath returns a valid default path for an index and path type.  The value will be:
// $SPLUNK_DB/{indexName}/{pathType}
func defaultIndexPath(indexName IndexName, pathType string) IndexPath {
	// path.Join always uses forward slashes, which is intended, to be valid in SplunkCloud
	return IndexPath(path.Join("$SPLUNK_DB", string(indexName), pathType))
}

// firstIndexPath returns the first non-zero IndexPath.  If no non-zero IndexPaths are given, returns ok=false.
func firstIndexPath(indexPaths ...IndexPath) (found IndexPath, ok bool) {
	for _, found = range indexPaths {
		if found != "" {
			ok = true
			return
		}
	}

	return
}

func firstIndexPathString(indexPaths ...IndexPath) (found string, ok bool) {
	foundIndexPath, ok := firstIndexPath(indexPaths...)
	found = string(foundIndexPath)

	return
}
