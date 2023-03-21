resource "synthetics_create_api_check_v2" "api_v2_foo_check" {
  test {
    active = true
    device_id = 1  
    frequency = 5
    location_ids = ["aws-us-east-1"]
    name = "Terraform - Api V2 Checkaroo"
    scheduling_strategy = "round_robin"
    requests {
        configuration {
          body = "\\'{\"alert_name\":\"the service is down\",\"url\":\"https://foo.com/bar\"}\\'\n"
          headers = {
            "Accept": "application/json"
            "x-foo": "bar"
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
