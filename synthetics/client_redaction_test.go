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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func TestCaCertificateRequestDetailsRedactsSensitiveValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cacerts" {
			t.Fatalf("expected /cacerts request, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"cacert":{"id":123,"name":"Terraform - CA Certificate V2"}}`))
	}))
	defer server.Close()

	client := sc2.NewConfigurableClient("secret-token", "test", sc2.NewClientArgs(30, server.URL))

	_, details, err := client.CreateCaCertificateV2(&sc2.CaCertificateV2Input{
		CaCert: sc2.CaCertificateInput{
			Name:          "Terraform - CA Certificate V2",
			Content:       "private-ca-material",
			FileExtension: "pem",
			Filename:      "ca.pem",
		},
	})
	if err != nil {
		t.Fatalf("CreateCaCertificateV2 returned error: %v", err)
	}

	if strings.Contains(details.RequestBody, "secret-token") {
		t.Fatalf("expected API token to be redacted, got %q", details.RequestBody)
	}
	if strings.Contains(details.RequestBody, "private-ca-material") {
		t.Fatalf("expected CA certificate content to be redacted, got %q", details.RequestBody)
	}
	if !strings.Contains(details.RequestBody, "X-Sf-Token: [REDACTED]") {
		t.Fatalf("expected redacted API token header, got %q", details.RequestBody)
	}
	if !strings.Contains(details.RequestBody, `"content":"[REDACTED]"`) {
		t.Fatalf("expected redacted CA certificate content, got %q", details.RequestBody)
	}
}
