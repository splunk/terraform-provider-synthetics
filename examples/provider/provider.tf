terraform {
  required_providers {
    synthetics = {
      version = "1.0.1"
      source  = "splunk/synthetics"
    }
  }
}

provider "synthetics" {
  product = "observability"
  realm = "us1"
  #apikey = "this-is-my-api-key"
}
