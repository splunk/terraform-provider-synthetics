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

const newVariableV2Config = `
resource "synthetics_create_variable_v2" "variable_v2_foo" {
	provider = synthetics.synthetics
  variable {
    description = "The most awesome variable. Full of snakes."
    value = "barv3v3"
    name = "acceptance-variable-terraform-test"
    secret = false  
  }    
}
`

const updatedVariableV2Config = `
resource "synthetics_create_variable_v2" "variable_v2_foo" {
	provider = synthetics.synthetics
  variable {
    description = "The most awesome variable. Full of snakes and birbs."
    value = "barv3v3"
    name = "acceptance-variable-terraform-test"
		// Any change to 'secret' will force re-creation of the resource
    secret = false  
  }    
}
`

func TestAccCreateUpdateVariableV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: providerConfig + newVariableV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.0.description", "The most awesome variable. Full of snakes."),
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.0.value", "barv3v3"),
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.0.name", "acceptance-variable-terraform-test"),
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.0.secret", "false"),
				),
			},
			{
				ResourceName:      "synthetics_create_variable_v2.variable_v2_foo",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_variable_v2.variable_v2_foo"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: providerConfig + updatedVariableV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.0.description", "The most awesome variable. Full of snakes and birbs."),
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.0.value", "barv3v3"),
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.0.name", "acceptance-variable-terraform-test"),
					resource.TestCheckResourceAttr("synthetics_create_variable_v2.variable_v2_foo", "variable.0.secret", "false"),
				),
			},
		},
	})
}
