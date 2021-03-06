---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "splunkconfig_app_package Resource - terraform-provider-splunkconfig"
subcategory: ""
description: |-
  Create a tarball for an app. Generated app.conf's version will be automatically incremented when app content changes.
---

# splunkconfig_app_package (Resource)

Create a tarball for an app. Generated app.conf's version will be automatically incremented when app content changes.

## Example Usage

```terraform
terraform {
  required_providers {
    splunkconfig = {
      source = "splunk/splunkconfig"
    }
  }
}

provider "splunkconfig" {
  configuration = <<EOF
indexes:
  - name: proxy
  - name: web

apps:
  - id: indexes
    name: indexes
    indexes: true
EOF
}

resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes"
  path   = "./"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **app_id** (String) ID of the app
- **path** (String) Path in which to create the app file

### Read-Only

- **base_version** (String) Version of the app, directly from the provider
- **effective_version** (String) Version of the app, accounting for patch count
- **files** (List of Object) File content of the app (see [below for nested schema](#nestedatt--files))
- **patch_count** (Number) Number of patches to the app since setting/changing its version
- **tarball_path** (String) Full path of the generated tarball

<a id="nestedatt--files"></a>
### Nested Schema for `files`

Read-Only:

- **content** (String)
- **path** (String)


