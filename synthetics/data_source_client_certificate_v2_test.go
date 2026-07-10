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

import "testing"

func TestClientCertificateV2DataSourceIsMetadataOnly(t *testing.T) {
	dataSource := dataSourceClientCertificateV2()

	for _, key := range []string{"id", "name", "description", "domain", "expires_at", "created_at", "created_by", "updated_at", "updated_by"} {
		if _, ok := dataSource.Schema[key]; !ok {
			t.Fatalf("client certificate data source schema missing %q", key)
		}
	}
	if _, ok := dataSource.Schema["public_key"]; ok {
		t.Fatal("client certificate data source must not expose public_key")
	}
	if _, ok := dataSource.Schema["private_key"]; ok {
		t.Fatal("client certificate data source must not expose private_key")
	}
}
