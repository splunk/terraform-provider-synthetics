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

package syntheticsclientv2

import (
	"strings"
	"testing"
)

func TestSanitizeRequestDumpRedactsCaCertificateContentAndToken(t *testing.T) {
	requestDump := strings.Join([]string{
		"POST /v2/synthetics/cacerts HTTP/1.1",
		"Host: api.example.signalfx.com",
		"X-Sf-Token: secret-token",
		"Content-Type: application/json",
		"",
		`{"caCertificate":{"content":"private-ca-material","nested":{"Content":"nested-private-material"}}}`,
	}, "\r\n")

	sanitized := sanitizeRequestDump("/cacerts", []byte(requestDump))

	if strings.Contains(sanitized, "secret-token") {
		t.Fatalf("expected API token to be redacted, got %q", sanitized)
	}
	if strings.Contains(sanitized, "private-ca-material") || strings.Contains(sanitized, "nested-private-material") {
		t.Fatalf("expected CA certificate content to be redacted, got %q", sanitized)
	}
	if !strings.Contains(sanitized, `X-Sf-Token: <REDACTED>`) {
		t.Fatalf("expected redacted token header, got %q", sanitized)
	}
	if strings.Count(sanitized, `<REDACTED>`) != 1 || strings.Count(sanitized, `\u003cREDACTED\u003e`) != 2 {
		t.Fatalf("expected token and CA content fields to be redacted, got %q", sanitized)
	}
}

func TestSanitizeRequestDumpOnlyRedactsCaCertificateContentForCaCertificateEndpoints(t *testing.T) {
	requestDump := strings.Join([]string{
		"POST /v2/synthetics/tests/ssl HTTP/1.1",
		"Host: api.example.signalfx.com",
		"X-Sf-Token: secret-token",
		"Content-Type: application/json",
		"",
		`{"content":"non-ca-payload"}`,
	}, "\r\n")

	sanitized := sanitizeRequestDump("/tests/ssl", []byte(requestDump))

	if strings.Contains(sanitized, "secret-token") {
		t.Fatalf("expected API token to be redacted, got %q", sanitized)
	}
	if !strings.Contains(sanitized, "non-ca-payload") {
		t.Fatalf("expected non-CA content to remain, got %q", sanitized)
	}
}
