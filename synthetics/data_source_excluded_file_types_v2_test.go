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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func TestDataSourceExcludedFileTypesV2Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/excluded_file_types" {
			t.Fatalf("path = %q, want /excluded_file_types", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		_, _ = w.Write([]byte(`{"excludedFileTypes":["chartbeat","google_analytics"]}`))
	}))
	defer server.Close()

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceExcludedFileTypesV2().Schema, map[string]interface{}{})

	diags := dataSourceExcludedFileTypesV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceExcludedFileTypesV2Read() diagnostics = %#v", diags)
	}
	if d.Id() != "global_excluded_file_types_synthetics" {
		t.Fatalf("id = %q, want global_excluded_file_types_synthetics", d.Id())
	}

	values := d.Get("excluded_file_types").(*schema.Set)
	if !values.Contains("chartbeat") {
		t.Fatalf("excluded_file_types missing chartbeat: %#v", values.List())
	}
	if !values.Contains("google_analytics") {
		t.Fatalf("excluded_file_types missing google_analytics: %#v", values.List())
	}
}

func TestAccDataSourceExcludedFileTypesV2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "synthetics_excluded_file_types_v2_check" "types" {
  provider = synthetics.synthetics
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.synthetics_excluded_file_types_v2_check.types", "id"),
					resource.TestCheckTypeSetElemAttr("data.synthetics_excluded_file_types_v2_check.types", "excluded_file_types.*", "google_analytics"),
				),
			},
		},
	})
}
