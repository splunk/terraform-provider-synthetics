resource "synthetics_create_port_check_v2" "port_v2_foo_check" {
  test {
    name = "Terraform - PORT V2 Checkaroo"
    # type = "port"
    port = 8080
    protocol = "udp"
    host = "www.splunk.com"
    location_ids = ["aws-us-west-2"]
    frequency = 5
    scheduling_strategy = "concurrent"
    custom_properties {
			key = "key"
			value = "value"
		}
    active = true 
  }    
}