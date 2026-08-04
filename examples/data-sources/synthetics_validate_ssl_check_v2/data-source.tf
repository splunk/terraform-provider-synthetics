# Create-style validation: checks the payload as if creating a new test.
data "synthetics_validate_ssl_check_v2" "new_ssl_check" {
  test {
    name                 = "Terraform - SSL V2 Check"
    active               = false
    frequency            = 5
    location_ids         = ["aws-us-east-1"]
    scheduling_strategy  = "round_robin"
    host                 = "www.splunk.com"
    port                 = 443
    server_name          = "www.splunk.com"
    allow_self_signed    = false
    allow_untrusted_root = false
  }
}

# Update-style validation: checks the payload as if updating an existing test.
data "synthetics_validate_ssl_check_v2" "existing_ssl_check" {
  test_id = 1650
  test {
    name                 = "Terraform - SSL V2 Check"
    active               = false
    frequency            = 5
    location_ids         = ["aws-us-east-1"]
    scheduling_strategy  = "round_robin"
    host                 = "www.splunk.com"
    port                 = 443
    server_name          = "www.splunk.com"
    allow_self_signed    = false
    allow_untrusted_root = false
  }
}

output "new_ssl_check_valid" {
  value = data.synthetics_validate_ssl_check_v2.new_ssl_check.valid
}

output "new_ssl_check_field_errors" {
  value = data.synthetics_validate_ssl_check_v2.new_ssl_check.field_errors
}
