resource "synthetics_create_api_check_v2" "full_api_v2_foo_check" {
  test {
    active = true
    device_id = 1  
    frequency = 5
    location_ids = ["aws-us-east-1"]
    name = "2 Terraform-Api V2 Checkaroo"
    scheduling_strategy = "round_robin"
    requests {
      configuration {
        body = "\\'{\"alert_name\":\"the service is down\",\"url\":\"https://foo.com/bar\"}\\'\n"
        headers = {
          "Accept": "application/json"
          "x-foo": "bar-foo"
        }
        name = "Get products"
        request_method = "GET"
        url = "https://dummyjson.com/products"
      }
      setup {
        name = "Extract from response body"
        type = "extract_json"
        source = "{{response.body}}"
        extractor = "extractosd"
        variable = "extractsetupvar"
      }
      setup {
        name = "Save Response Body"
        type = "save"
        value = "{{response.body}}"
        variable = "savesetupvar"
      }
      setup {
        name = "JS Run"
        type = "javascript"
        code = "js code"
        variable = "jsvarsetup"
      }
      validations {
        actual = "{{response.code}}"
        comparator = "equals"
        expected = 200
        name = "My validation step"
        type = "assert_numeric"
      }
      validations {
        name = "Extract from response body"
        type = "extract_json"
        source = "{{response.body}}"
        extractor = "js.extractor"
        variable = "extractjvar"
      }
      validations {
        name = "JavaScript run"
        type = "javascript"
        code = "codetorun"
        variable = "jscodevar"
      }
      validations {
        name = "Save response body"
        type = "save"
        value = "{{response.body}}"
        variable = "saverespvar"
      }
    }
    requests {
      configuration {
        body = "\\'{\"bad_alert\":\"the service is over\",\"url\":\"https://foo2.com/bar\"}\\'\n"
        headers = {
          "Accept": "application/json"
          "x-foo": "bar2-foo1"
        }
        name = "2nd Get products"
        request_method = "GET"
        url = "https://dummyjson.com/products1"
      }
      setup {
        name = "Extract from response body"
        type = "extract_json"
        source = "{{response.body}}"
        extractor = "extractosd"
        variable = "extractsetupvar"
      }
      setup {
        name = "Save Response Body"
        type = "save"
        value = "{{response.body}}"
        variable = "savesetupvar"
      }
      setup {
        name = "JS Run"
        type = "javascript"
        code = "js code"
        variable = "jsvarsetup"
      }
      validations {
        actual = "{{response.code}}"
        comparator = "equals"
        expected = 200
        name = "My validation step"
        type = "assert_numeric"
      }
      validations {
        name = "Extract from response body"
        type = "extract_json"
        source = "{{response.body}}"
        extractor = "js.extractor"
        variable = "extractjvar"
      }
      validations {
        name = "JavaScript run"
        type = "javascript"
        code = "codetorun"
        variable = "jscodevar"
      }
      validations {
        name = "Save response body"
        type = "save"
        value = "{{response.body}}"
        variable = "saverespvar"
      }
    }
  }
}