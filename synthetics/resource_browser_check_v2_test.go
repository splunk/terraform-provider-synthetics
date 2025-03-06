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
)

const newBrowserCheckV2Config = `
resource "synthetics_create_variable_v2" "variable_v2_foo" {
  provider = synthetics.synthetics
  variable {
    description = "The most awesome variable. Full of snakes."
    value = "barv3v3"
    name = "acceptance-variable-terraform-test"
    secret = false
  }
}

resource "synthetics_create_browser_check_v2" "browser_v2_foo_check" {
  provider = synthetics.synthetics
  depends_on = [synthetics_create_variable_v2.variable_v2_foo]
  test {
    active = true
    device_id = 1
    frequency = 5
    location_ids = ["aws-us-east-1"]
    automatic_retries = 1
    name = "01-acceptance-Terraform-Browser-V2"
    scheduling_strategy = "round_robin"
    custom_properties {
      key = "key"
      value = "value"
    }
    advanced_settings {
      verify_certificates = true
      user_agent = "Mozilla/5.0 (X11; Linux x86_64; Splunk Synthetics) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
      collect_interactive_metrics = false
      authentication {
        username = "batmab"
        password = "{{env.acceptance-variable-terraform-test}}"
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
      cookies {
        key = "yes"
        value = "no"
        domain = "zodiak.com"
        path = "/Edlesley"
      }
      chrome_flags {
        name = "--proxy-server"
        value = "my-proxy-server:80"
      }
      chrome_flags {
        name = "--proxy-bypass-list"
        value = "127.0.0.1:8080"
      }
      host_overrides {
        source = "asdasd.com"
        target = "whost.com"
        keep_host_header = false
      }
      host_overrides {
        source = "92.2.2.2"
        target = "91.1.1.1"
        keep_host_header = true
      }
    }
    transactions {
      name = "First Synthetic transaction"
      steps {
        name                 = "01 Go to URL"
        type                 = "go_to_url"
        url                  = "https://www.splunk.com"
      }
      steps {
        name                 = "02 fill in fieldz"
        selector             = "beep"
        selector_type        = "id"
        type                 = "enter_value"
        value                = "{{env.acceptance-variable-terraform-test}}"
        wait_for_nav_timeout = null
      }
      steps {
        name                 = "03 click"
        selector             = "clicky"
        selector_type        = "id"
        type                 = "click_element"
        wait_for_nav         = true
      }
      steps {
        name                 = "04 accept---Alert"
        type                 = "accept_alert"
      }
      steps {
        name                 = "05 Select-val-text"
        option_selector      = "sdad"
        option_selector_type = "text"
        selector             = "textzz"
        selector_type        = "id"
        type                 = "select_option"
        wait_for_nav         = false
        wait_for_nav_timeout = 1
      }
    }
    transactions {
      name = "2nd Synthetic transaction"
      steps {
        name                 = "Go to other URL"
        type                 = "go_to_url"
        url                  = "https://www.splunk.com"
      }
      steps {
        name                 = "fill in more fields field"
        selector             = "beep"
        selector_type        = "id"
        type                 = "enter_value"
        value                = "{{env.acceptance-variable-terraform-test}}"
      }
      steps {
        name                 = "assert element visible"
        type                 = "assert_element_visible"
        selector             = "beep"
        selector_type        = "id"
        max_wait_time        = 1000
      }
      steps {
        name                 = "assert element visible no max wait time"
        type                 = "assert_element_visible"
        selector             = "beep"
        selector_type        = "id"
      }
    }
    transactions {
      name = "3nd transaction - wait_for_nav tests"
      steps {
        name                 = "Go to other URL"
        type                 = "go_to_url"
        url                  = "https://www.splunk.com"
      }
      steps {
        name                 = "wait_for_nav true -> false with default timeout"
        selector             = "#buy"
        selector_type        = "id"
        type                 = "click_element"
        wait_for_nav         = true
      }
      steps {
        name                 = "wait_for_nav false -> true with default timeout"
        selector             = "#buy"
        selector_type        = "id"
        type                 = "click_element"
        wait_for_nav         = false
      }
      steps {
        name                 = "wait_for_nav true ->false with custom timeout"
        selector             = "#buy"
        selector_type        = "id"
        type                 = "click_element"
        wait_for_nav         = true
        wait_for_nav_timeout = 2010
      }
      steps {
        name                 = "wait_for_nav false ->true with custom timeout"
        selector             = "#buy"
        selector_type        = "id"
        type                 = "click_element"
        wait_for_nav         = false
        wait_for_nav_timeout = 2020
      }
      steps {
        name                 = "wait_for_nav default"
        selector             = "#buy"
        selector_type        = "id"
        type                 = "click_element"
      }
      steps {
        name                 = "wait_for_nav true with default->custom timeout"
        selector             = "#buy"
        selector_type        = "id"
        type                 = "click_element"
        wait_for_nav         = true
      }
      steps {
        name                 = "wait_for_nav true with custom->default timeout"
        selector             = "#buy"
        selector_type        = "id"
        type                 = "click_element"
        wait_for_nav         = true
        wait_for_nav_timeout = 2030
      }
    }
  }
}
`

const updatedBrowserCheckV2Config = `
resource "synthetics_create_variable_v2" "variable_v2_foo" {
    provider = synthetics.synthetics
  variable {
    description = "The most awesome variable. Full of snakes."
    value = "barv3v3"
    name = "acceptance-variable-terraform-test"
    secret = false
  }
}
resource "synthetics_create_browser_check_v2" "browser_v2_foo_check" {
  provider = synthetics.synthetics
  depends_on = [synthetics_create_variable_v2.variable_v2_foo]
  test {
    active = false
    device_id = 2
    frequency = 15
    location_ids = ["aws-us-west-1"]
    automatic_retries = 0
    name = "01-acceptance-updated-Terraform-Browser-V2"
    scheduling_strategy = "concurrent"
    custom_properties {
      key = "beepkey"
      value = "boop value 2"
    }
    advanced_settings {
      verify_certificates = false
      user_agent = "Jozilla/5.0"
      collect_interactive_metrics = false
      authentication {
        username = "batmantis"
        password = "{{env.acceptance-variable-terraform-test}}"
      }
      headers {
        name = "superstar-machine-show"
        value = "\"taking it too the stars\""
        domain = "davidcrossed.batman.com"
      }
      cookies {
        key = "sda2"
        value = "sda2"
        domain = "asd2.com"
        path = "/asd2"
      }
      cookies {
        key = "yes"
        value = "no"
        domain = "zodiak.com"
        path = "/Edlesley"
      }
      chrome_flags {
        name = "--proxy-server"
        value = "foo:80"
      }
      chrome_flags {
        name = "--proxy-bypass-list"
        value = "*google.com"
      }
      host_overrides {
        source = "asdasd.com"
        target = "whost.com"
        keep_host_header = false
      }
      host_overrides {
        source = "92.2.2.2"
        target = "91.1.1.1"
        keep_host_header = true
      }
    }
    transactions {
      name = "01 First Synthetic transaction"
      steps {
        name                 = "01 Go to URL"
        type                 = "go_to_url"
        url                  = "https://www.splunk.com"
      }
      steps {
        name                 = "06 Select-Val-Val"
        option_selector      = "{{env.acceptance-variable-terraform-test}}"
        option_selector_type = "value"
        selector             = "valz"
        selector_type        = "id"
        type                 = "select_option"
        wait_for_nav         = false
      }
      steps {
        name                 = "07 Select-Val-Index"
        option_selector      = "{{env.acceptance-variable-terraform-test}}"
        option_selector_type = "index"
        selector             = "selectionz"
        selector_type        = "id"
        type                 = "select_option"
        wait_for_nav         = false
      }
      steps {
        name                 = "08 Save as text"
        selector             = "beepval"
        selector_type        = "link"
        type                 = "store_variable_from_element"
        variable_name        = "{{env.terraform-test-foo-301}}"
      }
      steps {
        name                 = "08.5 Wait"
        duration             = 4234
        type                 = "wait"
      }
      steps {
        name                 = "09 Save JS2 return Val"
        type                 = "store_variable_from_javascript"
        value                = "sdasds"
        variable_name        = "{{env.terraform-test-foo-301}}"
        wait_for_nav         = true
      }
      steps {
        name                 = "010 Run JS"
        type                 = "run_javascript"
        value                = "beeeeeeep"
        wait_for_nav         = true
        wait_for_nav_timeout = 1000
      }
    }
    transactions {
      name = "2nd Synthetic transaction"
      steps {
        name                 = "Go to other URL"
        type                 = "go_to_url"
        url                  = "https://www.splunk.com"
      }
      steps {
        name                 = "fill in more fields field"
        selector             = "beep"
        selector_type        = "id"
        type                 = "enter_value"
        value                = "{{env.acceptance-variable-terraform-test}}"
        wait_for_nav_timeout = 60
      }
      steps {
        name                 = "assert element visible"
        type                 = "assert_element_visible"
        selector             = "beep"
        selector_type        = "id"
      }
      steps {
        name                 = "assert element visible with max wait time"
        type                 = "assert_element_visible"
        selector             = "beep"
        selector_type        = "id"
        max_wait_time        = "20000"
      }
    }
  }
}
`

func TestAccCreateUpdateBrowserCheckV2(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create It
			{
				Config: providerConfig + newBrowserCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.active", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.device_id", "1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.frequency", "5"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.automatic_retries", "1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.location_ids.0", "aws-us-east-1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.name", "01-acceptance-Terraform-Browser-V2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.scheduling_strategy", "round_robin"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.custom_properties.0.key", "key"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.custom_properties.0.value", "value"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.verify_certificates", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.user_agent", "Mozilla/5.0 (X11; Linux x86_64; Splunk Synthetics) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.collect_interactive_metrics", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.authentication.0.username", "batmab"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.authentication.0.password", "{{env.acceptance-variable-terraform-test}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.headers.0.name", "superstar-machine"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.headers.0.value", "\"taking it too the staaaaars\""),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.headers.0.domain", "asdasd.batman.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.0.key", "sda"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.0.value", "sda"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.0.domain", "asd.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.0.path", "/asd"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.1.key", "yes"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.1.value", "no"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.1.domain", "zodiak.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.1.path", "/Edlesley"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.chrome_flags.0.name", "--proxy-bypass-list"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.chrome_flags.0.value", "127.0.0.1:8080"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.chrome_flags.1.name", "--proxy-server"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.chrome_flags.1.value", "my-proxy-server:80"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.0.source", "asdasd.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.0.target", "whost.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.0.keep_host_header", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.1.source", "92.2.2.2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.1.target", "91.1.1.1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.1.keep_host_header", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.name", "First Synthetic transaction"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.name", "01 Go to URL"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.type", "go_to_url"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.url", "https://www.splunk.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.name", "02 fill in fieldz"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.selector", "beep"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.type", "enter_value"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.value", "{{env.acceptance-variable-terraform-test}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.wait_for_nav_timeout_default", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.name", "03 click"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.selector", "clicky"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.type", "click_element"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.wait_for_nav", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.name", "04 accept---Alert"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.type", "accept_alert"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.name", "05 Select-val-text"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.option_selector", "sdad"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.option_selector_type", "text"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.selector", "textzz"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.type", "select_option"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.wait_for_nav_timeout", "1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.name", "2nd Synthetic transaction"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.name", "Go to other URL"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.type", "go_to_url"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.url", "https://www.splunk.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.name", "fill in more fields field"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.selector", "beep"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.type", "enter_value"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.value", "{{env.acceptance-variable-terraform-test}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.wait_for_nav_timeout_default", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.name", "assert element visible"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.type", "assert_element_visible"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.selector", "beep"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.max_wait_time", "1000"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.name", "assert element visible no max wait time"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.type", "assert_element_visible"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.selector", "beep"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.max_wait_time", "0"),
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
				Config: providerConfig + updatedBrowserCheckV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.#", "1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.active", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.device_id", "2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.frequency", "15"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.automatic_retries", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.location_ids.0", "aws-us-west-1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.name", "01-acceptance-updated-Terraform-Browser-V2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.scheduling_strategy", "concurrent"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.custom_properties.0.key", "beepkey"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.custom_properties.0.value", "boop value 2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.verify_certificates", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.user_agent", "Jozilla/5.0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.collect_interactive_metrics", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.authentication.0.username", "batmantis"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.authentication.0.password", "{{env.acceptance-variable-terraform-test}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.headers.0.name", "superstar-machine-show"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.headers.0.value", "\"taking it too the stars\""),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.headers.0.domain", "davidcrossed.batman.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.0.key", "sda2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.0.value", "sda2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.0.domain", "asd2.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.0.path", "/asd2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.1.key", "yes"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.1.value", "no"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.1.domain", "zodiak.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.cookies.1.path", "/Edlesley"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.chrome_flags.0.name", "--proxy-bypass-list"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.chrome_flags.0.value", "*google.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.chrome_flags.1.name", "--proxy-server"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.chrome_flags.1.value", "foo:80"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.0.source", "asdasd.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.0.target", "whost.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.0.keep_host_header", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.1.source", "92.2.2.2"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.1.target", "91.1.1.1"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.advanced_settings.0.host_overrides.1.keep_host_header", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.name", "01 First Synthetic transaction"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.name", "01 Go to URL"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.type", "go_to_url"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.0.url", "https://www.splunk.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.name", "06 Select-Val-Val"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.option_selector", "{{env.acceptance-variable-terraform-test}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.option_selector_type", "value"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.selector", "valz"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.type", "select_option"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.1.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.name", "07 Select-Val-Index"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.option_selector", "{{env.acceptance-variable-terraform-test}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.option_selector_type", "index"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.selector", "selectionz"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.type", "select_option"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.2.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.name", "08 Save as text"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.selector", "beepval"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.selector_type", "link"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.type", "store_variable_from_element"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.variable_name", "{{env.terraform-test-foo-301}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.3.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.name", "08.5 Wait"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.duration", "4234"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.type", "wait"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.4.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.5.name", "09 Save JS2 return Val"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.5.type", "store_variable_from_javascript"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.5.value", "sdasds"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.5.variable_name", "{{env.terraform-test-foo-301}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.5.wait_for_nav", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.5.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.5.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.6.name", "010 Run JS"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.6.type", "run_javascript"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.6.value", "beeeeeeep"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.6.wait_for_nav", "true"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.6.wait_for_nav_timeout", "1000"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.0.steps.6.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.name", "2nd Synthetic transaction"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.name", "Go to other URL"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.type", "go_to_url"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.0.url", "https://www.splunk.com"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.name", "fill in more fields field"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.selector", "beep"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.type", "enter_value"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.value", "{{env.acceptance-variable-terraform-test}}"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.wait_for_nav_timeout", "60"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.1.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.name", "assert element visible"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.type", "assert_element_visible"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.selector", "beep"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.2.max_wait_time", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.name", "assert element visible with max wait time"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.type", "assert_element_visible"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.selector", "beep"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.selector_type", "id"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.wait_for_nav", "false"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.wait_for_nav_timeout", "0"),
					resource.TestCheckResourceAttr("synthetics_create_browser_check_v2.browser_v2_foo_check", "test.0.transactions.1.steps.3.max_wait_time", "20000"),
				),
			},
		},
	})
}
