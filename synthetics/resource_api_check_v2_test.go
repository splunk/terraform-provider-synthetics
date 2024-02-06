// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package synthetics

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const newApiCheckV2Config = `
resource "synthetics_create_api_check_v2" "api_v2_foo_check" {
	provider = synthetics.synthetics
	test {
		active = true
		device_id = 1  
		frequency = 5
		location_ids = ["aws-us-east-1"]
		name = "2 Terraform-Api V2 Acceptance Checkaroo"
		scheduling_strategy = "round_robin"
		custom_properties {
			key = "key"
			value = "value"
		}
		requests {
			configuration {
				name = "Get products"
				body = "\\'{\"alert_name\":\"the service is down\",\"url\":\"https://foo.com/bar\"}\\'\n"
				headers = {
					"Accept": "application/json"
					"x-foo": "bar-foo"
				}
				request_method = "GET"
				url = "https://dummyjson.com/products"
			}
			setup {
				name = "Extract from response body 01"
				type = "extract_json"
				source = "{{response.body}}"
				extractor = "extractosd"
				variable = "extractsetupvar"
			}
			setup {
				name = "Save Response Body 02"
				type = "save"
				value = "{{response.body}}"
				variable = "savesetupvar"
			}
			validations {
				name = "My assert validation step 01"
				actual = "{{response.code}}"
				comparator = "equals"
				expected = 200
				type = "assert_numeric"
			}
			validations {
				name = "Save response body 02"
				type = "save"
				value = "{{response.body}}"
				variable = "saverespvar"
			}
		}
		requests {
			configuration {
				name = "What about 2nd Get products?"
				body = "\\'{\"bad_alert\":\"the service is over\",\"url\":\"https://foo2.com/bar\"}\\'\n"
				headers = {
					"x-foo": "bar2-foo1"
				}
				request_method = "GET"
				url = "https://dummyjson.com/products1"
			}
			setup {
				name = "Extract from response body 01-02"
				type = "extract_json"
				source = "{{response.body}}"
				extractor = "extractosd"
				variable = "extractsetupvar"
			}
			setup {
				name = "Save Response Body 02-02"
				type = "save"
				value = "{{response.body}}"
				variable = "savesetupvar"
			}
			validations {
				name = "My Assert validation step 01-02"
				actual = "{{response.code}}"
				comparator = "equals"
				expected = 200
				type = "assert_numeric"
			}
			validations {
				name = "Save response body 02-02"
				type = "save"
				value = "{{response.body}}"
				variable = "saverespvar"
			}
		}
	}
}
`

const updatedApiCheckV2Config = `
resource "synthetics_create_api_check_v2" "api_v2_foo_check" {
	provider = synthetics.synthetics
	test {
		active = false
		device_id = 2  
		frequency = 15
		location_ids = ["aws-us-west-1"]
		name = "2 Terraform-Api V2 Acceptance Checkaroo Updated"
		scheduling_strategy = "concurrent"
		custom_properties {
			key = "beepkey"
			value = "boopvalue"
		}
		requests {
			configuration {
				name = "Get productz"
				body = "\\'{\"alert_name\":\"the service is down\",\"url\":\"https://foo.com/bar\"}\\'\n"
				headers = {
					"Accept": "application/xml"
					"x-foo": "bar-food"
				}
				request_method = "POST"
				url = "https://dummyjson.com/products2"
			}
			setup {
				name = "Save Response Body 02"
				type = "save"
				value = "{{response.body}}"
				variable = "savesetupvar"
			}
			setup {
				name = "Extract from response body updated 01"
				type = "extract_json"
				source = "{{response.body}}"
				extractor = "extractosd-updated"
				variable = "extractsetupvar-updated"
			}
			validations {
				name = "Save response body 02"
				type = "save"
				value = "{{response.body}}"
				variable = "saverespvar"
			}
			validations {
				name = "My assert validation step 01"
				actual = "{{response.code}}"
				comparator = "equals"
				expected = 200
				type = "assert_numeric"
			}
		}
		requests {
			configuration {
				body = "\\'{\"super_alert\":\"the service is over\",\"url\":\"https://foo2.com/bar\"}\\'\n"
				headers = {
					"x-foo": "bar2-foo1"
				}
				name = "What about 2nd Get elevensies?"
				request_method = "GET"
				url = "https://dummyjson.com/products4"
			}
			setup {
				name = "Save Response Body 02"
				type = "save"
				value = "{{response.body}}"
				variable = "savesetupvar"
			}
			setup {
				name = "Extract from response body 01"
				type = "extract_json"
				source = "{{response.body}}"
				extractor = "extractosd"
				variable = "extractsetupvar"
			}
			validations {
				name = "Save response body 02"
				type = "save"
				value = "{{response.body}}"
				variable = "saverespvar"
			}
			validations {
				name = "My Assert validation step 01"
				actual = "{{response.code}}"
				comparator = "equals"
				expected = 200
				type = "assert_numeric"
			}
		}
	}
}
`

func TestAccCreateUpdateApiCheckV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: providerConfig + newApiCheckV2Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.active", "true"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.device_id", "1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.frequency", "5"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.location_ids.0", "aws-us-east-1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.name", "2 Terraform-Api V2 Acceptance Checkaroo"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.scheduling_strategy", "round_robin"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.custom_properties.0.key", "key"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.custom_properties.0.value", "value"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.name", "Get products"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.body", "\\'{\"alert_name\":\"the service is down\",\"url\":\"https://foo.com/bar\"}\\'\n"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.headers.Accept", "application/json"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.headers.x-foo", "bar-foo"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.headers.%", "2"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.request_method", "GET"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.url", "https://dummyjson.com/products"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.name", "Extract from response body 01"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.type", "extract_json"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.source", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.extractor", "extractosd"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.variable", "extractsetupvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.name", "Save Response Body 02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.type", "save"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.value", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.variable", "savesetupvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.name", "My assert validation step 01"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.actual", "{{response.code}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.comparator", "equals"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.expected", "200"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.type", "assert_numeric"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.name", "Save response body 02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.type", "save"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.value", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.variable", "saverespvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.name", "What about 2nd Get products?"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.body", "\\'{\"bad_alert\":\"the service is over\",\"url\":\"https://foo2.com/bar\"}\\'\n"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.headers.x-foo", "bar2-foo1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.headers.%", "1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.request_method", "GET"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.url", "https://dummyjson.com/products1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.name", "Extract from response body 01-02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.type", "extract_json"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.source", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.extractor", "extractosd"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.variable", "extractsetupvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.name", "Save Response Body 02-02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.type", "save"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.value", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.variable", "savesetupvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.name", "My Assert validation step 01-02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.actual", "{{response.code}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.comparator", "equals"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.expected", "200"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.type", "assert_numeric"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.name", "Save response body 02-02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.type", "save"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.value", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.variable", "saverespvar"),
				),
			},
			// grab state
			{
				ResourceName:      "synthetics_create_api_check_v2.api_v2_foo_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_api_check_v2.api_v2_foo_check"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: providerConfig + updatedApiCheckV2Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.active", "false"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.device_id", "2"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.frequency", "15"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.location_ids.0", "aws-us-west-1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.name", "2 Terraform-Api V2 Acceptance Checkaroo Updated"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.scheduling_strategy", "concurrent"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.custom_properties.0.key", "beepkey"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.custom_properties.0.value", "boopvalue"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.name", "Get productz"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.body", "\\'{\"alert_name\":\"the service is down\",\"url\":\"https://foo.com/bar\"}\\'\n"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.headers.Accept", "application/xml"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.headers.x-foo", "bar-food"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.headers.%", "2"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.request_method", "POST"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.configuration.0.url", "https://dummyjson.com/products2"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.name", "Save Response Body 02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.type", "save"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.value", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.0.variable", "savesetupvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.name", "Extract from response body updated 01"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.type", "extract_json"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.source", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.extractor", "extractosd-updated"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.setup.1.variable", "extractsetupvar-updated"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.name", "Save response body 02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.type", "save"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.value", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.0.variable", "saverespvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.name", "My assert validation step 01"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.actual", "{{response.code}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.comparator", "equals"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.expected", "200"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.0.validations.1.type", "assert_numeric"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.name", "What about 2nd Get elevensies?"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.body", "\\'{\"super_alert\":\"the service is over\",\"url\":\"https://foo2.com/bar\"}\\'\n"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.headers.x-foo", "bar2-foo1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.headers.%", "1"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.request_method", "GET"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.configuration.0.url", "https://dummyjson.com/products4"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.name", "Save Response Body 02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.type", "save"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.value", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.0.variable", "savesetupvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.name", "Extract from response body 01"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.type", "extract_json"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.source", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.extractor", "extractosd"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.setup.1.variable", "extractsetupvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.name", "Save response body 02"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.type", "save"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.value", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.0.variable", "saverespvar"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.name", "My Assert validation step 01"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.actual", "{{response.code}}"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.comparator", "equals"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.expected", "200"),
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "test.0.requests.1.validations.1.type", "assert_numeric"),
				),
			},
		},
	})
}
