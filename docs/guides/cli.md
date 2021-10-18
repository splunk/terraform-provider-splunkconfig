---
page_title: "Command Line Tools"
subcategory: ""
description: |-
  Included command line tools
---

# Placement

Tools are included in this provider under the `tools/` directory of the provider.

After performing `terraform init`, this directory should exist, relative to the initialized root module, as:

```
.terraform/providers/registry.terraform.io/splunk/splunkconfig/<version>/<architecture>/tools
```

# Usage

## template-lookup-csv

Print YAML for a Lookup from CSV content.

### Arguments

- **csvFilename** (required) Path to CSV content. The CSV file must include a header row with field names.
- **lookupName** (required) Name to give the generated Lookup.
