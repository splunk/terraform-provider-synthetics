# Splunk Synthetics Terraform Provider

[![Release](https://img.shields.io/github/v/release/splunk/terraform-provider-synthetics)](https://github.com/splunk/terraform-provider-synthetics/releases)
[![CI Checks](https://img.shields.io/github/actions/workflow/status/splunk/terraform-provider-synthetics/ci.yml?branch=main&label=CI)](https://github.com/splunk/terraform-provider-synthetics/actions/workflows/ci.yml?query=branch%3Amain)
[![Build](https://img.shields.io/github/actions/workflow/status/splunk/terraform-provider-synthetics/ci.yml?branch=main&job=Build&label=build)](https://github.com/splunk/terraform-provider-synthetics/actions/workflows/ci.yml?query=branch%3Amain)
[![License](https://img.shields.io/github/license/splunk/terraform-provider-synthetics)](https://github.com/splunk/terraform-provider-synthetics/blob/main/LICENSE)

This repository is a **beta** Terraform provider for [Splunk Synthetics in Splunk Observability](https://docs.splunk.com/Observability/synthetics/intro-synthetics.html). It currently contains CRUD operations for API Checks, Real Browser Checks, Port Checks, HTTP Checks, SSL Tests, CA Certificates, and Variables.

**NOTE:** The client expects a valid Splunk Observability API token defined in the provider config (`apikey`) or as an environment variable named `OBSERVABILITY_API_TOKEN` (E.G. `export OBSERVABILITY_API_TOKEN="This_is_my_api_token"`)

## Installation

Whenever possible install from the official Terraform Registry:  
https://registry.terraform.io/providers/splunk/synthetics/latest

To install this provider locally follow the directions for installing [In-House Providers](https://www.terraform.io/docs/cloud/run/install-software.html#in-house-providers).

## Examples

see ./examples/ for examples of Splunk Synthetics resources and datasources.

## Import Existing Tests

Use `terraform import` as normally described in the [Terraform docs](https://developer.hashicorp.com/terraform/cli/import/usage) to bring the resource into your state file. Using the check id number as the identifier.

### Example: Import browser check 496 to state file
```
terraform import synthetics_create_browser_check_v2.browser_v2_foo_check 496
```

To rebuild your configuration file more easily use the datasource for the check in question. This will pull the entire configuration of the check for rebuilding the configuration in your tf files and comparing against a `terraform plan` command.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) — the version pinned in [`.go-version`](./.go-version) (currently Go 1.26.7)

## Contributions
Contributions are welcome and encouraged!

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details on contributing to this repository.

Before your contribution can be accepted, you will be asked to sign our
[Splunk Contributor License Agreement (CLA)](https://github.com/splunk/cla-agreement/blob/main/CLA.md).

To agree to the CLA and COC please comment these in **separate individual messages** on your PR:

CLA:
```
I have read the CLA Document and I hereby sign the CLA
```

Code of Conduct:
```
I have read the Code of Conduct and I hereby accept the Terms
```

### Local development

1. Update your `~/.terraformrc` and add `dev_overrides` config:

```
provider_installation {
  dev_overrides {
    "splunk/synthetics" = "/<path-to-terraform-provider-synthetics-repository>"
  }
  direct {}
}
```

2. Build the provider binary with `make build`. See GNUmakefile for more usefull commands.
3. Use the provider in your testing Terraform config and use `terraform plan/apply/destroy` to make sure your changes work as expected.

``` init.tf
terraform {
  required_providers {
    synthetics = {
      version = "3.0.0"
      source  = "splunk/synthetics"
    }
  }
}

provider "synthetics" {
  realm   = "us1"
  apikey  = "<token>"
}

# your TF resources
resource "synthetics_create_browser_check_v2" "test" { ... }
```

Before opening a pull request, run the same checks CI runs:

```shell
make fmtcheck      # gofmt -l, fails on formatting drift
make vet           # go vet -mod=vendor ./...
make lint          # golangci-lint v2.12.2, requires golangci-lint on PATH
make test          # go test -mod=vendor ./...
make test-race     # go test -mod=vendor -race ./...
make vendor-check  # go mod tidy && go mod vendor, fails on vendor/ drift
make govulncheck   # requires govulncheck on PATH
```

`golangci-lint`, `govulncheck`, and `actionlint` (for linting `.github/workflows/`) are not
vendored; install the versions CI pins with:

```shell
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
go install golang.org/x/vuln/cmd/govulncheck@v1.7.0
go install github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
```
