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
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sc2 "github.com/splunk/syntheticsclient/v3/syntheticsclientv2"
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

func TestProviderConfigureDefaultClient(t *testing.T) {
	// apiurl is omitted below so Provider().Schema's EnvDefaultFunc would otherwise pull
	// an ambient API_URL and make this test depend on the environment it runs in.
	t.Setenv("API_URL", "")

	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"apikey": "token123",
		"realm":  "us1",
	})

	got, diags := providerConfigure(context.Background(), d)
	if diags.HasError() {
		t.Fatalf("providerConfigure() diags = %#v, want none", diags)
	}
	c, ok := got.(*sc2.Client)
	if !ok {
		t.Fatalf("providerConfigure() = %T, want *sc2.Client", got)
	}
	if !strings.Contains(c.String(), "https://api.us1.signalfx.com/v2/synthetics") {
		t.Fatalf("client URL = %s, want the default us1 endpoint", c.String())
	}
}

func TestProviderConfigureApiUrlOverride(t *testing.T) {
	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"apikey": "token123",
		"realm":  "us1",
		"apiurl": "https://api.custom.example.com/",
	})

	got, diags := providerConfigure(context.Background(), d)
	if diags.HasError() {
		t.Fatalf("providerConfigure() diags = %#v, want none", diags)
	}
	c, ok := got.(*sc2.Client)
	if !ok {
		t.Fatalf("providerConfigure() = %T, want *sc2.Client", got)
	}
	if !strings.Contains(c.String(), "https://api.custom.example.com/v2/synthetics") {
		t.Fatalf("client URL = %s, want the apiurl override with trailing slash trimmed", c.String())
	}
}

func TestProviderConfigureEmptyCredentialsStillReturnsDefaultClient(t *testing.T) {
	// Every field is omitted below, so ambient OBSERVABILITY_API_TOKEN/REALM/API_URL
	// would otherwise be pulled in via EnvDefaultFunc, defeating the "empty credentials"
	// premise this test is named for.
	t.Setenv("OBSERVABILITY_API_TOKEN", "")
	t.Setenv("REALM", "")
	t.Setenv("API_URL", "")

	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})

	got, diags := providerConfigure(context.Background(), d)
	if diags.HasError() {
		t.Fatalf("providerConfigure() diags = %#v, want none", diags)
	}
	if _, ok := got.(*sc2.Client); !ok {
		t.Fatalf("providerConfigure() = %T, want *sc2.Client", got)
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
