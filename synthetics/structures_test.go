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
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func TestBuildSelectorsFromStep_multipleSelectors(t *testing.T) {
	step := map[string]interface{}{
		"name": "Click submit",
		"selectors": []interface{}{
			map[string]interface{}{"type": "css", "value": ".primary"},
			map[string]interface{}{"type": "id", "value": "submit-btn"},
		},
	}

	got, err := buildSelectorsFromStep(step)
	if err != nil {
		t.Fatalf("buildSelectorsFromStep() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(selectors) = %d, want 2", len(got))
	}
	if got[0].Type != "css" || got[0].Value != ".primary" {
		t.Fatalf("first selector = %+v, want css/.primary", got[0])
	}
	if got[1].Type != "id" || got[1].Value != "submit-btn" {
		t.Fatalf("second selector = %+v, want id/submit-btn", got[1])
	}
}

func TestBuildSelectorsFromStep_legacyFields(t *testing.T) {
	step := map[string]interface{}{
		"selector_type": "id",
		"selector":      "checkout-btn",
	}

	got, err := buildSelectorsFromStep(step)
	if err != nil {
		t.Fatalf("buildSelectorsFromStep() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(selectors) = %d, want 1", len(got))
	}
	if got[0].Type != "id" || got[0].Value != "checkout-btn" {
		t.Fatalf("selector = %+v", got[0])
	}
}

func TestDropStaleStepSelectors_legacyUpdateWins(t *testing.T) {
	step := map[string]interface{}{
		"name":          "06 Select-Val-Val",
		"selector_type": "id",
		"selector":      "valz",
		"selectors": []interface{}{
			map[string]interface{}{"type": "id", "value": "beep"},
		},
	}
	dropStaleStepSelectors(step)
	got, err := buildSelectorsFromStep(step)
	if err != nil {
		t.Fatalf("buildSelectorsFromStep() error = %v", err)
	}
	if len(got) != 1 || got[0].Value != "valz" {
		t.Fatalf("selectors = %+v, want id/valz after dropping stale state", got)
	}
}

func TestBuildSelectorsFromStep_prefersSelectorsList(t *testing.T) {
	step := map[string]interface{}{
		"selector_type": "id",
		"selector":      "ignored",
		"selectors": []interface{}{
			map[string]interface{}{"type": "css", "value": ".btn"},
		},
	}

	got, err := buildSelectorsFromStep(step)
	if err != nil {
		t.Fatalf("buildSelectorsFromStep() error = %v", err)
	}
	if len(got) != 1 || got[0].Value != ".btn" {
		t.Fatalf("selectors = %+v, want only .btn from selectors list", got)
	}
}

func TestFlattenStepsData_singleSelectorUsesSelectorsBlock(t *testing.T) {
	steps := []sc2.StepsV2{{
		Name:      "click",
		Type:      "click_element",
		Selectors: []sc2.Selector{{Type: "id", Value: "#order"}},
	}}

	got := flattenStepsData(&steps)
	if len(got) != 1 {
		t.Fatalf("len(steps) = %d, want 1", len(got))
	}
	step := got[0].(map[string]interface{})
	selectors, ok := step["selectors"].([]interface{})
	if !ok || len(selectors) != 1 {
		t.Fatalf("selectors = %#v, want one block", step["selectors"])
	}
	sel := selectors[0].(map[string]interface{})
	if sel["type"] != "id" || sel["value"] != "#order" {
		t.Fatalf("selector block = %#v, want id/#order", sel)
	}
	if _, ok := step["selector"]; ok {
		t.Fatal("expected no legacy selector field for single API selector")
	}
	if _, ok := step["selector_type"]; ok {
		t.Fatal("expected no legacy selector_type for single API selector")
	}
}

func TestStepSelectorInputsEquivalent_legacyAndSelectorsBlock(t *testing.T) {
	oldIn := stepSelectorInput{
		selectorType: "id",
		selector:     "#order",
	}
	newIn := stepSelectorInput{
		selectors: []sc2.Selector{{Type: "id", Value: "#order"}},
	}
	if !stepSelectorInputsEquivalent(oldIn, newIn) {
		t.Fatal("expected legacy and single selectors block to be equivalent")
	}
	if !migratingFromLegacyToSelectors(oldIn, newIn) {
		t.Fatal("expected legacy state to selectors config to be a migration")
	}
}

func TestMigratingFromLegacyToSelectors_afterApplyNotMigration(t *testing.T) {
	in := stepSelectorInput{
		selectors: []sc2.Selector{{Type: "id", Value: "#order"}},
	}
	if migratingFromLegacyToSelectors(in, in) {
		t.Fatal("selectors in state and config should not be treated as migration")
	}
}

func TestFlattenStepsData_multipleSelectorsExposed(t *testing.T) {
	steps := []sc2.StepsV2{{
		Name: "click",
		Type: "click_element",
		Selectors: []sc2.Selector{
			{Type: "css", Value: ".primary"},
			{Type: "id", Value: "submit-btn"},
		},
	}}

	got := flattenStepsData(&steps)
	step := got[0].(map[string]interface{})
	selectors, ok := step["selectors"].([]interface{})
	if !ok || len(selectors) != 2 {
		t.Fatalf("selectors = %#v, want 2 blocks", step["selectors"])
	}
	if _, ok := step["selector"]; ok {
		t.Fatal("expected no legacy selector fields for multiple selectors")
	}
	if _, ok := step["selector_type"]; ok {
		t.Fatal("expected no legacy selector_type for multiple selectors")
	}
}

func TestBuildSelectorsFromStep_tooManySelectors(t *testing.T) {
	selectors := make([]interface{}, browserCheckV2MaxSelectors+1)
	for i := range selectors {
		selectors[i] = map[string]interface{}{"type": "id", "value": "x"}
	}
	step := map[string]interface{}{"selectors": selectors}

	_, err := buildSelectorsFromStep(step)
	if err == nil {
		t.Fatal("expected error for too many selectors")
	}
}

func TestBuildExcludedFilesV2Data(t *testing.T) {
	input := testExcludedFileSet(
		map[string]interface{}{"type": "google_analytics", "regex": ""},
		map[string]interface{}{"type": "future_api_owned_type", "regex": ""},
		map[string]interface{}{"type": "custom", "regex": "cdn\\.example\\.com"},
		map[string]interface{}{"type": "all_except", "regex": "assets\\.example\\.com"},
	)

	got, err := buildExcludedFilesV2Data(input)
	if err != nil {
		t.Fatalf("buildExcludedFilesV2Data() error = %v", err)
	}

	want := []sc2.ExcludedFile{
		{Type: "google_analytics"},
		{Type: "future_api_owned_type"},
		{Type: "custom", Regex: "cdn\\.example\\.com"},
		{Type: "all_except", Regex: "assets\\.example\\.com"},
	}
	if !sameExcludedFiles(got, want) {
		t.Fatalf("excluded files = %#v, want %#v", got, want)
	}
}

func TestBuildExcludedFilesV2DataNilAndEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input *schema.Set
	}{
		{name: "nil"},
		{name: "empty", input: testExcludedFileSet()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildExcludedFilesV2Data(tt.input)
			if err != nil {
				t.Fatalf("buildExcludedFilesV2Data() error = %v", err)
			}
			if got == nil {
				t.Fatal("buildExcludedFilesV2Data() returned nil slice")
			}
			if len(got) != 0 {
				t.Fatalf("buildExcludedFilesV2Data() len = %d, want 0", len(got))
			}
		})
	}
}

func TestBuildExcludedFilesV2DataValidation(t *testing.T) {
	tests := []struct {
		name    string
		item    map[string]interface{}
		wantErr string
	}{
		{
			name:    "empty type",
			item:    map[string]interface{}{"type": "", "regex": ""},
			wantErr: "type must not be empty",
		},
		{
			name:    "custom missing regex",
			item:    map[string]interface{}{"type": "custom", "regex": ""},
			wantErr: "regex must be set when type is custom",
		},
		{
			name:    "all_except missing regex",
			item:    map[string]interface{}{"type": "all_except", "regex": "  "},
			wantErr: "regex must be set when type is all_except",
		},
		{
			name:    "invalid regex",
			item:    map[string]interface{}{"type": "custom", "regex": "[a-z"},
			wantErr: "regex is not valid RE2 syntax",
		},
		{
			name:    "predefined rejects regex",
			item:    map[string]interface{}{"type": "google_analytics", "regex": "google"},
			wantErr: "regex is only supported when type is custom or all_except",
		},
		{
			name:    "api owned type rejects whitespace regex",
			item:    map[string]interface{}{"type": "future_api_owned_type", "regex": " "},
			wantErr: "regex is only supported when type is custom or all_except",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildExcludedFilesV2Data(testExcludedFileSet(tt.item))
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestFlattenExcludedFilesV2Data(t *testing.T) {
	got := flattenExcludedFilesV2Data([]sc2.ExcludedFile{
		{Type: "google_analytics"},
		{Type: "custom", Regex: "cdn\\.example\\.com"},
	})

	want := []interface{}{
		map[string]interface{}{"type": "google_analytics"},
		map[string]interface{}{"type": "custom", "regex": "cdn\\.example\\.com"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("flattened excluded files = %#v, want %#v", got, want)
	}
}

func TestBuildAdvancedSettingsDataSendsEmptyExcludedFiles(t *testing.T) {
	settings := schema.NewSet(schema.HashResource(browserCheckV2AdvancedSettingsResource(false)), []interface{}{
		map[string]interface{}{
			"user_agent":                  "",
			"verify_certificates":         true,
			"collect_interactive_metrics": false,
			"authentication":              schema.NewSet(schema.HashString, nil),
			"chrome_flags":                schema.NewSet(schema.HashString, nil),
			"cookies":                     schema.NewSet(schema.HashString, nil),
			"headers":                     schema.NewSet(schema.HashString, nil),
			"host_overrides":              schema.NewSet(schema.HashString, nil),
			"excluded_files":              schema.NewSet(schema.HashString, nil),
		},
	})

	got, err := buildAdvancedSettingsData(settings)
	if err != nil {
		t.Fatalf("buildAdvancedSettingsData() error = %v", err)
	}
	body, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if !strings.Contains(string(body), `"excludedFiles":[]`) {
		t.Fatalf("advanced settings JSON = %s, want excludedFiles empty array", string(body))
	}
}

func TestBrowserV2AdvancedSettingsSensitiveFields(t *testing.T) {
	resourceAdvancedSettings := resourceBrowserCheckV2().Schema["test"].Elem.(*schema.Resource).Schema["advanced_settings"].Elem.(*schema.Resource)
	assertBrowserV2AdvancedSettingsSensitiveFields(t, "resource", resourceAdvancedSettings)

	computedAdvancedSettings := browserCheckV2AdvancedSettingsResource(true)
	assertBrowserV2AdvancedSettingsSensitiveFields(t, "computed helper", computedAdvancedSettings)

	dataSourceAdvancedSettings := dataSourceBrowserCheckV2().Schema["test"].Elem.(*schema.Resource).Schema["advanced_settings"].Elem.(*schema.Resource)
	assertBrowserV2AdvancedSettingsSensitiveFields(t, "data source", dataSourceAdvancedSettings)
}

func assertBrowserV2AdvancedSettingsSensitiveFields(t *testing.T, name string, advancedSettings *schema.Resource) {
	t.Helper()

	assertNestedFieldSensitive(t, name, advancedSettings, "authentication", "password", true)
	assertNestedFieldSensitive(t, name, advancedSettings, "cookies", "value", true)
	assertNestedFieldSensitive(t, name, advancedSettings, "headers", "value", true)
	assertNestedFieldSensitive(t, name, advancedSettings, "chrome_flags", "value", false)
}

func assertNestedFieldSensitive(t *testing.T, name string, advancedSettings *schema.Resource, setName string, fieldName string, want bool) {
	t.Helper()

	setSchema, ok := advancedSettings.Schema[setName]
	if !ok {
		t.Fatalf("%s advanced_settings.%s schema missing", name, setName)
	}
	nestedResource, ok := setSchema.Elem.(*schema.Resource)
	if !ok {
		t.Fatalf("%s advanced_settings.%s Elem = %T, want *schema.Resource", name, setName, setSchema.Elem)
	}
	fieldSchema, ok := nestedResource.Schema[fieldName]
	if !ok {
		t.Fatalf("%s advanced_settings.%s.%s schema missing", name, setName, fieldName)
	}
	if fieldSchema.Sensitive != want {
		t.Fatalf("%s advanced_settings.%s.%s Sensitive = %t, want %t", name, setName, fieldName, fieldSchema.Sensitive, want)
	}
}

func testExcludedFileSet(items ...map[string]interface{}) *schema.Set {
	values := make([]interface{}, 0, len(items))
	for _, item := range items {
		values = append(values, item)
	}
	return schema.NewSet(schema.HashResource(browserCheckV2ExcludedFileResource(false)), values)
}

func sameExcludedFiles(got []sc2.ExcludedFile, want []sc2.ExcludedFile) bool {
	if len(got) != len(want) {
		return false
	}
	remaining := append([]sc2.ExcludedFile(nil), want...)
	for _, gotItem := range got {
		found := false
		for i, wantItem := range remaining {
			if gotItem == wantItem {
				remaining = append(remaining[:i], remaining[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return len(remaining) == 0
}

func TestBuildSslCheckV2DataSetsNullableFieldsValidationsAndCustomProperties(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSslCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                 "ssl-create",
				"type":                 "ssl",
				"active":               false,
				"frequency":            5,
				"scheduling_strategy":  "round_robin",
				"location_ids":         []interface{}{"aws-us-east-1"},
				"automatic_retries":    1,
				"host":                 "example.com",
				"port":                 443,
				"server_name":          "tls.example.com",
				"allow_self_signed":    true,
				"allow_untrusted_root": true,
				"ca_certificate_id":    123,
				"validations": []interface{}{
					map[string]interface{}{
						"name":       "response code",
						"type":       "assert_numeric",
						"actual":     "{{response.code}}",
						"comparator": "equals",
						"expected":   "200",
					},
				},
				"custom_properties": []interface{}{
					map[string]interface{}{
						"key":   "owner",
						"value": "synthetics",
					},
				},
			},
		},
	})

	got := buildSslCheckV2Data(d)

	if got.Test.ServerName == nil || *got.Test.ServerName != "tls.example.com" {
		t.Fatalf("ServerName = %#v, want tls.example.com", got.Test.ServerName)
	}
	if got.Test.CaCertificateID == nil || *got.Test.CaCertificateID != 123 {
		t.Fatalf("CaCertificateID = %#v, want 123", got.Test.CaCertificateID)
	}
	if len(got.Test.Validations) != 1 {
		t.Fatalf("Validations = %#v, want one validation", got.Test.Validations)
	}
	validation := got.Test.Validations[0]
	if validation.Name != "response code" || validation.Type != "assert_numeric" || validation.Actual != "{{response.code}}" || validation.Comparator != "equals" || validation.Expected != "200" {
		t.Fatalf("Validation = %#v, want response code assertion", validation)
	}
	if validation.Extractor != "" || validation.Source != "" || validation.Variable != "" || validation.Value != "" || validation.Code != "" {
		t.Fatalf("Validation includes unsupported SSL fields: %#v", validation)
	}
	if len(got.Test.Customproperties) != 1 || got.Test.Customproperties[0].Key != "owner" || got.Test.Customproperties[0].Value != "synthetics" {
		t.Fatalf("Customproperties = %#v, want owner=synthetics", got.Test.Customproperties)
	}
}

func TestBuildValidationsDataToleratesMissingGenericValidationFields(t *testing.T) {
	got := buildValidationsData([]interface{}{
		map[string]interface{}{
			"name":       "response code",
			"type":       "assert_numeric",
			"actual":     "{{response.code}}",
			"comparator": "equals",
			"expected":   "200",
		},
	})

	if len(got) != 1 {
		t.Fatalf("Validations = %#v, want one validation", got)
	}
	if got[0].Name != "response code" || got[0].Type != "assert_numeric" || got[0].Actual != "{{response.code}}" || got[0].Comparator != "equals" || got[0].Expected != "200" {
		t.Fatalf("Validation = %#v, want response code assertion", got[0])
	}
}

func TestBuildSslCheckV2UpdateDataSendsExplicitNullsForAbsentNullableFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSslCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                 "ssl-update",
				"type":                 "ssl",
				"active":               false,
				"frequency":            5,
				"scheduling_strategy":  "round_robin",
				"location_ids":         []interface{}{"aws-us-east-1"},
				"automatic_retries":    1,
				"host":                 "example.com",
				"port":                 443,
				"allow_self_signed":    false,
				"allow_untrusted_root": false,
			},
		},
	})

	got := buildSslCheckV2UpdateData(d)

	if got.Test.ServerName == nil {
		t.Fatal("ServerName = nil, want explicit nullable null")
	}
	if got.Test.ServerName.Value != nil {
		t.Fatalf("ServerName.Value = %#v, want nil", *got.Test.ServerName.Value)
	}
	if got.Test.CaCertificateID == nil {
		t.Fatal("CaCertificateID = nil, want explicit nullable null")
	}
	if got.Test.CaCertificateID.Value != nil {
		t.Fatalf("CaCertificateID.Value = %#v, want nil", *got.Test.CaCertificateID.Value)
	}
}

func TestFlattenSslCheckV2ReadPreservesNullableServerNameAndCaCertificateID(t *testing.T) {
	serverName := "tls.example.com"
	caCertificateID := 123
	check := &sc2.SslCheckV2Response{}
	check.Test.Name = "ssl-read"
	check.Test.Type = "ssl"
	check.Test.Active = true
	check.Test.Frequency = 5
	check.Test.SchedulingStrategy = "round_robin"
	check.Test.LocationIds = []string{"aws-us-east-1"}
	check.Test.Host = "example.com"
	check.Test.Port = 443
	check.Test.ServerName = &serverName
	check.Test.CaCertificateID = &caCertificateID

	got := flattenSslCheckV2Read(check)
	if len(got) != 1 {
		t.Fatalf("len(flattened) = %d, want 1", len(got))
	}
	test := got[0].(map[string]interface{})
	if test["server_name"] != serverName {
		t.Fatalf("server_name = %#v, want %q", test["server_name"], serverName)
	}
	if test["ca_certificate_id"] != caCertificateID {
		t.Fatalf("ca_certificate_id = %#v, want %d", test["ca_certificate_id"], caCertificateID)
	}
}

func TestBuildCaCertificateV2DataRequiresContent(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCaCertificateV2().Schema, map[string]interface{}{
		"ca_certificate": []interface{}{
			map[string]interface{}{
				"name":           "ca-create",
				"description":    "test CA",
				"file_extension": "pem",
				"filename":       "test.pem",
			},
		},
	})

	_, err := buildCaCertificateV2Data(d)
	if err == nil {
		t.Fatal("expected error when CA certificate content is missing")
	}
}

func TestCaCertificateV2ContentSchemaIsRequiredAndSensitive(t *testing.T) {
	certSchema := resourceCaCertificateV2().Schema["ca_certificate"].Elem.(*schema.Resource).Schema
	contentSchema := certSchema["content"]
	if !contentSchema.Required {
		t.Fatal("content schema must be required for CA certificate resources")
	}
	if !contentSchema.Sensitive {
		t.Fatal("content schema must be sensitive for CA certificate resources")
	}
}

func TestBuildCaCertificateV2UpdateDataDoesNotSendRedactedContent(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCaCertificateV2().Schema, map[string]interface{}{
		"ca_certificate": []interface{}{
			map[string]interface{}{
				"name":           "ca-update",
				"description":    "updated CA",
				"content":        caCertificateRedactedContent,
				"file_extension": "pem",
				"filename":       "updated.pem",
			},
		},
	})

	got := buildCaCertificateV2UpdateData(d)

	if got.CaCert.Content != nil {
		t.Fatalf("Content = %#v, want nil when content is redacted", *got.CaCert.Content)
	}
	if got.CaCert.Description == nil || *got.CaCert.Description != "updated CA" {
		t.Fatalf("Description = %#v, want updated CA", got.CaCert.Description)
	}
}

func TestFlattenCaCertificateV2ReadPreservesExistingStateContent(t *testing.T) {
	check := &sc2.CaCertificateV2Response{}
	check.CaCert.ID = 123
	check.CaCert.Name = "ca-read"
	check.CaCert.Description = "test CA"
	check.CaCert.Content = caCertificateRedactedContent
	check.CaCert.FileExtension = "pem"
	check.CaCert.Filename = "test.pem"

	got := flattenCaCertificateV2Read(check, "existing-sensitive-content")
	if len(got) != 1 {
		t.Fatalf("len(flattened) = %d, want 1", len(got))
	}
	caCertificate := got[0].(map[string]interface{})
	if caCertificate["content"] != "existing-sensitive-content" {
		t.Fatalf("content = %#v, want existing state content", caCertificate["content"])
	}
}
