// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
)

// ACL represents the permissions configuration for a knowledge object.
type ACL struct {
	App     string
	Owner   string
	Sharing Sharing
	Read    RoleNames
	Write   RoleNames
}

// validate returns an error if the ACL is invalid. It is invalid if it:
// * has an invalid Sharing value
// * has Read or Write set, but not both
// * has invalid Read RoleNames
// * has invalid Write Rolenames
func (acl ACL) validate() error {
	if err := acl.Sharing.validate(); err != nil {
		return fmt.Errorf("ACL is invalid, has invalid Sharing: %s", err)
	}

	if acl.hasReadOrWrite() != acl.hasReadAndWrite() {
		return fmt.Errorf("ACL is invalid: Read and Write need to both be configured or unconfigured")
	}

	if err := acl.Read.validate(); err != nil {
		return fmt.Errorf("ACL is invalid, has invalid Read: %s", err)
	}

	if err := acl.Write.validate(); err != nil {
		return fmt.Errorf("ACL is invalid, has invalid Write: %s", err)
	}

	return nil
}

// hasReadOrWrite returns True if either Read or Write is non-nil.
func (acl ACL) hasReadOrWrite() bool {
	return acl.Read != nil || acl.Write != nil
}

// hasReadAndWrite returns True if both Read and Write are non-nil.
func (acl ACL) hasReadAndWrite() bool {
	return acl.Read != nil && acl.Write != nil
}

// stanzaValues returns the StanzaValues for the ACL object.
func (acl ACL) stanzaValues() StanzaValues {
	values := StanzaValues{}

	// hasReadAndWrite() would also work here, assuming the acl has been validated, because they are both true or false
	if acl.hasReadOrWrite() {
		values["access"] = fmt.Sprintf("read : %s, write : %s",
			acl.Read.metaAccessValue(),
			acl.Write.metaAccessValue())
	}

	if Sharing(acl.Sharing.metaValue()) != "" {
		values["export"] = string(acl.Sharing.metaValue())
	}

	return values
}
