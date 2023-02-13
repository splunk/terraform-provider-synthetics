terraform {
  required_providers {
    synthetics = {
      version = "1.0.0"
      source  = "splunk.com/splunk/synthetics"
    }
  }
}

provider "synthetics" {
  product = "rigor"
}

//Pull a Check as a datasource

# data "synthetics_check" "datasource_check" {
#   id = 270949
# }

# output "datasource_check" {
#   value = data.synthetics_check.datasource_check
# }

//=================================
//=================================
//=================================

//Create an HTTP Check
# resource "synthetics_create_http_check" "http_check" {
#   name = "Uptime Checkaroo"
#   frequency = 15
#   url = "https://www.google.com"
#   success_criteria {
#     action_type = "presence_of_text"
#     comparison_string = "Images"
#   }
#   success_criteria {
#     action_type = "presence_of_text"
#     comparison_string = "About"
#   }
#   check_connection {
#     download_bandwidth = 20010
#     upload_bandwidth = 5010
#     latency = 23
#     packet_loss = 0.1
#   }
#   tags = ["beep","boiiiiip","haaaai"]
#   locations = [1,2,158,48,22]
#   round_robin = false
#   integrations = [80,82]
#   auto_retry = true
#   enabled = false
#   http_method = "GET"
#   notifications {
#     sms = true
#     email = false
#     call = true
#     notify_after_failure_count = 10
#     notify_on_location_failure = false
#     muted = true
#     notify_who {
#       sms = false
#       email = true
#       custom_user_email = "example_personasdasdasd@splunk.com"
#       call = false
#     }
#     notify_who {
#       type = "user"
#       id = 18100
#       sms = true
#       email = true
#     }
#     escalations {
#       sms = true
#       email = true
#       call = false
#       after_minutes = 30
#       notify_who {
#         type = "user"
#         id = 18100
#       }
#     }
#   }
# }

# output "http_check" {
#   value = synthetics_create_http_check.http_check
# }



# // Create a Browser Check †
# // †: steps and javascript_files currently not available through public API endpoints
# resource "synthetics_create_browser_check" "browser_check" {
#   name = "BROWSERMANIA"
#   frequency = 60
#   type = "real_browser"
#   url = "https://www.google.com"
#   viewport {
#     width = 800    
#     height = 600
#   }
#   check_connection {
#     download_bandwidth = 0
#     upload_bandwidth = 0
#     latency = 0
#     packet_loss = 0
#   }
#   tags = ["beelp","boiiiiip","haaaai"]
#   locations = [6]
#   integrations = [80,71]  
#   cookies {
#     key = "beep"
#     value = "boop"
#     domain = "google.com"
#     path = "/"
#   }
#   cookies {
#     key = "bam"
#     value = "botch"
#     domain = "google.com"
#     path = "/"
#   }
#   dns_overrides {
#     original_domain = "new.domain.com"
#     original_host = "123.123.123.123"
#   }
#   threshold_monitors {
#     matcher = "*.google.com"
#     metric_name = "first_byte_time_ms"
#     comparison_type = "less_than"
#     value = 9999
#   }
#   excluded_files {
#     pattern = ".+\\.clicktale\\.net"
#     preset_name = "clicktale"
#     exclusion_type = "preset"
#   }  
#   round_robin = false
#   auto_retry = true
#   enabled = false
#   notifications {
#     sms = true
#     email = false
#     call = true
#     notify_after_failure_count = 10
#     notify_on_location_failure = false
#     muted = true
#     notify_who {
#       sms = false
#       email = true
#       custom_user_email = "21example_person@splunk.com"
#       call = false
#     }
#     notify_who {
#       type = "user"
#       id = 18100
#       sms = true
#       email = true
#     }
#     escalations {
#       sms = true
#       email = true
#       call = false
#       after_minutes = 30
#       notify_who {
#         type = "user"
#         id = 18100
#       }
#     }
#   }
# }

# output "browser_check" {
#   value = synthetics_create_browser_check.browser_check
# }
