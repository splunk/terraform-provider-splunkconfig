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
	"testing"
)

func TestIndexesPlaceholder_validate(t *testing.T) {
	tests := validatorTestCases{
		// empty placeholder needs to be valid
		{
			IndexesPlaceholder{},
			false,
		},
		// has only valid indexes
		{
			IndexesPlaceholder{
				Indexes: Indexes{
					Index{Name: "valid_index"},
				},
			},
			false,
		},
		// has invalid index
		{
			IndexesPlaceholder{
				Indexes: Indexes{
					Index{Name: "invalid index"},
				},
			},
			true,
		},
	}

	tests.test(t)
}

func TestIndexesPlaceholder_selectedIndexes(t *testing.T) {
	tests := []struct {
		inputPlaceholder IndexesPlaceholder
		inputIndexes     Indexes
		wantIndexes      Indexes
	}{
		// placeholder configured to import indexes
		{
			IndexesPlaceholder{Import: true},
			Indexes{Index{Name: "imported_index"}},
			Indexes{Index{Name: "imported_index"}},
		},
		// placeholder doesn't import indexes, also doesn't define any
		{
			IndexesPlaceholder{},
			Indexes{Index{Name: "imported_index"}},
			Indexes(nil),
		},
		// placeholder doesn't import indexes, uses its own
		{
			IndexesPlaceholder{Indexes: Indexes{Index{Name: "specified_index"}}},
			Indexes{Index{Name: "imported_index"}},
			Indexes{Index{Name: "specified_index"}},
		},
	}

	for _, test := range tests {
		gotIndexes := test.inputPlaceholder.selectedIndexes(test.inputIndexes)
		message := fmt.Sprintf("%#v.selectedIndexes(%#v)", test.inputPlaceholder, test.inputIndexes)
		testEqual(gotIndexes, test.wantIndexes, message, t)
	}
}

func TestIndexesPlaceholder_UnmarshalYAML(t *testing.T) {
	tests := yamlUnmarshallerTestCases{
		// empty definition isn't valid
		{
			&IndexesPlaceholder{},
			"",
			&IndexesPlaceholder{},
			true,
		},
		// explicit IndexesPlaceholder
		{
			&IndexesPlaceholder{},
			"{indexes: [{name: myindex}]}",
			&IndexesPlaceholder{Indexes: Indexes{Index{Name: "myindex"}}},
			false,
		},
		// list of indexes
		{
			&IndexesPlaceholder{},
			"[{name: myindex}]",
			&IndexesPlaceholder{Indexes: Indexes{Index{Name: "myindex"}}},
			false,
		},
		// boolean
		{
			&IndexesPlaceholder{},
			"true",
			&IndexesPlaceholder{Import: true},
			false,
		},
	}

	tests.test(t)
}
