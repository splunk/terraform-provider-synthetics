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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

// A downtime configuration requires start_time in the future and test_ids that exist in
// the org, so a disposable port check is created alongside it and referenced by generated
// id — the same shape as testAccDowntimeConfigurationV2Config.
func testAccDataSourceDowntimeConfigurationV2Config(name, startTime, endTime string) string {
	return fmt.Sprintf(`
resource "synthetics_create_port_check_v2" "downtime_v2_datasource_test_fixture" {
  provider = synthetics.synthetics
  test {
    active              = true
    frequency           = 5
    location_ids        = ["aws-us-west-2"]
    scheduling_strategy = "round_robin"
    name                = "%[1]s-test-fixture"
    port                = 8081
    protocol            = "tcp"
    host                = "www.splunk.com"
  }
}

resource "synthetics_create_downtime_configuration_v2" "downtime_configuration_v2_datasource_fixture" {
  provider = synthetics.synthetics
  downtime_configuration {
    name        = "%[1]s"
    description = "Terraform acceptance downtime configuration data source fixture"
    rule        = "augment_data"
    start_time  = "%[2]s"
    end_time    = "%[3]s"
    test_ids    = [tonumber(synthetics_create_port_check_v2.downtime_v2_datasource_test_fixture.id)]
  }
}

data "synthetics_downtime_configuration_v2_check" "downtime_configuration_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_downtime_configuration_v2.downtime_configuration_v2_datasource_fixture]

  downtime_configuration {
    id = tonumber(synthetics_create_downtime_configuration_v2.downtime_configuration_v2_datasource_fixture.id)
  }
}

data "synthetics_downtime_configurations_v2_check" "downtime_configurations_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_downtime_configuration_v2.downtime_configuration_v2_datasource_fixture]
}
`, name, startTime, endTime)
}

func TestAccDataSourceDowntimeConfigurationV2(t *testing.T) {
	name := testAccUniqueName("acceptance-terraform-downtime-configuration-v2-datasource")
	startTime := time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05.000Z")
	endTime := time.Now().Add(48 * time.Hour).Format("2006-01-02T15:04:05.000Z")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// This test creates two fixtures; both must be gone after destroy.
		CheckDestroy: testAccCheckFixturesDestroyed(map[string]func(string) (*sc2.RequestDetails, error){
			"synthetics_create_downtime_configuration_v2": testAccIntLookup(func(id int) (*sc2.RequestDetails, error) {
				_, details, err := testAccProvider.Meta().(*sc2.Client).GetDowntimeConfigurationV2(id)
				return details, err
			}),
			"synthetics_create_port_check_v2": testAccIntLookup(func(id int) (*sc2.RequestDetails, error) {
				_, details, err := testAccProvider.Meta().(*sc2.Client).GetPortCheckV2(id)
				return details, err
			}),
		}),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourceDowntimeConfigurationV2Config(name, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					// Singleton read.
					resource.TestCheckResourceAttrPair(
						"data.synthetics_downtime_configuration_v2_check.downtime_configuration_v2_datasource", "id",
						"synthetics_create_downtime_configuration_v2.downtime_configuration_v2_datasource_fixture", "id",
					),
					resource.TestCheckResourceAttr("data.synthetics_downtime_configuration_v2_check.downtime_configuration_v2_datasource", "downtime_configuration.#", "1"),
					resource.TestCheckResourceAttr("data.synthetics_downtime_configuration_v2_check.downtime_configuration_v2_datasource", "downtime_configuration.0.name", name),
					resource.TestCheckResourceAttr("data.synthetics_downtime_configuration_v2_check.downtime_configuration_v2_datasource", "downtime_configuration.0.description", "Terraform acceptance downtime configuration data source fixture"),
					resource.TestCheckResourceAttr("data.synthetics_downtime_configuration_v2_check.downtime_configuration_v2_datasource", "downtime_configuration.0.rule", "augment_data"),
					resource.TestCheckResourceAttr("data.synthetics_downtime_configuration_v2_check.downtime_configuration_v2_datasource", "downtime_configuration.0.start_time", startTime),
					resource.TestCheckResourceAttr("data.synthetics_downtime_configuration_v2_check.downtime_configuration_v2_datasource", "downtime_configuration.0.end_time", endTime),

					// List read must contain the fixture.
					resource.TestCheckResourceAttrSet("data.synthetics_downtime_configurations_v2_check.downtime_configurations_v2_datasource", "id"),
					testAccCheckDataSourceListContains(
						"data.synthetics_downtime_configurations_v2_check.downtime_configurations_v2_datasource",
						"downtime_configurations", "name", name,
					),
				),
			},
		},
	})
}
