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
	sc "github.com/splunk/syntheticsclient/syntheticsclient"
)

func TestAccHttpCheckBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccHttpCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: rigorConfig + testAccHttpCheckConfigBasic("ineffective test", "https://www.google.com", 15),
				Check: resource.ComposeTestCheckFunc(
					testAccHttpCheckExists("synthetics_create_http_check.http_check"),
					resource.TestCheckResourceAttr("synthetics_create_http_check.http_check", "name", "ineffective test"),
					resource.TestCheckResourceAttr("synthetics_create_http_check.http_check", "url", "https://www.google.com"),
					resource.TestCheckResourceAttr("synthetics_create_http_check.http_check", "frequency", "15"),
				),
			},
			{
				ResourceName:      "synthetics_create_http_check.http_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_http_check.http_check"),
			},
			{
				Config: rigorConfig + testAccHttpCheckConfigBasic("updated test", "https://about.google/", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccHttpCheckExists("synthetics_create_http_check.http_check"),
					resource.TestCheckResourceAttr("synthetics_create_http_check.http_check", "name", "updated test"),
					resource.TestCheckResourceAttr("synthetics_create_http_check.http_check", "url", "https://about.google/"),
					resource.TestCheckResourceAttr("synthetics_create_http_check.http_check", "frequency", "5"),
				),
			},
		},
	})
}

func testAccHttpCheckConfigBasic(name string, url string, frequency int) string {
	return fmt.Sprintf(`
resource "synthetics_create_http_check" "http_check" {
		provider = synthetics.rigor
    name = "%s"
    url = "%s"  
    frequency = %d
}
`, name, url, frequency)
}

func testAccHttpCheckExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Check Id set")
		}

		return nil
	}
}

func testAccHttpCheckDestroy(s *terraform.State) error {
	token := os.Getenv("API_ACCESS_TOKEN")
	client := sc.NewClient(token)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_http_check":
			checkId, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error converting check id: %s", err)
			}
			check, _, err := client.GetCheck(checkId)
			if check.ID != checkId || err != nil {
				return fmt.Errorf("Found deleted check %s", rs.Primary.ID)
			}
		default:
			return fmt.Errorf("Unexpected resource of type: %s", rs.Type)
		}
	}

	return nil
}

func testAccStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		return rs.Primary.Attributes["id"], nil
	}
}
