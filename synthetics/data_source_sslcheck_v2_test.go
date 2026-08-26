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
	sc2 "github.com/splunk/syntheticsclient/v3/syntheticsclientv2"
)

func testAccDataSourceSslCheckV2Config(name string) string {
	return fmt.Sprintf(`
resource "synthetics_create_ssl_check_v2" "ssl_v2_datasource_fixture" {
  provider = synthetics.synthetics
  test {
    active               = false
    frequency            = 5
    location_ids         = ["aws-us-east-1"]
    name                 = %[1]q
    host                 = "www.splunk.com"
    port                 = 443
    server_name          = "www.splunk.com"
    allow_self_signed    = false
    allow_untrusted_root = false
    scheduling_strategy  = "round_robin"
    automatic_retries    = 1
  }
}

data "synthetics_ssl_v2_check" "ssl_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_ssl_check_v2.ssl_v2_datasource_fixture]

  test {
    id = tonumber(synthetics_create_ssl_check_v2.ssl_v2_datasource_fixture.id)
  }
}
`, name)
}

func TestAccDataSourceSslCheckV2(t *testing.T) {
	name := testAccUniqueName("acceptance-terraform-ssl-v2-datasource")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckFixturesDestroyed(map[string]func(string) (*sc2.RequestDetails, error){
			"synthetics_create_ssl_check_v2": testAccIntLookup(func(id int) (*sc2.RequestDetails, error) {
				_, details, err := testAccProvider.Meta().(*sc2.Client).GetSslCheckV2(id)
				return details, err
			}),
		}),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourceSslCheckV2Config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.synthetics_ssl_v2_check.ssl_v2_datasource", "id",
						"synthetics_create_ssl_check_v2.ssl_v2_datasource_fixture", "id",
					),
					resource.TestCheckResourceAttr("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.#", "1"),
					resource.TestCheckResourceAttr("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.0.name", name),
					resource.TestCheckResourceAttr("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.0.active", "false"),
					resource.TestCheckResourceAttr("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.0.frequency", "5"),
					resource.TestCheckResourceAttr("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.0.host", "www.splunk.com"),
					resource.TestCheckResourceAttr("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.0.port", "443"),
					resource.TestCheckResourceAttr("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.0.server_name", "www.splunk.com"),
					resource.TestCheckResourceAttr("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.0.allow_self_signed", "false"),
					resource.TestCheckResourceAttrSet("data.synthetics_ssl_v2_check.ssl_v2_datasource", "test.0.type"),
				),
			},
		},
	})
}
