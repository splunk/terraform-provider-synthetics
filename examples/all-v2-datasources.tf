# terraform {
#   required_providers {
#     synthetics = {
#       version = "1.0.4"
#       source  = "splunk/synthetics"
#     }
#   }
# }

# provider "synthetics" {
#   product = "observability"
#   realm = "us1"
#   #apikey = "this-is-my-api-key"
# }

# //Pull a Check as a datasource

# data "synthetics_variable_v2_check" "datasource_check_variable" {
#   variable {
#     id = 246
#   }
# }

# output "datasource_check_variable" {
#   value = data.synthetics_variable_v2_check.datasource_check_variable
# }

# data "synthetics_port_v2_check" "datasource_check_port" {
#   test {
#     id = 1650
#   }
# }

# output "datasource_check_port" {
#   value = data.synthetics_port_v2_check.datasource_check_port
# }

# data "synthetics_http_v2_check" "datasource_check_http" {
#   test {
#     id = 490
#   }
# }

# output "datasource_check_http" {
#   value = data.synthetics_http_v2_check.datasource_check_http
# }


# data "synthetics_api_v2_check" "datasource_check_api" {
#   test {
#     id = 489
#   }
# }

# output "datasource_check_api" {
#   value = data.synthetics_api_v2_check.datasource_check_api
# }

data "synthetics_browser_v2_check" "datasource_check_browser" {
  test {
    id = 5056
  }
}

output "datasource_check_browser" {
  value = data.synthetics_browser_v2_check.datasource_check_browser
}

# data "synthetics_location_v2_check" "datasource_location" {
#   location {
#     id = "aws-af-south-1"
#   }
# }

# output "datasource_location" {
#   value = data.synthetics_location_v2_check.datasource_location
# }

# data "synthetics_locations_v2_check" "datasource_locations" {
#   locations {
#   }
# }

# output "datasource_locations" {
#   value = data.synthetics_locations_v2_check.datasource_locations
# }

# data "synthetics_variables_v2_check" "datasource_variables" {
#   variables {
#   }
# }

# output "datasource_variables" {
#   value = data.synthetics_variables_v2_check.datasource_variables
# }

# data "synthetics_devices_v2_check" "datasource_devices" {
#   devices {
#   }
# }

# output "datasource_devices" {
#   value = data.synthetics_devices_v2_check.datasource_devices
# }
