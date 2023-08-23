resource "synthetics_create_browser_check_v2" "long_browser_v2_foo_check" {
  test {
    active = true
    device_id = 1  
    frequency = 15
    location_ids = ["aws-us-east-1"]
    name = "0011aTerraform-Browser V2 Checkaroo"
    scheduling_strategy = "round_robin"
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
        value                = "{{env.beep-var}}"
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
        wait_for_nav         = false
      }
      steps {
        name                 = "05 Select-val-text"
        option_selector      = "sdad"
        option_selector_type = "text"
        selector             = "textzz"
        selector_type        = "id"
        type                 = "select_option"
        wait_for_nav         = false
      }
      steps {
        name                 = "06 Select-Val-Val"
        option_selector      = "{{env.beep-var}}"
        option_selector_type = "value"
        selector             = "valz"
        selector_type        = "id"
        type                 = "select_option"
        wait_for_nav         = false
      }
      steps {
        name                 = "07 Select-Val-Index"
        option_selector      = "{{env.beep-var}}"
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
        wait_for_nav         = false
      }
      steps {
        name                 = "08.5 Wait"
        duration             = 4234
        type                 = "wait"
        wait_for_nav         = false
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
        value                = "{{env.beep-var}}"
      }
    }
    advanced_settings {
      verify_certificates = true
      user_agent = "Mozilla/5.0 (X11; Linux x86_64; Splunk Synthetics) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"
      collect_interactive_metrics = false
      authentication {
        username = "batmab"
        password = "{{env.beep-var}}"
      }
      headers {
        name = "superstar-machine"
        value = "\"taking it too the staaaaars\""
        domain = "asdasd.batmab.com"
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
  }    
}