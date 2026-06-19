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
	"log"
	"net/http"
	"strconv"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCaCertificateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCaCertificateV2Create,
		ReadContext:   resourceCaCertificateV2Read,
		UpdateContext: resourceCaCertificateV2Update,
		DeleteContext: resourceCaCertificateV2Delete,

		Schema: map[string]*schema.Schema{
			"ca_certificate": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: caCertificateV2ResourceSchema(),
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func caCertificateV2ResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"content": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"file_extension": {
			Type:     schema.TypeString,
			Required: true,
		},
		"filename": {
			Type:     schema.TypeString,
			Required: true,
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
	}
}

func caCertificateV2DataSourceSchema() map[string]*schema.Schema {
	s := caCertificateV2ResourceSchema()
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
	s["content"] = &schema.Schema{
		Type:      schema.TypeString,
		Computed:  true,
		Sensitive: true,
	}
	s["file_extension"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	s["filename"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	return s
}

func caCertificateV2ListDataSourceSchema() map[string]*schema.Schema {
	s := caCertificateV2DataSourceSchema()
	s["id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	return s
}

func resourceCaCertificateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	caCertificateData, err := buildCaCertificateV2Data(d)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.CreateCaCertificateV2(&caCertificateData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.CaCert.ID))

	return append(diags, resourceCaCertificateV2Read(ctx, d, meta)...)
}

func resourceCaCertificateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	caCertificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	existingContent := caCertificateContentFromState(d)
	o, r, err := c.GetCaCertificateV2(caCertificateID)
	if r != nil && r.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] CA certificate exists in state but not in API. Removing resource from state.")
		return diags
	}
	if err != nil {
		statusCode := 0
		if r != nil {
			statusCode = r.StatusCode
		}
		log.Println("[WARN] Synthetics API error for CA certificate.", caCertificateID, err.Error(), statusCode)
		return diag.FromErr(err)
	}
	if err := d.Set("ca_certificate", flattenCaCertificateV2Read(o, existingContent)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCaCertificateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	caCertificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	caCertificateData := buildCaCertificateV2UpdateData(d)
	_, _, err = c.UpdateCaCertificateV2(caCertificateID, &caCertificateData)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCaCertificateV2Read(ctx, d, meta)
}

func resourceCaCertificateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	caCertificateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.DeleteCaCertificateV2(caCertificateID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
