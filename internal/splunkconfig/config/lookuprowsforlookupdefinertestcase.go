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

// lookupRowsForLookupDefinerTestCase defines a single test case for a lookupRowsForLookupDefiner object.
type lookupRowsForLookupDefinerTestCase struct {
	definer  lookupRowsForLookupDefiner
	lookup   Lookup
	wantRows LookupRows
}

// test performs the test defined in a lookupRowsForLookupDefinerTestCase.
func (testCase lookupRowsForLookupDefinerTestCase) test(t *testing.T) {
	gotRows := testCase.definer.lookupRowsForLookup(testCase.lookup)
	message := fmt.Sprintf("%T{%+v}.lookupRowsForLookup(%T{%+v}) = %v, want %v", testCase.definer, testCase.definer, testCase.lookup, testCase.lookup, gotRows, testCase.wantRows)
	testEqual(gotRows, testCase.wantRows, message, t)
}
