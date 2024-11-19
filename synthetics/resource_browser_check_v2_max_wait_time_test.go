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
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type maxWaitTimeTestValues struct {
	maxWaitTime                int
	expectedMaxWaitTime        string
	expectedMaxWaitTimeDefault string
}

func maxWaitTimeTestConfig(name string, maxWaitTime int) string {
	var maxWaitTimeStr string
	if maxWaitTime == 0 {
		maxWaitTimeStr = "null"
	} else {
		maxWaitTimeStr = strconv.Itoa(maxWaitTime)
	}

	return fmt.Sprintf(`
resource "synthetics_create_browser_check_v2" "example" {
  provider = synthetics.synthetics
  test {
    device_id = 1
    frequency = 5
    location_ids = ["aws-us-east-1"]
    name = "%s"
    advanced_settings {
      verify_certificates = true
    }
    transactions {
      name = "First Synthetic transaction"
      steps {
        name                 = "01 Go to URL"
        type                 = "go_to_url"
        url                  = "https://www.splunk.com"
      }
      steps {
        name                 = "assert element visible"
        type                 = "assert_element_visible"
        selector             = "beep"
        selector_type        = "id"
        max_wait_time        = %s
      }
    }
  }
}
`, name, maxWaitTimeStr)
}

func TestAccCreateUpdateBrowserCheckV2MaxWaitTime(t *testing.T) {

	maxWaitTimeTestCases := []struct {
		name   string
		apply1 maxWaitTimeTestValues
		apply2 maxWaitTimeTestValues
	}{
		{
			name: "default value",
			apply1: maxWaitTimeTestValues{
				maxWaitTime:                0,
				expectedMaxWaitTime:        "0",
				expectedMaxWaitTimeDefault: "true",
			},
			apply2: maxWaitTimeTestValues{
				maxWaitTime:                0,
				expectedMaxWaitTime:        "0",
				expectedMaxWaitTimeDefault: "true",
			},
		},
		{
			name: "custom value -> default",
			apply1: maxWaitTimeTestValues{
				maxWaitTime:                5000,
				expectedMaxWaitTime:        "5000",
				expectedMaxWaitTimeDefault: "false",
			},
			apply2: maxWaitTimeTestValues{
				maxWaitTime:                0,
				expectedMaxWaitTime:        "0",
				expectedMaxWaitTimeDefault: "true",
			},
		},
		{
			name: "custom value matching default -> default",
			apply1: maxWaitTimeTestValues{
				maxWaitTime:                10000,
				expectedMaxWaitTime:        "10000",
				expectedMaxWaitTimeDefault: "false",
			},
			apply2: maxWaitTimeTestValues{
				maxWaitTime:                0,
				expectedMaxWaitTime:        "0",
				expectedMaxWaitTimeDefault: "true",
			},
		},
		{
			name: "default -> custom value matching default",
			apply1: maxWaitTimeTestValues{
				maxWaitTime:                0,
				expectedMaxWaitTime:        "0",
				expectedMaxWaitTimeDefault: "true",
			},
			apply2: maxWaitTimeTestValues{
				maxWaitTime:                10000,
				expectedMaxWaitTime:        "10000",
				expectedMaxWaitTimeDefault: "false",
			},
		},
	}

	for _, tc := range maxWaitTimeTestCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + maxWaitTimeTestConfig(tc.name, tc.apply1.maxWaitTime),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.max_wait_time", tc.apply1.expectedMaxWaitTime),
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.max_wait_time_default", tc.apply1.expectedMaxWaitTimeDefault),
						),
					},
					{
						Config: providerConfig + maxWaitTimeTestConfig(tc.name, tc.apply2.maxWaitTime),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.max_wait_time", tc.apply2.expectedMaxWaitTime),
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.max_wait_time_default", tc.apply2.expectedMaxWaitTimeDefault),
						),
					},
				},
			})
		})
	}
}
