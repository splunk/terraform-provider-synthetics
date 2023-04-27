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
#   # apikey = "this-is-my-api-key"
# }


# //Create a V2 Location
# resource "synthetics_create_location_v2" "location_v2_foo" {
#   location {
#     id = "private-aws-awesome-east"
#     label = "awesome aws east location"
#     country = "IE"
#   }    
# }

# output "location_v2_foo" {
#   value = synthetics_create_location_v2.location_v2_foo
# }


# //Create a V2 Variable
# resource "synthetics_create_variable_v2" "variable_v2_foo" {
#   variable {
#     description = "The most awesome variable. Full of snakes and spiders."
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
#     frequency = 15
#     location_ids = ["aws-us-east-1"]
#     name = "Terraform - Browser V2 Checkaroo"
#     scheduling_strategy = "round_robin"
#     url_protocol = "https://"
#     start_url = "www.splunk.com"
#     transactions {
#       name = "First Synthetic transaction"
#       steps {
#         name                 = "01 Go to URL"
#         type                 = "go_to_url"
#         url                  = "https://www.splunk.com"
#         wait_for_nav         = true
#       }
#       steps {
#         name                 = "02 fill in fieldz"
#         selector             = "beep"
#         selector_type        = "id"
#         type                 = "enter_value"
#         value                = "{{env.beep-var}}"
#         wait_for_nav         = false
#       }
#       steps {
#         name                 = "03 click"
#         selector             = "clicky"
#         selector_type        = "id"
#         type                 = "click_element"
#         wait_for_nav         = true
#       }
#       steps {
#         name                 = "04 accept---Alert"
#         type                 = "accept_alert"
#         wait_for_nav         = false
#       }
#       steps {
#         name                 = "05 Select-val-text"
#         option_selector      = "sdad"
#         option_selector_type = "text"
#         selector             = "textzz"
#         selector_type        = "id"
#         type                 = "select_option"
#         wait_for_nav         = false
#       }
#       steps {
#         name                 = "06 Select-Val-Val"
#         option_selector      = "{{env.beep-var}}"
#         option_selector_type = "value"
#         selector             = "valz"
#         selector_type        = "id"
#         type                 = "select_option"
#         wait_for_nav         = false
#       }
#       steps {
#         name                 = "07 Select-Val-Index"
#         option_selector      = "{{env.beep-var}}"
#         option_selector_type = "index"
#         selector             = "selectionz"
#         selector_type        = "id"
#         type                 = "select_option"
#         wait_for_nav         = false
#       }
#       steps {
#         name                 = "08 Save as text"
#         selector             = "beepval"
#         selector_type        = "link"
#         type                 = "store_variable_from_element"
#         variable_name        = "{{env.terraform-test-foo-301}}"
#         wait_for_nav         = false
#       }
#       steps {
#         name                 = "09 Save JS2 return Val"
#         type                 = "store_variable_from_javascript"
#         value                = "sdasds"
#         variable_name        = "{{env.terraform-test-foo-301}}"
#         wait_for_nav         = true
#       }
#       steps {
#         name                 = "010 Run JS"
#         type                 = "run_javascript"
#         value                = "beeeeeeep"
#         wait_for_nav         = true
#       }
#     }
#     transactions {
#       name = "2nd Synthetic transaction"
#       steps {
#         name                 = "Go to other URL"
#         type                 = "go_to_url"
#         url                  = "https://www.splunk.com"
#         wait_for_nav         = true
#       }
#       steps {
#         name                 = "fill in more fields field"
#         selector             = "beep"
#         selector_type        = "id"
#         type                 = "enter_value"
#         value                = "{{env.beep-var}}"
#         wait_for_nav         = false
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


