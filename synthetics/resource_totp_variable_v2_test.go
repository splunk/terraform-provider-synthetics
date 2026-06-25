package synthetics

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
