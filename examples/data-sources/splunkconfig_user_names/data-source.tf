terraform {
  required_providers {
    splunkconfig = {
      source = "splunk/splunkconfig"
    }
  }
}

provider "splunkconfig" {
  configuration = <<EOF
users:
  - name: larry
  - name: moe
  - name: curly
EOF
}

data "splunkconfig_user_names" "users" {}

output "users" {
  value = data.splunkconfig_user_names.users.user_names
}
