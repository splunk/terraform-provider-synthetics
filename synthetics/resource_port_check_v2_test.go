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

const newPortCheckV2Config = `
resource "synthetics_create_port_check_v2" "port_v2_foo_check" {
  test {
    name = "Terraform - PORT V2 Checkaroo"
    port = 8080
    protocol = "udp"
    host = "www.splunk.com"
    location_ids = ["aws-us-east-1","aws-ap-northeast-3"]
    frequency = 5
    scheduling_strategy = "concurrent"
    active = true 
  }    
}
`

const updatedPortCheckV2Config = `
resource "synthetics_create_port_check_v2" "port_v2_foo_check" {
  test {
    name = "Terraform - PORT V2 Checkaroo"
    port = 8080
    protocol = "udp"
    host = "www.splunk.com"
    location_ids = ["aws-us-east-1","aws-ap-northeast-3"]
    frequency = 15
    scheduling_strategy = "concurrent"
    active = true 
  }    
}
`

func TestAccCreateUpdatePortCheckV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPortCheckV2Destroy,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: newPortCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdatePortCheckV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "frequency", "5"),
				),
			},
			{
				ResourceName:      "synthetics_create_port_check_v2.port_v2_foo_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_port_check_v2.port_v2_foo_check"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: updatedPortCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdatePortCheckV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "frequency", "15"),
				),
			},
		},
	})
}

func testAccCreateUpdatePortCheckV2ResourceExists(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_port_check_v2":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetPortCheckV2(checkId)
			if strconv.Itoa(check.Test.ID) != rs.Primary.ID || err != nil {
				return fmt.Errorf("Error finding port check v2 %s: %s", rs.Primary.ID, err)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}
	return nil
}

func testAccPortCheckV2Destroy(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_port_check_v2":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetPortCheckV2(checkId)
			if check.Test.ID != checkId || err != nil {
				return fmt.Errorf("Found deleted check %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

