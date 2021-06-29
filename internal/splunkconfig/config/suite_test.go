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

func TestSuite_validate(t *testing.T) {
	tests := validatorTestCases{
		{
			// simple valid configuration
			Suite{
				Indexes: Indexes{
					Index{Name: "index_a"},
					Index{Name: "index_b"},
				},
				Roles: Roles{
					Role{Name: "role_a"},
					Role{Name: "role_b"},
				},
			},
			false,
		},
		{
			// invalid Indexes
			Suite{
				Indexes: Indexes{
					Index{Name: "index_a"},
					Index{Name: "index_a"},
				},
			},
			true,
		},
		{
			// invalid Roles
			Suite{
				Roles: Roles{
					Role{Name: "role_a"},
					Role{Name: "role_a"},
				},
			},
			true,
		},
		{
			// Index references RoleName that doesn't exist
			Suite{
				Indexes: Indexes{
					Index{
						Name:               "index_a",
						SearchRolesAllowed: RoleNames{"role_a", "role_b"},
					},
				},
				Roles: Roles{
					Role{Name: "role_a"},
				},
			},
			true,
		},
	}

	tests.test(t)
}

func TestSuite_NewSuiteFromYAML(t *testing.T) {
	yamlContent := `
# this is a comment to make sure comments work
anchors:
  retention_days_1: &retention_days_1
    frozenTimePeriod:
      seconds: 86400
indexes:
  - name: index_a
    <<: *retention_days_1
    srchRolesAllowed:
      - role_a
roles:
  - name: role_a
    srchIndexesAllowed:
      - existing_index
users:
  - name: test_account
    password: test_password
`

	want := Suite{
		Indexes: Indexes{
			Index{Name: "index_a", FrozenTime: TimePeriod{Seconds: 86400}, SearchRolesAllowed: RoleNames{"role_a"}},
		},
		Roles: Roles{
			Role{Name: "role_a", SearchIndexesAllowed: IndexNames{"existing_index"}},
		},
		Users: Users{
			User{Name: "test_account", Password: "test_password"},
		},
	}

	got, err := NewSuiteFromYAML([]byte(yamlContent))
	// because Anchors is an interface{}, it can really be anything, and so we don't want to bother testing it
	// we only care that the defined anchors were correctly used by other components
	// so we set Anchors to nil here, to pretend it didn't ever get set
	got.Anchors = nil
	if err != nil {
		t.Errorf("NewSuiteFromYAML error: %s", err)
	}

	message := fmt.Sprintf("NewSuiteFromYAML(%q)", yamlContent)
	testEqual(got, want, message, t)
}
