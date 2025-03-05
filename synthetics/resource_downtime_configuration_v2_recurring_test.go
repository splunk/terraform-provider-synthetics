// Copyright 2024 Splunk, Inc.
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

// start_time must be in the future and no more than one year in the future
// end_time must be after the start_time and no more than one year after the start_time
// test_ids must be existing test_ids in the org
const newRecurringDowntimeConfigurationV2Config = `
resource "synthetics_create_downtime_configuration_v2" "downtime_configuration_v2_foo_recurring" {
	provider = synthetics.synthetics
  downtime_configuration {
    name = "acceptance-downtime-configuration-recurring-terraform-test"
    description = "The most awesome recurring downtime_configuration. Full of snakes."
    rule = "augment_data"
    start_time = "2025-03-25T17:13:00.000Z"
    end_time = "2025-03-25T18:13:00.000Z"
    test_ids = [1523512]
    timezone = "America/New_York"
    recurrence {
      repeats {
        type = "daily"
      }
      end {
        type = "on"
        value = "2025-04-25"
      }
    }
  }
}
`

func TestAccCreateRecurringDowntimeConfigurationV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: providerConfig + newRecurringDowntimeConfigurationV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.description", "The most awesome recurring downtime_configuration. Full of snakes."),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.rule", "augment_data"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.name", "acceptance-downtime-configuration-recurring-terraform-test"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.start_time", "2025-03-25T17:13:00.000Z"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.end_time", "2025-03-25T18:13:00.000Z"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.timezone", "America/New_York"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.recurrence.0.repeats.0.type", "daily"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.recurrence.0.end.0.type", "on"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.recurrence.0.end.0.value", "2025-04-25"),
				),
			},
			{
				ResourceName:      "synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring"),
				ImportStateVerify: true,
			},
		},
	})
}
