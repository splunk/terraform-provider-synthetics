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
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

// The fixture is deliberately secret = false. Asserting the returned `value` is only
// acceptable because this is a non-secret variable holding a non-sensitive literal.
func testAccDataSourceVariableV2Config(name string) string {
	return fmt.Sprintf(`
resource "synthetics_create_variable_v2" "variable_v2_datasource_fixture" {
  provider = synthetics.synthetics
  variable {
    name        = %[1]q
    description = "Terraform acceptance variable data source fixture"
    value       = "datasource-fixture-value"
    secret      = false
  }
}

data "synthetics_variable_v2_check" "variable_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_variable_v2.variable_v2_datasource_fixture]

  variable {
    id = tonumber(synthetics_create_variable_v2.variable_v2_datasource_fixture.id)
  }
}

data "synthetics_variables_v2_check" "variables_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_variable_v2.variable_v2_datasource_fixture]
}
`, name)
}

func TestAccDataSourceVariableV2(t *testing.T) {
	name := testAccUniqueName("acceptance-terraform-variable-v2-datasource")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckFixturesDestroyed(map[string]func(string) (*sc2.RequestDetails, error){
			"synthetics_create_variable_v2": testAccIntLookup(func(id int) (*sc2.RequestDetails, error) {
				_, details, err := testAccProvider.Meta().(*sc2.Client).GetVariableV2(id)
				return details, err
			}),
		}),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourceVariableV2Config(name),
				Check: resource.ComposeTestCheckFunc(
					// Singleton read.
					resource.TestCheckResourceAttrPair(
						"data.synthetics_variable_v2_check.variable_v2_datasource", "id",
						"synthetics_create_variable_v2.variable_v2_datasource_fixture", "id",
					),
					resource.TestCheckResourceAttr("data.synthetics_variable_v2_check.variable_v2_datasource", "variable.#", "1"),
					resource.TestCheckResourceAttr("data.synthetics_variable_v2_check.variable_v2_datasource", "variable.0.name", name),
					resource.TestCheckResourceAttr("data.synthetics_variable_v2_check.variable_v2_datasource", "variable.0.description", "Terraform acceptance variable data source fixture"),
					resource.TestCheckResourceAttr("data.synthetics_variable_v2_check.variable_v2_datasource", "variable.0.value", "datasource-fixture-value"),
					resource.TestCheckResourceAttr("data.synthetics_variable_v2_check.variable_v2_datasource", "variable.0.secret", "false"),

					// List read must contain the fixture.
					resource.TestCheckResourceAttrSet("data.synthetics_variables_v2_check.variables_v2_datasource", "id"),
					testAccCheckDataSourceListContains(
						"data.synthetics_variables_v2_check.variables_v2_datasource",
						"variables", "name", name,
					),
				),
			},
		},
	})
}
