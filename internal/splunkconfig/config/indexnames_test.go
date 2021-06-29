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
	"testing"
)

// IndexNames.validate should return an error when expected.
func TestIndexNames_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			IndexNames{"fine", "alsofine"},
			false,
		},
		{
			IndexNames{"fine", "!fine", "also_fine"},
			true,
		},
		{
			IndexNames{"duplicate", "duplicate"},
			true,
		},
	}

	tests.test(t)
}

func TestIndexNames_deduplicatedSorted(t *testing.T) {
	indexNames := IndexNames{
		"a",
		"b",
		"c",
		"a0",
		"b0",
		"_internal",
		"a",
	}

	expectedIndexNames := IndexNames{
		"_internal",
		"a",
		"a0",
		"b",
		"b0",
		"c",
	}

	got := indexNames.deduplicatedSorted()
	message := fmt.Sprintf("%T{+%v}.deduplicatedSorted()", indexNames, indexNames)
	testEqual(got, expectedIndexNames, message, t)
}
