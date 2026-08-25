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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use the HASHICUPS_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
variable "observability_token" {
	description = "API token for observability"
}
variable "realm" {
	description = "Splunk Observability realm"
}

provider "synthetics" {
	alias = "synthetics"
	product = "observability"
	realm = var.realm
	apikey = var.observability_token
}
`
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"synthetics": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func TestProviderContainsRecentV2ResourcesAndDataSources(t *testing.T) {
	provider := Provider()

	expectedResources := []string{
		"synthetics_create_ssl_check_v2",
		"synthetics_create_ca_certificate_v2",
		"synthetics_create_client_certificate_v2",
		"synthetics_create_totp_variable_v2",
	}
	for _, name := range expectedResources {
		if _, ok := provider.ResourcesMap[name]; !ok {
			t.Fatalf("Provider().ResourcesMap missing %q", name)
		}
	}

	expectedDataSources := []string{
		"synthetics_ssl_v2_check",
		"synthetics_ca_certificate_v2_check",
		"synthetics_ca_certificates_v2_check",
		"synthetics_client_certificate_v2_check",
		"synthetics_client_certificates_v2_check",
		"synthetics_excluded_file_types_v2_check",
		"synthetics_totp_variable_v2_check",
		"synthetics_totp_variables_v2_check",
	}
	for _, name := range expectedDataSources {
		if _, ok := provider.DataSourcesMap[name]; !ok {
			t.Fatalf("Provider().DataSourcesMap missing %q", name)
		}
	}
}

func testAccStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.Attributes["id"], nil
	}
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("TF_VAR_observability_token"); err == "" {
		t.Fatal("TF_VAR_observability_token environment variable must be set for acceptance tests")
	}
	if err := os.Getenv("TF_VAR_realm"); err == "" {
		t.Fatal("TF_VAR_realm environment variable must be set for acceptance tests")
	}
}
