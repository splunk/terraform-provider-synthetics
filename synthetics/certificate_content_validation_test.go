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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testCertificateDER(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	return der
}

// testCertificatesPEM returns count CERTIFICATE blocks followed by nothing else.
func testCertificatesPEM(t *testing.T, count int) []byte {
	t.Helper()
	block := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: testCertificateDER(t)})
	return bytes.Repeat(block, count)
}

func base64Of(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func TestValidateCertificateContentSize(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr string
	}{
		{name: "a single certificate", value: base64Of(testCertificatesPEM(t, 1))},
		{name: "many certificates within the size limit", value: base64Of(testCertificatesPEM(t, 50))},
		{name: "content at the limit", value: base64Of(bytes.Repeat([]byte("a"), maxCertificateContentBytes))},
		{
			name:    "content one byte over the limit",
			value:   base64Of(bytes.Repeat([]byte("a"), maxCertificateContentBytes+1)),
			wantErr: "content is too large (257 KB); the maximum is 256 KB",
		},
		{name: "content that is not Base64", value: "not base64!"},
		{name: "a value that is not a string", value: 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings, errs := validateCertificateContentSize(tt.value, "content")
			if len(warnings) != 0 {
				t.Fatalf("unexpected warnings: %v", warnings)
			}
			if tt.wantErr == "" {
				if len(errs) != 0 {
					t.Fatalf("unexpected errors: %v", errs)
				}
				return
			}
			if len(errs) != 1 || errs[0].Error() != tt.wantErr {
				t.Fatalf("errors = %v, want %q", errs, tt.wantErr)
			}
		})
	}
}

func diagnosticsText(diags diag.Diagnostics) string {
	var parts []string
	for _, d := range diags {
		parts = append(parts, d.Summary+" "+d.Detail)
	}
	return strings.Join(parts, "\n")
}

func clientCertificateRawConfig(publicKeyContent string) map[string]interface{} {
	return map[string]interface{}{
		"client_certificate": []interface{}{map[string]interface{}{
			"name":   "chain-test",
			"domain": "example.com",
			"public_key": []interface{}{map[string]interface{}{
				"content":        publicKeyContent,
				"filename":       "chain.pem",
				"file_extension": "pem",
			}},
			"private_key": []interface{}{map[string]interface{}{
				"content":        base64Of([]byte("key")),
				"filename":       "client.key",
				"file_extension": "pem",
			}},
		}},
	}
}

func TestClientCertificateV2ResourceValidatesContentSize(t *testing.T) {
	resource := resourceClientCertificateV2()

	atLimit := terraform.NewResourceConfigRaw(clientCertificateRawConfig(base64Of(bytes.Repeat([]byte("a"), maxCertificateContentBytes))))
	if diags := resource.Validate(atLimit); diags.HasError() {
		t.Fatalf("content at the limit was rejected: %s", diagnosticsText(diags))
	}

	overLimit := terraform.NewResourceConfigRaw(clientCertificateRawConfig(base64Of(bytes.Repeat([]byte("a"), maxCertificateContentBytes+1))))
	diags := resource.Validate(overLimit)
	if !diags.HasError() || !strings.Contains(diagnosticsText(diags), "is too large (257 KB); the maximum is 256 KB") {
		t.Fatalf("content over the limit was not rejected as expected: %s", diagnosticsText(diags))
	}
}

func TestCaCertificateV2ResourceValidatesContentSize(t *testing.T) {
	resource := resourceCaCertificateV2()
	rawConfig := func(content string) map[string]interface{} {
		return map[string]interface{}{
			"ca_certificate": []interface{}{map[string]interface{}{
				"name":           "bundle-test",
				"content":        content,
				"filename":       "bundle.pem",
				"file_extension": "pem",
			}},
		}
	}

	atLimit := terraform.NewResourceConfigRaw(rawConfig(base64Of(bytes.Repeat([]byte("a"), maxCertificateContentBytes))))
	if diags := resource.Validate(atLimit); diags.HasError() {
		t.Fatalf("content at the limit was rejected: %s", diagnosticsText(diags))
	}

	overLimit := terraform.NewResourceConfigRaw(rawConfig(base64Of(bytes.Repeat([]byte("a"), maxCertificateContentBytes+1))))
	diags := resource.Validate(overLimit)
	if !diags.HasError() || !strings.Contains(diagnosticsText(diags), "is too large (257 KB); the maximum is 256 KB") {
		t.Fatalf("content over the limit was not rejected as expected: %s", diagnosticsText(diags))
	}
}

func TestCertificateExpiresAtIsDescribed(t *testing.T) {
	schemas := map[string][2]string{
		"client resource":    {clientCertificateV2ResourceSchema()["expires_at"].Description, clientCertificateExpiresAtDescription},
		"client data source": {dataSourceClientCertificateV2().Schema["expires_at"].Description, clientCertificateExpiresAtDescription},
		"ca resource":        {caCertificateV2ResourceSchema()["expires_at"].Description, caCertificateExpiresAtDescription},
		"ca data source":     {caCertificateV2DataSourceSchema()["expires_at"].Description, caCertificateExpiresAtDescription},
	}
	for name, descriptions := range schemas {
		if descriptions[0] != descriptions[1] {
			t.Errorf("%s expires_at description = %q", name, descriptions[0])
		}
	}
}
