package synthetics

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const acceptanceCaCertificateContent = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURRVENDQWltZ0F3SUJBZ0lVS01CTmgrZmtHa2RwRTRLZXVGUnpnc2o4ck5jd0RRWUpLb1pJaHZjTkFRRUwKQlFBd01ERXVNQ3dHQTFVRUF3d2xkR1Z5Y21GbWIzSnRMWEJ5YjNacFpHVnlMWE41Ym5Sb1pYUnBZM010ZEdWegpkQzFqWVRBZUZ3MHlOakEyTVRnd056QTBOVGxhRncweU5qQTJNVGt3TnpBME5UbGFNREF4TGpBc0JnTlZCQU1NCkpYUmxjbkpoWm05eWJTMXdjbTkyYVdSbGNpMXplVzUwYUdWMGFXTnpMWFJsYzNRdFkyRXdnZ0VpTUEwR0NTcUcKU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQ3ZVOVpFb2xBVzdaQ05ZTkpnNC9EaXlpdENKSHg5eDJKTApXVjVTYVBRK2JmOGRCeEtHMExRWDJTZXd3SEU2U3Vzc2hQc2x3ZExuOUtUaVExMWFULzlJZTRkb3pnTi9rV0pjClhIY0huK1hNSkxZVTZaYkFDcllSb1JZWHJYbHJIeU05YzdDeHYxUFJxaFF4bjRqR2hRaDViQjhRZmdXNkh6V1UKOEtmeldxU0dsUDZ1RXJPQ0l3S1VIaUNvSVcxdWZwUWloUEs4enVyZUt4bVNlQ1dORnVsUE1Ud0hQT285UGJzRApRb3J5dHEwZUZ6NkN0Wnhrd1RlWnJNU1M4QnJvZmlKLzM4bXNBTFdzQjZzSXBwV0hXelVPMEF4UXd4RUJ5OGNXCnBwS2ptSzI5QjZ0bFZ0bElBczMzU3htL1VtWVpnYnk0ZWpRL21JSmpTSmRhSWxPaWxoZGZBZ01CQUFHalV6QlIKTUIwR0ExVWREZ1FXQkJSN2NYQm5DMGxwRFBjb2NtdDF2a0pta2RnYTFqQWZCZ05WSFNNRUdEQVdnQlI3Y1hCbgpDMGxwRFBjb2NtdDF2a0pta2RnYTFqQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01BMEdDU3FHU0liM0RRRUJDd1VBCkE0SUJBUUNBZDR6MVFIMUZLY2FIVEowVDhHY2VkbWVudXVJV3BtNFdIUHZuMkxjcDAxU09CcUJmVk1TVnRGdVkKWDBoeldzbDVYOEZUODFiNnp1SGk3TGFwQ2RtMWNTWDJuVDJVZTJOVndndGRsSUxsV0d2dlBuWVVtOTl3NnFZWgpJMUdkM3diRjNXb0xvNGpCNWJsTWdzZHNrK1VIMXljYmNzdkRJMWdseXU3SGtBTERkaVhsbXFlbnBoTzh2RjVQCk9aM2NNRE1aZXF2cnUvb2djNGdaMUI5ejJXc1MrOWhmbUFyeTJDUEJET2lOZEdGMEhtWmdrSGlUSnQzWHIzSkoKZkZPRmxxTjlwdnlNWlorbmZZcTZVNzFMdE53RWZCWmgwVW1pbUtld2h6WFlQdkVMdm1jVm1Jc1pHMjNFbDE0cQo4cU5Pdmt1SEI4eVZmZUhaSGdPeFBwL2xIUURTCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"

const newCaCertificateV2Config = `
resource "synthetics_create_ca_certificate_v2" "ca_certificate_v2" {
  provider = synthetics.synthetics
  ca_certificate {
    name           = "acceptance-Terraform-CA-Certificate-V2"
    description    = "Terraform acceptance CA certificate"
    content        = "` + acceptanceCaCertificateContent + `"
    file_extension = "pem"
    filename       = "terraform-ca.pem"
  }
}
`

const updatedCaCertificateV2Config = `
resource "synthetics_create_ca_certificate_v2" "ca_certificate_v2" {
  provider = synthetics.synthetics
  ca_certificate {
    name           = "acceptance-Terraform-CA-Certificate-V2"
    description    = "Terraform acceptance CA certificate updated"
    content        = "` + acceptanceCaCertificateContent + `"
    file_extension = "pem"
    filename       = "terraform-ca-updated.pem"
  }
}
`

func TestAccCreateUpdateCaCertificateV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + newCaCertificateV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_ca_certificate_v2.ca_certificate_v2", "ca_certificate.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_ca_certificate_v2.ca_certificate_v2", "ca_certificate.0.name", "acceptance-Terraform-CA-Certificate-V2"),
					resource.TestCheckResourceAttr("synthetics_create_ca_certificate_v2.ca_certificate_v2", "ca_certificate.0.description", "Terraform acceptance CA certificate"),
					resource.TestCheckResourceAttr("synthetics_create_ca_certificate_v2.ca_certificate_v2", "ca_certificate.0.file_extension", "pem"),
				),
			},
			{
				ResourceName:      "synthetics_create_ca_certificate_v2.ca_certificate_v2",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_ca_certificate_v2.ca_certificate_v2"),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ca_certificate.0.content",
				},
			},
			{
				Config: providerConfig + updatedCaCertificateV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_ca_certificate_v2.ca_certificate_v2", "ca_certificate.0.description", "Terraform acceptance CA certificate updated"),
					resource.TestCheckResourceAttr("synthetics_create_ca_certificate_v2.ca_certificate_v2", "ca_certificate.0.filename", "terraform-ca-updated.pem"),
				),
			},
		},
	})
}
