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
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func TestSslCheckV2AcceptanceConfigsOmitUnsupportedCertificateValidations(t *testing.T) {
	for name, config := range map[string]string{
		"create": newSslCheckV2Config,
		"update": updatedSslCheckV2Config,
	} {
		if strings.Contains(config, "{{certificate.days_until_expiration}}") {
			t.Fatalf("%s config uses unsupported certificate metric validation actual", name)
		}
		if strings.Contains(config, "validations {") {
			t.Fatalf("%s config must not send unsupported SSL certificate validations", name)
		}
	}
}

func TestSslCheckV2LastRunIDSchemaUsesString(t *testing.T) {
	for name, testSchema := range map[string]map[string]*schema.Schema{
		"resource":    sslCheckV2ResourceTestSchema(),
		"data source": sslCheckV2DataSourceTestSchema(),
	} {
		lastRunIDSchema := testSchema["last_run_id"]
		if lastRunIDSchema.Type != schema.TypeString {
			t.Errorf("%s last_run_id type = %v, want TypeString", name, lastRunIDSchema.Type)
		}
	}
}

func TestSslCheckV2ValidationSchemaExposesOnlyControllerAllowedFields(t *testing.T) {
	schema := sslCheckV2ValidationSchema()
	allowed := map[string]bool{
		"actual":     true,
		"comparator": true,
		"expected":   true,
		"name":       true,
		"type":       true,
	}

	if len(schema) != len(allowed) {
		t.Fatalf("SSL validation schema fields = %d, want %d: %#v", len(schema), len(allowed), schema)
	}
	for key := range schema {
		if !allowed[key] {
			t.Fatalf("SSL validation schema exposes unsupported field %q", key)
		}
	}
}
