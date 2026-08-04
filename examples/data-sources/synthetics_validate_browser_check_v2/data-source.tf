# Create-style validation: checks the payload as if creating a new test.
data "synthetics_validate_browser_check_v2" "new_browser_check" {
  test {
    active               = true
    device_id            = 1
    frequency            = 15
    location_ids         = ["aws-us-east-1"]
    name                 = "Terraform-Browser V2 Checkaroo"
    scheduling_strategy  = "round_robin"
    advanced_settings {
      verify_certificates = true
    }
    transactions {
      name = "First Synthetic transaction"
      steps {
        name = "01 Go to URL"
        type = "go_to_url"
        url  = "https://www.splunk.com"
      }
    }
  }
}

# Update-style validation: checks the payload as if updating an existing test.
data "synthetics_validate_browser_check_v2" "existing_browser_check" {
  test_id = 1650
  test {
    active               = true
    device_id            = 1
    frequency            = 15
    location_ids         = ["aws-us-east-1"]
    name                 = "Terraform-Browser V2 Checkaroo"
    scheduling_strategy  = "round_robin"
    advanced_settings {
      verify_certificates = true
    }
    transactions {
      name = "First Synthetic transaction"
      steps {
        name = "01 Go to URL"
        type = "go_to_url"
        url  = "https://www.splunk.com"
      }
    }
  }
}

output "new_browser_check_valid" {
  value = data.synthetics_validate_browser_check_v2.new_browser_check.valid
}

output "new_browser_check_field_errors" {
  value = data.synthetics_validate_browser_check_v2.new_browser_check.field_errors
}
