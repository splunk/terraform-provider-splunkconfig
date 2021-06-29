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

import "fmt"

// ExplicitInt defines an integer value that can be explicitly zero.
type ExplicitInt struct {
	Value    int
	Explicit bool
}

// ExplicitlySetInt returns an ExplicitInt, with Explicit=true, for a given value.
func ExplicitlySetInt(value int) ExplicitInt {
	return ExplicitInt{
		Value:    value,
		Explicit: true,
	}
}

// UnmarshalYAML implements custom unmarshalling for an ExplicitInt.  It enables an ExplicitInt to be unmarshalled from
// these types of content:
// * {Value: 0, Explicit: true}    # an actual ExplicitInt
// * 0                             # a bare integer
func (explicitInt *ExplicitInt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type realExplicitInt ExplicitInt

	unmarshalledExplicitInt := realExplicitInt{}
	if err := unmarshal(&unmarshalledExplicitInt); err == nil {
		*explicitInt = ExplicitInt(unmarshalledExplicitInt)
		return nil
	}

	var unmarshalledInt int
	if err := unmarshal(&unmarshalledInt); err == nil {
		*explicitInt = ExplicitlySetInt(unmarshalledInt)
		return nil
	}

	return fmt.Errorf("unable to unmarshal ExplicitInt")
}
