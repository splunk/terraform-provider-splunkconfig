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
