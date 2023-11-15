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

const newLocationV2Config = `
resource "synthetics_create_location_v2" "location_v2_foo" {
	provider = synthetics.synthetics
  location {
    id = "private-bacon"
    label = "awesome aws bacon east location part1"
  }    
}
`

func TestAccCreateLocationV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: providerConfig + newLocationV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_location_v2.location_v2_foo", "location.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_location_v2.location_v2_foo", "location.0.id", "private-bacon"),
					resource.TestCheckResourceAttr("synthetics_create_location_v2.location_v2_foo", "location.0.label", "awesome aws bacon east location part1"),
				),
			},
			{
				ResourceName:      "synthetics_create_location_v2.location_v2_foo",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_location_v2.location_v2_foo"),
				ImportStateVerify: true,
			},
			// Locations are immutable and can not be "updated" so no 2nd step to update the test
		},
	})
}
