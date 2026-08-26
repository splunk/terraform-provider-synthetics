terraform {
  required_providers {
    synthetics = {
      version = "3.0.0"
      source  = "splunk/synthetics"
    }
  }
}

provider "synthetics" {
  realm = "us1"
  #apikey = "this-is-my-api-key"
}
