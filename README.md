# Splunk Synthetics Terraform Provider

This repository is a **beta* Terraform provider for [Splunk Synthetics in Splunk Observability](https://docs.splunk.com/Observability/synthetics/intro-synthetics.html). It currently contains CRUD operations for API Checks, Real Browser Checks, Port Checks, HTTP Checks, and Variables.

**NOTE:** The client expects a valid Splunk Observability API token defined as an environment variable named `OBSERVABILITY_API_TOKEN` (E.G. `export OBSERVABILITY_API_TOKEN="This_is_my_api_token"`)

### Rigor Classic (V1)
Rigor Classic endpoints and CRUD operations are still available by setting the provider's `product` setting to `rigor`
```
provider "synthetics" {
  product = "rigor"
}
```
**NOTE:** The Rigor Classic client expects a valid Rigor Monitoring API token defined as an environment variable named `API_ACCESS_TOKEN` (E.G. `export API_ACCESS_TOKEN="This_is_my_api_token"`)

## Installation

Whenever possible install from the official Terraform Registry:  
https://registry.terraform.io/providers/splunk/synthetics/latest

To install this provider locally follow the directions for installing [In-House Providers](https://www.terraform.io/docs/cloud/run/install-software.html#in-house-providers).

## Examples

see ./examples/ for examples of Splunk Synthetics resources and datasources.
see ./examples/rigor/ for examples of Rigor Classic resources and datasources

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.18