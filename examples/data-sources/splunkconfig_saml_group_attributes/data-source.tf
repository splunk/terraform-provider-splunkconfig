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
    roles: [admin]
  - name: splunk_users
    roles: [user]
EOF
}

data "splunkconfig_saml_group_attributes" "admin" {
  saml_group_name = "splunk_admins"
}

output "admin_roles" {
  value = data.splunkconfig_saml_group_attributes.admin.roles
}
