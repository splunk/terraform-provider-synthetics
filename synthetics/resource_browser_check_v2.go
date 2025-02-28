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
	"regexp"
	"strconv"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBrowserCheckV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBrowserCheckV2Create,
		ReadContext:   resourceBrowserCheckV2Read,
		UpdateContext: resourceBrowserCheckV2Update,
		DeleteContext: resourceBrowserCheckV2Delete,

		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeList,
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
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "round_robin",
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`(^concurrent$|^round_robin$)`), "Setting must match concurrent or round_robin"),
						},
						"url_protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_url": {
							Type:     schema.TypeString,
							Optional: true,
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
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_agent": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"verify_certificates": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"collect_interactive_metrics": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
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
									"chrome_flags": {
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
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`), "Setting must be a valid domain"),
												},
												"path": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\/`), "Setting must be a valid path starting with /"),
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
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`), "Setting must be a valid domain"),
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
						"transactions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"steps": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Unique steps for the transaction. See official [API documentation](https://dev.splunk.com/observability/reference/api/synthetics_browser/latest#endpoint-createbrowsertest) as the source of truth for descriptions and options for these values.",
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
												"option_selector_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"option_selector": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"variable_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"duration": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"wait_for_nav": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"wait_for_nav_timeout": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.All(validation.IntAtLeast(1), validation.IntAtMost(20000)),
												},
												"wait_for_nav_timeout_default": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"max_wait_time": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.All(validation.IntAtLeast(1), validation.IntAtMost(90000)),
												},
												"max_wait_time_default": {
													Type:     schema.TypeBool,
													Computed: true,
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
						"custom_properties": {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z](\w|-|_){1,128}$`), "custom_properties key must start with a letter and only consist of alphanumeric and underscore characters with no whitespace"),
									},
									"value": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9](\w|-|_){1,128}$`), "custom_properties value can only consist of alphanumeric and underscore characters with no whitespace"),
									},
								},
							},
						},
						"automatic_retries": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
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
	o, _, err := c.CreateBrowserCheckV2(&checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.Test.ID))

	resourceBrowserCheckV2Read(ctx, d, meta)

	return diags
}

func resourceBrowserCheckV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	o, r, err := c.GetBrowserCheckV2(checkID)

	if r.StatusCode == http.StatusNotFound {
		d.SetId("")
		log.Println("[WARN] Resource exists in state but not in API. Removing resource from state.")
		return diags
	}

	if err != nil {
		log.Println("[WARN] Synthetics API error.", checkID, err.Error(), r.StatusCode)
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] GET BROWSER BODY: ", o)
	if err := d.Set("test", flattenBrowserV2Read(o)); err != nil {
		return diag.FromErr(err)
	}

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

	log.Println("[DEBUG] Delete check response data: ", resp)
	d.SetId("")

	return diags
}

func resourceBrowserCheckV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc2.Client)

	checkID := d.Id()

	log.Println("[DEBUG] UPDATE BROWSER CHECK ID: ", checkID)

	checkData := processBrowserCheckV2Items(d)

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.UpdateBrowserCheckV2(checkIdString, &checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] UPDATE BODY: ", o)

	return resourceBrowserCheckV2Read(ctx, d, meta)
}

func processBrowserCheckV2Items(d *schema.ResourceData) sc2.BrowserCheckV2Input {

	var check = buildBrowserV2Data(d)
	log.Println("[DEBUG] BROWSER V2 CHECK OUTPUT: ", check)
	return check
}
