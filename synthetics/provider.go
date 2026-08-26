// Copyright 2021 Splunk, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package synthetics

import (
	"context"
	"regexp"
	"strings"

	sc2 "github.com/splunk/syntheticsclient/v3/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{

			"product": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "observability",
				Description: "Must be `observability`. Retained for compatibility with existing configurations.",
				Deprecated:  "product is no longer required now that this provider supports only Splunk Observability Synthetics; it will be removed in a future major release.",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^observability$`),
					"product must be observability; this provider supports only Splunk Observability Synthetics",
				),
			},
			"apikey": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Splunk Observability API Key. Will pull from `OBSERVABILITY_API_TOKEN` environment variable if available.",
				DefaultFunc: schema.EnvDefaultFunc("OBSERVABILITY_API_TOKEN", nil),
			},
			"realm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Splunk Observability Realm (E.G. `us1`). Will pull from `REALM` environment variable if available.",
				DefaultFunc: schema.EnvDefaultFunc("REALM", nil),
			},
			"apiurl": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Splunk Observability Realm API Endpoint (E.G. `https://api.<REALM>.signalfx.com`). Will pull from `API_URL` environment variable if available.",
				DefaultFunc: schema.EnvDefaultFunc("API_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"synthetics_create_api_check_v2":              resourceApiCheckV2(),
			"synthetics_create_browser_check_v2":          resourceBrowserCheckV2(),
			"synthetics_create_http_check_v2":             resourceHttpCheckV2(),
			"synthetics_create_port_check_v2":             resourcePortCheckV2(),
			"synthetics_create_ssl_check_v2":              resourceSslCheckV2(),
			"synthetics_create_ca_certificate_v2":         resourceCaCertificateV2(),
			"synthetics_create_client_certificate_v2":     resourceClientCertificateV2(),
			"synthetics_create_totp_variable_v2":          resourceTotpVariableV2(),
			"synthetics_create_variable_v2":               resourceVariableV2(),
			"synthetics_create_location_v2":               resourceLocationV2(),
			"synthetics_create_downtime_configuration_v2": resourceDowntimeConfigurationV2(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"synthetics_api_v2_check":                     dataSourceApiCheckV2(),
			"synthetics_browser_v2_check":                 dataSourceBrowserCheckV2(),
			"synthetics_http_v2_check":                    dataSourceHttpCheckV2(),
			"synthetics_port_v2_check":                    dataSourcePortCheckV2(),
			"synthetics_ssl_v2_check":                     dataSourceSslCheckV2(),
			"synthetics_ca_certificate_v2_check":          dataSourceCaCertificateV2(),
			"synthetics_ca_certificates_v2_check":         dataSourceCaCertificatesV2(),
			"synthetics_client_certificate_v2_check":      dataSourceClientCertificateV2(),
			"synthetics_client_certificates_v2_check":     dataSourceClientCertificatesV2(),
			"synthetics_excluded_file_types_v2_check":     dataSourceExcludedFileTypesV2(),
			"synthetics_totp_variable_v2_check":           dataSourceTotpVariableV2(),
			"synthetics_totp_variables_v2_check":          dataSourceTotpVariablesV2(),
			"synthetics_variable_v2_check":                dataSourceVariableV2(),
			"synthetics_variables_v2_check":               dataSourceVariablesV2(),
			"synthetics_location_v2_check":                dataSourceLocationV2(),
			"synthetics_locations_v2_check":               dataSourceLocationsV2(),
			"synthetics_devices_v2_check":                 dataSourceDevicesV2(),
			"synthetics_downtime_configuration_v2_check":  dataSourceDowntimeConfigurationV2(),
			"synthetics_downtime_configurations_v2_check": dataSourceDowntimeConfigurationsV2(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("apikey").(string)
	realm := d.Get("realm").(string)
	apiurl := d.Get("apiurl").(string)

	var diags diag.Diagnostics

	if token != "" && realm != "" && apiurl != "" {
		args := sc2.NewClientArgs(
			30,
			strings.TrimSuffix(apiurl, "/")+"/v2/synthetics",
		)
		c := sc2.NewConfigurableClient(token, realm, args)
		return c, diags
	}

	// Default client (no apiurl override)
	c := sc2.NewClient(token, realm)

	return c, diags
}
