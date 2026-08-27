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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v3/syntheticsclientv2"
)

func TestDataSourceChromeFlagsRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chrome_flags" {
			t.Fatalf("path = %q, want /chrome_flags", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		_, _ = w.Write([]byte(`{"chromeFlags":[{"name":"--disable-http2","label":"Disable h2","description":"Disables HTTP/2.0.","acceptsValue":false},{"name":"--proxy-server","label":"Proxy server","description":"Use a specified proxy server.","acceptsValue":true}]}`))
	}))
	defer server.Close()

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceChromeFlags().Schema, map[string]interface{}{})

	diags := dataSourceChromeFlagsRead(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceChromeFlagsRead() diagnostics = %#v", diags)
	}
	if d.Id() != "global_chrome_flags_synthetics" {
		t.Fatalf("id = %q, want global_chrome_flags_synthetics", d.Id())
	}

	values := d.Get("chrome_flags").(*schema.Set)
	if values.Len() != 2 {
		t.Fatalf("chrome_flags length = %d, want 2: %#v", values.Len(), values.List())
	}

	found := false
	for _, v := range values.List() {
		flag := v.(map[string]interface{})
		if flag["name"] == "--proxy-server" {
			found = true
			if flag["label"] != "Proxy server" {
				t.Errorf("label = %v, want Proxy server", flag["label"])
			}
			if flag["description"] != "Use a specified proxy server." {
				t.Errorf("description = %v, want %q", flag["description"], "Use a specified proxy server.")
			}
			if flag["accepts_value"] != true {
				t.Errorf("accepts_value = %v, want true", flag["accepts_value"])
			}
		}
	}
	if !found {
		t.Fatalf("chrome_flags missing --proxy-server: %#v", values.List())
	}
}

func TestDataSourceChromeFlagsReadReturnsErrorOnFailure(t *testing.T) {
	unreachableClient := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, "http://127.0.0.1:1"))
	d := schema.TestResourceDataRaw(t, dataSourceChromeFlags().Schema, map[string]interface{}{})

	diags := dataSourceChromeFlagsRead(context.Background(), d, unreachableClient)
	if !diags.HasError() {
		t.Fatal("expected a connection error")
	}
}

// The chrome flag catalog is account-global and cannot be created by Terraform, so there
// is no fixture to look for. Assertions are limited to a stable representative shape: the
// collection is non-empty and some element carries every field the flattener populates.
const testAccDataSourceChromeFlagsConfig = `
data "synthetics_chrome_flags" "flags" {
  provider = synthetics.synthetics
}
`

func TestAccDataSourceChromeFlags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccDataSourceChromeFlagsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.synthetics_chrome_flags.flags", "id", "global_chrome_flags_synthetics"),
					testAccCheckDataSourceCollectionNotEmpty(
						"data.synthetics_chrome_flags.flags", "chrome_flags",
					),
					testAccCheckDataSourceElemFieldsNonEmpty(
						"data.synthetics_chrome_flags.flags",
						"chrome_flags", "name", "label", "description",
					),
				),
			},
		},
	})
}
