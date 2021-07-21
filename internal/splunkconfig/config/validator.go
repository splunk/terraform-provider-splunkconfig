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
	"reflect"
)

// validator objects can validate themselves.
type validator interface {
	validate() error
}

type uniqueValidator interface {
	validator
	uid() string
}

func allValidNoDuplicates(validators []uniqueValidator) error {
	seen := make(map[string]bool)

	for _, v := range validators {
		if err := v.validate(); err != nil {
			return err
		}

		uniqueID := v.uid()

		if seen[uniqueID] {
			return fmt.Errorf("duplicate %T: %s", v, uniqueID)
		}
		seen[uniqueID] = true
	}

	return nil
}

// uniqueValidators returns a list of uniqueValidator objects from a list of objects that adhere to the uniqueValidator interface.
func uniqueValidators(list interface{}) []uniqueValidator {
	listValue := reflect.ValueOf(list)

	listLength := listValue.Len()
	validators := make([]uniqueValidator, listLength)

	for i := range validators {
		iValue := listValue.Index(i)

		validators[i] = iValue.Interface().(uniqueValidator)
	}

	return validators
}
