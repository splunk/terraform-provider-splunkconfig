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

// RolesPlaceholder represents a set of Roles or an intent to import Roles from elsewhere.
type RolesPlaceholder struct {
	Roles  Roles
	Import bool
}

// validate returns an error if RolesPlaceholder is invalid.  It is invalid if its Roles are invalid.
func (rolesPlaceholder RolesPlaceholder) validate() error {
	if err := rolesPlaceholder.Roles.validate(); err != nil {
		return fmt.Errorf("RolesPlaceholder invalid, invalid Roles: %s", err)
	}

	return nil
}

// selectedRoles returns the candidateRoles if RolesPlaceholder.Import is true.  Otherwise it returns
// RolesPlaceholder.Roles.
func (rolesPlaceholder RolesPlaceholder) selectedRoles(candidateRoles Roles) Roles {
	if rolesPlaceholder.Import {
		return candidateRoles
	}

	return rolesPlaceholder.Roles
}

// UnmarshalYAML implements custom unmarshalling for a RolesPlaceholder.  It enables a RolesPlaceholder to be
// unmarshalled from these types of content:
// * {roles: [{name: my_role}]        # explicitly define its roles as part of the RolesPlaceholder structure
// * [{name: my_role}]                # provide a list of roles directly
// * true                             # configure RolesPlaceholder to import roles instead
func (rolesPlaceholder *RolesPlaceholder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// realRolesPlaceholder only exists inside this function, and is used to allow attempting to unmarshal into what
	// is really just a RolesPlaceholder directly.  Attempting to unmarshal(&RolesPlaceholder) from inside this
	// function will result in infinite recursion back into this function, so we need another type to attempt that
	// unmarshalling.
	type realRolesPlaceholder RolesPlaceholder

	// first try to unmarshal into (effectively) an actual RolesPlaceholder
	unmarshalledRolesPlaceholder := realRolesPlaceholder{}
	if err := unmarshal(&unmarshalledRolesPlaceholder); err == nil {
		*rolesPlaceholder = RolesPlaceholder(unmarshalledRolesPlaceholder)
		return nil
	}

	// then try to unmarshal into a Roles object, to be embedded in the placeholder
	unmarshalledRoles := Roles{}
	if err := unmarshal(&unmarshalledRoles); err == nil {
		rolesPlaceholderFromRoles := RolesPlaceholder{
			Roles: unmarshalledRoles,
		}
		*rolesPlaceholder = rolesPlaceholderFromRoles

		return nil
	}

	// and finally try to unmarshal into a boolean, to be embedded in the placeholder
	unmarshalledRolesPlaceholderBool := false
	if err := unmarshal(&unmarshalledRolesPlaceholderBool); err == nil {
		rolesPlaceHolderFromBool := RolesPlaceholder{
			Import: unmarshalledRolesPlaceholderBool,
		}
		*rolesPlaceholder = rolesPlaceHolderFromBool

		return nil
	}

	// if none of the above unmarshal attempts succeed, return an error
	return fmt.Errorf("unable to unmarshall RolesPlaceholder from YAML")
}
