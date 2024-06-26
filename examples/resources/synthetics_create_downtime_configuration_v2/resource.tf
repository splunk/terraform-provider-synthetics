resource "synthetics_create_downtime_configuration_v2" "downtime_configuration_v2_foo" {
  downtime_configuration {
    name = "acceptance-downtime-configuration-terraform-test"
    description = "The most awesome downtime_configuration. Full of snakes."
    rule = "augment_data"
    start_time = "2024-07-01T17:13:00.000Z"
    end_time = "2024-08-01T17:13:00.000Z"
    test_ids = [932826] 
  }
}