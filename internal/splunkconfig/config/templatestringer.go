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
	"bytes"
	"log"
	"text/template"
)

// templateStringer objects implement templateString() which returns its template string.
type templateStringer interface {
	templateString() string
}

// templateContent returns a templateStringer's templated content.
func templateContent(t templateStringer) string {
	textTemplate := template.Must(template.New("ConfTemplate").Parse(t.templateString()))
	buf := new(bytes.Buffer)
	if err := textTemplate.Execute(buf, t); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}
