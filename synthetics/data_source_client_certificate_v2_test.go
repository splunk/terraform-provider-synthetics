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

func TestClientCertificateV2DataSourceIsMetadataOnly(t *testing.T) {
	dataSource := dataSourceClientCertificateV2()

	for _, key := range []string{"id", "name", "description", "domain", "expires_at", "created_at", "created_by", "updated_at", "updated_by"} {
		if _, ok := dataSource.Schema[key]; !ok {
			t.Fatalf("client certificate data source schema missing %q", key)
		}
	}
	if _, ok := dataSource.Schema["public_key"]; ok {
		t.Fatal("client certificate data source must not expose public_key")
	}
	if _, ok := dataSource.Schema["private_key"]; ok {
		t.Fatal("client certificate data source must not expose private_key")
	}
}

// The client certificate data sources are metadata-only by design: the API returns
// certificate and private key material redacted. These assertions cover the metadata
// fields and positively assert that no key material reaches state. No assertion or
// failure message here interpolates certificate or key content.
func TestAccDataSourceClientCertificateV2(t *testing.T) {
	name := testAccUniqueName("terraform-client-cert-datasource")
	publicKey, privateKey := testAccClientCertificateV2Material(t)

	config := testAccClientCertificateV2Config(
		name,
		"Terraform acceptance client certificate data source fixture",
		"client.crt", "client.key",
		publicKey, privateKey,
	) + `
data "synthetics_client_certificate_v2_check" "mtls" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_client_certificate_v2.mtls]

  id = tonumber(synthetics_create_client_certificate_v2.mtls.id)
}

data "synthetics_client_certificates_v2_check" "all" {
  provider   = synthetics.synthetics
  depends_on = [synthetics_create_client_certificate_v2.mtls]
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check: resource.ComposeTestCheckFunc(
					// Singleton read: metadata only.
					resource.TestCheckResourceAttrPair(
						"data.synthetics_client_certificate_v2_check.mtls", "id",
						"synthetics_create_client_certificate_v2.mtls", "id",
					),
					resource.TestCheckResourceAttr("data.synthetics_client_certificate_v2_check.mtls", "name", name),
					resource.TestCheckResourceAttr("data.synthetics_client_certificate_v2_check.mtls", "description", "Terraform acceptance client certificate data source fixture"),
					resource.TestCheckResourceAttr("data.synthetics_client_certificate_v2_check.mtls", "domain", "api.example.com"),
					resource.TestCheckResourceAttrSet("data.synthetics_client_certificate_v2_check.mtls", "expires_at"),
					testAccCheckDataSourceAttrsAbsent(
						"data.synthetics_client_certificate_v2_check.mtls",
						"public_key", "private_key", "content", "password",
					),

					// List read must contain the fixture, still metadata only.
					resource.TestCheckResourceAttrSet("data.synthetics_client_certificates_v2_check.all", "id"),
					testAccCheckDataSourceListContains(
						"data.synthetics_client_certificates_v2_check.all",
						"client_certificates", "name", name,
					),
					testAccCheckDataSourceAttrsAbsent(
						"data.synthetics_client_certificates_v2_check.all",
						"public_key", "private_key", "content", "password",
					),
				),
			},
		},
	})
}
