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

	sc2 "github.com/splunk/syntheticsclient/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sc "github.com/splunk/syntheticsclient/syntheticsclient"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{

			"product": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "One of: `observability` or `rigor`",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`(^observability$|^rigor$)`), "product setting must match either observability or rigor (v1.0.0+)"),
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
				Description: "Splunk Observability Realm (E.G. `us1`). Will pull from `REALM` environment variable if available. For Rigor use realm rigor",
				DefaultFunc: schema.EnvDefaultFunc("REALM", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"synthetics_create_http_check":       resourceHttpCheck(),
			"synthetics_create_browser_check":    resourceBrowserCheck(),
			"synthetics_create_api_check_v2":     resourceApiCheckV2(),
			"synthetics_create_browser_check_v2": resourceBrowserCheckV2(),
			"synthetics_create_http_check_v2":    resourceHttpCheckV2(),
			"synthetics_create_port_check_v2":    resourcePortCheckV2(),
			"synthetics_create_variable_v2":      resourceVariableV2(),
			"synthetics_create_location_v2":      resourceLocationV2(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"synthetics_check":              dataSourceCheck(),
			"synthetics_api_v2_check":       dataSourceApiCheckV2(),
			"synthetics_browser_v2_check":   dataSourceBrowserCheckV2(),
			"synthetics_http_v2_check":      dataSourceHttpCheckV2(),
			"synthetics_port_v2_check":      dataSourcePortCheckV2(),
			"synthetics_variable_v2_check":  dataSourceVariableV2(),
			"synthetics_variables_v2_check": dataSourceVariablesV2(),
			"synthetics_location_v2_check":  dataSourceLocationV2(),
			"synthetics_locations_v2_check": dataSourceLocationsV2(),
			"synthetics_devices_v2_check":   dataSourceDevicesV2(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("apikey").(string)
	realm := d.Get("realm").(string)
	product := d.Get("product").(string)

	var diags diag.Diagnostics

	if product == "observability" {
		if token != "" && realm != "" {
			c := sc2.NewClient(token, realm)

			return c, diags
		}

		c := sc2.NewClient(token, realm)

		return c, diags
	} else {
		if product == "rigor" && token != "" {
			c := sc.NewClient(token)

			return c, diags
		}

		c := sc.NewClient(token)

		return c, diags
	}
}
