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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

const newHttpCheckV2Config = `
resource "synthetics_create_http_check_v2" "http_v2_foo_check" {
	provider = synthetics.synthetics
  test {
    active = true 
    frequency = 5
    location_ids = ["aws-us-east-1","aws-us-west-2"]
    name = "01-acceptance-Terraform-HTTP-V2"
    type = "http"
    url = "https://www.splunk.com"
    automatic_retries = 1
    scheduling_strategy = "round_robin"
		custom_properties {
			key = "key"
			value = "value"
		}
    request_method = "POST"
    verify_certificates = true
    user_agent = "Another User of Agents"
    body = "Beepbeepbeep"
		headers {
			name = "Synthetic_transaction_1"
			value = "batmantis is a mantis not a man"
		}
		headers {
			name = "back_transaction_1"
			value = "peeko"
		}
    validations {
        name = "My validation step 01"
        actual = "{{response.code}}"
        comparator = "equals"
        expected = 200
        type = "assert_numeric"
    }
    validations {
        name = "02 My validation step"
        actual = "{{response.body}}"
        comparator = "does_not_equal"
        expected = "11"
        type = "assert_string"
    }
  }    
}
`

const updatedHttpCheckV2Config = `
resource "synthetics_create_http_check_v2" "http_v2_foo_check" {
	provider = synthetics.synthetics
  test {
    active = false 
    frequency = 15
    location_ids = ["aws-us-west-2"]
    name = "01-acceptance-updated-Terraform-HTTP-V2"
    type = "http"
    url = "https://www.duckduckgo.com"
    automatic_retries = 0
    scheduling_strategy = "concurrent"
		custom_properties {
			key = "beepkey"
			value = "boopvalue"
		}
    request_method = "PUT"
    verify_certificates = false
    user_agent = "Another User of Agents and snake oil"
    body = "boopboopboop"
		headers {
			name = "Synthetic_transaction_01"
			value = "batmantis is a mantis not a man. Man."
		}
		headers {
			name = "back_transaction_01"
			value = "peekoboot"
		}
		validations {
				name = "002 My validation step"
				actual = "{{response.body}}"
				comparator = "matches"
				expected = "12221"
				type = "assert_string"
		}
		validations {
				name = "My validation step 001"
				actual = "{{response.code}}"
				comparator = "does_not_equal"
				expected = 400
				type = "assert_numeric"
		}
  }    
}
`

func TestAccCreateUpdateHttpCheckV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: providerConfig + newHttpCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.active", "true"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.frequency", "5"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.automatic_retries", "1"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.location_ids.0", "aws-us-east-1"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.location_ids.1", "aws-us-west-2"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.name", "01-acceptance-Terraform-HTTP-V2"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.type", "http"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.url", "https://www.splunk.com"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.scheduling_strategy", "round_robin"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.custom_properties.0.key", "key"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.custom_properties.0.value", "value"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.request_method", "POST"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.verify_certificates", "true"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.user_agent", "Another User of Agents"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.body", "Beepbeepbeep"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.headers.0.name", "Synthetic_transaction_1"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.headers.0.value", "batmantis is a mantis not a man"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.headers.1.name", "back_transaction_1"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.headers.1.value", "peeko"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.name", "My validation step 01"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.actual", "{{response.code}}"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.comparator", "equals"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.expected", "200"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.type", "assert_numeric"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.name", "02 My validation step"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.actual", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.comparator", "does_not_equal"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.expected", "11"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.type", "assert_string"),
				),
			},
			{
				ResourceName:      "synthetics_create_http_check_v2.http_v2_foo_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_http_check_v2.http_v2_foo_check"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: providerConfig + updatedHttpCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.active", "false"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.frequency", "15"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.automatic_retries", "0"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.location_ids.0", "aws-us-west-2"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.name", "01-acceptance-updated-Terraform-HTTP-V2"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.type", "http"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.url", "https://www.duckduckgo.com"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.scheduling_strategy", "concurrent"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.custom_properties.0.key", "beepkey"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.custom_properties.0.value", "boopvalue"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.request_method", "PUT"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.verify_certificates", "false"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.user_agent", "Another User of Agents and snake oil"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.body", "boopboopboop"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.headers.0.name", "Synthetic_transaction_01"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.headers.0.value", "batmantis is a mantis not a man. Man."),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.headers.1.name", "back_transaction_01"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.headers.1.value", "peekoboot"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.name", "002 My validation step"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.actual", "{{response.body}}"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.comparator", "matches"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.expected", "12221"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.0.type", "assert_string"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.name", "My validation step 001"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.actual", "{{response.code}}"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.comparator", "does_not_equal"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.expected", "400"),
					resource.TestCheckResourceAttr("synthetics_create_http_check_v2.http_v2_foo_check", "test.0.validations.1.type", "assert_numeric"),
				),
			},
		},
	})
}

func TestBuildHttpV2DataIncludesCertificateID(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceHttpCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{map[string]interface{}{
			"name":                "api-example",
			"type":                "http",
			"url":                 "https://api.example.com/health",
			"active":              true,
			"frequency":           5,
			"automatic_retries":   0,
			"scheduling_strategy": "round_robin",
			"request_method":      "GET",
			"body":                "",
			"location_ids":        []interface{}{"aws-us-east-1"},
			"user_agent":          "",
			"verify_certificates": true,
			"headers":             []interface{}{},
			"validations":         []interface{}{},
			"custom_properties":   []interface{}{},
			"certificate_id":      123,
		}},
	})

	input := buildHttpV2Data(d)

	if input.Test.CertificateID == nil {
		t.Fatal("CertificateID is nil, want nullable int")
	}
	if input.Test.CertificateID.Value == nil {
		t.Fatal("CertificateID.Value is nil, want 123")
	}
	if *input.Test.CertificateID.Value != 123 {
		t.Fatalf("CertificateID.Value = %d, want 123", *input.Test.CertificateID.Value)
	}
}

func TestFlattenHttpV2ReadIncludesCertificateID(t *testing.T) {
	response := &sc2.HttpCheckV2Response{}
	response.Test.CertificateID = sc2.NewNullableInt(123)

	flattened := flattenHttpV2Read(response)
	testBlock := flattened[0].(map[string]interface{})
	if got := testBlock["certificate_id"]; got != 123 {
		t.Fatalf("certificate_id = %#v, want 123", got)
	}
}

func TestFlattenHttpV2ReadResetsClearedCertificateIDInState(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceHttpCheckV2().Schema, map[string]interface{}{
		"test": []interface{}{map[string]interface{}{
			"name":                "api-example",
			"type":                "http",
			"url":                 "https://api.example.com/health",
			"active":              true,
			"frequency":           5,
			"automatic_retries":   0,
			"scheduling_strategy": "round_robin",
			"request_method":      "GET",
			"location_ids":        []interface{}{"aws-us-west-2"},
			"verify_certificates": true,
			"headers":             []interface{}{},
			"validations":         []interface{}{},
			"custom_properties":   []interface{}{},
			"certificate_id":      123,
		}},
	})

	response := &sc2.HttpCheckV2Response{}
	response.Test.Name = "api-example"
	response.Test.Type = "http"
	response.Test.URL = "https://api.example.com/health"
	response.Test.Active = true
	response.Test.Frequency = 5
	response.Test.SchedulingStrategy = "round_robin"
	response.Test.RequestMethod = "GET"
	response.Test.LocationIds = []string{"aws-us-west-2"}
	response.Test.Verifycertificates = true

	if err := d.Set("test", flattenHttpV2Read(response)); err != nil {
		t.Fatalf("set flattened HTTP test: %s", err)
	}

	testBlock := d.Get("test").(*schema.Set).List()[0].(map[string]interface{})
	if got := testBlock["certificate_id"]; got != 0 {
		t.Fatalf("certificate_id = %#v, want 0", got)
	}
}

func TestBuildHttpV2UpdateClearsRemovedCertificateID(t *testing.T) {
	oldTest := map[string]interface{}{"certificate_id": 123}
	newTest := map[string]interface{}{}

	input := buildHttpV2CertificateIDForUpdate(oldTest, newTest)

	if input == nil {
		t.Fatal("certificate ID update is nil, want null value")
	}
	if input.Value != nil {
		t.Fatalf("certificate ID update value = %#v, want nil", *input.Value)
	}
}
