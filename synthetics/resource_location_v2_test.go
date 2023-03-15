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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sc "syntheticsclientv2"
)

const newLocationV2Config = `
resource "synthetics_create_location_v2" "location_v2_foo" {
  location {
    id = "aws-awesome-west-2"
    label = "awesome aws location"
    default = false
    type = "private"
    country = "US"
  }    
}
`

const updatedLocationV2Config = `
resource "synthetics_create_location_v2" "location_v2_foo" {
  location {
    id = "aws-awesome-west"
    label = "awesome aws location. Now with snakes!"
    default = false
    type = "private"
    country = "UK"
  }    
}
`

func TestAccCreateUpdateLocationV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccLocationV2Destroy,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: newLocationV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateLocationV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_location_v2.location_v2_foo", "value", "barv3v3"),
				),
			},
			{
				ResourceName:      "synthetics_create_location_v2.location_v2_foo",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_location_v2.location_v2_foo"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: updatedLocationV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateLocationV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_location_v2.location_v2_foo", "value", "barv3v322"),
				),
			},
		},
	})
}

func testAccCreateUpdateLocationV2ResourceExists(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "signalfx_alert_muting_rule":
			var locationId = rs.Primary.ID
			location, _, err := client.GetLocationV2(locationId)
			if location.Location.ID != rs.Primary.ID || err != nil {
				return fmt.Errorf("Error finding location v2 %s: %s", rs.Primary.ID, err)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}
	return nil
}

func testAccLocationV2Destroy(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_location_v2":
			var locationId = rs.Primary.ID
			location, _, err := client.GetLocationV2(locationId)
			if location.Location.ID != locationId || err != nil {
				return fmt.Errorf("Found deleted location %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

