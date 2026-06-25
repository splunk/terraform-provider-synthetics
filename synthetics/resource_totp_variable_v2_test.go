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
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCreateTotpVariableV2(t *testing.T) {
	name := "terraform_totp_login_mfa"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccTotpVariableV2Config(name, "login MFA", "JBSWY3DPEHPK3PXP", 6, 30),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.name", name),
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.description", "login MFA"),
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.digits", "6"),
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.interval", "30"),
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.hmac_digest", "sha1"),
				),
			},
			{
				ResourceName:            "synthetics_create_totp_variable_v2.login_mfa",
				ImportState:             true,
				ImportStateIdFunc:       testAccStateIdFunc("synthetics_create_totp_variable_v2.login_mfa"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"totp_variable.0.secret"},
			},
			{
				Config: providerConfig + testAccTotpVariableV2Config(name, "updated login MFA", "JBSWY3DPEHPK3PXQ", 8, 45),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.description", "updated login MFA"),
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.digits", "8"),
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.interval", "45"),
				),
			},
		},
	})
}

func TestAccBrowserCheckV2WithTotpVariableReference(t *testing.T) {
	name := "terraform_totp_browser_mfa"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccBrowserCheckV2TotpReferenceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_totp_variable_v2.login_mfa", "totp_variable.0.name", name),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_totp_check", "test.0.transactions.0.steps.1.value", fmt.Sprintf("{{totp.%s}}", name)),
				),
			},
		},
	})
}

func TestAccTotpVariableV2DataSources(t *testing.T) {
	name := "terraform_totp_datasource_mfa"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccTotpVariableV2DataSourcesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.synthetics_totp_variable_v2_check.login_mfa",
						"id",
						"synthetics_create_totp_variable_v2.login_mfa",
						"id",
					),
					resource.TestCheckResourceAttr("data.synthetics_totp_variable_v2_check.login_mfa", "totp_variable.#", "1"),
					resource.TestCheckResourceAttr("data.synthetics_totp_variable_v2_check.login_mfa", "totp_variable.0.name", name),
					resource.TestCheckResourceAttr("data.synthetics_totp_variable_v2_check.login_mfa", "totp_variable.0.description", "data source MFA"),
					resource.TestCheckResourceAttr("data.synthetics_totp_variable_v2_check.login_mfa", "totp_variable.0.digits", "6"),
					resource.TestCheckResourceAttr("data.synthetics_totp_variable_v2_check.login_mfa", "totp_variable.0.interval", "30"),
					resource.TestCheckResourceAttr("data.synthetics_totp_variable_v2_check.login_mfa", "totp_variable.0.hmac_digest", "sha1"),
					resource.TestCheckNoResourceAttr("data.synthetics_totp_variable_v2_check.login_mfa", "totp_variable.0.secret"),
					resource.TestCheckResourceAttrSet("data.synthetics_totp_variables_v2_check.all", "id"),
					resource.TestCheckNoResourceAttr("data.synthetics_totp_variables_v2_check.all", "totp_variables.0.secret"),
					testAccCheckTotpVariablesDataSourceContains("data.synthetics_totp_variables_v2_check.all", name),
				),
			},
		},
	})
}

func testAccTotpVariableV2Config(name, description, secret string, digits, interval int) string {
	return fmt.Sprintf(`
resource "synthetics_create_totp_variable_v2" "login_mfa" {
  provider = synthetics.synthetics
  totp_variable {
    name        = %[1]q
    description = %[2]q
    secret      = %[3]q
    digits      = %[4]d
    interval    = %[5]d
    hmac_digest = "sha1"
  }
}
`, name, description, secret, digits, interval)
}

func testAccBrowserCheckV2TotpReferenceConfig(name string) string {
	return testAccTotpVariableV2Config(name, "browser MFA", "JBSWY3DPEHPK3PXP", 6, 30) + fmt.Sprintf(`
resource "synthetics_create_browser_check_v2" "browser_v2_totp_check" {
  provider = synthetics.synthetics
  depends_on = [synthetics_create_totp_variable_v2.login_mfa]

  test {
    active              = true
    device_id           = 1
    frequency           = 5
    location_ids        = ["aws-us-west-2"]
    name                = "terraform-browser-totp-reference"
    scheduling_strategy = "round_robin"

    advanced_settings {
      verify_certificates         = true
      collect_interactive_metrics = false
    }

    transactions {
      name = "Login"

      steps {
        name = "Go to login"
        type = "go_to_url"
        url  = "https://example.com/login"
      }

      steps {
        name          = "Enter MFA code"
        selector      = "mfa-code"
        selector_type = "id"
        type          = "enter_value"
        value         = "{{totp.%[1]s}}"
      }
    }
  }
}
`, name)
}

func testAccTotpVariableV2DataSourcesConfig(name string) string {
	return testAccTotpVariableV2Config(name, "data source MFA", "JBSWY3DPEHPK3PXP", 6, 30) + `
data "synthetics_totp_variable_v2_check" "login_mfa" {
  provider = synthetics.synthetics
  depends_on = [synthetics_create_totp_variable_v2.login_mfa]

  totp_variable {
    id = synthetics_create_totp_variable_v2.login_mfa.totp_variable[0].id
  }
}

data "synthetics_totp_variables_v2_check" "all" {
  provider = synthetics.synthetics
  depends_on = [synthetics_create_totp_variable_v2.login_mfa]
}
`
}

func testAccCheckTotpVariablesDataSourceContains(dataSourceName string, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source %s not found", dataSourceName)
		}

		count, ok := rs.Primary.Attributes["totp_variables.#"]
		if !ok || count == "0" {
			return fmt.Errorf("%s has no totp_variables entries", dataSourceName)
		}

		for key, value := range rs.Primary.Attributes {
			if strings.HasPrefix(key, "totp_variables.") && strings.HasSuffix(key, ".secret") {
				return fmt.Errorf("%s exposed secret attribute %s", dataSourceName, key)
			}
			if strings.HasPrefix(key, "totp_variables.") && strings.HasSuffix(key, ".name") && value == expectedName {
				return nil
			}
		}

		return fmt.Errorf("%s did not include TOTP variable named %q", dataSourceName, expectedName)
	}
}
