resource "splunkconfig_app_package" "indexes" {
  app_id = "indexes"
  path = "."
}

resource "splunkconfig_app_package" "lookups" {
  app_id = "lookups"
  path = "."
}
