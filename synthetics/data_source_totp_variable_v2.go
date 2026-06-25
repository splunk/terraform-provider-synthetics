package synthetics

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func dataSourceTotpVariableV2() *schema.Resource {
	return &schema.Resource{
		Description: "Reads Synthetics TOTP variable metadata. The TOTP secret is not returned by the API and is not exposed by this data source.",
		ReadContext: dataSourceTotpVariableV2Read,
		Schema: map[string]*schema.Schema{
			"totp_variable": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: totpVariableV2DataSourceSchema(),
				},
			},
		},
	}
}

func totpVariableV2DataSourceSchema() map[string]*schema.Schema {
	s := totpVariableV2ResourceSchema()
	s["id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	s["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	s["description"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	delete(s, "secret")
	s["digits"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s["interval"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s["hmac_digest"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	return s
}

func dataSourceTotpVariableV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sc2.Client)
	var diags diag.Diagnostics

	totpVariableID := totpVariableIDFromList(d.Get("totp_variable"))
	totpVariable, _, err := c.GetTotpVariableV2(totpVariableID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("totp_variable", flattenTotpVariableV2Data(totpVariable)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(totpVariable.Totp.ID))
	return diags
}
