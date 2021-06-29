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

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// helper test function to return a list of tests for list values
func testCheckResourceAttrList(name string, key string, values []string) resource.TestCheckFunc {
	// testCheckFuncs is 1 longer than values.  we add a "next index is missing" test
	testCheckFuncs := make([]resource.TestCheckFunc, len(values)+1)

	for i, value := range values {
		keyWithIndex := fmt.Sprintf("%s.%d", key, i)
		testCheckFuncs[i] = resource.TestCheckResourceAttr(name, keyWithIndex, value)
	}

	// add a check that the next index is unset, to ensure list length equality
	keyWithIndex := fmt.Sprintf("%s.%d", key, len(values))
	testCheckFuncs[len(values)] = resource.TestCheckNoResourceAttr(name, keyWithIndex)

	return resource.ComposeTestCheckFunc(testCheckFuncs...)
}
