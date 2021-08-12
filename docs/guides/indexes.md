---
page_title: "Indexes Example"
subcategory: ""
description: |-
  Working with indexes
---

# Working with indexes

## Example terraform

```
resource "splunkconfig_app_package" "index_lookups" {
  app_id = "index_lookups"
  path   = path.root
}

resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes"
  path   = path.root
}
```

## Initial splunkconfig

```
lookups:
  - name: index_contacts
    fields:
      - name: index
      - name: contact

apps:
  - name: Index Lookups
    id: index_lookups
    lookups: [index_contacts]
    version: 1.0.0
  - name: Indexes
    id: indexes
    indexes: true
    version: 1.0.0

indexes:
  - name: web
    frozenTimePeriod: {days: 365}
```

## Add an index

### Altered splunkconfig

```
lookups:
  - name: index_contacts
    fields:
      - name: index
      - name: contact

apps:
  - name: Index Lookups
    id: index_lookups
    lookups: [index_contacts]
    version: 1.0.0
  - name: Indexes
    id: indexes
    indexes: true
    version: 1.0.0

indexes:
  - name: web
    frozenTimePeriod: {days: 365}
  - name: proxy                     # new index added
    frozenTimePeriod: {days: 365}
```

### terraform apply output
```
splunkconfig_app_package.indexes: Refreshing state... [id=indexes]
splunkconfig_app_package.index_lookups: Refreshing state... [id=index_lookups]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # splunkconfig_app_package.indexes will be updated in-place
  ~ resource "splunkconfig_app_package" "indexes" {
      ~ effective_version = "1.0.0" -> "1.0.1"
      ~ files             = [
          ~ {
              ~ content = <<-EOT
                    [ui]
                    is_visible = false
                    label = Indexes
                    
                    [launcher]
                    author = 
                    description = 
                  - version = 1.0.0
                  + version = 1.0.1
                    
                    [package]
                    check_for_updates = false
                    id = indexes
                    
                EOT
                # (1 unchanged element hidden)
            },
          ~ {
              ~ content = <<-EOT
                  + [proxy]
                  + coldPath = $SPLUNK_DB/proxy/colddb
                  + frozenTimePeriodInSecs = 31536000
                  + homePath = $SPLUNK_DB/proxy/db
                  + thawedPath = $SPLUNK_DB/proxy/thaweddb
                  + 
                    [web]
                    coldPath = $SPLUNK_DB/web/colddb
                    frozenTimePeriodInSecs = 31536000
                    homePath = $SPLUNK_DB/web/db
                    thawedPath = $SPLUNK_DB/web/thaweddb
                    
                EOT
                # (1 unchanged element hidden)
            },
        ]
        id                = "indexes"
      ~ patch_count       = 0 -> 1
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
splunkconfig_app_package.indexes: Modifying... [id=indexes]
splunkconfig_app_package.indexes: Modifications complete after 0s [id=indexes]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## Index adds a lookup row

### Altered splunkconfig

```
lookups:
  - name: index_contacts
    fields:
      - name: index
      - name: contact

apps:
  - name: Index Lookups
    id: index_lookups
    lookups: [index_contacts]
    version: 1.0.0
  - name: Indexes
    id: indexes
    indexes: true
    version: 1.0.0

indexes:
  - name: web
    frozenTimePeriod: {days: 365}
    lookup_rows:
      - lookup_name: index_contacts                # this index adds a row to the lookup
        values: {contact: web_admins@example.com}  # "index" field is automatically set
  - name: proxy
    frozenTimePeriod: {days: 365}
```

### terraform apply output

```
splunkconfig_app_package.index_lookups: Refreshing state... [id=index_lookups]
splunkconfig_app_package.indexes: Refreshing state... [id=indexes]

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  ~ update in-place

Terraform will perform the following actions:

  # splunkconfig_app_package.index_lookups will be updated in-place
  ~ resource "splunkconfig_app_package" "index_lookups" {
      ~ effective_version = "1.0.0" -> "1.0.1"
      ~ files             = [
          ~ {
              ~ content = <<-EOT
                    [ui]
                    is_visible = false
                    label = Index Lookups
                    
                    [launcher]
                    author = 
                    description = 
                  - version = 1.0.0
                  + version = 1.0.1
                    
                    [package]
                    check_for_updates = false
                    id = index_lookups
                    
                EOT
                # (1 unchanged element hidden)
            },
          ~ {
              ~ content = <<-EOT
                    index,contact
                  + web,web_admins@example.com
                EOT
                # (1 unchanged element hidden)
            },
        ]
        id                = "index_lookups"
      ~ patch_count       = 0 -> 1
        # (4 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
splunkconfig_app_package.index_lookups: Modifying... [id=index_lookups]
splunkconfig_app_package.index_lookups: Modifications complete after 0s [id=index_lookups]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```
