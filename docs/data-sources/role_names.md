---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "splunkconfig_role_names Data Source - terraform-provider-splunkconfig"
subcategory: ""
description: |-
  Return Role Names from the Splunk Configuration
---

# splunkconfig_role_names (Data Source)

Return Role Names from the Splunk Configuration

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
roles:
  - name: admin
  - name: user
EOF
}

data "splunkconfig_role_names" "roles" {}

output "roles" {
  value = data.splunkconfig_role_names.roles.role_names
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- **role_names** (List of String) List of Role Names in the Splunk Configuration


