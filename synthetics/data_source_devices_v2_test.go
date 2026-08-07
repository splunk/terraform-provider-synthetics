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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The device catalog is account-global and cannot be created by Terraform, so there is no
// fixture to look for. Assertions are limited to stable representative shape: the
// collection is non-empty and some element carries every field the flattener populates.
const testAccDataSourceDevicesV2Config = `
data "synthetics_devices_v2_check" "devices_v2_datasource" {
  provider = synthetics.synthetics
}
`

func TestAccDataSourceDevicesV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourceDevicesV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.synthetics_devices_v2_check.devices_v2_datasource", "id", "global_devices_synthetics"),
					testAccCheckDataSourceCollectionNotEmpty(
						"data.synthetics_devices_v2_check.devices_v2_datasource", "devices",
					),
					testAccCheckDataSourceElemFieldsNonEmpty(
						"data.synthetics_devices_v2_check.devices_v2_datasource",
						"devices", "id", "label", "user_agent", "viewport_height", "viewport_width",
					),
				),
			},
		},
	})
}
