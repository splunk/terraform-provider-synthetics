terraform {
  required_providers {
    synthetics = {
      version = "1.0.1"
      source  = "splunk.com/splunk/synthetics"
    }
  }
}

provider "synthetics" {
  product = "observability"
  realm = "us1"
  #apikey = "this-is-my-api-key"
}

//Pull a Check as a datasource

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

# data "synthetics_browser_v2_check" "datasource_check_browser" {
#   test {
#     id = 496
#   }
# }

# output "datasource_check_browser" {
#   value = data.synthetics_browser_v2_check.datasource_check_browser
# }

# data "synthetics_location_v2_check" "datasource_location" {
#   location {
#     id = "aws-us-east-1"
#   }
# }

# output "datasource_location" {
#   value = data.synthetics_location_v2_check.datasource_location
# }

data "synthetics_locations_v2_check" "datasource_locations" {
  locations {
  }
}

output "datasource_locations" {
  value = data.synthetics_locations_v2_check.datasource_locations
}


//=================================
//=================================
//=================================

//Create a V2 Variable
# resource "synthetics_create_variable_v2" "variable_v2_foo" {
#   variable {
#     description = "The most awesome variable. Full of snakes."
#     value = "barv3--oopsasdasd"
#     // Once created name and secret can not be changed and will result in a 422 from the API
#     // unless the variable is deleted and re-created
#     name = "terraform-test121"
#     secret = false  
#   }    
# }

  
# output "variable_v2_foo" {
#   value = synthetics_create_variable_v2.variable_v2_foo
# }

# //Create a Http V2 Check
# resource "synthetics_create_http_check_v2" "http_v2_foo_check" {
#   test {
#     active = true 
#     frequency = 5
#     location_ids = ["aws-us-east-1","aws-ap-northeast-3"]
#     name = "Terraform - HTTP V2 Checkaroo"
#     type = "http"
#     url = "https://www.splunk.com"
#     scheduling_strategy = "round_robin"
#     request_method = "GET"
#     body = null
#     headers {
#       name = "Synthetic_transaction_1"
#       value = "batman is the man"
#     }
#     headers {
#       name = "back_transaction_1"
#       value = "peeko"
#     }
#   }    
# }

  
# output "http_v2_foo_check" {
#   value = synthetics_create_http_check_v2.http_v2_foo_check
# }

# //Create a Port V2 Check
# resource "synthetics_create_port_check_v2" "port_v2_foo_check" {
#   test {
#     name = "Terraform - PORT V2 Checkaroo"
#     # type = "port"
#     port = 8080
#     protocol = "udp"
#     host = "www.splunk.com"
#     location_ids = ["aws-us-west-2"]
#     frequency = 5
#     scheduling_strategy = "concurrent"
#     active = true 
#   }    
# }

  
# output "port_v2_foo_check" {
#   value = synthetics_create_port_check_v2.port_v2_foo_check
# }

# //Create a Browser V2 Check
# resource "synthetics_create_browser_check_v2" "browser_v2_foo_check" {
#   test {
#     active = true
#     device_id = 1  
#     frequency = 5
#     location_ids = ["aws-us-east-1"]
#     name = "Terraform - Browser V2 Checkaroo"
#     scheduling_strategy = "round_robin"
#     url_protocol = "https://"
#     start_url = "www.splunk.com"
#     transactions {
#       name = "Synthetic transaction 1"
#       steps {
#         name = "Go to URL"
#         type = "go_to_url"
#         url = "https://www.splunk.com"
#         action = "go_to_url"
#         wait_for_nav = true
#         options {
#           url = "https://www.splunk.com"
#         }
#       }
#       steps {
#         name = "New step"
#         type = "click_element"
#         selector_type = "id"
#         wait_for_nav = false
#         selector = "\"free-splunk-click-mobile\""
#       }
#       steps {
#         name = "New step"
#         type = "click_element"
#         selector_type = "id"
#         wait_for_nav = false
#         selector = "login-button"
#       }
#     }
#     transactions {
#       name = "New synthetic transaction"
#       steps {
#         name = "New step"
#         type = "go_to_url"
#         wait_for_nav = true
#         action = "go_to_url"
#         url = "https://www.batman.com"
#       }
#     }
#     advanced_settings {
#       verify_certificates = false
#       user_agent = "Mozilla/5.0 (X11; Linux x86_64; Splunk Synthetics) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
#       authentication {
#         username = "batmab"
#         password = "{{env.beep-var}}"
#       }
#       headers {
#         name = "superstar-machine"
#         value = "\"taking it too the staaaaars\""
#         domain = "asdasd.batman.com"
#       }
#       cookies {
#         key = "sda"
#         value = "sda"
#         domain = "asd.com"
#         path = "/asd"
#       }
#       cookies {
#         key = "yes"
#         value = "no"
#         domain = "zodiak.com"
#         path = "/Edlesley"
#       }
#       host_overrides {
#         source = "asdasd.com"
#         target = "whost.com"
#         keep_host_header = false
#       }
#       host_overrides {
#         source = "92.2.2.2"
#         target = "91.1.1.1"
#         keep_host_header = true
#       }
#     }
#   }    
# }

  
# output "browser_v2_foo_check" {
#   value = synthetics_create_browser_check_v2.browser_v2_foo_check
# }

# //Create an API V2 Check
# resource "synthetics_create_api_check_v2" "api_v2_foo_check" {
#   test {
#     active = true
#     device_id = 1  
#     frequency = 5
#     location_ids = ["aws-us-east-1"]
#     name = "Terraform - Api V2 Checkaroo"
#     scheduling_strategy = "round_robin"
#     requests {
#         configuration {
#           body = "\\'{\"alert_name\":\"the service is down\",\"url\":\"https://foo.com/bar\"}\\'\n"
#           headers = {
#             "Accept": "application/json"
#             "x-foo": "bar"
#           }
#           name = "Get products"
#           request_method = "GET"
#           url = "https://dummyjson.com/products"
#         }
#         setup {
#             extractor = "$.foo"
#             name = "First setup step"
#             source = "{\\'foo\\': \\'bar\\'}"
#             type = "extract_json"
#             variable = "myVariable"
#           }
#         validations {
#             actual = "{{response.code}}"
#             comparator = "equals"
#             expected = 200
#             name = "My validation step"
#             type = "assert_numeric"
#           }
#       }
#   }
# }

  
# output "api_v2_foo_check" {
#   value = synthetics_create_api_check_v2.api_v2_foo_check
# }


