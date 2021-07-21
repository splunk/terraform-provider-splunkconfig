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
    importRoles: [user]
  - name: user
EOF
}

data "splunkconfig_role_attributes" "admin" {
  role_name = "admin"
}

output "admin_imported_roles" {
  value = data.splunkconfig_role_attributes.admin.imported_roles
}
