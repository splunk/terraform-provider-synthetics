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
	"fmt"
)

// maxCertificateContentBytes mirrors the limit the Synthetics API enforces on uploaded certificate content (a client
// certificate chain or a CA bundle). Checking it here reports the problem at plan time instead of as an API error
// during apply.
const maxCertificateContentBytes = 256 * 1024

const certificateExpiresAtDescription = "When the certificate expires, in UTC. For a chain or bundle of certificates, the earliest expiry among them."

// validateCertificateContentSize rejects certificate content larger than the API accepts.
func validateCertificateContentSize(value interface{}, key string) ([]string, []error) {
	content, ok := value.(string)
	if !ok {
		return nil, nil
	}
	// Content that is not Base64 is left for the API to reject.
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, nil
	}
	if len(data) > maxCertificateContentBytes {
		return nil, []error{fmt.Errorf("%s is too large (%d KB); the maximum is %d KB", key, (len(data)+1023)/1024, maxCertificateContentBytes/1024)}
	}
	return nil, nil
}
