resource "synthetics_create_location_v2" "location_v2_foo" {
  location {
    id = "private-aws-awesome-east"
    label = "awesome aws east location"
  }    
}