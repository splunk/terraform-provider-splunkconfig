---
page_title: "SAML Groups Example"
subcategory: ""
description: |-
  Working with SAML groups
---

# Working with SAML groups

## Example terraform

```
data "splunkconfig_role_names" "found_roles" {}

data "splunkconfig_role_attributes" "found_role" {
  for_each = toset(data.splunkconfig_role_names.found_roles.role_names)

  role_name = each.key
}

resource "splunk_authorization_roles" "deployed_role" {
  for_each = data.splunkconfig_role_attributes.found_role

  name = each.value.role_name
}

data "splunkconfig_saml_group_names" "found_saml_groups" {}

data "splunkconfig_saml_group_attributes" "found_saml_group" {
  for_each = toset(data.splunkconfig_saml_group_names.found_saml_groups.saml_group_names)

  saml_group_name = each.key
}

resource "splunk_admin_saml_groups" "deployed_saml_group" {
  for_each = data.splunkconfig_saml_group_attributes.found_saml_group

  name  = each.value.saml_group_name
  roles = each.value.roles

  depends_on = [splunk_authorization_roles.deployed_role]
}
```

## Initial splunkconfig

```
saml_groups:
  - name: IT-Web-Admins
    roles: [user]

roles:
  - name: web_users
```

## Add SAML group to role

### Altered splunkconfig

```
saml_groups:
  - name: IT-Web-Admins
    roles: [user]

roles:
  - name: web_users
    saml_groups: [IT-Web-Admins]   # reverse-assign role to the IT-Web-Admins SAML group
```

### terraform apply output

```
splunk_authorization_roles.deployed_role["web_users"]: Refreshing state... [id=web_users]
splunk_admin_saml_groups.deployed_saml_group["IT-Web-Admins"]: Refreshing state... [id=IT-Web-Admins]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # splunk_admin_saml_groups.deployed_saml_group["IT-Web-Admins"] will be updated in-place
  ~ resource "splunk_admin_saml_groups" "deployed_saml_group" {
        id    = "IT-Web-Admins"
        name  = "IT-Web-Admins"
      ~ roles = [
            "user",
          + "web_users",
        ]
    }

Plan: 0 to add, 1 to change, 0 to destroy.
splunk_admin_saml_groups.deployed_saml_group["IT-Web-Admins"]: Modifying... [id=IT-Web-Admins]
splunk_admin_saml_groups.deployed_saml_group["IT-Web-Admins"]: Modifications complete after 0s [id=IT-Web-Admins]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```
