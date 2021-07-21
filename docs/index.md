---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "splunkconf Provider"
subcategory: ""
description: |-
  Define and query a Splunk Configuration to be used by other Terraform providers or other automation
---

# splunkconf Provider

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

### Configuration

- **configuration** (String) YAML content containing the abstracted configruation. Either this or `configuration_file` must be set.
- **configuration_file** (String) Full path to YAML file containing the abstracted configuration. Either this or `configuration` must be set.