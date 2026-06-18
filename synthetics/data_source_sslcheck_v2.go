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

func dataSourceSslCheckV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSslCheckV2Read,
		Schema: map[string]*schema.Schema{
			"test": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: sslCheckV2DataSourceTestSchema(),
				},
			},
		},
	}
}

func dataSourceSslCheckV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sc2.Client)

	var diags diag.Diagnostics

	checkID := flattenIdData(d.Get("test"))

	check, _, err := c.GetSslCheckV2(checkID)
	if err != nil {
		return diag.FromErr(err)
	}

	checkTest := flattenSslCheckV2Data(check)
	if err := d.Set("test", checkTest); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprint(check.Test.ID))
	return diags
}
