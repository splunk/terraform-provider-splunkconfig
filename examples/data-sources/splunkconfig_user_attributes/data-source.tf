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
    email: larry@example.com
  - name: moe
    email: moe@example.com
  - name: curly
    email: curly@example.com
EOF
}

data "splunkconfig_user_attributes" "larry" {
  user_name = "larry"
}

output "larry_email" {
  value = data.splunkconfig_user_attributes.larry.email
}
