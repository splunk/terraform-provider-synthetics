resource "synthetics_create_http_check_v2" "http_v2_foo_check" {
  test {
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
      value = "peeko"
    }
  }    
}