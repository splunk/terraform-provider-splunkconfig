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
