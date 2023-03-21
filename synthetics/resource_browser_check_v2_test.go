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
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sc "github.com/splunk/syntheticsclient/syntheticsclientv2"
)

const newBrowserCheckV2Config = `
resource "synthetics_create_browser_check_v2" "browser_v2_foo_check" {
  test {
    active = true
    device_id = 1  
    frequency = 5
    location_ids = ["aws-us-east-1"]
    name = "Terraform - Browser V2 Checkaroo"
    scheduling_strategy = "round_robin"
    url_protocol = "https://"
    start_url = "www.splunk.com"
    business_transactions {
      name = "Synthetic transaction 1"
      steps {
        name = "Go to URL"
        type = "go_to_url"
        url = "https://www.splunk.com"
        action = "go_to_url"
        wait_for_nav = true
        options {
          url = "https://www.splunk.com"
        }
      }
    }
    business_transactions {
      name = "New synthetic transaction"
      steps {
        name = "New step"
        type = "go_to_url"
        wait_for_nav = true
        action = "go_to_url"
        url = "https://www.batman.com"
      }
    }
    advanced_settings {
      verify_certificates = false
      user_agent = "Mozilla/5.0 (X11; Linux x86_64; Splunk Synthetics) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
      authentication {
        username = "batmab"
        password = "{{env.beep-var}}"
      }
      headers {
        name = "superstar-machine"
        value = "\"taking it too the staaaaars\""
        domain = "asdasd.batman.com"
      }
      cookies {
        key = "sda"
        value = "sda"
        domain = "asd.com"
        path = "/asd"
      }
      host_overrides {
        source = "asdasd.com"
        target = "whost.com"
        keep_host_header = false
      }
    }
  }    
}
`

const updatedBrowserCheckV2Config = `
resource "synthetics_create_browser_check_v2" "browser_v2_foo_check" {
  test {
    active = true
    device_id = 1  
    frequency = 15
    location_ids = ["aws-us-east-1"]
    name = "Terraform - Browser V2 Checkaroo"
    scheduling_strategy = "round_robin"
    url_protocol = "https://"
    start_url = "www.splunk.com"
    business_transactions {
      name = "Synthetic Super transaction 1"
      steps {
        name = "Go to URL"
        type = "go_to_url"
        url = "https://www.splunk.com"
        action = "go_to_url"
        wait_for_nav = true
        options {
          url = "https://www.splunk.com"
        }
      }
    }
    business_transactions {
      name = "New synthetic transaction"
      steps {
        name = "New step"
        type = "go_to_url"
        wait_for_nav = true
        action = "go_to_url"
        url = "https://www.batman.com"
      }
    }
    advanced_settings {
      verify_certificates = false
      user_agent = "Mozilla/5.0 (X11; Linux x86_64; Splunk Synthetics) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
      authentication {
        username = "batmab"
        password = "{{env.beep-var}}"
      }
      headers {
        name = "superstar-machine"
        value = "\"taking it too the staaaaars\""
        domain = "asdasd.batman.com"
      }
      cookies {
        key = "sda"
        value = "sda"
        domain = "asd.com"
        path = "/asd"
      }
      host_overrides {
        source = "asdasd.com"
        target = "whost.com"
        keep_host_header = false
      }
    }
  }    
}
`

func TestAccCreateUpdateBrowserCheckV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccBrowserCheckV2Destroy,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: newBrowserCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateBrowserCheckV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "frequency", "5"),
				),
			},
			{
				ResourceName:      "synthetics_create_browser_check_v2.browser_v2_foo_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_browser_check_v2.browser_v2_foo_check"),
				ImportStateVerify: true,
			},
			// Update It
			{
				Config: updatedBrowserCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCreateUpdateBrowserCheckV2ResourceExists,
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "frequency", "15"),
				),
			},
		},
	})
}

func testAccCreateUpdateBrowserCheckV2ResourceExists(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_browser_check_v2":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetBrowserCheckV2(checkId)
			if strconv.Itoa(check.Test.ID) != rs.Primary.ID || err != nil {
				return fmt.Errorf("Error finding browser check v2 %s: %s", rs.Primary.ID, err)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}
	return nil
}

func testAccBrowserCheckV2Destroy(s *terraform.State) error {
	token := os.Getenv("OBSERVABILITY_API_TOKEN")
	realm := os.Getenv("REALM")
	client := sc.NewClient(token, realm)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_browser_check_v2":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetBrowserCheckV2(checkId)
			if check.Test.ID != checkId || err != nil {
				return fmt.Errorf("Found deleted check %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

