package synthetics

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func dataSourceTotpVariablesV2() *schema.Resource {
	return &schema.Resource{
		Description: "Reads Synthetics TOTP variable metadata. TOTP secrets are not returned by the API and are not exposed by this data source.",
		ReadContext: dataSourceTotpVariablesV2Read,
		Schema: map[string]*schema.Schema{
			"totp_variables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: totpVariableV2ListDataSourceSchema(),
				},
			},
		},
	}
}

func totpVariableV2ListDataSourceSchema() map[string]*schema.Schema {
	s := totpVariableV2DataSourceSchema()
	s["id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	return s
}

func dataSourceTotpVariablesV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sc2.Client)
	var diags diag.Diagnostics

	totpVariables, _, err := c.GetTotpVariablesV2()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("totp_variables", flattenTotpVariablesV2Data(totpVariables.Totps)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("global_totp_variables_synthetics")
	return diags
}
