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

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCaCertificatesV2() *schema.Resource {
	return &schema.Resource{
		Description: "Reads Synthetics CA certificates. CA certificate content is read into Terraform state even when marked sensitive. Use encrypted, access-controlled remote state and do not commit private CA material to source control.",
		ReadContext: dataSourceCaCertificatesV2Read,
		Schema: map[string]*schema.Schema{
			"ca_certificates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: caCertificateV2ListDataSourceSchema(),
				},
			},
		},
	}
}

func dataSourceCaCertificatesV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sc2.Client)

	var diags diag.Diagnostics

	caCertificates, _, err := c.GetCaCertificatesV2()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ca_certificates", flattenCaCertificatesV2Data(caCertificates.CaCerts)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("global_ca_certificates_synthetics")
	return diags
}
