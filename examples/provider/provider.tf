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
