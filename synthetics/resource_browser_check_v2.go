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
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc2 "syntheticsclientv2"
)

func resourceBrowserCheckV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBrowserCheckV2Create,
		ReadContext:   resourceBrowserCheckV2Read,
		UpdateContext: resourceBrowserCheckV2Update,
		DeleteContext: resourceBrowserCheckV2Delete,

		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"frequency": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"scheduling_strategy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"location_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"device_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"advanced_settings": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_agent": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"verify_certificates": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"authentication": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"username": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"password": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"cookies": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"domain": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"path": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"headers": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"domain": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"host_overrides": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"target": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"keep_host_header": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"business_transactions": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"steps": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"url": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"action": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"selector_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"selector": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"wait_for_nav": {
													Type:     schema.TypeBool,
													Required: true,
												},
												"options": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"url": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},
								},
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

func resourceBrowserCheckV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkData := processBrowserCheckV2Items(d)

	o, req, err := c.CreateBrowserCheckV2(&checkData)
	log.Printf("[WARN] ^^^^^^^^^^^^^^^^ CREATE REQUEST BODY JSON^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	log.Println(o)
	log.Println(req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.Test.ID))

	resourceBrowserCheckV2Read(ctx, d, meta)

	return diags
	// return nil
}

func resourceBrowserCheckV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	check, req, err := c.GetBrowserCheckV2(checkID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[WARN] ^^^^^^^^^^^^^^^^GET REQUEST BODY JSON^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	log.Println(req)
	log.Printf("[DEBUG] G***************************************G: ")
	log.Println("[DEBUG] GET check response data: ", check)
	log.Printf("[DEBUG] #########################################: ")
	log.Printf("[DEBUG] $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$: ")

	return diags
}

func resourceBrowserCheckV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	var diags diag.Diagnostics

	checkID := d.Id()

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.DeleteBrowserCheckV2(checkIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] #########################################: ")
	log.Printf("[DEBUG] $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$: ")
	log.Println("[DEBUG] Delete check response data: ", resp)
	d.SetId("")

	return diags
}

func resourceBrowserCheckV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	checkID := d.Id()

	log.Printf("[WARN] ^^^^^^^^^^^^^^^^UPDATE!! CHECK ID!!!^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	log.Println(checkID)

	checkData := processBrowserCheckV2Items(d)

	log.Printf("[WARN] ^^^^^^^^^^^^^^^^UPDATE!! CHECK DATA#$#^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	log.Println(checkID)

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, req, err := c.UpdateBrowserCheckV2(checkIdString, &checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[WARN] ^^^^^^^^^^^^^^^^UPDATE BODY JSON^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	log.Println(req)

	log.Printf("[DEBUG] #########################################: ")
	log.Printf("[DEBUG] $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$: ")
	log.Println("[DEBUG] Update check response data: ", o)
	d.Set("test.updated_at", time.Now().Format(time.RFC850))

	return resourceBrowserCheckV2Read(ctx, d, meta)
}

func processBrowserCheckV2Items(d *schema.ResourceData) sc2.BrowserCheckV2Input {

	log.Printf("[WARN] *****&&*&  PRE OUTPUT ****************")
	log.Println(d)
	//These MUST exist or the request will fail
	var check = buildBrowserV2Data(d)
	// check.Test.Active = d.Get("test.active").(bool)
	// check.Test.Deviceid = d.Get("test.device_id").(int)
	// check.Test.Frequency = d.Get("test.frequency").(int)
	// check.Test.Locationids = buildLocationIdData(d)   //Start tomorry
	// check.Test.Name = d.Get("test.name").(string)
	// //check.Test.Requests = buildRequestsData(d)
	// check.Test.Schedulingstrategy = d.Get("test.scheduling_strategy").(string)
	log.Printf("[WARN] *****&&*& CHECK OUTPUT ****************")
	log.Println(check)
	return check
}
