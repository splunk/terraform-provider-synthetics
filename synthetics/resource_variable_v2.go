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

func resourceVariableV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVariableV2Create,
		ReadContext:   resourceVariableV2Read,
		UpdateContext: resourceVariableV2Update,
		DeleteContext: resourceVariableV2Delete,

		Schema: map[string]*schema.Schema{
			"variable": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"secret": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceVariableV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkData := processVariableV2Items(d)

	o, _, err := c.CreateVariableV2(&checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.Variable.ID))

	resourceVariableV2Read(ctx, d, meta)

	return diags
	// return nil
}

func resourceVariableV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	variableID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	variable, r, err := c.GetVariableV2(variableID)

	if r.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] Resource exists in state but not in API. Removing resource from state.")
		return diags
	}
	if err != nil {
		log.Println("[WARN] Synthetics API error.", variableID, err.Error(), r.StatusCode)
		return diag.FromErr(err)
	}
	log.Println("DEBUG] GET variable response data: ", variable)
	if err := d.Set("variable", flattenVariableV2Read(variable)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceVariableV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	variableID := d.Id()

	variableIdString, err := strconv.Atoi(variableID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.DeleteVariableV2(variableIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Delete variable response data: ", resp)
	d.SetId("")

	return diags
}

func resourceVariableV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	variableID := d.Id()

	variableData := processVariableV2Items(d)

	variableIdString, err := strconv.Atoi(variableID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.UpdateVariableV2(variableIdString, &variableData)
	if err != nil {
		log.Println("[ERROR] Variable failed to update. Dumping request data: ", o)
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Update variable response data: ", o)
	return resourceVariableV2Read(ctx, d, meta)
}

func processVariableV2Items(d *schema.ResourceData) sc2.VariableV2Input {

	log.Println("[DEBUG] Process Variable Resource Data: ", d)

	var check = buildVariableV2Data(d)

	log.Println("[DEBUG] Processed Variable Resource Data OUTPUT: ", check)
	return check
}
