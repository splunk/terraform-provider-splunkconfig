terraform {
  required_providers {
    splunkconfig = {
      source = "splunk/splunkconfig"
    }
  }
}

provider "splunkconfig" {
  configuration = <<EOF
saml_groups:
  - name: splunk_admins
  - name: splunk_users
EOF
}

data "splunkconfig_saml_group_names" "groups" {}

output "saml_groups" {
  value = data.splunkconfig_saml_group_names.groups.saml_group_names
}
