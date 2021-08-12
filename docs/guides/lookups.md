---
page_title: "Lookups Example"
subcategory: ""
description: |-
  Working with lookups
---

# Working with lookups

## Example terraform

```
resource "splunkconfig_app_package" "http_lookups" {
  app_id = "http_lookups"
  path   = path.root
}

resource "splunkconfig_app_package" "role_lookups" {
  app_id = "role_lookups"
  path   = path.root
}
```

## Initial splunkconfig

```
lookups:
  - name: http_status_codes
    fields:
      - name: code
      - name: description
    rows:
      - values: {code: 200, description: OK}
      - values: {code: 404, description: Not Found}
  - name: role_contacts
    fields:
      - name: role
      - name: contact

apps:
  - name: HTTP Lookups
    id: http_lookups
    lookups: [http_status_codes]
    version: 1.0.0
  - name: Role Lookups
    id: role_lookups
    lookups: [role_contacts]
    version: 1.0.0

roles:
  - name: web_admins
```

## Add a row

### Altered splunkconfig

```
lookups:
  - name: http_status_codes
    fields:
      - name: code
      - name: description
    rows:
      - values: {code: 200, description: OK}
      - values: {code: 401, description: Unauthorized}     # new row added here
      - values: {code: 404, description: Not Found}
  - name: role_contacts
    fields:
      - name: role
      - name: contact

apps:
  - name: HTTP Lookups
    id: http_lookups
    lookups: [http_status_codes]
    version: 1.0.0
  - name: Role Lookups
    id: role_lookups
    lookups: [role_contacts]
    version: 1.0.0

roles:
  - name: web_admins
```

### terraform apply output

```
splunkconfig_app_package.role_lookups: Refreshing state... [id=role_lookups]
splunkconfig_app_package.http_lookups: Refreshing state... [id=http_lookups]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # splunkconfig_app_package.http_lookups will be updated in-place
  ~ resource "splunkconfig_app_package" "http_lookups" {
      ~ effective_version = "1.0.0" -> "1.0.1"
      ~ files             = [
          ~ {
              ~ content = <<-EOT
                    [ui]
                    is_visible = false
                    label = HTTP Lookups
                    
                    [launcher]
                    author = 
                    description = 
                  - version = 1.0.0
                  + version = 1.0.1
                    
                    [package]
                    check_for_updates = false
                    id = http_lookups
                    
                EOT
                # (1 unchanged element hidden)
            },
          ~ {
              ~ content = <<-EOT
                    code,description
                    200,OK
                  + 401,Unauthorized
                    404,Not Found
                EOT
                # (1 unchanged element hidden)
            },
        ]
        id                = "http_lookups"
      ~ patch_count       = 0 -> 1
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
splunkconfig_app_package.http_lookups: Modifying... [id=http_lookups]
splunkconfig_app_package.http_lookups: Modifications complete after 0s [id=http_lookups]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## Major version change

### Altered splunkconfig

```
lookups:
  - name: http_status_codes
    fields:
      - name: code
      - name: description
    rows:
      - values: {code: 200, description: 200 OK}            # potentially breaking change to rows
      - values: {code: 401, description: 401 Unauthorized}  # justifies a major version bump to the
      - values: {code: 404, description: 404 Not Found}     # app containing this lookup
  - name: role_contacts
    fields:
      - name: role
      - name: contact

apps:
  - name: HTTP Lookups
    id: http_lookups
    lookups: [http_status_codes]
    version: 2.0.0                                          # explicity increase version
  - name: Role Lookups
    id: role_lookups
    lookups: [role_contacts]
    version: 1.0.0

roles:
  - name: web_admins
```

### terraform apply output

```
splunkconfig_app_package.role_lookups: Refreshing state... [id=role_lookups]
splunkconfig_app_package.http_lookups: Refreshing state... [id=http_lookups]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # splunkconfig_app_package.http_lookups will be updated in-place
  ~ resource "splunkconfig_app_package" "http_lookups" {
      ~ base_version      = "1.0.0" -> "2.0.0"
      ~ effective_version = "1.0.1" -> "2.0.0"
      ~ files             = [
          ~ {
              ~ content = <<-EOT
                    [ui]
                    is_visible = false
                    label = HTTP Lookups
                    
                    [launcher]
                    author = 
                    description = 
                  - version = 1.0.1
                  + version = 2.0.0
                    
                    [package]
                    check_for_updates = false
                    id = http_lookups
                    
                EOT
                # (1 unchanged element hidden)
            },
          ~ {
              ~ content = <<-EOT
                    code,description
                  - 200,OK
                  - 401,Unauthorized
                  - 404,Not Found
                  + 200,200 OK
                  + 401,401 Unauthorized
                  + 404,404 Not Found
                EOT
                # (1 unchanged element hidden)
            },
        ]
        id                = "http_lookups"
      ~ patch_count       = 1 -> 0
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
splunkconfig_app_package.http_lookups: Modifying... [id=http_lookups]
splunkconfig_app_package.http_lookups: Modifications complete after 0s [id=http_lookups]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## Role adds a row to the lookup

### Altered splunkconfig

```
lookups:
  - name: http_status_codes
    fields:
      - name: code
      - name: description
    rows:
      - values: {code: 200, description: 200 OK}
      - values: {code: 401, description: 401 Unauthorized}
      - values: {code: 404, description: 404 Not Found}
  - name: role_contacts
    fields:
      - name: role
      - name: contact

apps:
  - name: HTTP Lookups
    id: http_lookups
    lookups: [http_status_codes]
    version: 2.0.0
  - name: Role Lookups
    id: role_lookups
    lookups: [role_contacts]
    version: 1.0.0

roles:
  - name: web_admins
    lookup_rows:
      - lookup_name: role_contacts                  # this role adds a row to the lookup
        values: {contact: web_admins@example.com}   # "role" field is automatically set
```

### terraform apply output

```
splunkconfig_app_package.role_lookups: Refreshing state... [id=role_lookups]
splunkconfig_app_package.http_lookups: Refreshing state... [id=http_lookups]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # splunkconfig_app_package.role_lookups will be updated in-place
  ~ resource "splunkconfig_app_package" "role_lookups" {
      ~ effective_version = "1.0.0" -> "1.0.1"
      ~ files             = [
          ~ {
              ~ content = <<-EOT
                    [ui]
                    is_visible = false
                    label = Role Lookups
                    
                    [launcher]
                    author = 
                    description = 
                  - version = 1.0.0
                  + version = 1.0.1
                    
                    [package]
                    check_for_updates = false
                    id = role_lookups
                    
                EOT
                # (1 unchanged element hidden)
            },
          ~ {
              ~ content = <<-EOT
                    role,contact
                  + web_admins,web_admins@example.com
                EOT
                # (1 unchanged element hidden)
            },
        ]
        id                = "role_lookups"
      ~ patch_count       = 0 -> 1
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
splunkconfig_app_package.role_lookups: Modifying... [id=role_lookups]
splunkconfig_app_package.role_lookups: Modifications complete after 0s [id=role_lookups]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```
