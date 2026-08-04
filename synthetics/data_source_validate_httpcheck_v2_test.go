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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func httpValidateTestData() map[string]interface{} {
	return map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                "http-validate",
				"type":                "http",
				"url":                 "https://example.com",
				"active":              true,
				"frequency":           5,
				"scheduling_strategy": "round_robin",
				"request_method":      "GET",
				"body":                "",
				"location_ids":        []interface{}{"aws-us-east-1"},
				"user_agent":          "",
				"verify_certificates": true,
				"headers":             []interface{}{},
				"validations":         []interface{}{},
				"custom_properties":   []interface{}{},
				"automatic_retries":   0,
			},
		},
	}
}

func TestDataSourceValidateHttpCheckV2ReadCreateStyleSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tests/http/validate" {
			t.Fatalf("path = %q, want /tests/http/validate", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want POST", r.Method)
		}
		_, _ = w.Write([]byte(`{"valid":true,"message":"Test is valid","details":[]}`))
	}))
	defer server.Close()

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceValidateHttpCheckV2().Schema, httpValidateTestData())

	diags := dataSourceValidateHttpCheckV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceValidateHttpCheckV2Read() diagnostics = %#v", diags)
	}
	if !d.Get("valid").(bool) {
		t.Fatalf("valid = %#v, want true", d.Get("valid"))
	}
}

func TestDataSourceValidateHttpCheckV2ReadCreateStyleFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tests/http/validate" {
			t.Fatalf("path = %q, want /tests/http/validate", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want POST", r.Method)
		}
		_, _ = w.Write([]byte(`{"valid":false,"message":"Test is invalid","details":{"url":["can't be blank"]}}`))
	}))
	defer server.Close()

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceValidateHttpCheckV2().Schema, httpValidateTestData())

	diags := dataSourceValidateHttpCheckV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceValidateHttpCheckV2Read() diagnostics = %#v", diags)
	}
	if d.Get("valid").(bool) {
		t.Fatalf("valid = %#v, want false", d.Get("valid"))
	}

	fieldErrors := d.Get("field_errors").([]interface{})
	if len(fieldErrors) != 1 {
		t.Fatalf("field_errors = %#v, want one entry", fieldErrors)
	}
	entry := fieldErrors[0].(map[string]interface{})
	if entry["field"] != "url" {
		t.Fatalf("field_errors[0].field = %#v, want url", entry["field"])
	}
}

func TestDataSourceValidateHttpCheckV2ReadUpdateStyleSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tests/http/4001/validate" {
			t.Fatalf("path = %q, want /tests/http/4001/validate", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Fatalf("method = %q, want PUT", r.Method)
		}
		_, _ = w.Write([]byte(`{"valid":true,"message":"Test is valid","details":[]}`))
	}))
	defer server.Close()

	testData := httpValidateTestData()
	testData["test_id"] = 4001

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceValidateHttpCheckV2().Schema, testData)

	diags := dataSourceValidateHttpCheckV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceValidateHttpCheckV2Read() diagnostics = %#v", diags)
	}
	if !d.Get("valid").(bool) {
		t.Fatalf("valid = %#v, want true", d.Get("valid"))
	}
}

func TestDataSourceValidateHttpCheckV2ReadUpdateStyleFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tests/http/4002/validate" {
			t.Fatalf("path = %q, want /tests/http/4002/validate", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Fatalf("method = %q, want PUT", r.Method)
		}
		_, _ = w.Write([]byte(`{"valid":false,"message":"Test is invalid","details":{"request_method":["is not included in the list"]}}`))
	}))
	defer server.Close()

	testData := httpValidateTestData()
	testData["test_id"] = 4002

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceValidateHttpCheckV2().Schema, testData)

	diags := dataSourceValidateHttpCheckV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceValidateHttpCheckV2Read() diagnostics = %#v", diags)
	}
	if d.Get("valid").(bool) {
		t.Fatalf("valid = %#v, want false", d.Get("valid"))
	}

	fieldErrors := d.Get("field_errors").([]interface{})
	if len(fieldErrors) != 1 {
		t.Fatalf("field_errors = %#v, want one entry", fieldErrors)
	}
	entry := fieldErrors[0].(map[string]interface{})
	if entry["field"] != "request_method" {
		t.Fatalf("field_errors[0].field = %#v, want request_method", entry["field"])
	}
}
