terraform {
  required_providers {
    splunkconfig = {
      source = "splunk/splunkconfig"
    }
  }
}

provider "splunkconfig" {
  configuration = <<EOF
indexes:
  - name: proxy
  - name: web

apps:
  - id: indexes
    name: indexes
    indexes: true
EOF
}

resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes"
  path   = "./"
}
