package synthetics

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc "github.com/splunk/syntheticsclient/syntheticsclient"
)

func resourceHttpCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHttpCheckCreate,
		ReadContext:   resourceHttpCheckRead,
		UpdateContext: resourceHttpCheckUpdate,
		DeleteContext: resourceHttpCheckDelete,

		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frequency": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notifications": {
				Type:     schema.TypeSet,
				Optional: true,
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
						"notify_after_failure_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"notify_on_location_failure": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"muted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"notify_who": {
							Type:     schema.TypeSet,
							Optional: true,
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
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"duration_in_minutes": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"escalations": {
							Type:     schema.TypeSet,
							Optional: true,
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
									"after_minutes": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"is_repeat": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"notify_who": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeInt,
													Optional: true,
												},
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
											},
										},
									},
									"notification_window": {
										Type:     schema.TypeSet,
										Computed: true,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_time": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"end_time": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"duration_in_minutes": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"time_zone": {
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
			"http_request_headers": {
				Type:     schema.TypeList,
				Optional: true,
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
				Default:  true,
				Optional: true,
			},
			"auto_retry": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"integrations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"http_request_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"blackout_periods": {
				Type:     schema.TypeSet,
				Optional: true,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"success_criteria": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"comparison_string": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"check_connection": {
				// `connection` is a reserved field name in Terraform
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"download_bandwidth": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"upload_bandwidth": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"latency": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"packet_loss": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceHttpCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	check := sc.HttpCheckInput{}

	checkData := processHttpCheckItems(d, check)

	o, _, err := c.CreateHttpCheck(&checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(o.ID))

	resourceHttpCheckRead(ctx, d, meta)

	return diags
	// return nil
}

func resourceHttpCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	checkID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	check, _, err := c.GetHttpCheck(checkID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] GET check response data: ", check)

	return diags
}

func resourceHttpCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc.Client)

	var diags diag.Diagnostics

	checkID := d.Id()

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := c.DeleteHttpCheck(checkIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Delete check response data: ", resp)
	d.SetId("")

	return diags
}

func resourceHttpCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*sc.Client)

	checkID := d.Id()

	check := sc.HttpCheckInput{}

	checkData := processHttpCheckItems(d, check)

	checkIdString, err := strconv.Atoi(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	o, _, err := c.UpdateHttpCheck(checkIdString, &checkData)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] Update check response data: ", o)
	d.Set("last_updated", time.Now().Format(time.RFC850))

	// return nil
	return resourceHttpCheckRead(ctx, d, meta)
}

func processHttpCheckItems(d *schema.ResourceData, check sc.HttpCheckInput) sc.HttpCheckInput {

	//These MUST exist of the request will fail
	check.Name = d.Get("name").(string)
	check.Frequency = d.Get("frequency").(int)
	check.URL = d.Get("url").(string)

	check.AutoRetry = d.Get("auto_retry").(bool)
	check.Enabled = d.Get("enabled").(bool)
	check.HTTPRequestBody = d.Get("http_request_body").(string)
	check.HTTPMethod = d.Get("http_method").(string)
	check.RoundRobin = d.Get("round_robin").(bool)

	check.SuccessCriteria = buildSuccessCriteriaData(d)
	check.Connection = buildConnectionData(d)
	check.Tags = buildTagsData(d)
	check.Locations = buildLocationsData(d)
	check.Integrations = buildIntegrationsData(d)
	check.Notifications = buildNotificationsData(check.Notifications, d)

	return check
}
