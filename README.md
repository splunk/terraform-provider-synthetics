# Splunk Synthetics Terraform Provider

This repository is an **alpha* Terraform provider for [Splunk Synthetics (formerly Rigor™)](https://monitoring.rigor.com/). It currently contains CRUD operations for HTTP (Uptime) Checks and Real Browser Checks with some caveats:

 - Currently Real Browser Checks cannot have `steps` or `javascript_files` added via Public API and thus are not included in this provider.
 - Integrations are not managed by this provider and must be setup in the UI and referenced with their ID number.
 - Private Locations are not managed by this provider and must be setup in the UI and referenced with their ID number.
 - Excluding custom files is currently not supported. All preset file exclusions are included and working.
 
This repo and the companion [Synthetics Golang client](https://github.com/splunk/syntheticsclient) are not DRY and are specifically verbose for code auditing and teaching reasons.     

**NOTE:** The client expects a valid Synthetics API token defined as an environment variable named `API_ACCESS_TOKEN` (E.G. `export API_ACCESS_TOKEN="This_is_my_api_token"`)

## Installation

Whenever possible install from the official Terraform Registry:  
https://registry.terraform.io/providers/splunk/synthetics/latest

To install this provider locally follow the directions for installing [In-House Providers](https://www.terraform.io/docs/cloud/run/install-software.html#in-house-providers).

## Examples

see ./examples/ for current examples of HTTP and Browser Checks

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15