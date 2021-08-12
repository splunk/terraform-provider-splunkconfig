---
page_title: "Roles Example"
subcategory: ""
description: |-
  Working with roles
---

# Working with roles

## Example terraform

```
data "splunkconfig_role_names" "found_roles" {}

data "splunkconfig_role_attributes" "found_role" {
  for_each = toset(data.splunkconfig_role_names.found_roles.role_names)

  role_name = each.key
}

resource "splunk_authorization_roles" "deployed_role" {
  for_each = data.splunkconfig_role_attributes.found_role

  name           = each.value.role_name
  imported_roles = each.value.imported_roles
}
```

## Initial splunkconfig
```
roles:
  - name: web_admins
  - name: itsec_admins
```

## Edit a role

### Altered splunkconfig
```
roles:
  - name: web_admins
  - name: itsec_admins
    importRoles:
      - web_admins      # added an imported role
```

### terraform apply output
```
splunk_authorization_roles.deployed_role["web_admins"]: Refreshing state... [id=web_admins]
splunk_authorization_roles.deployed_role["itsec_admins"]: Refreshing state... [id=itsec_admins]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # splunk_authorization_roles.deployed_role["itsec_admins"] will be updated in-place
  ~ resource "splunk_authorization_roles" "deployed_role" {
        id                                    = "itsec_admins"
      ~ imported_roles                        = [
          + "web_admins",
        ]
        name                                  = "itsec_admins"
        # (9 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
splunk_authorization_roles.deployed_role["itsec_admins"]: Modifying... [id=itsec_admins]
splunk_authorization_roles.deployed_role["itsec_admins"]: Modifications complete after 1s [id=itsec_admins]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```
