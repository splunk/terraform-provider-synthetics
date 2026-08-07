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

// A disposable private location is used for the singleton read so the test does not
// depend on any particular public location existing in the org. Location ids are
// validated against privateLocationIDPattern, which forbids digits, so the id comes
// from testAccUniquePrivateLocationID rather than a numeric suffix.
func testAccDataSourceLocationV2Config(id, label string) string {
	return fmt.Sprintf(`
resource "synthetics_create_location_v2" "location_v2_datasource_fixture" {
  provider = synthetics.synthetics
  location {
    id    = %[1]q
    label = %[2]q
  }
}

data "synthetics_location_v2_check" "location_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_location_v2.location_v2_datasource_fixture]

  location {
    id = synthetics_create_location_v2.location_v2_datasource_fixture.id
  }
}

data "synthetics_locations_v2_check" "locations_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_location_v2.location_v2_datasource_fixture]
}
`, id, label)
}

func TestAccDataSourceLocationV2(t *testing.T) {
	id := testAccUniquePrivateLocationID("acc-datasource")
	label := "Terraform acceptance location data source fixture"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourceLocationV2Config(id, label),
				Check: resource.ComposeTestCheckFunc(
					// Singleton read of the disposable private location.
					resource.TestCheckResourceAttrPair(
						"data.synthetics_location_v2_check.location_v2_datasource", "id",
						"synthetics_create_location_v2.location_v2_datasource_fixture", "id",
					),
					resource.TestCheckResourceAttr("data.synthetics_location_v2_check.location_v2_datasource", "location.#", "1"),
					resource.TestCheckResourceAttr("data.synthetics_location_v2_check.location_v2_datasource", "location.0.id", id),
					resource.TestCheckResourceAttr("data.synthetics_location_v2_check.location_v2_datasource", "location.0.label", label),
					resource.TestCheckResourceAttrSet("data.synthetics_location_v2_check.location_v2_datasource", "location.0.type"),

					// List read: the fixture must appear, and the account-global fields
					// are asserted only for stable non-empty shape.
					resource.TestCheckResourceAttrSet("data.synthetics_locations_v2_check.locations_v2_datasource", "id"),
					testAccCheckDataSourceListContains(
						"data.synthetics_locations_v2_check.locations_v2_datasource",
						"locations", "id", id,
					),
					testAccCheckDataSourceElemFieldsNonEmpty(
						"data.synthetics_locations_v2_check.locations_v2_datasource",
						"locations", "id", "label", "type",
					),
					resource.TestCheckResourceAttrSet("data.synthetics_locations_v2_check.locations_v2_datasource", "default_location_ids.#"),
				),
			},
		},
	})
}
