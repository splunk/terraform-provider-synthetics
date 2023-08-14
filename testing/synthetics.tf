terraform {
  required_providers {
    synthetics = {
      version = "1.0.4"
      source  = "splunk/synthetics"
    }
  }
}

provider "synthetics" {
  product = "observability"
  realm = "us1"
  apikey = "mQ_EjIKlOxtdY_UTkJafdw"
}

resource "synthetics_create_http_check_v2" "http_test" {
    active = true 
    frequency = 5
    location_ids = ["aws-us-east-1","aws-ap-northeast-3"]
    name = "Terraform - HTTP V2 Checkaroo"
    type = "http"
    url = "https://www.splunk.com"
    scheduling_strategy = "round_robin"
    request_method = "GET"
    body = null
    headers {
      name = "Synthetic_transaction_1"
      value = "batman is the man"
    }
    headers {
      name = "back_transaction_1"
      value = "test"
    }
}    

# resource "synthetics_create_api_check_v2" "test" {
#   test {
#     active = true
#     device_id = 1  
#     frequency = 5
#     location_ids = ["aws-us-east-1"]
#     name = "Terraform - test"
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
#         validations {
#             actual = "{{response.dns_time}}"
#             comparator = "is_less_than"
#             expected = 2000
#             name = "My validation step"
#             type = "assert_numeric"
#           }
#       }
#   }
# }

resource "synthetics_create_api_check_v2" "new" {
  test {
     active = true
    device_id = 1  
    frequency = 5
    location_ids = ["aws-us-east-1"]
    name = "Terraform - new schema"
    scheduling_strategy = "round_robin"
      requests {
        configuration {
          body = "\\'{\"alert_name\":\"the service is down\",\"url\":\"https://foo.com/bar\"}\\'\n"
          headers = {
            "Accept": "application/json"
            "xFoo": "bar"
          }
          name = "Get products"
          request_method = "GET"
          url = "https://dummyjson.com/products"
        }
        setup {
          extractor = "$.foo"
          name = "First setup step"
          source = "{\\'foo\\': \\'bar\\'}"
          type = "extract_json"
          variable = "myVariable"
        }
        validations {
          actual = "{{response.code}}"
          comparator = "equals"
          expected = 200
          name = "My validation step"
          type = "assert_numeric"
        }
      }
  } 
}