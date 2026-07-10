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
	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"
)

func dataSourceClientCertificatesV2() *schema.Resource {
	return &schema.Resource{
		Description: "Reads Synthetics client certificate metadata. Certificate contents and private key passwords are not exposed by this data source. " + clientCertificateAuthMaterialDescription,
		ReadContext: dataSourceClientCertificatesV2Read,
		Schema: map[string]*schema.Schema{
			"client_certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"description": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"domain": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"expires_at": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_at": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"updated_at": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"updated_by": {
						Type:     schema.TypeString,
						Computed: true,
					},
				}},
			},
		},
	}
}

func dataSourceClientCertificatesV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sc2.Client)

	response, _, err := c.GetClientCertificatesV2()
	if err != nil {
		return diag.FromErr(err)
	}

	certificates := make([]interface{}, 0, len(response.Certificates))
	for _, certificate := range response.Certificates {
		certificates = append(certificates, flattenClientCertificateMetadata(certificate))
	}

	if err := d.Set("client_certificates", certificates); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("client-certificates")
	return nil
}
