// Copyright 2024 Splunk, Inc.
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

func resourceDowntimeConfigurationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDowntimeConfigurationV2Create,
		ReadContext:   resourceDowntimeConfigurationV2Read,
		UpdateContext: resourceDowntimeConfigurationV2Update,
		DeleteContext: resourceDowntimeConfigurationV2Delete,

		Schema: map[string]*schema.Schema{
			"downtime_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"rule": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"status": {
							Type:     schema.TypeString,
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
						"tests_updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"test_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
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

func resourceDowntimeConfigurationV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	downtimeConfigData := processDowntimeConfigurationV2Items(d)

	o, _, err := c.CreateDowntimeConfigurationV2(&downtimeConfigData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.DowntimeConfiguration.ID))

	resourceDowntimeConfigurationV2Read(ctx, d, meta)

	return diags
	// return nil
}

func resourceDowntimeConfigurationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	downtimeConfigurationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	downtimeConfiguration, r, err := c.GetDowntimeConfigurationV2(downtimeConfigurationID)

	if r.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] Resource exists in state but not in API. Removing resource from state.")
		return diags
	}
	if err != nil {
		log.Println("[WARN] Synthetics API error.", downtimeConfigurationID, err.Error(), r.StatusCode)
		return diag.FromErr(err)
	}
	log.Println("DEBUG] GET downtime_configuration response data: ", downtimeConfiguration)
	if err := d.Set("downtime_configuration", flattenDowntimeConfigurationV2Read(downtimeConfiguration)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDowntimeConfigurationV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	DowntimeConfigurationID := d.Id()

	DowntimeConfigurationIdString, err := strconv.Atoi(DowntimeConfigurationID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.DeleteDowntimeConfigurationV2(DowntimeConfigurationIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Delete downtime_configuration response data: ", resp)
	d.SetId("")

	return diags
}

func resourceDowntimeConfigurationV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	DowntimeConfigurationID := d.Id()

	DowntimeConfigurationData := processDowntimeConfigurationV2Items(d)

	DowntimeConfigurationIdString, err := strconv.Atoi(DowntimeConfigurationID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.UpdateDowntimeConfigurationV2(DowntimeConfigurationIdString, &DowntimeConfigurationData)
	if err != nil {
		log.Println("[ERROR] downtime_configuration failed to update. Dumping request data: ", o)
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Update downtime_configuration response data: ", o)
	return resourceDowntimeConfigurationV2Read(ctx, d, meta)
}

func processDowntimeConfigurationV2Items(d *schema.ResourceData) sc2.DowntimeConfigurationV2Input {

	log.Println("[DEBUG] Process downtime_configuration Resource Data: ", d)

	var downtimeConfig = buildDowntimeConfigurationV2Data(d)

	log.Println("[DEBUG] Processed downtime_configuration Resource Data OUTPUT: ", downtimeConfig)
	return downtimeConfig
}
