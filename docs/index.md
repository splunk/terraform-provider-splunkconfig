---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "splunkconfig Provider"
subcategory: ""
description: |-
  Define and query a Splunk Configuration to be used by other Terraform providers or other automation
---

# splunkconfig Provider

## Description

Define and query a Splunk Configuration to be used by other Terraform providers or other automation

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
  configuration_file = "./splunkconfig.yml"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **configuration** (String) YAML content containing the abstracted configuration. Exactly one of `configuration`, `configuration_file`, or `configuration_path` must be set.
- **configuration_file** (String) Full path to YAML file containing the abstracted configuration. Exactly one of `configuration`, `configuration_file`, or `configuration_path` must be set.
- **configuration_path** (String) Full path to directory containing one or more YAML files containing the abstracted configuration. Exactly one of `configuration`, `configuration_file`, or `configuration_path` must be set.
