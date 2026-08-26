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

// Reuses acceptanceCaCertificateContent from resource_ca_certificate_v2_test.go rather
// than minting new CA material.
func testAccDataSourceCaCertificateV2Config(name string) string {
	return fmt.Sprintf(`
resource "synthetics_create_ca_certificate_v2" "ca_certificate_v2_datasource_fixture" {
  provider = synthetics.synthetics
  ca_certificate {
    name           = %[1]q
    description    = "Terraform acceptance CA certificate data source fixture"
    content        = %[2]q
    file_extension = "pem"
    filename       = "terraform-ca-datasource.pem"
  }
}

data "synthetics_ca_certificate_v2_check" "ca_certificate_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_ca_certificate_v2.ca_certificate_v2_datasource_fixture]

  ca_certificate {
    id = tonumber(synthetics_create_ca_certificate_v2.ca_certificate_v2_datasource_fixture.id)
  }
}

data "synthetics_ca_certificates_v2_check" "ca_certificates_v2_datasource" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_ca_certificate_v2.ca_certificate_v2_datasource_fixture]
}
`, name, acceptanceCaCertificateContent)
}

func TestAccDataSourceCaCertificateV2(t *testing.T) {
	name := testAccUniqueName("acceptance-terraform-ca-certificate-v2-datasource")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckFixturesDestroyed(map[string]func(string) (*sc2.RequestDetails, error){
			"synthetics_create_ca_certificate_v2": testAccIntLookup(func(id int) (*sc2.RequestDetails, error) {
				_, details, err := testAccProvider.Meta().(*sc2.Client).GetCaCertificateV2(id)
				return details, err
			}),
		}),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourceCaCertificateV2Config(name),
				Check: resource.ComposeTestCheckFunc(
					// Singleton read.
					resource.TestCheckResourceAttrPair(
						"data.synthetics_ca_certificate_v2_check.ca_certificate_v2_datasource", "id",
						"synthetics_create_ca_certificate_v2.ca_certificate_v2_datasource_fixture", "id",
					),
					resource.TestCheckResourceAttr("data.synthetics_ca_certificate_v2_check.ca_certificate_v2_datasource", "ca_certificate.#", "1"),
					resource.TestCheckResourceAttr("data.synthetics_ca_certificate_v2_check.ca_certificate_v2_datasource", "ca_certificate.0.name", name),
					resource.TestCheckResourceAttr("data.synthetics_ca_certificate_v2_check.ca_certificate_v2_datasource", "ca_certificate.0.description", "Terraform acceptance CA certificate data source fixture"),
					resource.TestCheckResourceAttr("data.synthetics_ca_certificate_v2_check.ca_certificate_v2_datasource", "ca_certificate.0.file_extension", "pem"),
					resource.TestCheckResourceAttr("data.synthetics_ca_certificate_v2_check.ca_certificate_v2_datasource", "ca_certificate.0.filename", "terraform-ca-datasource.pem"),

					// List read must contain the fixture. Asserted by name, never by content.
					resource.TestCheckResourceAttrSet("data.synthetics_ca_certificates_v2_check.ca_certificates_v2_datasource", "id"),
					testAccCheckDataSourceListContains(
						"data.synthetics_ca_certificates_v2_check.ca_certificates_v2_datasource",
						"ca_certificates", "name", name,
					),
				),
			},
		},
	})
}
