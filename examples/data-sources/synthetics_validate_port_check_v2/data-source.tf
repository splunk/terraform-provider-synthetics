# Create-style validation: checks the payload as if creating a new test.
data "synthetics_validate_port_check_v2" "new_port_check" {
  test {
    name                 = "Terraform - PORT V2 Checkaroo"
    port                 = 8080
    protocol             = "udp"
    host                 = "www.splunk.com"
    location_ids         = ["aws-us-west-2"]
    frequency            = 5
    scheduling_strategy  = "concurrent"
    active               = true
  }
}

# Update-style validation: checks the payload as if updating an existing test.
data "synthetics_validate_port_check_v2" "existing_port_check" {
  test_id = 1650
  test {
    name                 = "Terraform - PORT V2 Checkaroo"
    port                 = 8080
    protocol             = "udp"
    host                 = "www.splunk.com"
    location_ids         = ["aws-us-west-2"]
    frequency            = 5
    scheduling_strategy  = "concurrent"
    active               = true
  }
}

output "new_port_check_valid" {
  value = data.synthetics_validate_port_check_v2.new_port_check.valid
}

output "new_port_check_field_errors" {
  value = data.synthetics_validate_port_check_v2.new_port_check.field_errors
}
