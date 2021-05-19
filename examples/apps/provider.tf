terraform {
  required_providers {
    splunkconfig = {
      source = "cd.splunkdev.com/sas/splunkconfig"
    }
  }
}

provider "splunkconfig" {
  configuration_file = "splunkconfig.yml"
}
