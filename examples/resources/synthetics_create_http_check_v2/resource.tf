resource "synthetics_create_http_check_v2" "http_v2_foo_check" {
  test {
    active = true 
    frequency = 10
    location_ids = ["aws-us-east-1","aws-ap-northeast-3"]
    name = "Terraform1 - HTTP V2 Checkaroo"
    type = "http"
    url = "https://www.splunk.com"
    scheduling_strategy = "round_robin"
    custom_properties {
			key = "key"
			value = "value"
		}
    request_method = "GET"
    verify_certificates = true
    user_agent = "Another User of Agents"
    body = null
    headers {
      name = "Synthetic_transaction_1"
      value = "batmab is the mab"
    }
    headers {
      name = "back_transaction_1"
      value = "peeko"
    }
  }    
}