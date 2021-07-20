package main

import (
	"github.com/greatestusername-splunk/terraform-provider-synthetics/synthetics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return synthetics.Provider()
		},
	})
}
