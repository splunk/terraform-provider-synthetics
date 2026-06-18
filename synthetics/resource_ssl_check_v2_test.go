package synthetics

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const newSslCheckV2Config = `
resource "synthetics_create_ssl_check_v2" "ssl_v2_check" {
  provider = synthetics.synthetics
  test {
    active               = false
    frequency            = 5
    location_ids         = ["aws-us-east-1"]
    automatic_retries    = 1
    scheduling_strategy  = "round_robin"
    name                 = "acceptance-Terraform-SSL-V2"
    host                 = "www.splunk.com"
    port                 = 443
    server_name          = "www.splunk.com"
    allow_self_signed    = false
    allow_untrusted_root = false

    custom_properties {
      key   = "env"
      value = "terraform-acceptance"
    }

    validations {
      name       = "Certificate expires later"
      type       = "assert_numeric"
      actual     = "{{certificate.days_until_expiration}}"
      comparator = "is_greater_than"
      expected   = "7"
    }
  }
}
`

const updatedSslCheckV2Config = `
resource "synthetics_create_ssl_check_v2" "ssl_v2_check" {
  provider = synthetics.synthetics
  test {
    active               = false
    frequency            = 15
    location_ids         = ["aws-us-east-1"]
    automatic_retries    = 0
    scheduling_strategy  = "concurrent"
    name                 = "acceptance-updated-Terraform-SSL-V2"
    host                 = "example.com"
    port                 = 443
    allow_self_signed    = false
    allow_untrusted_root = false

    custom_properties {
      key   = "env"
      value = "terraform-acceptance-updated"
    }

    validations {
      name       = "Certificate expires later"
      type       = "assert_numeric"
      actual     = "{{certificate.days_until_expiration}}"
      comparator = "is_greater_than"
      expected   = "7"
    }
  }
}
`

func TestAccCreateUpdateSslCheckV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + newSslCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.active", "false"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.frequency", "5"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.host", "www.splunk.com"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.port", "443"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.server_name", "www.splunk.com"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.allow_self_signed", "false"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.allow_untrusted_root", "false"),
				),
			},
			{
				ResourceName:      "synthetics_create_ssl_check_v2.ssl_v2_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_ssl_check_v2.ssl_v2_check"),
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + updatedSslCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.frequency", "15"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.scheduling_strategy", "concurrent"),
					resource.TestCheckResourceAttr("synthetics_create_ssl_check_v2.ssl_v2_check", "test.0.host", "example.com"),
				),
			},
		},
	})
}
