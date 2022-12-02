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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc "github.com/splunk/syntheticsclient/syntheticsclient"
	sc2 "syntheticsclientv2"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// Add a setting for Realm

			// Add a setting to chose o11y or classic
			"apikey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_ACCESS_TOKEN", nil),
			},
			"realm": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("REALM", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"synthetics_create_http_check":    resourceHttpCheck(),
			"synthetics_create_browser_check": resourceBrowserCheck(),
			"synthetics_create_api_check_v2": resourceApiCheckV2(),
			"synthetics_create_browser_check_v2": resourceBrowserCheckV2(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"synthetics_check": dataSourceCheck(),
			"synthetics_api_v2_check": dataSourceApiCheckV2(),
			"synthetics_browser_v2_check": dataSourceBrowserCheckV2(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("apikey").(string)
	realm := d.Get("realm").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// If realm is defined we'll be using the V2 endpoints from Splunk Observability
	// Otherwise use Rigor Classic endpoints
	if realm != "" {
		if token != "" {
			c := sc2.NewClient(token, realm)
	
			return c, diags
		}
	
		c := sc2.NewClient(token, realm)

		return c, diags
	} else {
		if token != "" {
			c := sc.NewClient(token)
	
			return c, diags
		}
	
		c := sc.NewClient(token)

		return c, diags
	}
}
