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
	"fmt"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBrowserCheckV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrowserCheckV2Read,
		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"frequency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scheduling_strategy": {
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
						"location_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"advanced_settings": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication": {
										Type:     schema.TypeSet,
										Computed: true,
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
										Computed: true,
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
										Computed: true,
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
										Computed: true,
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
									"user_agent": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"verify_certificates": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"business_transactions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"steps": {
										Type:     schema.TypeSet,
										Computed: true,
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
												"wait_for_nav": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"options": {
													Type:     schema.TypeSet,
													Computed: true,
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
						"transactions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"steps": {
										Type:     schema.TypeSet,
										Computed: true,
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
												"wait_for_nav": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"options": {
													Type:     schema.TypeSet,
													Computed: true,
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
						"device": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"label": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_agent": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"network_connection": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"download_bandwidth": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"latency": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"packet_loss": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"upload_bandwidth": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"viewport_height": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"viewport_width": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBrowserCheckV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*sc2.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID := flattenIdData(d.Get("test"))

	check, _, err := c.GetBrowserCheckV2(checkID)
	println(check)
	if err != nil {
		return diag.FromErr(err)
	}

	checkTest := flattenBrowserV2Data(check)
	if err := d.Set("test", checkTest); err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprint(check.Test.ID)
	d.SetId(id)
	return diags
}
