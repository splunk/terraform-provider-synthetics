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
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// start_time must be in the future and no more than one year in the future
// end_time must be after the start_time and no more than one year after the start_time
// test_ids must be existing test_ids in the org, so a fixture port check is created
// in the same config and its generated ID is referenced instead of a hardcoded value.
func testAccRecurringDowntimeConfigurationV2Config(name, description, startTime, endTime, recurrenceEndDate string) string {
	return fmt.Sprintf(`
resource "synthetics_create_port_check_v2" "downtime_v2_recurring_fixture" {
	provider = synthetics.synthetics
  test {
    active = true
    frequency = 5
    location_ids = ["aws-us-west-2"]
    scheduling_strategy = "round_robin"
    name = "%[1]s-fixture"
    port = 8081
    protocol = "tcp"
    host = "www.splunk.com"
  }
}

resource "synthetics_create_downtime_configuration_v2" "downtime_configuration_v2_foo_recurring" {
	provider = synthetics.synthetics
  downtime_configuration {
    name = "%[1]s"
    description = "%[2]s"
    rule = "augment_data"
    start_time = "%[3]s"
    end_time = "%[4]s"
    test_ids = [tonumber(synthetics_create_port_check_v2.downtime_v2_recurring_fixture.id)]
    timezone = "America/New_York"
    recurrence {
      repeats {
        type = "daily"
      }
      end {
        type = "on"
        value = "%[5]s"
      }
    }
  }
}
`, name, description, startTime, endTime, recurrenceEndDate)
}

func TestAccCreateRecurringDowntimeConfigurationV2(t *testing.T) {

	name := fmt.Sprintf("acceptance-downtime-configuration-recurring-terraform-test-%d", time.Now().UnixNano())
	description := "The most awesome recurring downtime_configuration. Full of snakes."
	startTime := time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05.000Z")
	endTime := time.Now().Add(25 * time.Hour).Format("2006-01-02T15:04:05.000Z")
	recurrenceEndDate := time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02")

	// name is ForceNew; keep it constant across the update step below so this exercises
	// resourceDowntimeConfigurationV2Update rather than a destroy/recreate.
	updatedDescription := "The most awesome recurring downtime_configuration, now even more updated. Still full of snakes."
	updatedRecurrenceEndDate := time.Now().Add(45 * 24 * time.Hour).Format("2006-01-02")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create the fixture test and the recurring downtime configuration referencing it.
			{
				Config: providerConfig + testAccRecurringDowntimeConfigurationV2Config(name, description, startTime, endTime, recurrenceEndDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.description", description),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.rule", "augment_data"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.name", name),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.start_time", startTime),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.end_time", endTime),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.timezone", "America/New_York"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.recurrence.0.repeats.0.type", "daily"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.recurrence.0.end.0.type", "on"),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.recurrence.0.end.0.value", recurrenceEndDate),
				),
			},
			{
				ResourceName:      "synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring"),
				ImportStateVerify: true,
			},
			// Update the recurring downtime configuration in place - exercises
			// resourceDowntimeConfigurationV2Update on a config with recurrence/repeats/end set.
			{
				Config: providerConfig + testAccRecurringDowntimeConfigurationV2Config(name, updatedDescription, startTime, endTime, updatedRecurrenceEndDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.description", updatedDescription),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.name", name),
					resource.TestCheckResourceAttr("synthetics_create_downtime_configuration_v2.downtime_configuration_v2_foo_recurring", "downtime_configuration.0.recurrence.0.end.0.value", updatedRecurrenceEndDate),
				),
			},
		},
	})
}
