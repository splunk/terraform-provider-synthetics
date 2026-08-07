// Copyright 2026 Splunk, Inc.
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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourcePortCheckV2Config(name string) string {
	return fmt.Sprintf(`
resource "synthetics_create_port_check_v2" "port_v2_datasource_fixture" {
  provider = synthetics.synthetics
  test {
    active              = true
    frequency           = 5
    location_ids        = ["aws-us-east-1"]
    name                = %[1]q
    host                = "www.splunk.com"
    port                = 8081
    protocol            = "tcp"
    scheduling_strategy = "round_robin"
    automatic_retries   = 1
  }
}

data "synthetics_port_v2_check" "port_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_port_check_v2.port_v2_datasource_fixture]

  test {
    id = tonumber(synthetics_create_port_check_v2.port_v2_datasource_fixture.id)
  }
}
`, name)
}

func TestAccDataSourcePortCheckV2(t *testing.T) {
	name := testAccUniqueName("acceptance-terraform-port-v2-datasource")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourcePortCheckV2Config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.synthetics_port_v2_check.port_v2_datasource", "id",
						"synthetics_create_port_check_v2.port_v2_datasource_fixture", "id",
					),
					resource.TestCheckResourceAttr("data.synthetics_port_v2_check.port_v2_datasource", "test.#", "1"),
					resource.TestCheckResourceAttr("data.synthetics_port_v2_check.port_v2_datasource", "test.0.name", name),
					resource.TestCheckResourceAttr("data.synthetics_port_v2_check.port_v2_datasource", "test.0.active", "true"),
					resource.TestCheckResourceAttr("data.synthetics_port_v2_check.port_v2_datasource", "test.0.frequency", "5"),
					resource.TestCheckResourceAttr("data.synthetics_port_v2_check.port_v2_datasource", "test.0.host", "www.splunk.com"),
					resource.TestCheckResourceAttr("data.synthetics_port_v2_check.port_v2_datasource", "test.0.port", "8081"),
					resource.TestCheckResourceAttr("data.synthetics_port_v2_check.port_v2_datasource", "test.0.protocol", "tcp"),
					resource.TestCheckResourceAttrSet("data.synthetics_port_v2_check.port_v2_datasource", "test.0.type"),
				),
			},
		},
	})
}
