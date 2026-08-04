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

func sslValidateTestData() map[string]interface{} {
	return map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                "ssl-validate",
				"active":              true,
				"frequency":           5,
				"scheduling_strategy": "round_robin",
				"location_ids":        []interface{}{"aws-us-east-1"},
				"host":                "example.com",
				"port":                443,
			},
		},
	}
}

func TestDataSourceValidateSslCheckV2ReadCreateStyleSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tests/ssl/validate" {
			t.Fatalf("path = %q, want /tests/ssl/validate", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want POST", r.Method)
		}
		_, _ = w.Write([]byte(`{"valid":true,"message":"Test is valid","details":[]}`))
	}))
	defer server.Close()

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceValidateSslCheckV2().Schema, sslValidateTestData())

	diags := dataSourceValidateSslCheckV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceValidateSslCheckV2Read() diagnostics = %#v", diags)
	}
	if !d.Get("valid").(bool) {
		t.Fatalf("valid = %#v, want true", d.Get("valid"))
	}
	if len(d.Get("field_errors").([]interface{})) != 0 {
		t.Fatalf("field_errors = %#v, want empty", d.Get("field_errors"))
	}
}

func TestDataSourceValidateSslCheckV2ReadCreateStyleFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tests/ssl/validate" {
			t.Fatalf("path = %q, want /tests/ssl/validate", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method = %q, want POST", r.Method)
		}
		_, _ = w.Write([]byte(`{"valid":false,"message":"Test is invalid","details":{"host":["can't be blank"]}}`))
	}))
	defer server.Close()

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceValidateSslCheckV2().Schema, sslValidateTestData())

	diags := dataSourceValidateSslCheckV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceValidateSslCheckV2Read() diagnostics = %#v", diags)
	}
	if d.Get("valid").(bool) {
		t.Fatalf("valid = %#v, want false", d.Get("valid"))
	}

	fieldErrors := d.Get("field_errors").([]interface{})
	if len(fieldErrors) != 1 {
		t.Fatalf("field_errors = %#v, want one entry", fieldErrors)
	}
	entry := fieldErrors[0].(map[string]interface{})
	if entry["field"] != "host" {
		t.Fatalf("field_errors[0].field = %#v, want host", entry["field"])
	}
	messages := entry["messages"].([]interface{})
	if len(messages) != 1 || messages[0] != "can't be blank" {
		t.Fatalf("field_errors[0].messages = %#v, want [can't be blank]", messages)
	}
}

func TestDataSourceValidateSslCheckV2ReadUpdateStyleSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tests/ssl/1655/validate" {
			t.Fatalf("path = %q, want /tests/ssl/1655/validate", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Fatalf("method = %q, want PUT", r.Method)
		}
		_, _ = w.Write([]byte(`{"valid":true,"message":"Test is valid","details":[]}`))
	}))
	defer server.Close()

	testData := sslValidateTestData()
	testData["test_id"] = 1655

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceValidateSslCheckV2().Schema, testData)

	diags := dataSourceValidateSslCheckV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceValidateSslCheckV2Read() diagnostics = %#v", diags)
	}
	if !d.Get("valid").(bool) {
		t.Fatalf("valid = %#v, want true", d.Get("valid"))
	}
}

func TestDataSourceValidateSslCheckV2ReadUpdateStyleFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tests/ssl/1656/validate" {
			t.Fatalf("path = %q, want /tests/ssl/1656/validate", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Fatalf("method = %q, want PUT", r.Method)
		}
		_, _ = w.Write([]byte(`{"valid":false,"message":"Test is invalid","details":{"port":["is not included in the list"]}}`))
	}))
	defer server.Close()

	testData := sslValidateTestData()
	testData["test_id"] = 1656

	client := sc2.NewConfigurableClient("token", "test", sc2.NewClientArgs(30, server.URL))
	d := schema.TestResourceDataRaw(t, dataSourceValidateSslCheckV2().Schema, testData)

	diags := dataSourceValidateSslCheckV2Read(context.Background(), d, client)
	if diags.HasError() {
		t.Fatalf("dataSourceValidateSslCheckV2Read() diagnostics = %#v", diags)
	}
	if d.Get("valid").(bool) {
		t.Fatalf("valid = %#v, want false", d.Get("valid"))
	}

	fieldErrors := d.Get("field_errors").([]interface{})
	if len(fieldErrors) != 1 {
		t.Fatalf("field_errors = %#v, want one entry", fieldErrors)
	}
	entry := fieldErrors[0].(map[string]interface{})
	if entry["field"] != "port" {
		t.Fatalf("field_errors[0].field = %#v, want port", entry["field"])
	}
}
