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

const newHttpCheckV2Config = `
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
`

const updatedHttpCheckV2Config = `
resource "synthetics_create_http_check_v2" "http_v2_foo_check" {
  test {
    active = true 
    frequency = 15
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
`

func TestAccCreateUpdateHttpCheckV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccHttpCheckV2Destroy,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: newHttpCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateHttpCheckV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "frequency", "5"),
				),
			},
			{
				ResourceName:      "synthetics_create_http_check_v2.http_v2_foo_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_http_check_v2.http_v2_foo_check"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: updatedHttpCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateHttpCheckV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "frequency", "15"),
				),
			},
		},
	})
}

func testAccCreateUpdateHttpCheckV2ResourceExists(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_http_check_v2":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetHttpCheckV2(checkId)
			if strconv.Itoa(check.Test.ID) != rs.Primary.ID || err != nil {
				return fmt.Errorf("Error finding http check v2 %s: %s", rs.Primary.ID, err)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}
	return nil
}

func testAccHttpCheckV2Destroy(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_http_check_v2":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetHttpCheckV2(checkId)
			if check.Test.ID != checkId || err != nil {
				return fmt.Errorf("Found deleted check %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

