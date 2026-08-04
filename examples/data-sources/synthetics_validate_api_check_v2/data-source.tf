# Create-style validation: checks the payload as if creating a new test.
data "synthetics_validate_api_check_v2" "new_api_check" {
  test {
    active               = true
    device_id            = 1
    frequency            = 5
    location_ids         = ["aws-us-east-1"]
    name                 = "Terraform-Api V2 Checkaroo"
    scheduling_strategy  = "round_robin"
    requests {
      configuration {
        name           = "Get products"
        request_method = "GET"
        url            = "https://dummyjson.com/products"
      }
    }
  }
}

# Update-style validation: checks the payload as if updating an existing test.
data "synthetics_validate_api_check_v2" "existing_api_check" {
  test_id = 1650
  test {
    active               = true
    device_id            = 1
    frequency            = 5
    location_ids         = ["aws-us-east-1"]
    name                 = "Terraform-Api V2 Checkaroo"
    scheduling_strategy  = "round_robin"
    requests {
      configuration {
        name           = "Get products"
        request_method = "GET"
        url            = "https://dummyjson.com/products"
      }
    }
  }
}

output "new_api_check_valid" {
  value = data.synthetics_validate_api_check_v2.new_api_check.valid
}

output "new_api_check_field_errors" {
  value = data.synthetics_validate_api_check_v2.new_api_check.field_errors
}
