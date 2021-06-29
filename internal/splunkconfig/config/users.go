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
	"reflect"
	"sort"
)

// Users is a list of User objects.
type Users []User

// validate returns an error if any of users' members are invalid.
func (users Users) validate() error {
	return allValidNoDuplicates(uniqueValidators(users))
}

// Names returns a list of user names for each User in users, sorted by Name.
func (users Users) Names() []string {
	uids := uidsOfUIDers(users)
	sort.Strings(uids)

	return uids
}

// WithName returns the User object with the given Name. Returns ok=false if not found.
func (users Users) WithName(name string) (found User, ok bool) {
	foundUIDer, ok := withUID(users, name)
	if !ok {
		return
	}

	foundValue := reflect.ValueOf(foundUIDer)
	found = foundValue.Interface().(User)

	return
}
