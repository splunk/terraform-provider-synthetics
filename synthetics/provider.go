package synthetics

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc "github.com/splunk/syntheticsclient/syntheticsclient"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_ACCESS_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"synthetics_create_http_check":    resourceHttpCheck(),
			"synthetics_create_browser_check": resourceBrowserCheck(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"synthetics_check": dataSourceCheck(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("apikey").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if token != "" {
		c := sc.NewClient(token)

		return c, diags
	}

	c := sc.NewClient(token)

	return c, diags
}
