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

const newVariableV2Config = `
resource "synthetics_create_variable_v2" "variable_v2_foo" {
  variable {
    description = "The most awesome variable. Full of snakes."
    value = "barv3v3"
    name = "terraform-test"
    secret = false  
  }    
}
`

const updatedVariableV2Config = `
resource "synthetics_create_variable_v2" "variable_v2_foo" {
  variable {
    description = "The most awesome variable. Full of snakes."
    value = "barv3v322"
    name = "terraform-test"
    secret = false  
  }    
}
`

func TestAccCreateUpdateVariableV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVariableV2Destroy,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: newVariableV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateVariableV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "value", "barv3v3"),
				),
			},
			{
				ResourceName:      "synthetics_create_variable_v2.variable_v2_foo",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_variable_v2.variable_v2_foo"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: updatedVariableV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateVariableV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "value", "barv3v322"),
				),
			},
		},
	})
}

func testAccCreateUpdateVariableV2ResourceExists(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "signalfx_alert_muting_rule":
			variableId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting variable id: %s", err)
			}
			variable, _, err := client.GetVariableV2(variableId)
			if strconv.Itoa(variable.Variable.ID) != rs.Primary.ID || err != nil {
				return fmt.Errorf("Error finding variable v2 %s: %s", rs.Primary.ID, err)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}
	return nil
}

func testAccVariableV2Destroy(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_variable_v2":
			variableId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting variable id: %s", err)
			}
			variable, _, err := client.GetVariableV2(variableId)
			if variable.Variable.ID != variableId || err != nil {
				return fmt.Errorf("Found deleted variable %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

