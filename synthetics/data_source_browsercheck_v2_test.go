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

// The top-level resource id is used rather than test[0].id: the block is a TypeSet on
// most check resources and HCL cannot index a set, so this keeps one idiom throughout.
func testAccDataSourceBrowserCheckV2Config(name string) string {
	return fmt.Sprintf(`
resource "synthetics_create_browser_check_v2" "browser_v2_datasource_fixture" {
  provider = synthetics.synthetics
  test {
    active              = true
    device_id           = 1
    frequency           = 5
    location_ids        = ["aws-us-east-1"]
    name                = %[1]q
    scheduling_strategy = "round_robin"
    automatic_retries   = 1

    advanced_settings {
      verify_certificates         = true
      collect_interactive_metrics = false
    }

    transactions {
      name = "First Synthetic transaction"

      steps {
        name = "01 Go to URL"
        type = "go_to_url"
        url  = "https://www.splunk.com"
      }
    }
  }
}

data "synthetics_browser_v2_check" "browser_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_browser_check_v2.browser_v2_datasource_fixture]

  test {
    id = tonumber(synthetics_create_browser_check_v2.browser_v2_datasource_fixture.id)
  }
}
`, name)
}

func TestAccDataSourceBrowserCheckV2(t *testing.T) {
	name := testAccUniqueName("acceptance-terraform-browser-v2-datasource")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourceBrowserCheckV2Config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.synthetics_browser_v2_check.browser_v2_datasource", "id",
						"synthetics_create_browser_check_v2.browser_v2_datasource_fixture", "id",
					),
					resource.TestCheckResourceAttr("data.synthetics_browser_v2_check.browser_v2_datasource", "test.#", "1"),
					resource.TestCheckResourceAttr("data.synthetics_browser_v2_check.browser_v2_datasource", "test.0.name", name),
					resource.TestCheckResourceAttr("data.synthetics_browser_v2_check.browser_v2_datasource", "test.0.active", "true"),
					resource.TestCheckResourceAttr("data.synthetics_browser_v2_check.browser_v2_datasource", "test.0.frequency", "5"),
					resource.TestCheckResourceAttr("data.synthetics_browser_v2_check.browser_v2_datasource", "test.0.scheduling_strategy", "round_robin"),
					resource.TestCheckResourceAttrSet("data.synthetics_browser_v2_check.browser_v2_datasource", "test.0.type"),
					testAccCheckDataSourceListContains(
						"data.synthetics_browser_v2_check.browser_v2_datasource",
						"test.0.transactions", "name", "First Synthetic transaction",
					),
				),
			},
		},
	})
}
