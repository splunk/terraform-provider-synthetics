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
						"name":       "expires",
						"type":       "certificate",
						"comparator": "less_than",
						"value":      "30",
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
	if len(got.Test.Validations) != 1 || got.Test.Validations[0].Name != "expires" || got.Test.Validations[0].Value != "30" {
		t.Fatalf("Validations = %#v, want expires validation", got.Test.Validations)
	}
	if len(got.Test.Customproperties) != 1 || got.Test.Customproperties[0].Key != "owner" || got.Test.Customproperties[0].Value != "synthetics" {
		t.Fatalf("Customproperties = %#v, want owner=synthetics", got.Test.Customproperties)
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
