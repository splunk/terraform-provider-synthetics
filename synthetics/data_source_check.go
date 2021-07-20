package synthetics

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc "github.com/splunk/syntheticsclient/syntheticsclient"
)

func dataSourceCheck() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChecksRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"frequency": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"paused": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"muted": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"links": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"self": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"self_html": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"metrics": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"last_run": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_code": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"last_message": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"last_response_time": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"last_run_at": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"last_failure_at": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"last_alert_at": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"has_failure": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"has_location_failure": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"viewport": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"height": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"dns_overrides": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"original_domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"original_host": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"browser": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "chrome",
						},
					},
				},
			},
			"steps": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"item_method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"how": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"what": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"variable_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"position": {
							Type:     schema.TypeInt,
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
			"javascript_files": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"threshold_monitors": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matcher": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								switch v {
								case
									"first_byte_time_ms",
									"dom_interactive_time_ms",
									"dom_load_time_ms",
									"dom_complete_time_ms",
									"start_render_ms",
									"onload_time_ms",
									"visually_complete_ms",
									"fully_loaded_time_ms",
									"first_paint_time_ms",
									"first_contentful_paint_time_ms",
									"first_meaningful_paint_time_ms",
									"first_interactive_time_ms",
									"first_cpu_idle_time_ms",
									"first_request_dns_time_ms",
									"first_request_connect_time_ms",
									"first_request_ssl_time_ms",
									"first_request_send_time_ms",
									"first_request_wait_time_ms",
									"first_request_receive_time_ms",
									"speed_index",
									"requests",
									"content_bytes",
									"html_files",
									"html_bytes",
									"image_files",
									"image_bytes",
									"javascript_files",
									"javascript_bytes",
									"css_files",
									"css_bytes",
									"video_files",
									"video_bytes",
									"font_files",
									"font_bytes",
									"other_files",
									"other_bytes",
									"client_errors",
									"connection_errors",
									"server_errors",
									"errors",
									"run_count",
									"success_count",
									"failure_count",
									"lighthouse_performance_score",
									"availability",
									"downtime",
									"total_blocking_time_ms",
									"largest_contentful_paint_time_ms",
									"cumulative_layout_shift":
									return
								}
								errs = append(errs, fmt.Errorf("%s is not a valid metric_name. Please check your input against API docs at https://monitoring-api.rigor.com", v))
								return
							},
						},
						"comparison_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								switch v {
								case
									"less_than",
									"equals",
									"greater_than":
									return
								}
								errs = append(errs, fmt.Errorf("%s is not a valid comparison_type. Please check your input against API docs at https://monitoring-api.rigor.com", v))
								return
							},
						},
						"value": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"excluded_files": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"preset_name": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								switch v {
								case
									"chartbeat",
									"clicktale",
									"comscore",
									"coremetrics",
									"crazy-egg",
									"eloqua",
									"gomez",
									"google-analytics",
									"hubspot",
									"liveperson",
									"mixpanel",
									"omniture",
									"optimizely",
									"pardot",
									"quantcast",
									"spectate",
									"tealium",
									"white-ops":
									return
								}
								errs = append(errs, fmt.Errorf("%s is not a valid preset_name. Please check your input against API docs at https://monitoring-api.rigor.com", v))
								return
							},
						},
						"exclusion_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								switch v {
								case
									"preset",
									"custom",
									"all_except":
									return
								}
								errs = append(errs, fmt.Errorf("%s is not a valid exclusion_type. Please check your input against API docs at https://monitoring-api.rigor.com", v))
								return
							},
						},
					},
				},
			},
			"notifications": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sms": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"call": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"email": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"notify_after_failure_count": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"notify_on_location_failure": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"muted": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"notify_who": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sms": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"call": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"email": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"custom_user_email": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"links": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"self_html": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
										Optional: true,
									},
								},
							},
						},
						"notification_windows": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"duration_in_minutes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"escalations": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sms": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"call": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"email": {
										Type:     schema.TypeBool,
										Computed: true,
										Optional: true,
									},
									"after_minutes": {
										Type:     schema.TypeInt,
										Computed: true,
										Optional: true,
									},
									"is_repeat": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"notify_who": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"custom_user_email": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"links": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"self_html": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"notification_window": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"end_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"duration_in_minutes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"time_zone": {
													Type:     schema.TypeString,
													Computed: true,
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
			"response_time_monitor_milliseconds": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"http_request_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_agent": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"round_robin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_retry": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"integrations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_request_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"blackout_periods": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repeat_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration_in_minutes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_repeat": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"monthly_repeat_type": {
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
					},
				},
			},
			"locations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"world_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"success_criteria": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_type": {
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
						"comparison_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"check_connection": {
				// `connection` is a reserved field name in Terraform
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"download_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"upload_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"latency": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"packet_loss": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceChecksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sc.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID := d.Get("id").(int)

	check, _, err := c.GetCheck(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", check.Name)
	d.Set("frequency", check.Frequency)
	d.Set("muted", check.Muted)
	d.Set("paused", check.Paused)
	d.Set("created_at", check.CreatedAt)
	d.Set("type", check.Type)
	d.Set("updated_at", check.UpdatedAt)
	d.Set("response_time_monitor_milliseconds", check.ResponseTimeMonitorMilliseconds)
	d.Set("round_robin", check.RoundRobin)
	d.Set("auto_retry", check.AutoRetry)
	d.Set("enabled", check.Enabled)
	d.Set("url", check.URL)
	d.Set("http_request_body", check.HTTPRequestBody)
	d.Set("http_method", check.HTTPMethod)

	checkLinks := flattenLinkData(&check.Links)
	if err := d.Set("links", checkLinks); err != nil {
		return diag.FromErr(err)
	}

	checkStatus := flattenStatusData(&check.Status)
	if err := d.Set("status", checkStatus); err != nil {
		return diag.FromErr(err)
	}

	checkTags := flattenTagsData(&check.Tags)
	if err := d.Set("tags", checkTags); err != nil {
		return diag.FromErr(err)
	}

	checkBlackout := flattenBlackoutData(&check.BlackoutPeriods)
	if err := d.Set("blackout_periods", checkBlackout); err != nil {
		return diag.FromErr(err)
	}

	checkNotifications := flattenNotificationsData(&check.Notifications)
	if err := d.Set("notifications", checkNotifications); err != nil {
		return diag.FromErr(err)
	}

	checkConnection := flattenConnectionData(&check.Connection)
	if err := d.Set("check_connection", checkConnection); err != nil {
		return diag.FromErr(err)
	}

	checkIntegrations := flattenIntegrationsData(&check.Integrations)
	if err := d.Set("integrations", checkIntegrations); err != nil {
		return diag.FromErr(err)
	}

	checkLocations := flattenLocationsData(&check.Locations)
	if err := d.Set("locations", checkLocations); err != nil {
		return diag.FromErr(err)
	}

	checkSuccessCriteria := flattenSuccessCriteriaData(&check.SuccessCriteria)
	if err := d.Set("success_criteria", checkSuccessCriteria); err != nil {
		return diag.FromErr(err)
	}

	checkSteps := flattenStepData(check.Steps)
	if err := d.Set("steps", checkSteps); err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprint(check.ID)
	d.SetId(id)
	return diags
}
