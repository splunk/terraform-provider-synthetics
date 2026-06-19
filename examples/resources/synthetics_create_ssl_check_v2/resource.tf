resource "synthetics_create_ssl_check_v2" "ssl_v2_check" {
  test {
    name                 = "Terraform - SSL V2 Check"
    active               = false
    frequency            = 5
    location_ids         = ["aws-us-east-1"]
    automatic_retries    = 1
    scheduling_strategy  = "round_robin"
    host                 = "www.splunk.com"
    port                 = 443
    server_name          = "www.splunk.com"
    allow_self_signed    = false
    allow_untrusted_root = false

    custom_properties {
      key   = "env"
      value = "example"
    }
  }
}
