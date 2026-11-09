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
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// These limits mirror the ones the Synthetics API enforces on uploaded certificates. Checking them here reports
// the problem at plan time instead of as an API error during apply.
const (
	maxClientCertificateChainLength = 10
	maxCaCertificateContentBytes    = 256 * 1024
)

const certificateExpiresAtDescription = "When the certificate expires, in UTC. For a chain or bundle of certificates, the earliest expiry among them."

// validateClientCertificatePublicKeyContent rejects a client certificate chain with more certificates than the API accepts.
func validateClientCertificatePublicKeyContent(value interface{}, key string) ([]string, []error) {
	data, ok := decodeBase64Content(value)
	if !ok {
		return nil, nil
	}
	if count := countPEMCertificates(data); count > maxClientCertificateChainLength {
		return nil, []error{fmt.Errorf("%s contains %d certificates; the maximum is %d", key, count, maxClientCertificateChainLength)}
	}
	return nil, nil
}

// validateCaCertificateContent rejects CA certificate content larger than the API accepts.
func validateCaCertificateContent(value interface{}, key string) ([]string, []error) {
	data, ok := decodeBase64Content(value)
	if !ok {
		return nil, nil
	}
	if len(data) > maxCaCertificateContentBytes {
		return nil, []error{fmt.Errorf("%s is too large (%d KB); the maximum is %d KB", key, (len(data)+1023)/1024, maxCaCertificateContentBytes/1024)}
	}
	return nil, nil
}

// decodeBase64Content returns the decoded content. Content that is not a Base64 string is left for the API to reject.
func decodeBase64Content(value interface{}) ([]byte, bool) {
	content, ok := value.(string)
	if !ok {
		return nil, false
	}
	data, err := base64.StdEncoding.DecodeString(content)
	return data, err == nil
}

// countPEMCertificates counts the CERTIFICATE blocks in PEM data. Other blocks, such as a private key, and
// content that is not PEM (for example a single DER certificate) are not counted.
func countPEMCertificates(data []byte) int {
	count := 0
	for {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			return count
		}
		if block.Type == "CERTIFICATE" {
			count++
		}
	}
}
