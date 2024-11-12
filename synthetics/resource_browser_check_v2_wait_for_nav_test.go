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

type waitForNavTestValues struct {
	waitForNav                bool
	waitForNavTimeout         int
	expectedNav               string
	expectedNavTimeout        string
	expectedNavTimeoutDefault string
}

func waitForNavTestConfig(name string, waitForNav bool, waitForNavTimeout int) string {
	var waitForNavTimeoutStr string
	if waitForNavTimeout == 0 {
		waitForNavTimeoutStr = "null"
	} else {
		waitForNavTimeoutStr = strconv.Itoa(waitForNavTimeout)
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
		name                 = "click"
		selector             = "#order"
		selector_type        = "id"
		type                 = "click_element"
		wait_for_nav         = %t
		wait_for_nav_timeout = %s
	  }
	}
  }
}
`, name, waitForNav, waitForNavTimeoutStr)
}

func TestAccCreateUpdateBrowserCheckV2WaitForNav(t *testing.T) {

	waitForNavTestCases := []struct {
		name   string
		apply1 waitForNavTestValues
		apply2 waitForNavTestValues
	}{
		{
			name: "wait_for_nav true with default timeout",
			apply1: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         0,
				expectedNav:               "true",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
			apply2: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         0,
				expectedNav:               "true",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
		},
		{
			name: "wait_for_nav false with default timeout",
			apply1: waitForNavTestValues{
				waitForNav:                false,
				waitForNavTimeout:         0,
				expectedNav:               "false",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
			apply2: waitForNavTestValues{
				waitForNav:                false,
				waitForNavTimeout:         0,
				expectedNav:               "false",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
		},
		{
			name: "wait_for_nav true->false with default timeout",
			apply1: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         0,
				expectedNav:               "true",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
			apply2: waitForNavTestValues{
				waitForNav:                false,
				waitForNavTimeout:         0,
				expectedNav:               "false",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
		},
		{
			name: "wait_for_nav false->true with default timeout",
			apply1: waitForNavTestValues{
				waitForNav:                false,
				waitForNavTimeout:         0,
				expectedNav:               "false",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
			apply2: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         0,
				expectedNav:               "true",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
		},
		{
			name: "wait_for_nav true default->custom timeout",
			apply1: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         0,
				expectedNav:               "true",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
			apply2: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         1500,
				expectedNav:               "true",
				expectedNavTimeout:        "1500",
				expectedNavTimeoutDefault: "false",
			},
		},
		{
			name: "wait_for_nav true custom->default timeout",
			apply1: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         1500,
				expectedNav:               "true",
				expectedNavTimeout:        "1500",
				expectedNavTimeoutDefault: "false",
			},
			apply2: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         0,
				expectedNav:               "true",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
		},
		{
			name: "wait_for_nav true custom=default timeout -> default",
			apply1: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         2000,
				expectedNav:               "true",
				expectedNavTimeout:        "2000",
				expectedNavTimeoutDefault: "false",
			},
			apply2: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         0,
				expectedNav:               "true",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
		},
		{
			name: "wait_for_nav true default -> custom=default timeout",
			apply1: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         0,
				expectedNav:               "true",
				expectedNavTimeout:        "0",
				expectedNavTimeoutDefault: "true",
			},
			apply2: waitForNavTestValues{
				waitForNav:                true,
				waitForNavTimeout:         2000,
				expectedNav:               "true",
				expectedNavTimeout:        "2000",
				expectedNavTimeoutDefault: "false",
			},
		},
	}

	for _, tc := range waitForNavTestCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: providerConfig + waitForNavTestConfig(tc.name, tc.apply1.waitForNav, tc.apply1.waitForNavTimeout),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.wait_for_nav", tc.apply1.expectedNav),
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.wait_for_nav_timeout", tc.apply1.expectedNavTimeout),
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.wait_for_nav_timeout_default", tc.apply1.expectedNavTimeoutDefault),
						),
					},
					{
						Config: providerConfig + waitForNavTestConfig(tc.name, tc.apply2.waitForNav, tc.apply2.waitForNavTimeout),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.wait_for_nav", tc.apply2.expectedNav),
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.wait_for_nav_timeout", tc.apply2.expectedNavTimeout),
							resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.example", "test.0.transactions.0.steps.1.wait_for_nav_timeout_default", tc.apply2.expectedNavTimeoutDefault),
						),
					},
				},
			})
		})
	}
}
