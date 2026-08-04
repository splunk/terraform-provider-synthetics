# Create-style validation: checks the payload as if creating a new test.
data "synthetics_validate_http_check_v2" "new_http_check" {
  test {
    active               = true
    frequency            = 10
    location_ids         = ["aws-us-east-1", "aws-ap-northeast-3"]
    name                 = "Terraform1 - HTTP V2 Checkaroo"
    type                 = "http"
    url                  = "https://www.splunk.com"
    port                 = 443
    scheduling_strategy  = "round_robin"
    request_method       = "GET"
    verify_certificates  = true
    user_agent           = "Another User of Agents"
  }
}

# Update-style validation: checks the payload as if updating an existing test.
data "synthetics_validate_http_check_v2" "existing_http_check" {
  test_id = 1650
  test {
    active               = true
    frequency            = 10
    location_ids         = ["aws-us-east-1", "aws-ap-northeast-3"]
    name                 = "Terraform1 - HTTP V2 Checkaroo"
    type                 = "http"
    url                  = "https://www.splunk.com"
    port                 = 443
    scheduling_strategy  = "round_robin"
    request_method       = "GET"
    verify_certificates  = true
    user_agent           = "Another User of Agents"
  }
}

output "new_http_check_valid" {
  value = data.synthetics_validate_http_check_v2.new_http_check.valid
}

output "new_http_check_field_errors" {
  value = data.synthetics_validate_http_check_v2.new_http_check.field_errors
}
