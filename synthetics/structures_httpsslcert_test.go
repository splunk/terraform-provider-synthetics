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
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v3/syntheticsclientv2"
)

// ============================================================================
// flattenPortCheckV2Read tests (ZERO coverage)
// ============================================================================

// TestFlattenPortCheckV2ReadNilInput skipped - function assumes non-nil input

func TestFlattenPortCheckV2ReadMinimal(t *testing.T) {
	jsonData := `{"test":{"active":true,"automaticRetries":0,"locationIds":[],"customProperties":[]}}`
	var resp sc2.PortCheckV2Response
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenPortCheckV2Read(&resp)
	if len(got) != 1 {
		t.Errorf("len(got) = %d, want 1", len(got))
	}
}

func TestFlattenPortCheckV2ReadFull(t *testing.T) {
	jsonData := `{
		"test": {
			"name": "port-check",
			"type": "port",
			"protocol": "tcp",
			"host": "example.com",
			"port": 8080,
			"active": true,
			"frequency": 60,
			"automaticRetries": 1,
			"schedulingStrategy": "round_robin",
			"locationIds": ["aws-us-east-1"],
			"customProperties": [{"key":"env","value":"prod"}]
		}
	}`
	var resp sc2.PortCheckV2Response
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenPortCheckV2Read(&resp)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}

	gotMap := got[0].(map[string]interface{})
	if gotMap["name"] != "port-check" {
		t.Errorf("name = %v, want port-check", gotMap["name"])
	}
	if gotMap["host"] != "example.com" {
		t.Errorf("host = %v, want example.com", gotMap["host"])
	}
	if gotMap["port"] != 8080 {
		t.Errorf("port = %v, want 8080", gotMap["port"])
	}
}

// ============================================================================
// flattenPortCheckV2Data tests (ZERO coverage)
// ============================================================================

// TestFlattenPortCheckV2DataNilInput skipped - function assumes non-nil input

func TestFlattenPortCheckV2DataWithID(t *testing.T) {
	jsonData := `{
		"test": {
			"id": 123,
			"name": "port-check",
			"active": true,
			"host": "localhost",
			"port": 9000,
			"protocol": "udp",
			"locationIds": [],
			"customProperties": [],
			"automaticRetries": 0
		}
	}`
	var resp sc2.PortCheckV2Response
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenPortCheckV2Data(&resp)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}

	gotMap := got[0].(map[string]interface{})
	if id, ok := gotMap["id"]; !ok || id != 123 {
		t.Errorf("id = %v, want 123", id)
	}
}

// ============================================================================
// flattenSslCheckV2Data tests (ZERO coverage)
// ============================================================================

// TestFlattenSslCheckV2DataNilInput skipped - function assumes non-nil input

func TestFlattenSslCheckV2DataMinimal(t *testing.T) {
	jsonData := `{
		"test": {
			"name": "ssl-check",
			"host": "example.com",
			"port": 443,
			"active": true,
			"allowSelfSigned": false,
			"allowUntrustedRoot": false,
			"automaticRetries": 0,
			"locationIds": [],
			"customProperties": [],
			"validations": []
		}
	}`
	var resp sc2.SslCheckV2Response
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenSslCheckV2Data(&resp)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
}

// ============================================================================
// flattenCaCertificatesV2Data tests (ZERO coverage)
// ============================================================================

func TestFlattenCaCertificatesV2DataEmpty(t *testing.T) {
	// flattenCaCertificatesV2Data returns empty slice for empty input
	got := flattenCaCertificatesV2Data([]sc2.CaCertificate{})
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}

func TestFlattenCaCertificatesV2DataSingle(t *testing.T) {
	certs := []sc2.CaCertificate{
		{
			ID:   1,
			Name: "root-ca",
		},
	}
	got := flattenCaCertificatesV2Data(certs)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}

	gotMap := got[0].(map[string]interface{})
	if gotMap["name"] != "root-ca" {
		t.Errorf("name = %v, want root-ca", gotMap["name"])
	}
}

func TestFlattenCaCertificatesV2DataMultiple(t *testing.T) {
	certs := []sc2.CaCertificate{
		{ID: 1, Name: "ca-one"},
		{ID: 2, Name: "ca-two"},
		{ID: 3, Name: "ca-three"},
	}
	got := flattenCaCertificatesV2Data(certs)
	if len(got) != 3 {
		t.Errorf("len(got) = %d, want 3", len(got))
	}
}

// ============================================================================
// caCertificateContentFromState tests (ZERO coverage)
// ============================================================================

func TestCaCertificateContentFromStateEmpty(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCaCertificateV2().Schema, map[string]interface{}{})
	got := caCertificateContentFromState(d)
	if got != "" {
		t.Errorf("got = %q, want empty", got)
	}
}

func TestCaCertificateContentFromStateWithContent(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCaCertificateV2().Schema, map[string]interface{}{
		"ca_certificate": []interface{}{
			map[string]interface{}{
				"content": "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
			},
		},
	})
	got := caCertificateContentFromState(d)
	if got != "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----" {
		t.Errorf("got = %q, want certificate content", got)
	}
}

// ============================================================================
// buildClientCertificateV2Data tests (ZERO coverage)
// ============================================================================

func TestBuildClientCertificateV2DataEmpty(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceClientCertificateV2().Schema, map[string]interface{}{})
	got := buildClientCertificateV2Data(d)
	if got == nil {
		t.Errorf("got = nil, want non-nil")
	}
}

func TestBuildClientCertificateV2DataFull(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceClientCertificateV2().Schema, map[string]interface{}{
		"client_certificate": []interface{}{
			map[string]interface{}{
				"name":        "test-cert",
				"description": "Test certificate",
				"domain":      "example.com",
				"public_key": []interface{}{
					map[string]interface{}{
						"content":        "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----",
						"filename":       "public.key",
						"file_extension": "pem",
					},
				},
				"private_key": []interface{}{
					map[string]interface{}{
						"content":        "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
						"filename":       "private.key",
						"file_extension": "pem",
						"password":       "secret123",
					},
				},
			},
		},
	})
	got := buildClientCertificateV2Data(d)
	if got == nil {
		t.Fatalf("got = nil, want non-nil")
	}
	if got.Certificate.Name != "test-cert" {
		t.Errorf("name = %q, want test-cert", got.Certificate.Name)
	}
}

// ============================================================================
// buildClientCertificateV2UpdateData tests (ZERO coverage)
// ============================================================================

func TestBuildClientCertificateV2UpdateDataValid(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceClientCertificateV2().Schema, map[string]interface{}{
		"client_certificate": []interface{}{
			map[string]interface{}{
				"private_key": []interface{}{
					map[string]interface{}{
						"content":  "new-key",
						"password": "",
					},
				},
			},
		},
	})

	got, err := buildClientCertificateV2UpdateData(d)
	if err != nil {
		t.Errorf("buildClientCertificateV2UpdateData() error = %v", err)
	}
	if got == nil {
		t.Errorf("got = nil, want non-nil")
	}
}

// ============================================================================
// flattenClientCertificateMetadata tests (ZERO coverage)
// ============================================================================

func TestFlattenClientCertificateMetadataZero(t *testing.T) {
	jsonData := `{
		"id": 0,
		"name": "",
		"description": "",
		"domain": ""
	}`
	var cert sc2.ClientCertificate
	if err := json.Unmarshal([]byte(jsonData), &cert); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenClientCertificateMetadata(cert)
	if len(got) == 0 {
		t.Errorf("len(got) = 0, want non-empty map")
	}
}

func TestFlattenClientCertificateMetadataFull(t *testing.T) {
	jsonData := `{
		"id": 42,
		"name": "my-cert",
		"description": "My certificate",
		"domain": "example.com",
		"createdBy": "user1",
		"updatedBy": "user2",
		"createdAt": "2023-01-01T12:00:00Z",
		"updatedAt": "2023-02-01T12:00:00Z",
		"expiresAt": "2024-01-01T12:00:00Z"
	}`
	var cert sc2.ClientCertificate
	if err := json.Unmarshal([]byte(jsonData), &cert); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenClientCertificateMetadata(cert)
	if id, ok := got["id"].(int); !ok || id != 42 {
		t.Errorf("id = %v, want 42", got["id"])
	}
}

// ============================================================================
// buildPortCheckV2Data tests (ZERO coverage)
// ============================================================================

func TestBuildPortCheckV2DataEmpty(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourcePortCheckV2().Schema, map[string]interface{}{})
	got := buildPortCheckV2Data(d)
	// Should not panic
	_ = got
}

func TestBuildPortCheckV2DataFull(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourcePortCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                "port-test",
				"type":                "port",
				"host":                "example.com",
				"port":                8080,
				"protocol":            "tcp",
				"active":              true,
				"frequency":           60,
				"automatic_retries":   2,
				"scheduling_strategy": "round_robin",
				"location_ids":        []interface{}{"aws-us-east-1"},
				"custom_properties":   []interface{}{},
			},
		},
	})
	got := buildPortCheckV2Data(d)
	if got.Test.Name != "port-test" {
		t.Errorf("name = %q, want port-test", got.Test.Name)
	}
}

// ============================================================================
// flattenHttpV2Read partial coverage tests
// ============================================================================

func TestFlattenHttpV2ReadNilPort(t *testing.T) {
	jsonData := `{
		"test": {
			"active": true,
			"automaticRetries": 0,
			"verifyCertificates": false,
			"headers": [],
			"validations": [],
			"customProperties": [],
			"locationIds": []
		}
	}`
	var resp sc2.HttpCheckV2ResponseWithNullablePort
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenHttpV2Read(&resp)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
}

func TestFlattenHttpV2ReadWithPort(t *testing.T) {
	jsonData := `{
		"test": {
			"name": "http-check",
			"url": "http://example.com",
			"port": 8080,
			"active": true,
			"automaticRetries": 0,
			"headers": [],
			"validations": [],
			"customProperties": [],
			"locationIds": []
		}
	}`
	var resp sc2.HttpCheckV2ResponseWithNullablePort
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenHttpV2Read(&resp)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}

	gotMap := got[0].(map[string]interface{})
	if port, ok := gotMap["port"]; !ok || port != 8080 {
		t.Errorf("port = %v, want 8080", port)
	}
}

// ============================================================================
// flattenHttpV2Data partial coverage tests
// ============================================================================

func TestFlattenHttpV2DataNoUserAgent(t *testing.T) {
	jsonData := `{
		"test": {
			"active": true,
			"automaticRetries": 0,
			"headers": [],
			"validations": [],
			"customProperties": [],
			"locationIds": []
		}
	}`
	var resp sc2.HttpCheckV2ResponseWithNullablePort
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	got := flattenHttpV2Data(&resp)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}
}

// ============================================================================
// flattenSslCheckV2Read partial coverage tests
// ============================================================================

func TestFlattenSslCheckV2ReadWithServerName(t *testing.T) {
	sn := "secure.example.com"
	jsonData := `{
		"test": {
			"name": "ssl-check",
			"host": "example.com",
			"port": 443,
			"serverName": "secure.example.com",
			"active": true,
			"allowSelfSigned": false,
			"allowUntrustedRoot": false,
			"automaticRetries": 0,
			"locationIds": [],
			"customProperties": [],
			"validations": []
		}
	}`
	var resp sc2.SslCheckV2Response
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	resp.Test.ServerName = &sn

	got := flattenSslCheckV2Read(&resp)
	if len(got) != 1 {
		t.Fatalf("len(got) = %d, want 1", len(got))
	}

	gotMap := got[0].(map[string]interface{})
	if sn, ok := gotMap["server_name"]; !ok || sn != "secure.example.com" {
		t.Errorf("server_name = %v, want secure.example.com", sn)
	}
}

// ============================================================================
// flattenCaCertificateData partial coverage tests
// ============================================================================

func TestFlattenCaCertificateDataEmpty(t *testing.T) {
	cert := sc2.CaCertificate{}
	got := flattenCaCertificateData(cert, "")
	// Empty certificate returns empty map
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0 for empty certificate", len(got))
	}
}

func TestFlattenCaCertificateDataWithState(t *testing.T) {
	cert := sc2.CaCertificate{
		ID:   123,
		Name: "prod-ca",
	}
	existing := "-----BEGIN CERTIFICATE-----\nMIIC..."
	got := flattenCaCertificateData(cert, existing)
	if len(got) == 0 {
		t.Fatalf("len(got) = 0, want non-empty map")
	}
	if got["id"] != 123 {
		t.Errorf("id = %v, want 123", got["id"])
	}
}

// ============================================================================
// caCertificateContentForState partial coverage tests
// ============================================================================

func TestCaCertificateContentForStateAPIContent(t *testing.T) {
	got := caCertificateContentForState("new-content", "old-content")
	if got != "new-content" {
		t.Errorf("got = %q, want new-content", got)
	}
}

func TestCaCertificateContentForStateEmpty(t *testing.T) {
	got := caCertificateContentForState("", "existing-content")
	if got != "existing-content" {
		t.Errorf("got = %q, want existing-content", got)
	}
}

func TestCaCertificateContentForStateRedacted(t *testing.T) {
	got := caCertificateContentForState("<REDACTED>", "existing-content")
	if got != "existing-content" {
		t.Errorf("got = %q, want existing-content", got)
	}
}

// ============================================================================
// httpV2PortFromRawValue partial coverage tests
// ============================================================================

func TestHttpV2PortFromRawValueNull(t *testing.T) {
	val := cty.NullVal(cty.Object(map[string]cty.Type{}))
	_, ok := httpV2PortFromRawValue(val)
	if ok {
		t.Errorf("ok = true, want false for null value")
	}
}

func TestHttpV2PortFromRawValueUnknown(t *testing.T) {
	val := cty.UnknownVal(cty.Object(map[string]cty.Type{}))
	_, ok := httpV2PortFromRawValue(val)
	if ok {
		t.Errorf("ok = true, want false for unknown value")
	}
}

// ============================================================================
// buildHttpV2CertificateIDForUpdate partial coverage tests
// ============================================================================

func TestBuildHttpV2CertificateIDForUpdateNewID(t *testing.T) {
	oldTest := map[string]interface{}{"certificate_id": 0}
	newTest := map[string]interface{}{"certificate_id": 123}
	got := buildHttpV2CertificateIDForUpdate(oldTest, newTest)
	if got == nil {
		t.Errorf("got = nil, want non-nil")
	}
	if got != nil && (got.Value == nil || *got.Value != 123) {
		t.Errorf("value = %v, want 123", got.Value)
	}
}

func TestBuildHttpV2CertificateIDForUpdateRemoved(t *testing.T) {
	oldTest := map[string]interface{}{"certificate_id": 456}
	newTest := map[string]interface{}{"certificate_id": 0}
	got := buildHttpV2CertificateIDForUpdate(oldTest, newTest)
	if got == nil {
		t.Errorf("got = nil, want non-nil")
	}
	if got != nil && got.Value != nil {
		t.Errorf("value = %v, want nil", got.Value)
	}
}

func TestBuildHttpV2CertificateIDForUpdateNoChange(t *testing.T) {
	oldTest := map[string]interface{}{"certificate_id": 0}
	newTest := map[string]interface{}{"certificate_id": 0}
	got := buildHttpV2CertificateIDForUpdate(oldTest, newTest)
	if got != nil {
		t.Errorf("got = %v, want nil", got)
	}
}

// ============================================================================
// buildSslCheckV2UpdateData partial coverage tests
// ============================================================================

func TestBuildSslCheckV2UpdateDataClearFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSslCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{
			map[string]interface{}{
				"name":                 "ssl",
				"active":               true,
				"automatic_retries":    0,
				"frequency":            60,
				"scheduling_strategy":  "round_robin",
				"host":                 "example.com",
				"port":                 443,
				"allow_self_signed":    false,
				"allow_untrusted_root": false,
				"server_name":          "",
				"ca_certificate_id":    0,
				"location_ids":         []interface{}{},
				"custom_properties":    []interface{}{},
				"validations":          []interface{}{},
			},
		},
	})

	got := buildSslCheckV2UpdateData(d)
	if got.Test.ServerName == nil {
		t.Errorf("ServerName = nil, want explicit nullable")
	}
}

// ============================================================================
// buildCaCertificateV2Data partial coverage tests
// ============================================================================

func TestBuildCaCertificateV2DataNoBlock(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCaCertificateV2().Schema, map[string]interface{}{})
	_, err := buildCaCertificateV2Data(d)
	if err == nil {
		t.Errorf("expected error for missing ca_certificate block")
	}
}

func TestBuildCaCertificateV2DataValid(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCaCertificateV2().Schema, map[string]interface{}{
		"ca_certificate": []interface{}{
			map[string]interface{}{
				"name":           "test-ca",
				"description":    "Test CA",
				"content":        "-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----",
				"file_extension": "pem",
				"filename":       "test-ca.pem",
			},
		},
	})

	got, err := buildCaCertificateV2Data(d)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got.CaCert.Name != "test-ca" {
		t.Errorf("name = %q, want test-ca", got.CaCert.Name)
	}
}

// ============================================================================
// buildCaCertificateV2UpdateData partial coverage tests
// ============================================================================

func TestBuildCaCertificateV2UpdateDataEmpty(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceCaCertificateV2().Schema, map[string]interface{}{
		"ca_certificate": []interface{}{},
	})
	got := buildCaCertificateV2UpdateData(d)
	// Should not panic
	_ = got
}

// ============================================================================
// caCertificateStringField tests
// ============================================================================

func TestCaCertificateStringFieldValid(t *testing.T) {
	input := map[string]interface{}{"name": "test-ca"}
	got := caCertificateStringField(input, "name")
	if got != "test-ca" {
		t.Errorf("got = %q, want test-ca", got)
	}
}

func TestCaCertificateStringFieldMissing(t *testing.T) {
	input := map[string]interface{}{}
	got := caCertificateStringField(input, "missing")
	if got != "" {
		t.Errorf("got = %q, want empty", got)
	}
}

// ============================================================================
// sslStringField tests
// ============================================================================

func TestSslStringFieldValid(t *testing.T) {
	input := map[string]interface{}{"host": "example.com"}
	got := sslStringField(input, "host")
	if got != "example.com" {
		t.Errorf("got = %q, want example.com", got)
	}
}

func TestSslStringFieldWrongType(t *testing.T) {
	input := map[string]interface{}{"port": 443}
	got := sslStringField(input, "port")
	if got != "" {
		t.Errorf("got = %q, want empty", got)
	}
}

// ============================================================================
// sslIntField tests
// ============================================================================

func TestSslIntFieldValid(t *testing.T) {
	input := map[string]interface{}{"port": 443}
	got := sslIntField(input, "port")
	if got != 443 {
		t.Errorf("got = %d, want 443", got)
	}
}

func TestSslIntFieldWrongType(t *testing.T) {
	input := map[string]interface{}{"host": "example.com"}
	got := sslIntField(input, "host")
	if got != 0 {
		t.Errorf("got = %d, want 0", got)
	}
}

// ============================================================================
// sslInterfaceListField tests
// ============================================================================

func TestSslInterfaceListFieldValid(t *testing.T) {
	input := map[string]interface{}{
		"validations": []interface{}{
			map[string]interface{}{"name": "val1"},
			map[string]interface{}{"name": "val2"},
		},
	}
	got := sslInterfaceListField(input, "validations")
	if len(got) != 2 {
		t.Errorf("len(got) = %d, want 2", len(got))
	}
}

func TestSslInterfaceListFieldMissing(t *testing.T) {
	input := map[string]interface{}{}
	got := sslInterfaceListField(input, "validations")
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}

// ============================================================================
// flattenClientCertificateV2Read partial coverage tests
// ============================================================================

func TestFlattenClientCertificateV2ReadWithRedaction(t *testing.T) {
	jsonData := `{
		"id": 1,
		"name": "cert",
		"publicKey": {
			"content": "",
			"id": 10
		},
		"privateKey": {
			"content": "<REDACTED>",
			"password": "",
			"id": 20
		}
	}`
	var cert sc2.ClientCertificate
	if err := json.Unmarshal([]byte(jsonData), &cert); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	existing := map[string]interface{}{
		"public_key": []interface{}{
			map[string]interface{}{
				"content": "public-key-content",
			},
		},
		"private_key": []interface{}{
			map[string]interface{}{
				"content":  "private-key-content",
				"password": "secret123",
			},
		},
	}

	got := flattenClientCertificateV2Read(cert, existing)
	if len(got) == 0 {
		t.Fatalf("len(got) = 0, want non-empty map")
	}
}

// ============================================================================
// stateOrAPISecret partial coverage tests
// ============================================================================

func TestStateOrAPISecretAPIValue(t *testing.T) {
	existing := map[string]interface{}{"content": "old"}
	got := stateOrAPISecret(existing, "new", "content")
	if got != "new" {
		t.Errorf("got = %q, want new", got)
	}
}

func TestStateOrAPISecretEmpty(t *testing.T) {
	existing := map[string]interface{}{"content": "existing"}
	got := stateOrAPISecret(existing, "", "content")
	if got != "existing" {
		t.Errorf("got = %q, want existing", got)
	}
}

func TestStateOrAPISecretRedacted(t *testing.T) {
	existing := map[string]interface{}{"password": "old-password"}
	got := stateOrAPISecret(existing, "<REDACTED>", "password")
	if got != "old-password" {
		t.Errorf("got = %q, want old-password", got)
	}
}

// ============================================================================
// firstMapFromList partial coverage tests
// ============================================================================

func TestFirstMapFromListEmpty(t *testing.T) {
	value := []interface{}{}
	got := firstMapFromList(value)
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}

func TestFirstMapFromListWithMap(t *testing.T) {
	value := []interface{}{
		map[string]interface{}{"name": "test", "id": 42},
	}
	got := firstMapFromList(value)
	if got["name"] != "test" || got["id"] != 42 {
		t.Errorf("got = %v, want map with name=test, id=42", got)
	}
}

func TestFirstMapFromListNil(t *testing.T) {
	got := firstMapFromList(nil)
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}

// ============================================================================
// intSliceFromInterface partial coverage tests
// ============================================================================

func TestIntSliceFromInterfaceEmpty(t *testing.T) {
	got := intSliceFromInterface([]interface{}{})
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}

func TestIntSliceFromInterfaceValid(t *testing.T) {
	got := intSliceFromInterface([]interface{}{1, 2, 3, 4, 5})
	if len(got) != 5 {
		t.Fatalf("len(got) = %d, want 5", len(got))
	}
	for i, v := range got {
		if v != i+1 {
			t.Errorf("element %d = %d, want %d", i, v, i+1)
		}
	}
}

func TestIntSliceFromInterfaceMixed(t *testing.T) {
	got := intSliceFromInterface([]interface{}{1, "two", 3, nil, 5})
	if len(got) != 3 {
		t.Errorf("len(got) = %d, want 3 (only ints extracted)", len(got))
	}
}

func TestIntSliceFromInterfaceNil(t *testing.T) {
	got := intSliceFromInterface(nil)
	if len(got) != 0 {
		t.Errorf("len(got) = %d, want 0", len(got))
	}
}
