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

func TestAccBrowserCheckBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccBrowserCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBrowserCheckConfigBasic("ineffective browser test", "https://www.google.com", "real_browser", 15),
				Check: resource.ComposeTestCheckFunc(
					testAccBrowserCheckExists("synthetics_create_browser_check.browser_check"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check.browser_check", "name", "ineffective browser test"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check.browser_check", "url", "https://www.google.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check.browser_check", "type", "real_browser"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check.browser_check", "frequency", "15"),
				),
			},
			{
				ResourceName:      "synthetics_create_browser_check.browser_check",
				ImportState:       true,
				ImportStateIdFunc: testAccStateIdFunc("synthetics_create_browser_check.browser_check"),
			},
			{
				Config: testAccBrowserCheckConfigBasic("updated test", "https://about.google/", "real_browser", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccBrowserCheckExists("synthetics_create_browser_check.browser_check"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check.browser_check", "name", "updated test"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check.browser_check", "url", "https://about.google/"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check.browser_check", "type", "real_browser"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check.browser_check", "frequency", "5"),
				),
			},
		},
	})
}

func testAccBrowserCheckConfigBasic(name string, url string, checktype string, frequency int) string {
	check := fmt.Sprintf(`
resource "synthetics_create_browser_check" "browser_check" {
 	name = "%s"
 	url = "%s"	
 	type = "%s"
 	frequency = %d
}
`, name, url, checktype, frequency)

	return check
}

func testAccBrowserCheckExists(n string) resource.TestCheckFunc {
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

func testAccBrowserCheckDestroy(s *terraform.State) error {
	token := os.Getenv("API_ACCESS_TOKEN")
	client := sc.NewClient(token)
	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "synthetics_create_browser_check":
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
