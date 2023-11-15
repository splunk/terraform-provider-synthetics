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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider

const (
  // providerConfig is a shared configuration to combine with the actual
  // test configuration so the HashiCups client is properly configured.
  // It is also possible to use the HASHICUPS_ environment variables instead,
  // such as updating the Makefile and running the testing through that tool.
  providerConfig = `
variable "observability_token" {
	description = "API token for observability"
}
provider "synthetics" {
	alias = "synthetics"
	product = "observability"
	realm = "us1"
	apikey = var.observability_token
}
`
	rigorConfig = `
variable "rigor_token" {
	description = "API token for rigor"
}
provider "synthetics" {
	alias = "rigor"
	product = "rigor"
	realm = "us1"
	apikey = var.rigor_token
}	
`
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"synthetics": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("TF_VAR_rigor_token"); err == "" {
		t.Fatal("TF_VAR_rigor_token environment variable must be set for acceptance tests. Set to empty string if not testing v1 rigor resources.")
	}
	if err := os.Getenv("TF_VAR_observability_token"); err == "" {
		t.Fatal("TF_VAR_observability_token environment variable must be set for acceptance tests. Set to empty string if not testing v2 Observability resources")
	}
	if err := os.Getenv("REALM"); err == "" {
		t.Fatal("REALM environment variable must be set for acceptance tests. If testing v1 rigor resources any value can be provided")
	}
}
