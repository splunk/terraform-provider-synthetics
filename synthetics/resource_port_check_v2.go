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
	"regexp"
	"strconv"
	"time"

	sc2 "syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePortCheckV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePortCheckV2Create,
		ReadContext:   resourcePortCheckV2Read,
		UpdateContext: resourcePortCheckV2Update,
		DeleteContext: resourcePortCheckV2Delete,

		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "port",
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`(^tcp$|^udp$)`), "Setting must match tcp or udp"),
						},
						"host": {
							Type:     schema.TypeString,
							Required: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"frequency": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"scheduling_strategy": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "round_robin",
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`(^concurrent$|^round_robin$)`), "Setting must match concurrent or round_robin"),
						},
						"location_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func resourcePortCheckV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkData := processPortCheckV2Items(d)

	o, _, err := c.CreatePortCheckV2(&checkData)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.Test.ID))

	resourcePortCheckV2Read(ctx, d, meta)

	return diags
}

func resourcePortCheckV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.GetPortCheckV2(checkID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] GET PORT BODY: ", o)

	return diags
}

func resourcePortCheckV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	checkID := d.Id()

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.DeletePortCheckV2(checkIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Delete check response data: ", resp)
	d.SetId("")

	return diags
}

func resourcePortCheckV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	checkID := d.Id()

	checkData := processPortCheckV2Items(d)

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.UpdatePortCheckV2(checkIdString, &checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] UPDATE BODY: ", o)

	d.Set("test.updated_at", time.Now().Format(time.RFC850))

	return resourcePortCheckV2Read(ctx, d, meta)
}

func processPortCheckV2Items(d *schema.ResourceData) sc2.PortCheckV2Input {

	var check = buildPortCheckV2Data(d)
	log.Println("[DEBUG] PORT V2 CHECK OUTPUT: ", check)
	return check
}