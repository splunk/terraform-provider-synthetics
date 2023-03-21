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
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sc "github.com/splunk/syntheticsclient/syntheticsclientv2"
)

const newApiCheckV2Config = `
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
`

const updatedApiCheckV2Config = `
resource "synthetics_create_api_check_v2" "api_v2_foo_check" {
  test {
    active = true
    device_id = 1  
    frequency = 15
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
`

func TestAccCreateUpdateApiCheckV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccApiCheckV2Destroy,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: newApiCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateApiCheckV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "frequency", "5"),
				),
			},
			{
				ResourceName:      "synthetics_create_api_check_v2.api_v2_foo_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_api_check_v2.api_v2_foo_check"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: updatedApiCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateApiCheckV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_api_check_v2.api_v2_foo_check", "frequency", "15"),
				),
			},
		},
	})
}

func testAccCreateUpdateApiCheckV2ResourceExists(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_api_check_v2":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetApiCheckV2(checkId)
			if strconv.Itoa(check.Test.ID) != rs.Primary.ID || err != nil {
				return fmt.Errorf("Error finding Api check v2 %s: %s", rs.Primary.ID, err)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}
	return nil
}

func testAccApiCheckV2Destroy(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_api_check_v2":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetApiCheckV2(checkId)
			if check.Test.ID != checkId || err != nil {
				return fmt.Errorf("Found deleted check %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

