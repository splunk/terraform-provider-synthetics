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

const newPortCheckV2Config = `
resource "synthetics_create_port_check_v2" "port_v2_foo_check" {
	provider = synthetics.synthetics
  test {
    active = true 
    frequency = 5
    location_ids = ["aws-us-west-2", "aws-us-east-1"]
    scheduling_strategy = "round_robin"
		custom_properties {
			key = "key"
			value = "value"
		}
    name = "acceptance-Terraform-PORT-V2"
    port = 8081
    protocol = "udp"
    host = "www.splunk.com"
  }    
}
`

const updatedPortCheckV2Config = `
resource "synthetics_create_port_check_v2" "port_v2_foo_check" {
	provider = synthetics.synthetics
  test {
    active = false 
    frequency = 15
    location_ids = ["aws-us-east-1"]
    scheduling_strategy = "concurrent"
		custom_properties {
			key = "beepkey"
			value = "boopvalue"
		}
    name = "acceptance-updated-Terraform-PORT-V2"
    port = 8082
    protocol = "tcp"
    host = "www.duckduckgo.com"
  }    
}
`

func TestAccCreateUpdatePortCheckV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: providerConfig + newPortCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.active", "true"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.frequency", "5"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.location_ids.0", "aws-us-west-2"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.location_ids.1", "aws-us-east-1"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.scheduling_strategy", "round_robin"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.custom_properties.0.key", "key"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.custom_properties.0.value", "value"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.name", "acceptance-Terraform-PORT-V2"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.port", "8081"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.protocol", "udp"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.host", "www.splunk.com"),
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
				Config: providerConfig + updatedPortCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.active", "false"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.frequency", "15"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.location_ids.0", "aws-us-east-1"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.scheduling_strategy", "concurrent"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.custom_properties.0.key", "beepkey"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.custom_properties.0.value", "boopvalue"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.name", "acceptance-updated-Terraform-PORT-V2"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.port", "8082"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("synthetics_create_port_check_v2.port_v2_foo_check", "test.0.host", "www.duckduckgo.com"),
				),
			},
		},
	})
}
