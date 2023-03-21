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
    transactions {
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
      steps {
        name = "New step"
        type = "click_element"
        selector_type = "id"
        wait_for_nav = false
        selector = "\"free-splunk-click-mobile\""
      }
      steps {
        name = "New step"
        type = "click_element"
        selector_type = "id"
        wait_for_nav = false
        selector = "login-button"
      }
    }
    transactions {
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
