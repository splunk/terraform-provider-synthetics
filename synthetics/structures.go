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
	"fmt"
	"log"
	"strings"
	"time"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/gocty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sc "github.com/splunk/syntheticsclient/syntheticsclient"
)

func flattenIdData(test interface{}) int {

	test_schema := test.(*schema.Set)

	test_list := test_schema.List()
	test_map := test_list[0].(map[string]interface{})
	id := test_map["id"]
	return id.(int)
}

func flattenStringIdData(test interface{}) string {

	test_schema := test.(*schema.Set)

	test_list := test_schema.List()
	test_map := test_list[0].(map[string]interface{})
	id := test_map["id"]
	return id.(string)
}

func flattenApiV2Read(checkApiV2 *sc2.ApiCheckV2Response) []interface{} {
	apiV2 := make(map[string]interface{})

	apiV2["active"] = checkApiV2.Test.Active
	apiV2["automatic_retries"] = checkApiV2.Test.Automaticretries

	if checkApiV2.Test.Frequency != 0 {
		apiV2["frequency"] = checkApiV2.Test.Frequency
	}

	if checkApiV2.Test.Name != "" {
		apiV2["name"] = checkApiV2.Test.Name
	}

	if checkApiV2.Test.Schedulingstrategy != "" {
		apiV2["scheduling_strategy"] = checkApiV2.Test.Schedulingstrategy
	}

	apiV2["device_id"] = checkApiV2.Test.Deviceid

	locationIds := flattenLocationData(&checkApiV2.Test.Locationids)
	apiV2["location_ids"] = locationIds

	requests := flattenRequestData(&checkApiV2.Test.Requests)
	apiV2["requests"] = requests

	customProperties := flattenCustomProperties(&checkApiV2.Test.Customproperties)
	apiV2["custom_properties"] = customProperties

	log.Println("[DEBUG] apiv2 data: ", apiV2)

	return []interface{}{apiV2}
}

func findDeviceByID(devices []sc2.Device, deviceID int) *sc2.Device {
	for i := range devices {
		if devices[i].ID == deviceID {
			return &devices[i]
		}
	}
	return nil
}

func flattenDeviceFromID(deviceID int, devices []sc2.Device) []interface{} {
	if device := findDeviceByID(devices, deviceID); device != nil {
		return flattenDeviceData(device)
	}
	if deviceID != 0 {
		return flattenDeviceData(&sc2.Device{ID: deviceID})
	}
	return []interface{}{}
}

const browserCheckV2MaxSelectors = 10

func selectorsFromFields(selectorType, selector string) []sc2.Selector {
	if selectorType == "" || selector == "" {
		return nil
	}
	return []sc2.Selector{{Type: selectorType, Value: selector}}
}

// stepSelectorInput holds legacy and selectors-block fields for one browser step.
type stepSelectorInput struct {
	selectorType string
	selector     string
	selectors    []sc2.Selector
}

func parseSelectorsList(raw interface{}) []sc2.Selector {
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	out := make([]sc2.Selector, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		selType := stepStringField(m, "type")
		selValue := stepStringField(m, "value")
		if selType == "" || selValue == "" {
			continue
		}
		out = append(out, sc2.Selector{Type: selType, Value: selValue})
	}
	return out
}

// resolveSingleSelector returns the single selector implied by the step, using the
// same precedence as buildSelectorsFromStep. ok is false when there are multiple
// selectors or no selector at all.
func (in stepSelectorInput) resolveSingleSelector() (selectorType, selector string, ok bool) {
	if len(in.selectors) > 1 {
		return "", "", false
	}
	if len(in.selectors) == 1 {
		return in.selectors[0].Type, in.selectors[0].Value, true
	}
	if in.selectorType != "" && in.selector != "" {
		return in.selectorType, in.selector, true
	}
	return "", "", false
}

func stepSelectorInputsEquivalent(a, b stepSelectorInput) bool {
	aType, aVal, aOK := a.resolveSingleSelector()
	bType, bVal, bOK := b.resolveSingleSelector()
	return aOK && bOK && aType == bType && aVal == bVal
}

// stepSelectorRepresentationDiffers is true when state and config encode the same
// single selector using different Terraform fields (legacy vs selectors block).
func stepSelectorRepresentationDiffers(a, b stepSelectorInput) bool {
	aUsesList := len(a.selectors) > 0
	bUsesList := len(b.selectors) > 0
	aUsesLegacy := a.selector != "" || a.selectorType != ""
	bUsesLegacy := b.selector != "" || b.selectorType != ""
	return aUsesList != bUsesList || aUsesLegacy != bUsesLegacy
}

// migratingFromLegacyToSelectors is true when state still has legacy fields and
// config has moved to a single selectors block (the one-time migration shape).
func migratingFromLegacyToSelectors(state, config stepSelectorInput) bool {
	stateUsesLegacy := state.selector != "" || state.selectorType != ""
	stateUsesList := len(state.selectors) > 0
	configUsesList := len(config.selectors) == 1
	configUsesLegacy := config.selector != "" || config.selectorType != ""
	return stateUsesLegacy && !stateUsesList && configUsesList && !configUsesLegacy
}

func flattenSelectorsData(selectors []sc2.Selector) []interface{} {
	if len(selectors) == 0 {
		return nil
	}
	cls := make([]interface{}, len(selectors))
	for i, sel := range selectors {
		cls[i] = map[string]interface{}{
			"type":  sel.Type,
			"value": sel.Value,
		}
	}
	return cls
}

func stepStringField(step map[string]interface{}, key string) string {
	if v, ok := step[key].(string); ok {
		return v
	}
	return ""
}

func buildSelectorsFromStep(step map[string]interface{}) ([]sc2.Selector, error) {
	if raw, ok := step["selectors"]; ok && raw != nil {
		list, ok := raw.([]interface{})
		if ok && len(list) > 0 {
			if len(list) > browserCheckV2MaxSelectors {
				return nil, fmt.Errorf(
					"step %q has %d selectors; maximum is %d",
					stepStringField(step, "name"),
					len(list),
					browserCheckV2MaxSelectors,
				)
			}
			result := make([]sc2.Selector, len(list))
			for i, item := range list {
				m, ok := item.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("step %q: invalid selector at index %d", stepStringField(step, "name"), i)
				}
				selType := stepStringField(m, "type")
				selValue := stepStringField(m, "value")
				if selType == "" || selValue == "" {
					return nil, fmt.Errorf("step %q: selector at index %d requires type and value", stepStringField(step, "name"), i)
				}
				result[i] = sc2.Selector{Type: selType, Value: selValue}
			}
			return result, nil
		}
	}
	return selectorsFromFields(
		stepStringField(step, "selector_type"),
		stepStringField(step, "selector"),
	), nil
}

func flattenApiV2Data(checkApiV2 *sc2.ApiCheckV2Response, devices []sc2.Device) []interface{} {
	apiV2 := make(map[string]interface{})

	apiV2["active"] = checkApiV2.Test.Active
	apiV2["automatic_retries"] = checkApiV2.Test.Automaticretries

	if checkApiV2.Test.Createdat.IsZero() {
	} else {
		apiV2["created_at"] = checkApiV2.Test.Createdat.String()
	}

	if checkApiV2.Test.Updatedat.IsZero() {
	} else {
		apiV2["updated_at"] = checkApiV2.Test.Updatedat.String()
	}

	if checkApiV2.Test.Createdby != "" {
		apiV2["created_by"] = checkApiV2.Test.Createdby
	}

	if checkApiV2.Test.Updatedby != "" {
		apiV2["updated_by"] = checkApiV2.Test.Updatedby
	}

	if checkApiV2.Test.Lastrunat.IsZero() {
	} else {
		apiV2["last_run_at"] = checkApiV2.Test.Lastrunat.String()
	}

	if checkApiV2.Test.Lastrunstatus != "" {
		apiV2["last_run_status"] = checkApiV2.Test.Lastrunstatus
	}

	if checkApiV2.Test.Frequency != 0 {
		apiV2["frequency"] = checkApiV2.Test.Frequency
	}

	if checkApiV2.Test.ID != 0 {
		apiV2["id"] = checkApiV2.Test.ID
	}

	if checkApiV2.Test.Name != "" {
		apiV2["name"] = checkApiV2.Test.Name
	}

	if checkApiV2.Test.Schedulingstrategy != "" {
		apiV2["scheduling_strategy"] = checkApiV2.Test.Schedulingstrategy
	}

	if checkApiV2.Test.Type != "" {
		apiV2["type"] = checkApiV2.Test.Type
	}

	apiV2["device"] = flattenDeviceFromID(checkApiV2.Test.Deviceid, devices)

	locationIds := flattenLocationData(&checkApiV2.Test.Locationids)
	apiV2["location_ids"] = locationIds

	requests := flattenRequestData(&checkApiV2.Test.Requests)
	apiV2["requests"] = requests

	customProperties := flattenCustomProperties(&checkApiV2.Test.Customproperties)
	apiV2["custom_properties"] = customProperties

	log.Println("[DEBUG] apiv2 data: ", apiV2)

	return []interface{}{apiV2}
}

func flattenVariableV2Read(checkVariableV2 *sc2.VariableV2Response) []interface{} {
	variableV2 := make(map[string]interface{})

	variableV2["name"] = checkVariableV2.Name

	variableV2["description"] = checkVariableV2.Description

	variableV2["value"] = checkVariableV2.Value

	variableV2["secret"] = checkVariableV2.Secret

	log.Println("[DEBUG] VARIABLE V2 data: ", variableV2)

	return []interface{}{variableV2}
}

func flattenVariableV2Data(checkVariableV2 *sc2.VariableV2Response) []interface{} {
	variableV2 := make(map[string]interface{})

	variableV2["name"] = checkVariableV2.Name

	variableV2["id"] = checkVariableV2.ID

	variableV2["description"] = checkVariableV2.Description

	variableV2["value"] = checkVariableV2.Value

	variableV2["secret"] = checkVariableV2.Secret

	if checkVariableV2.Createdat.IsZero() {
	} else {
		variableV2["created_at"] = checkVariableV2.Createdat.String()
	}

	if checkVariableV2.Updatedat.IsZero() {
	} else {
		variableV2["updated_at"] = checkVariableV2.Updatedat.String()
	}

	log.Println("[DEBUG] VARIABLE V2 data: ", variableV2)

	return []interface{}{variableV2}
}

func flattenVariablesV2Data(variables *[]sc2.Variable) []interface{} {
	if variables != nil {
		cls := make([]interface{}, len(*variables))

		for i, variable := range *variables {
			cl := make(map[string]interface{})

			cl["id"] = variable.ID
			cl["name"] = variable.Name
			cl["secret"] = variable.Secret
			cl["value"] = variable.Value
			cl["description"] = variable.Description
			cl["created_at"] = variable.Createdat.String()
			cl["updated_at"] = variable.Updatedat.String()

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func buildDowntimeConfigurationV2Data(d *schema.ResourceData) sc2.DowntimeConfigurationV2Input {
	var downtimeConfigurationV2 sc2.DowntimeConfigurationV2Input
	downtimeConfigurationV2Data := d.Get("downtime_configuration").(*schema.Set).List()
	var i = 0
	layout := "2006-01-02T15:04:05.000Z"
	for _, downtimeConfiguration := range downtimeConfigurationV2Data {
		if i < 1 {
			downtimeConfiguration := downtimeConfiguration.(map[string]interface{})
			downtimeConfigurationV2.Description = downtimeConfiguration["description"].(string)
			downtimeConfigurationV2.Name = downtimeConfiguration["name"].(string)
			downtimeConfigurationV2.Rule = downtimeConfiguration["rule"].(string)
			startTime, err := time.Parse(layout, downtimeConfiguration["start_time"].(string))
			if err != nil {
				_ = err
			}
			downtimeConfigurationV2.Starttime = startTime
			endTime, err := time.Parse(layout, downtimeConfiguration["end_time"].(string))
			if err != nil {
				_ = err
			}
			downtimeConfigurationV2.Endtime = endTime
			timezoneValue, ok := downtimeConfiguration["timezone"].(string)
			if ok {
				downtimeConfigurationV2.Timezone = &timezoneValue
			}
			downtimeConfigurationV2.Testids = buildTestIdData(downtimeConfiguration["test_ids"].([]interface{}))
			downtimeConfigurationV2.Recurrence = buildRecurrenceData(downtimeConfiguration["recurrence"].(*schema.Set))
			i++
		}
	}
	log.Println("[DEBUG]] build downtimeConfigurationV2 data: ", downtimeConfigurationV2)
	return downtimeConfigurationV2
}

func flattenDowntimeConfigurationV2Read(downtimeConfigurationV2 *sc2.DowntimeConfigurationV2Response) []interface{} {
	DowntimeConfigurationV2 := make(map[string]interface{})

	DowntimeConfigurationV2["name"] = downtimeConfigurationV2.Name

	if DowntimeConfigurationV2["description"] != "" {
		DowntimeConfigurationV2["description"] = downtimeConfigurationV2.Description
	}

	DowntimeConfigurationV2["rule"] = downtimeConfigurationV2.Rule

	DowntimeConfigurationV2["start_time"] = downtimeConfigurationV2.Starttime.Format("2006-01-02T15:04:05.000Z")

	DowntimeConfigurationV2["end_time"] = downtimeConfigurationV2.Endtime.Format("2006-01-02T15:04:05.000Z")

	DowntimeConfigurationV2["test_ids"] = downtimeConfigurationV2.Testids

	if downtimeConfigurationV2.Timezone != nil {
		DowntimeConfigurationV2["timezone"] = downtimeConfigurationV2.Timezone
	}

	if downtimeConfigurationV2.Recurrence != nil {
		DowntimeConfigurationV2["recurrence"] = flattenRecurrenceData(downtimeConfigurationV2.Recurrence)
	}

	log.Println("[DEBUG] DowntimeConfiguration V2 data: ", downtimeConfigurationV2)

	return []interface{}{DowntimeConfigurationV2}
}

func flattenDowntimeConfigurationV2Data(downtimeConfigurationV2 *sc2.DowntimeConfigurationV2Response) []interface{} {
	DowntimeConfigurationV2 := make(map[string]interface{})

	DowntimeConfigurationV2["name"] = downtimeConfigurationV2.Name

	DowntimeConfigurationV2["id"] = downtimeConfigurationV2.ID

	DowntimeConfigurationV2["description"] = downtimeConfigurationV2.Description

	DowntimeConfigurationV2["rule"] = downtimeConfigurationV2.Rule

	DowntimeConfigurationV2["start_time"] = downtimeConfigurationV2.Starttime.Format("2006-01-02T15:04:05.000Z")

	DowntimeConfigurationV2["end_time"] = downtimeConfigurationV2.Endtime.Format("2006-01-02T15:04:05.000Z")

	DowntimeConfigurationV2["status"] = downtimeConfigurationV2.Status

	if downtimeConfigurationV2.Createdat.IsZero() {
	} else {
		DowntimeConfigurationV2["created_at"] = downtimeConfigurationV2.Createdat.String()
	}

	if downtimeConfigurationV2.Updatedat.IsZero() {
	} else {
		DowntimeConfigurationV2["updated_at"] = downtimeConfigurationV2.Updatedat.String()
	}

	if downtimeConfigurationV2.Testsupdatedat.IsZero() {
	} else {
		DowntimeConfigurationV2["tests_updated_at"] = downtimeConfigurationV2.Testsupdatedat.String()
	}

	DowntimeConfigurationV2["test_count"] = downtimeConfigurationV2.Testcount

	if downtimeConfigurationV2.Timezone != nil {
		DowntimeConfigurationV2["timezone"] = downtimeConfigurationV2.Timezone
	}

	if downtimeConfigurationV2.Recurrence != nil {
		DowntimeConfigurationV2["recurrence"] = flattenRecurrenceData(downtimeConfigurationV2.Recurrence)
	}

	log.Println("[DEBUG] DowntimeConfiguration V2 data: ", downtimeConfigurationV2)

	return []interface{}{DowntimeConfigurationV2}
}

func flattenDowntimeConfigurationsV2Data(downtimeConfigurations *[]sc2.DowntimeConfiguration) []interface{} {
	if downtimeConfigurations != nil {
		cls := make([]interface{}, len(*downtimeConfigurations))

		for i, downtimeConfiguration := range *downtimeConfigurations {
			cl := make(map[string]interface{})

			cl["id"] = downtimeConfiguration.ID
			cl["name"] = downtimeConfiguration.Name
			cl["description"] = downtimeConfiguration.Description
			cl["rule"] = downtimeConfiguration.Rule
			cl["start_time"] = downtimeConfiguration.Starttime.Format("2006-01-02T15:04:05.000Z")
			cl["end_time"] = downtimeConfiguration.Endtime.Format("2006-01-02T15:04:05.000Z")
			cl["status"] = downtimeConfiguration.Status
			cl["created_at"] = downtimeConfiguration.Createdat.String()
			cl["updated_at"] = downtimeConfiguration.Updatedat.String()
			cl["tests_updated_at"] = downtimeConfiguration.Testsupdatedat.String()
			cl["test_count"] = downtimeConfiguration.Testcount
			if downtimeConfiguration.Timezone != nil {
				cl["timezone"] = downtimeConfiguration.Timezone
			}
			if downtimeConfiguration.Recurrence != nil {
				cl["recurrence"] = flattenRecurrenceData(downtimeConfiguration.Recurrence)
			}

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenRecurrenceData(recurrenceData *sc2.Recurrence) []interface{} {
	recurrence := make(map[string]interface{})

	recurrence["repeats"] = flattenRepeatsData(recurrenceData.Repeats)

	if recurrenceData.End != nil {
		recurrence["end"] = flattenEndData(recurrenceData.End)
	}

	return []interface{}{recurrence}
}

func flattenRepeatsData(repeatsData sc2.Repeats) []interface{} {
	repeats := make(map[string]interface{})

	repeats["type"] = repeatsData.Type

	if repeatsData.Customvalue != nil {
		repeats["custom_value"] = *repeatsData.Customvalue
	}
	if repeatsData.Customfrequency != nil {
		repeats["custom_frequency"] = *repeatsData.Customfrequency
	}

	return []interface{}{repeats}
}

func flattenEndData(endData *sc2.End) []interface{} {
	end := make(map[string]interface{})

	end["type"] = endData.Type

	end["value"] = endData.Value

	return []interface{}{end}
}

func flattenDevicesV2Data(devices *[]sc2.Device) []interface{} {
	if devices != nil {
		cls := make([]interface{}, len(*devices))

		for i, variable := range *devices {
			cl := make(map[string]interface{})

			cl["id"] = variable.ID
			cl["label"] = variable.Label
			cl["user_agent"] = variable.UserAgent
			Networkconnection := flattenNetworkConnectionData(&variable.Networkconnection)
			cl["network_connection"] = Networkconnection
			cl["viewport_height"] = variable.Viewportheight
			cl["viewport_width"] = variable.Viewportwidth

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenLocationsV2Data(locations *[]sc2.Location) []interface{} {
	if locations != nil {
		cls := make([]interface{}, len(*locations))

		for i, location := range *locations {
			cl := make(map[string]interface{})

			cl["id"] = location.ID
			cl["label"] = location.Label
			cl["default"] = location.Default
			cl["type"] = location.Type
			cl["country"] = location.Country

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenDefaultLocationData(checkLocations []string) []interface{} {
	if checkLocations != nil {
		cls := make([]interface{}, len(checkLocations))

		for i, checkLocations := range checkLocations {
			cls[i] = checkLocations
		}
		return cls
	}
	return make([]interface{}, 0)
}

func buildLocationV2Data(d *schema.ResourceData) sc2.LocationV2Input {
	locationData := d.Get("location").(*schema.Set).List()
	var location sc2.LocationV2Input
	for _, lo := range locationData {
		loc := lo.(map[string]interface{})
		location.ID = loc["id"].(string)
		location.Label = loc["label"].(string)
		location.Default = loc["default"].(bool)
		location.Type = loc["type"].(string)
		location.Country = loc["country"].(string)
	}
	return location
}

func flattenLocationV2Data(checkLocationV2 sc2.Location) []interface{} {
	locationV2 := make(map[string]interface{})

	locationV2["id"] = checkLocationV2.ID

	locationV2["label"] = checkLocationV2.Label

	locationV2["default"] = checkLocationV2.Default

	locationV2["type"] = checkLocationV2.Type

	locationV2["country"] = checkLocationV2.Country

	log.Println("[DEBUG] Location V2 data: ", locationV2)

	return []interface{}{locationV2}
}

func flattenLocationMetaV2Data(checkLocationV2 sc2.Meta) []interface{} {
	locationMetaV2 := make(map[string]interface{})

	locationMetaV2["active_test_ids"] = checkLocationV2.ActiveTestIds

	locationMetaV2["paused_test_ids"] = checkLocationV2.PausedTestIds

	log.Println("[DEBUG] Location Meta V2 data: ", locationMetaV2)

	return []interface{}{locationMetaV2}
}

func flattenBrowserV2Read(checkBrowserV2 *sc2.BrowserCheckV2Response) []interface{} {
	browserV2 := make(map[string]interface{})

	browserV2["active"] = checkBrowserV2.Test.Active
	browserV2["automatic_retries"] = checkBrowserV2.Test.Automaticretries

	browserV2["device_id"] = checkBrowserV2.Test.Deviceid

	if checkBrowserV2.Test.Frequency != 0 {
		browserV2["frequency"] = checkBrowserV2.Test.Frequency
	}

	if checkBrowserV2.Test.Name != "" {
		browserV2["name"] = checkBrowserV2.Test.Name
	}

	if checkBrowserV2.Test.Schedulingstrategy != "" {
		browserV2["scheduling_strategy"] = checkBrowserV2.Test.Schedulingstrategy
	}

	locationIds := flattenLocationData(&checkBrowserV2.Test.Locationids)
	browserV2["location_ids"] = locationIds

	advancedSettings := flattenAdvancedSettingsData(&checkBrowserV2.Test.Advancedsettings)
	browserV2["advanced_settings"] = advancedSettings

	transactions := flattenTransactionsData(&checkBrowserV2.Test.Transactions)
	browserV2["transactions"] = transactions

	customProperties := flattenCustomProperties(&checkBrowserV2.Test.Customproperties)
	browserV2["custom_properties"] = customProperties

	log.Println("[DEBUG] read browserv2 data: ", browserV2)

	return []interface{}{browserV2}
}

func flattenBrowserV2Data(checkBrowserV2 *sc2.BrowserCheckV2Response, devices []sc2.Device) []interface{} {
	browserV2 := make(map[string]interface{})

	browserV2["active"] = checkBrowserV2.Test.Active
	browserV2["automatic_retries"] = checkBrowserV2.Test.Automaticretries

	if checkBrowserV2.Test.Createdat.IsZero() {
	} else {
		browserV2["created_at"] = checkBrowserV2.Test.Createdat.String()
	}

	if checkBrowserV2.Test.Updatedat.IsZero() {
	} else {
		browserV2["updated_at"] = checkBrowserV2.Test.Updatedat.String()
	}

	if checkBrowserV2.Test.Createdby != "" {
		browserV2["created_by"] = checkBrowserV2.Test.Createdby
	}

	if checkBrowserV2.Test.Updatedby != "" {
		browserV2["updated_by"] = checkBrowserV2.Test.Updatedby
	}

	if checkBrowserV2.Test.Lastrunat.IsZero() {
	} else {
		browserV2["last_run_at"] = checkBrowserV2.Test.Lastrunat.String()
	}

	if checkBrowserV2.Test.Lastrunstatus != "" {
		browserV2["last_run_status"] = checkBrowserV2.Test.Lastrunstatus
	}

	if checkBrowserV2.Test.Frequency != 0 {
		browserV2["frequency"] = checkBrowserV2.Test.Frequency
	}

	if checkBrowserV2.Test.ID != 0 {
		browserV2["id"] = checkBrowserV2.Test.ID
	}

	if checkBrowserV2.Test.Name != "" {
		browserV2["name"] = checkBrowserV2.Test.Name
	}

	if checkBrowserV2.Test.Schedulingstrategy != "" {
		browserV2["scheduling_strategy"] = checkBrowserV2.Test.Schedulingstrategy
	}

	if checkBrowserV2.Test.Type != "" {
		browserV2["type"] = checkBrowserV2.Test.Type
	}

	locationIds := flattenLocationData(&checkBrowserV2.Test.Locationids)
	browserV2["location_ids"] = locationIds

	browserV2["device"] = flattenDeviceFromID(checkBrowserV2.Test.Deviceid, devices)

	advancedSettings := flattenAdvancedSettingsData(&checkBrowserV2.Test.Advancedsettings)
	browserV2["advanced_settings"] = advancedSettings

	businessTranscations := flattenBusinessTransactionsData(&checkBrowserV2.Test.Transactions)
	browserV2["transactions"] = businessTranscations

	customProperties := flattenCustomProperties(&checkBrowserV2.Test.Customproperties)
	browserV2["custom_properties"] = customProperties

	transcations := flattenTransactionsData(&checkBrowserV2.Test.Transactions)
	browserV2["transactions"] = transcations

	log.Println("[DEBUG] flatten browserv2 data: ", browserV2)

	return []interface{}{browserV2}
}

func flattenHttpV2Read(checkHttpV2 *sc2.HttpCheckV2ResponseWithNullablePort) []interface{} {
	httpV2 := make(map[string]interface{})

	if checkHttpV2.Test.Name != "" {
		httpV2["name"] = checkHttpV2.Test.Name
	}

	httpV2["active"] = checkHttpV2.Test.Active
	httpV2["automatic_retries"] = checkHttpV2.Test.Automaticretries

	if checkHttpV2.Test.Frequency != 0 {
		httpV2["frequency"] = checkHttpV2.Test.Frequency
	}

	if checkHttpV2.Test.SchedulingStrategy != "" {
		httpV2["scheduling_strategy"] = checkHttpV2.Test.SchedulingStrategy
	}

	if checkHttpV2.Test.Type != "" {
		httpV2["type"] = checkHttpV2.Test.Type
	}

	if checkHttpV2.Test.URL != "" {
		httpV2["url"] = checkHttpV2.Test.URL
	}

	if checkHttpV2.Test.Port.Value != nil {
		httpV2["port"] = *checkHttpV2.Test.Port.Value
	}

	if checkHttpV2.Test.RequestMethod != "" {
		httpV2["request_method"] = checkHttpV2.Test.RequestMethod
	}

	if checkHttpV2.Test.Body != "" {
		httpV2["body"] = checkHttpV2.Test.Body
	}

	httpV2["user_agent"] = checkHttpV2.Test.UserAgent

	httpV2["verify_certificates"] = checkHttpV2.Test.Verifycertificates

	locationIds := flattenLocationData(&checkHttpV2.Test.LocationIds)
	httpV2["location_ids"] = locationIds

	httpHeaders := flattenHttpHeadersData(&checkHttpV2.Test.HttpHeaders)
	httpV2["headers"] = httpHeaders

	validations := flattenValidationsData(&checkHttpV2.Test.Validations)
	httpV2["validations"] = validations

	customProperties := flattenCustomProperties(&checkHttpV2.Test.Customproperties)
	httpV2["custom_properties"] = customProperties

	log.Println("[DEBUG] httpV2 data: ", httpV2)

	return []interface{}{httpV2}
}

func flattenHttpV2Data(checkHttpV2 *sc2.HttpCheckV2ResponseWithNullablePort) []interface{} {
	httpV2 := make(map[string]interface{})

	if checkHttpV2.Test.ID != 0 {
		httpV2["id"] = checkHttpV2.Test.ID
	}

	if checkHttpV2.Test.Name != "" {
		httpV2["name"] = checkHttpV2.Test.Name
	}

	httpV2["active"] = checkHttpV2.Test.Active
	httpV2["automatic_retries"] = checkHttpV2.Test.Automaticretries

	if checkHttpV2.Test.Frequency != 0 {
		httpV2["frequency"] = checkHttpV2.Test.Frequency
	}

	if checkHttpV2.Test.CreatedAt.IsZero() {
	} else {
		httpV2["created_at"] = checkHttpV2.Test.CreatedAt.String()
	}

	if checkHttpV2.Test.UpdatedAt.IsZero() {
	} else {
		httpV2["updated_at"] = checkHttpV2.Test.UpdatedAt.String()
	}

	if checkHttpV2.Test.Createdby != "" {
		httpV2["created_by"] = checkHttpV2.Test.Createdby
	}

	if checkHttpV2.Test.Updatedby != "" {
		httpV2["updated_by"] = checkHttpV2.Test.Updatedby
	}

	if checkHttpV2.Test.Lastrunat.IsZero() {
	} else {
		httpV2["last_run_at"] = checkHttpV2.Test.Lastrunat.String()
	}

	if checkHttpV2.Test.Lastrunstatus != "" {
		httpV2["last_run_status"] = checkHttpV2.Test.Lastrunstatus
	}

	if checkHttpV2.Test.SchedulingStrategy != "" {
		httpV2["scheduling_strategy"] = checkHttpV2.Test.SchedulingStrategy
	}

	if checkHttpV2.Test.Type != "" {
		httpV2["type"] = checkHttpV2.Test.Type
	}

	if checkHttpV2.Test.URL != "" {
		httpV2["url"] = checkHttpV2.Test.URL
	}

	if checkHttpV2.Test.Port.Value != nil {
		httpV2["port"] = *checkHttpV2.Test.Port.Value
	}

	if checkHttpV2.Test.RequestMethod != "" {
		httpV2["request_method"] = checkHttpV2.Test.RequestMethod
	}

	if checkHttpV2.Test.Body != "" {
		httpV2["body"] = checkHttpV2.Test.Body
	}

	if checkHttpV2.Test.UserAgent != nil {
		httpV2["user_agent"] = checkHttpV2.Test.UserAgent
	}

	if checkHttpV2.Test.Verifycertificates {
		httpV2["verify_certificates"] = checkHttpV2.Test.Verifycertificates
	}

	locationIds := flattenLocationData(&checkHttpV2.Test.LocationIds)
	httpV2["location_ids"] = locationIds

	httpHeaders := flattenHttpHeadersData(&checkHttpV2.Test.HttpHeaders)
	httpV2["headers"] = httpHeaders

	validations := flattenValidationsData(&checkHttpV2.Test.Validations)
	httpV2["validations"] = validations

	customProperties := flattenCustomProperties(&checkHttpV2.Test.Customproperties)
	httpV2["custom_properties"] = customProperties

	log.Println("[DEBUG] httpV2 data: ", httpV2)

	return []interface{}{httpV2}
}

func flattenPortCheckV2Read(checkPortV2 *sc2.PortCheckV2Response) []interface{} {
	portV2 := make(map[string]interface{})

	if checkPortV2.Test.Name != "" {
		portV2["name"] = checkPortV2.Test.Name
	}

	portV2["active"] = checkPortV2.Test.Active
	portV2["automatic_retries"] = checkPortV2.Test.Automaticretries

	if checkPortV2.Test.Frequency != 0 {
		portV2["frequency"] = checkPortV2.Test.Frequency
	}

	if checkPortV2.Test.SchedulingStrategy != "" {
		portV2["scheduling_strategy"] = checkPortV2.Test.SchedulingStrategy
	}

	if checkPortV2.Test.Type != "" {
		portV2["type"] = checkPortV2.Test.Type
	}

	if checkPortV2.Test.Protocol != "" {
		portV2["protocol"] = checkPortV2.Test.Protocol
	}

	if checkPortV2.Test.Host != "" {
		portV2["host"] = checkPortV2.Test.Host
	}

	if checkPortV2.Test.Port != 0 {
		portV2["port"] = checkPortV2.Test.Port
	}

	locationIds := flattenLocationData(&checkPortV2.Test.LocationIds)
	portV2["location_ids"] = locationIds

	customProperties := flattenCustomProperties(&checkPortV2.Test.Customproperties)
	portV2["custom_properties"] = customProperties

	log.Println("[DEBUG] portv2 data: ", portV2)

	return []interface{}{portV2}
}

func flattenPortCheckV2Data(checkPortV2 *sc2.PortCheckV2Response) []interface{} {
	portV2 := make(map[string]interface{})

	if checkPortV2.Test.ID != 0 {
		portV2["id"] = checkPortV2.Test.ID
	}

	if checkPortV2.Test.Name != "" {
		portV2["name"] = checkPortV2.Test.Name
	}

	portV2["active"] = checkPortV2.Test.Active
	portV2["automatic_retries"] = checkPortV2.Test.Automaticretries

	if checkPortV2.Test.Frequency != 0 {
		portV2["frequency"] = checkPortV2.Test.Frequency
	}

	if checkPortV2.Test.CreatedAt.IsZero() {
	} else {
		portV2["created_at"] = checkPortV2.Test.CreatedAt.String()
	}

	if checkPortV2.Test.UpdatedAt.IsZero() {
	} else {
		portV2["updated_at"] = checkPortV2.Test.UpdatedAt.String()
	}

	if checkPortV2.Test.Createdby != "" {
		portV2["created_by"] = checkPortV2.Test.Createdby
	}

	if checkPortV2.Test.Updatedby != "" {
		portV2["updated_by"] = checkPortV2.Test.Updatedby
	}

	if checkPortV2.Test.Lastrunat.IsZero() {
	} else {
		portV2["last_run_at"] = checkPortV2.Test.Lastrunat.String()
	}

	if checkPortV2.Test.Lastrunstatus != "" {
		portV2["last_run_status"] = checkPortV2.Test.Lastrunstatus
	}

	if checkPortV2.Test.SchedulingStrategy != "" {
		portV2["scheduling_strategy"] = checkPortV2.Test.SchedulingStrategy
	}

	if checkPortV2.Test.Type != "" {
		portV2["type"] = checkPortV2.Test.Type
	}

	if checkPortV2.Test.Protocol != "" {
		portV2["protocol"] = checkPortV2.Test.Protocol
	}

	if checkPortV2.Test.Host != "" {
		portV2["host"] = checkPortV2.Test.Host
	}

	if checkPortV2.Test.Port != 0 {
		portV2["port"] = checkPortV2.Test.Port
	}

	locationIds := flattenLocationData(&checkPortV2.Test.LocationIds)
	portV2["location_ids"] = locationIds

	customProperties := flattenCustomProperties(&checkPortV2.Test.Customproperties)
	portV2["custom_properties"] = customProperties

	log.Println("[DEBUG] portv2 data: ", portV2)

	return []interface{}{portV2}
}

func flattenSslCheckV2Read(checkSslV2 *sc2.SslCheckV2Response) []interface{} {
	sslV2 := make(map[string]interface{})

	if checkSslV2.Test.ID != 0 {
		sslV2["id"] = checkSslV2.Test.ID
	}

	if checkSslV2.Test.Name != "" {
		sslV2["name"] = checkSslV2.Test.Name
	}

	sslV2["active"] = checkSslV2.Test.Active
	sslV2["automatic_retries"] = checkSslV2.Test.Automaticretries

	if checkSslV2.Test.Frequency != 0 {
		sslV2["frequency"] = checkSslV2.Test.Frequency
	}

	if checkSslV2.Test.CreatedAt.IsZero() {
	} else {
		sslV2["created_at"] = checkSslV2.Test.CreatedAt.String()
	}

	if checkSslV2.Test.UpdatedAt.IsZero() {
	} else {
		sslV2["updated_at"] = checkSslV2.Test.UpdatedAt.String()
	}

	if checkSslV2.Test.Createdby != "" {
		sslV2["created_by"] = checkSslV2.Test.Createdby
	}

	if checkSslV2.Test.Updatedby != "" {
		sslV2["updated_by"] = checkSslV2.Test.Updatedby
	}

	if checkSslV2.Test.Lastrunat.IsZero() {
	} else {
		sslV2["last_run_at"] = checkSslV2.Test.Lastrunat.String()
	}

	if checkSslV2.Test.Lastrunstatus != "" {
		sslV2["last_run_status"] = checkSslV2.Test.Lastrunstatus
	}

	if checkSslV2.Test.LastRunLocationId != "" {
		sslV2["last_run_location_id"] = checkSslV2.Test.LastRunLocationId
	}

	if checkSslV2.Test.LastRunId != 0 {
		sslV2["last_run_id"] = checkSslV2.Test.LastRunId
	}

	if checkSslV2.Test.LastRunCoreMetricsPublishedAt.IsZero() {
	} else {
		sslV2["last_run_core_metrics_published_at"] = checkSslV2.Test.LastRunCoreMetricsPublishedAt.String()
	}

	if checkSslV2.Test.SchedulingStrategy != "" {
		sslV2["scheduling_strategy"] = checkSslV2.Test.SchedulingStrategy
	}

	if checkSslV2.Test.Type != "" {
		sslV2["type"] = checkSslV2.Test.Type
	}

	if checkSslV2.Test.Host != "" {
		sslV2["host"] = checkSslV2.Test.Host
	}

	if checkSslV2.Test.Port != 0 {
		sslV2["port"] = checkSslV2.Test.Port
	}

	if checkSslV2.Test.ServerName != nil {
		sslV2["server_name"] = *checkSslV2.Test.ServerName
	}

	sslV2["allow_self_signed"] = checkSslV2.Test.AllowSelfSigned
	sslV2["allow_untrusted_root"] = checkSslV2.Test.AllowUntrustedRoot

	if checkSslV2.Test.CaCertificateID != nil {
		sslV2["ca_certificate_id"] = *checkSslV2.Test.CaCertificateID
	}

	locationIds := flattenLocationData(&checkSslV2.Test.LocationIds)
	sslV2["location_ids"] = locationIds

	validations := flattenSslValidationsData(&checkSslV2.Test.Validations)
	sslV2["validations"] = validations

	customProperties := flattenCustomProperties(&checkSslV2.Test.Customproperties)
	sslV2["custom_properties"] = customProperties

	return []interface{}{sslV2}
}

func flattenSslCheckV2Data(checkSslV2 *sc2.SslCheckV2Response) []interface{} {
	return flattenSslCheckV2Read(checkSslV2)
}

func flattenCaCertificateV2Read(resp *sc2.CaCertificateV2Response, existingContent string) []interface{} {
	return flattenCaCertificateV2Data(resp, existingContent)
}

func flattenCaCertificateV2Data(resp *sc2.CaCertificateV2Response, existingContent string) []interface{} {
	if resp == nil {
		return []interface{}{}
	}
	return []interface{}{flattenCaCertificateData(resp.CaCert, existingContent)}
}

func flattenCaCertificatesV2Data(caCertificates []sc2.CaCertificate) []interface{} {
	cls := make([]interface{}, len(caCertificates))
	for i, caCertificate := range caCertificates {
		cls[i] = flattenCaCertificateData(caCertificate, "")
	}
	return cls
}

func flattenCaCertificateData(caCertificate sc2.CaCertificate, existingContent string) map[string]interface{} {
	caCertificateData := make(map[string]interface{})

	if caCertificate.ID != 0 {
		caCertificateData["id"] = caCertificate.ID
	}
	if caCertificate.Name != "" {
		caCertificateData["name"] = caCertificate.Name
	}
	if caCertificate.Description != "" {
		caCertificateData["description"] = caCertificate.Description
	}
	if content := caCertificateContentForState(caCertificate.Content, existingContent); content != "" {
		caCertificateData["content"] = content
	}
	if caCertificate.FileExtension != "" {
		caCertificateData["file_extension"] = caCertificate.FileExtension
	}
	if caCertificate.Filename != "" {
		caCertificateData["filename"] = caCertificate.Filename
	}
	if !caCertificate.ExpiresAt.IsZero() {
		caCertificateData["expires_at"] = caCertificate.ExpiresAt.String()
	}
	if !caCertificate.CreatedAt.IsZero() {
		caCertificateData["created_at"] = caCertificate.CreatedAt.String()
	}
	if caCertificate.CreatedBy != "" {
		caCertificateData["created_by"] = caCertificate.CreatedBy
	}
	if !caCertificate.UpdatedAt.IsZero() {
		caCertificateData["updated_at"] = caCertificate.UpdatedAt.String()
	}
	if caCertificate.UpdatedBy != "" {
		caCertificateData["updated_by"] = caCertificate.UpdatedBy
	}

	return caCertificateData
}

func caCertificateContentForState(apiContent, existingContent string) string {
	if existingContent != "" && (apiContent == "" || apiContent == caCertificateRedactedContent) {
		return existingContent
	}
	return apiContent
}

func flattenRequestData(checkRequests *[]sc2.Requests) []interface{} {
	if checkRequests != nil {
		cls := make([]interface{}, len(*checkRequests))

		for i, checkRequests := range *checkRequests {
			cl := make(map[string]interface{})

			configuration := flattenConfigurationData(&checkRequests.Configuration)
			cl["configuration"] = configuration

			setup := flattenSetupData(&checkRequests.Setup)
			cl["setup"] = setup

			validations := flattenValidationsData(&checkRequests.Validations)
			cl["validations"] = validations

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenBusinessTransactionsData(checkBusinessTransactions *[]sc2.Transactions) []interface{} {
	if checkBusinessTransactions != nil {
		cls := make([]interface{}, len(*checkBusinessTransactions))

		for i, checkBusinessTransactions := range *checkBusinessTransactions {
			cl := make(map[string]interface{})

			cl["name"] = checkBusinessTransactions.Name

			steps := flattenStepsData(&checkBusinessTransactions.StepsV2)
			cl["steps"] = steps
			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenHttpHeadersData(checkHttpHeaders *[]sc2.HttpHeaders) []interface{} {
	if checkHttpHeaders != nil {
		cls := make([]interface{}, len(*checkHttpHeaders))

		for i, checkHttpHeaders := range *checkHttpHeaders {
			cl := make(map[string]interface{})

			cl["name"] = checkHttpHeaders.Name
			cl["value"] = checkHttpHeaders.Value

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenCustomProperties(checkCustomProperties *[]sc2.CustomProperties) []interface{} {
	if checkCustomProperties != nil {
		cls := make([]interface{}, len(*checkCustomProperties))

		for i, checkCustomProperties := range *checkCustomProperties {
			cl := make(map[string]interface{})

			cl["key"] = checkCustomProperties.Key
			cl["value"] = checkCustomProperties.Value

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenTransactionsData(checkTransactions *[]sc2.Transactions) []interface{} {
	if checkTransactions != nil {
		cls := make([]interface{}, len(*checkTransactions))

		for i, checkTransactions := range *checkTransactions {
			cl := make(map[string]interface{})

			cl["name"] = checkTransactions.Name

			steps := flattenStepsData(&checkTransactions.StepsV2)
			cl["steps"] = steps
			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenConfigurationData(checkConfiguration *sc2.Configuration) []interface{} {
	configuration := make(map[string]interface{})

	if checkConfiguration.Body != "" {
		configuration["body"] = checkConfiguration.Body
	}
	if checkConfiguration.Name != "" {
		configuration["name"] = checkConfiguration.Name
	}

	if checkConfiguration.RequestMethod != "" {
		configuration["request_method"] = checkConfiguration.RequestMethod
	}
	if checkConfiguration.URL != "" {
		configuration["url"] = checkConfiguration.URL
	}

	headers := flattenHeaderData(&checkConfiguration.Headers)
	configuration["headers"] = headers

	return []interface{}{configuration}
}

func flattenStepsData(checkSteps *[]sc2.StepsV2) []interface{} {
	if checkSteps != nil {
		cls := make([]interface{}, len(*checkSteps))

		for i, checkStep := range *checkSteps {
			cl := make(map[string]interface{})

			if checkStep.Name != "" {
				cl["name"] = checkStep.Name
			}

			if checkStep.Type != "" {
				cl["type"] = checkStep.Type
			}

			if checkStep.URL != "" {
				cl["url"] = checkStep.URL
			}

			cl["wait_for_nav"] = checkStep.WaitForNav

			if checkStep.WaitForNavTimeout != 0 && !checkStep.WaitForNavTimeoutDefault {
				cl["wait_for_nav_timeout"] = checkStep.WaitForNavTimeout
			}

			cl["wait_for_nav_timeout_default"] = checkStep.WaitForNavTimeoutDefault

			if checkStep.MaxWaitTime != 0 && !checkStep.MaxWaitTimeDefault {
				cl["max_wait_time"] = checkStep.MaxWaitTime
			}

			cl["max_wait_time_default"] = checkStep.MaxWaitTimeDefault

			// Persist all API selectors (including a single one) as selectors blocks so
			// config using selectors { } stays in sync after apply. Legacy fields are only
			// written when the API returns no selectors.
			if len(checkStep.Selectors) > 0 {
				if selectors := flattenSelectorsData(checkStep.Selectors); selectors != nil {
					cl["selectors"] = selectors
				}
			}

			if checkStep.OptionSelectorType != "" {
				cl["option_selector_type"] = checkStep.OptionSelectorType
			}

			if checkStep.OptionSelector != "" {
				cl["option_selector"] = checkStep.OptionSelector
			}

			if checkStep.VariableName != "" {
				cl["variable_name"] = checkStep.VariableName
			}

			if checkStep.Value != "" {
				cl["value"] = string(checkStep.Value)
			}

			if checkStep.Duration != 0 {
				cl["duration"] = checkStep.Duration
			}

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenChromeFlagsData(chromeFlags []sc2.ChromeFlag) []interface{} {
	if chromeFlags == nil {
		return []interface{}{}
	}

	var result []interface{}
	for _, flag := range chromeFlags {
		flagData := map[string]interface{}{
			"name":  flag.Name,
			"value": flag.Value,
		}
		result = append(result, flagData)
	}

	return result
}

func flattenSetupData(checkSetup *[]sc2.Setup) []interface{} {
	if checkSetup != nil {
		cls := make([]interface{}, len(*checkSetup))

		for i, checkSetup := range *checkSetup {
			cl := make(map[string]interface{})

			if checkSetup.Extractor != "" {
				cl["extractor"] = checkSetup.Extractor
			}

			if checkSetup.Name != "" {
				cl["name"] = checkSetup.Name
			}

			if checkSetup.Source != "" {
				cl["source"] = checkSetup.Source
			}

			if checkSetup.Type != "" {
				cl["type"] = checkSetup.Type
			}

			if checkSetup.Variable != "" {
				cl["variable"] = checkSetup.Variable
			}

			if checkSetup.Code != "" {
				cl["code"] = checkSetup.Code
			}

			if checkSetup.Value != "" {
				cl["value"] = checkSetup.Value
			}

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenCookiesData(checkCookies *[]sc2.Cookiesv2) []interface{} {
	if checkCookies != nil {
		cls := make([]interface{}, len(*checkCookies))

		for i, checkSetup := range *checkCookies {
			cl := make(map[string]interface{})

			if checkSetup.Key != "" {
				cl["key"] = checkSetup.Key
			}

			if checkSetup.Value != "" {
				cl["value"] = checkSetup.Value
			}

			if checkSetup.Domain != "" {
				cl["domain"] = checkSetup.Domain
			}

			if checkSetup.Path != "" {
				cl["path"] = checkSetup.Path
			}

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenBrowserHeadersData(checkBrowserHeaders *[]sc2.BrowserHeaders) []interface{} {
	if checkBrowserHeaders != nil {
		cls := make([]interface{}, len(*checkBrowserHeaders))

		for i, checkSetup := range *checkBrowserHeaders {
			cl := make(map[string]interface{})

			if checkSetup.Name != "" {
				cl["name"] = checkSetup.Name
			}

			if checkSetup.Value != "" {
				cl["value"] = checkSetup.Value
			}

			if checkSetup.Domain != "" {
				cl["domain"] = checkSetup.Domain
			}

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenHostOverridesData(checkHostOverrides *[]sc2.HostOverrides) []interface{} {
	if checkHostOverrides != nil {
		cls := make([]interface{}, len(*checkHostOverrides))

		for i, checkSetup := range *checkHostOverrides {
			cl := make(map[string]interface{})

			if checkSetup.Source != "" {
				cl["source"] = checkSetup.Source
			}

			if checkSetup.Target != "" {
				cl["target"] = checkSetup.Target
			}

			if checkSetup.KeepHostHeader {
				cl["keep_host_header"] = checkSetup.KeepHostHeader
			}

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenValidationsData(checkValidations *[]sc2.Validations) []interface{} {
	if checkValidations != nil {
		cls := make([]interface{}, len(*checkValidations))

		for i, checkValidations := range *checkValidations {
			cl := make(map[string]interface{})

			if checkValidations.Name != "" {
				cl["name"] = checkValidations.Name
			}

			if checkValidations.Type != "" {
				cl["type"] = checkValidations.Type
			}

			if checkValidations.Actual != "" {
				cl["actual"] = checkValidations.Actual
			}

			if checkValidations.Expected != "" {
				cl["expected"] = checkValidations.Expected
			}

			if checkValidations.Comparator != "" {
				cl["comparator"] = checkValidations.Comparator
			}

			if checkValidations.Extractor != "" {
				cl["extractor"] = checkValidations.Extractor
			}

			if checkValidations.Source != "" {
				cl["source"] = checkValidations.Source
			}

			if checkValidations.Variable != "" {
				cl["variable"] = checkValidations.Variable
			}

			if checkValidations.Value != "" {
				cl["value"] = checkValidations.Value
			}

			if checkValidations.Code != "" {
				cl["code"] = checkValidations.Code
			}

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenSslValidationsData(checkValidations *[]sc2.Validations) []interface{} {
	if checkValidations != nil {
		cls := make([]interface{}, len(*checkValidations))

		for i, checkValidation := range *checkValidations {
			cl := make(map[string]interface{})

			if checkValidation.Name != "" {
				cl["name"] = checkValidation.Name
			}

			if checkValidation.Type != "" {
				cl["type"] = checkValidation.Type
			}

			if checkValidation.Actual != "" {
				cl["actual"] = checkValidation.Actual
			}

			if checkValidation.Expected != "" {
				cl["expected"] = checkValidation.Expected
			}

			if checkValidation.Comparator != "" {
				cl["comparator"] = checkValidation.Comparator
			}

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenHeaderData(checkHeaders *sc2.Headers) map[string]interface{} {
	if checkHeaders != nil {
		cls := make(map[string]interface{}, len(*checkHeaders))

		for k, v := range *checkHeaders {
			cls[k] = v
		}
		return cls
	}
	return make(map[string]interface{}, 0)
}

func flattenLocationData(checkLocations *[]string) []interface{} {
	if checkLocations != nil {
		cls := make([]interface{}, len(*checkLocations))

		for i, checkLocations := range *checkLocations {
			cls[i] = checkLocations
		}
		return cls
	}
	return make([]interface{}, 0)
}

func flattenDeviceData(checkDevice *sc2.Device) []interface{} {
	device := make(map[string]interface{})

	if checkDevice.ID != 0 {
		device["id"] = checkDevice.ID
	}

	if checkDevice.Label != "" {
		device["label"] = checkDevice.Label
	}

	if checkDevice.UserAgent != "" {
		device["user_agent"] = checkDevice.UserAgent
	}

	if checkDevice.Viewportheight != 0 {
		device["viewport_height"] = checkDevice.Viewportheight
	}
	if checkDevice.Viewportwidth != 0 {
		device["viewport_width"] = checkDevice.Viewportwidth
	}

	Networkconnection := flattenNetworkConnectionData(&checkDevice.Networkconnection)
	device["network_connection"] = Networkconnection

	return []interface{}{device}
}

func flattenAdvancedSettingsData(advSettings *sc2.Advancedsettings) []interface{} {
	advancedSettings := make(map[string]interface{})

	if advSettings.Verifycertificates {
		advancedSettings["verify_certificates"] = advSettings.Verifycertificates
	}

	if advSettings.CollectInteractiveMetrics {
		advancedSettings["collect_interactive_metrics"] = advSettings.CollectInteractiveMetrics
	}

	if advSettings.UserAgent != nil {
		advancedSettings["user_agent"] = advSettings.UserAgent
	}

	if advSettings.Authentication != nil {
		Authentication := flattenAuthenticationData(advSettings.Authentication)
		advancedSettings["authentication"] = Authentication
	}

	Cookies := flattenCookiesData(&advSettings.Cookiesv2)
	advancedSettings["cookies"] = Cookies

	BrowserHeaders := flattenBrowserHeadersData(&advSettings.BrowserHeaders)
	advancedSettings["headers"] = BrowserHeaders

	HostOverRides := flattenHostOverridesData(&advSettings.HostOverrides)
	advancedSettings["host_overrides"] = HostOverRides

	ChromeFlags := flattenChromeFlagsData(advSettings.ChromeFlags)
	advancedSettings["chrome_flags"] = ChromeFlags

	return []interface{}{advancedSettings}
}

func flattenNetworkConnectionData(checkNetworkConnection *sc2.Networkconnection) []interface{} {
	networkConnection := make(map[string]interface{})

	networkConnection["description"] = checkNetworkConnection.Description
	networkConnection["download_bandwidth"] = checkNetworkConnection.Downloadbandwidth
	networkConnection["latency"] = checkNetworkConnection.Latency
	networkConnection["packet_loss"] = checkNetworkConnection.Packetloss
	networkConnection["upload_bandwidth"] = checkNetworkConnection.Uploadbandwidth

	return []interface{}{networkConnection}
}

func flattenAuthenticationData(checkAuthentications *sc2.Authentication) []interface{} {
	authentication := make(map[string]interface{})

	if checkAuthentications.Username != "" {
		authentication["username"] = checkAuthentications.Username
	}
	if checkAuthentications.Password != "" {
		authentication["password"] = checkAuthentications.Password
	}

	return []interface{}{authentication}
}

func buildApiV2Data(d *schema.ResourceData) sc2.ApiCheckV2Input {
	var apiv2 sc2.ApiCheckV2Input
	apiv2Data := d.Get("test").(*schema.Set).List()
	for _, api := range apiv2Data {
		api := api.(map[string]interface{})
		if api["name"].(string) != "" {
			apiv2.Test.Active = api["active"].(bool)
			apiv2.Test.Deviceid = api["device_id"].(int)
			apiv2.Test.Frequency = api["frequency"].(int)
			apiv2.Test.Automaticretries = api["automatic_retries"].(int)
			apiv2.Test.Locationids = buildLocationIdData(api["location_ids"].([]interface{}))
			apiv2.Test.Name = api["name"].(string)
			apiv2.Test.Requests = buildRequestsData(api["requests"].(([]interface{})))
			apiv2.Test.Schedulingstrategy = api["scheduling_strategy"].(string)
			apiv2.Test.Customproperties = buildCustomPropertiesData(api["custom_properties"].(*schema.Set))
		}
	}
	log.Println("[DEBUG] build apiv2 data: ", apiv2)
	return apiv2
}

func buildBrowserV2Data(d *schema.ResourceData) (sc2.BrowserCheckV2Input, error) {
	var browserv2 sc2.BrowserCheckV2Input
	browserv2Data := d.Get("test").([]interface{})
	for _, browser := range browserv2Data {
		browser := browser.(map[string]interface{})
		if browser["name"].(string) != "" {
			browserv2.Test.Active = browser["active"].(bool)
			browserv2.Test.DeviceID = browser["device_id"].(int)
			browserv2.Test.Frequency = browser["frequency"].(int)
			browserv2.Test.Automaticretries = browser["automatic_retries"].(int)
			browserv2.Test.LocationIds = buildLocationIdData(browser["location_ids"].([]interface{}))
			browserv2.Test.Name = browser["name"].(string)
			transactions, err := buildBusinessTransactionsData(browser["transactions"].([]interface{}))
			if err != nil {
				return browserv2, err
			}
			browserv2.Test.Transactions = transactions
			browserv2.Test.Schedulingstrategy = browser["scheduling_strategy"].(string)
			browserv2.Test.Advancedsettings = buildAdvancedSettingsData(browser["advanced_settings"].(*schema.Set))
			browserv2.Test.Customproperties = buildCustomPropertiesData(browser["custom_properties"].(*schema.Set))
		}
	}

	log.Println("[DEBUG] build browserv2 data:", browserv2)
	return browserv2, nil
}

func buildHttpV2Data(d *schema.ResourceData) sc2.HttpCheckV2InputWithNullablePort {
	var httpv2 sc2.HttpCheckV2InputWithNullablePort
	httpv2Data := d.Get("test").(*schema.Set).List()
	var i = 0
	for _, http := range httpv2Data {
		if i < 1 {
			http := http.(map[string]interface{})
			httpv2.Test.Name = http["name"].(string)
			httpv2.Test.Type = http["type"].(string)
			httpv2.Test.URL = http["url"].(string)
			httpv2.Test.Port = httpV2PortFromResourceData(d)
			httpv2.Test.LocationIds = buildLocationIdData(http["location_ids"].([]interface{}))
			httpv2.Test.Frequency = http["frequency"].(int)
			httpv2.Test.Automaticretries = http["automatic_retries"].(int)
			httpv2.Test.SchedulingStrategy = http["scheduling_strategy"].(string)
			httpv2.Test.Active = http["active"].(bool)
			httpv2.Test.RequestMethod = http["request_method"].(string)
			httpv2.Test.Body = http["body"].(string)
			httpv2.Test.Verifycertificates = http["verify_certificates"].(bool)
			userAgentString := http["user_agent"].(string)
			httpv2.Test.UserAgent = &userAgentString
			httpv2.Test.HttpHeaders = buildHttpHeadersData(http["headers"].(*schema.Set))
			httpv2.Test.Validations = buildValidationsData(http["validations"].([]interface{}))
			httpv2.Test.Customproperties = buildCustomPropertiesData(http["custom_properties"].(*schema.Set))
			i++
		}
	}
	log.Println("[DEBUG] build httpv2 data: ", httpv2)
	return httpv2
}

func httpV2PortFromResourceData(d *schema.ResourceData) sc2.NullableInt {
	port, ok := httpV2PortFromRawValue(d.GetRawConfig())
	if !ok {
		return *sc2.NewNullInt()
	}
	return port
}

func httpV2PortFromRawValue(raw cty.Value) (sc2.NullableInt, bool) {
	if raw.IsNull() || !raw.IsKnown() || !raw.Type().IsObjectType() || !raw.Type().HasAttribute("test") {
		return *sc2.NewNullInt(), false
	}

	test := raw.GetAttr("test")
	if test.IsNull() || !test.IsKnown() {
		return *sc2.NewNullInt(), false
	}

	it := test.ElementIterator()
	if !it.Next() {
		return *sc2.NewNullInt(), false
	}

	_, testBlock := it.Element()
	if testBlock.IsNull() || !testBlock.IsKnown() {
		return *sc2.NewNullInt(), false
	}

	if !testBlock.Type().IsObjectType() || !testBlock.Type().HasAttribute("port") {
		return *sc2.NewNullInt(), true
	}
	portValue := testBlock.GetAttr("port")
	if portValue.IsNull() {
		return *sc2.NewNullInt(), true
	}
	if !portValue.IsKnown() {
		return *sc2.NewNullInt(), false
	}

	var port int
	if err := gocty.FromCtyValue(portValue, &port); err != nil {
		return *sc2.NewNullInt(), false
	}
	return *sc2.NewNullableInt(port), true
}

func buildPortCheckV2Data(d *schema.ResourceData) sc2.PortCheckV2Input {
	var portv2 sc2.PortCheckV2Input
	portv2Data := d.Get("test").(*schema.Set).List()
	var i = 0
	for _, port := range portv2Data {
		if i < 1 {
			port := port.(map[string]interface{})
			portv2.Test.Name = port["name"].(string)
			portv2.Test.Type = port["type"].(string)
			portv2.Test.URL = port["url"].(string)
			portv2.Test.Port = port["port"].(int)
			portv2.Test.Protocol = port["protocol"].(string)
			portv2.Test.Host = port["host"].(string)
			portv2.Test.LocationIds = buildLocationIdData(port["location_ids"].([]interface{}))
			portv2.Test.Frequency = port["frequency"].(int)
			portv2.Test.Automaticretries = port["automatic_retries"].(int)
			portv2.Test.SchedulingStrategy = port["scheduling_strategy"].(string)
			portv2.Test.Active = port["active"].(bool)
			portv2.Test.Customproperties = buildCustomPropertiesData(port["custom_properties"].(*schema.Set))
			i++

		}
	}
	log.Println("[DEBUG] build portv2 data: ", portv2)
	return portv2
}

func buildSslCheckV2Data(d *schema.ResourceData) sc2.SslCheckV2Input {
	var sslv2 sc2.SslCheckV2Input
	sslv2Data := d.Get("test").(*schema.Set).List()
	var i = 0
	for _, ssl := range sslv2Data {
		if i < 1 {
			ssl := ssl.(map[string]interface{})
			sslv2.Test.Name = ssl["name"].(string)
			sslv2.Test.LocationIds = buildLocationIdData(ssl["location_ids"].([]interface{}))
			sslv2.Test.Frequency = ssl["frequency"].(int)
			sslv2.Test.SchedulingStrategy = ssl["scheduling_strategy"].(string)
			sslv2.Test.Active = ssl["active"].(bool)
			sslv2.Test.Customproperties = buildCustomPropertiesData(ssl["custom_properties"].(*schema.Set))
			sslv2.Test.Automaticretries = ssl["automatic_retries"].(int)
			sslv2.Test.Host = ssl["host"].(string)
			sslv2.Test.Port = ssl["port"].(int)
			if serverName := sslStringField(ssl, "server_name"); serverName != "" {
				sslv2.Test.ServerName = &serverName
			}
			sslv2.Test.AllowSelfSigned = ssl["allow_self_signed"].(bool)
			sslv2.Test.AllowUntrustedRoot = ssl["allow_untrusted_root"].(bool)
			if caCertificateID := sslIntField(ssl, "ca_certificate_id"); caCertificateID != 0 {
				sslv2.Test.CaCertificateID = &caCertificateID
			}
			sslv2.Test.Validations = buildValidationsData(sslInterfaceListField(ssl, "validations"))
			i++
		}
	}
	return sslv2
}

func buildSslCheckV2UpdateData(d *schema.ResourceData) sc2.SslCheckV2UpdateInput {
	var sslv2 sc2.SslCheckV2UpdateInput
	sslv2Data := d.Get("test").(*schema.Set).List()
	var i = 0
	for _, ssl := range sslv2Data {
		if i < 1 {
			ssl := ssl.(map[string]interface{})
			name := ssl["name"].(string)
			locationIds := buildLocationIdData(ssl["location_ids"].([]interface{}))
			frequency := ssl["frequency"].(int)
			schedulingStrategy := ssl["scheduling_strategy"].(string)
			active := ssl["active"].(bool)
			customProperties := buildCustomPropertiesData(ssl["custom_properties"].(*schema.Set))
			automaticRetries := ssl["automatic_retries"].(int)
			host := ssl["host"].(string)
			port := ssl["port"].(int)
			allowSelfSigned := ssl["allow_self_signed"].(bool)
			allowUntrustedRoot := ssl["allow_untrusted_root"].(bool)
			validations := buildValidationsData(sslInterfaceListField(ssl, "validations"))

			sslv2.Test.Name = &name
			sslv2.Test.LocationIds = &locationIds
			sslv2.Test.Frequency = &frequency
			sslv2.Test.SchedulingStrategy = &schedulingStrategy
			sslv2.Test.Active = &active
			sslv2.Test.Customproperties = &customProperties
			sslv2.Test.Automaticretries = &automaticRetries
			sslv2.Test.Host = &host
			sslv2.Test.Port = &port
			if serverName := sslStringField(ssl, "server_name"); serverName != "" {
				sslv2.Test.ServerName = sc2.NewNullableString(serverName)
			} else {
				sslv2.Test.ServerName = sc2.NewNullString()
			}
			sslv2.Test.AllowSelfSigned = &allowSelfSigned
			sslv2.Test.AllowUntrustedRoot = &allowUntrustedRoot
			if caCertificateID := sslIntField(ssl, "ca_certificate_id"); caCertificateID != 0 {
				sslv2.Test.CaCertificateID = sc2.NewNullableInt(caCertificateID)
			} else {
				sslv2.Test.CaCertificateID = sc2.NewNullInt()
			}
			sslv2.Test.Validations = &validations
			i++
		}
	}
	return sslv2
}

const caCertificateRedactedContent = "<REDACTED>"

func buildCaCertificateV2Data(d *schema.ResourceData) (sc2.CaCertificateV2Input, error) {
	var caCertificateV2 sc2.CaCertificateV2Input
	caCertificateV2Data := d.Get("ca_certificate").(*schema.Set).List()
	if len(caCertificateV2Data) == 0 {
		return caCertificateV2, fmt.Errorf("ca_certificate block is required")
	}

	caCertificate := caCertificateV2Data[0].(map[string]interface{})
	content := caCertificateStringField(caCertificate, "content")
	if content == "" || content == caCertificateRedactedContent {
		return caCertificateV2, fmt.Errorf("ca_certificate content is required")
	}

	caCertificateV2.CaCert.Name = caCertificateStringField(caCertificate, "name")
	caCertificateV2.CaCert.Description = caCertificateStringField(caCertificate, "description")
	caCertificateV2.CaCert.Content = content
	caCertificateV2.CaCert.FileExtension = caCertificateStringField(caCertificate, "file_extension")
	caCertificateV2.CaCert.Filename = caCertificateStringField(caCertificate, "filename")
	return caCertificateV2, nil
}

func buildCaCertificateV2UpdateData(d *schema.ResourceData) sc2.CaCertificateV2UpdateInput {
	var caCertificateV2 sc2.CaCertificateV2UpdateInput
	caCertificateV2Data := d.Get("ca_certificate").(*schema.Set).List()
	if len(caCertificateV2Data) == 0 {
		return caCertificateV2
	}

	caCertificate := caCertificateV2Data[0].(map[string]interface{})
	description := caCertificateStringField(caCertificate, "description")
	fileExtension := caCertificateStringField(caCertificate, "file_extension")
	filename := caCertificateStringField(caCertificate, "filename")
	content := caCertificateStringField(caCertificate, "content")

	caCertificateV2.CaCert.Description = &description
	caCertificateV2.CaCert.FileExtension = &fileExtension
	caCertificateV2.CaCert.Filename = &filename
	if content != "" && content != caCertificateRedactedContent {
		caCertificateV2.CaCert.Content = &content
	}

	return caCertificateV2
}

func caCertificateContentFromState(d *schema.ResourceData) string {
	caCertificateData, ok := d.Get("ca_certificate").(*schema.Set)
	if !ok || caCertificateData.Len() == 0 {
		return ""
	}

	caCertificate, ok := caCertificateData.List()[0].(map[string]interface{})
	if !ok {
		return ""
	}
	return caCertificateStringField(caCertificate, "content")
}

func caCertificateStringField(caCertificate map[string]interface{}, key string) string {
	if value, ok := caCertificate[key].(string); ok {
		return value
	}
	return ""
}

func sslStringField(ssl map[string]interface{}, key string) string {
	if value, ok := ssl[key].(string); ok {
		return value
	}
	return ""
}

func sslIntField(ssl map[string]interface{}, key string) int {
	if value, ok := ssl[key].(int); ok {
		return value
	}
	return 0
}

func sslInterfaceListField(ssl map[string]interface{}, key string) []interface{} {
	if value, ok := ssl[key].([]interface{}); ok {
		return value
	}
	return []interface{}{}
}

func buildVariableV2Data(d *schema.ResourceData) sc2.VariableV2Input {
	var variablev2 sc2.VariableV2Input
	variablev2Data := d.Get("variable").(*schema.Set).List()
	var i = 0
	for _, variable := range variablev2Data {
		if i < 1 {
			variable := variable.(map[string]interface{})
			variablev2.Description = variable["description"].(string)
			variablev2.Name = variable["name"].(string)
			variablev2.Secret = variable["secret"].(bool)
			variablev2.Value = variable["value"].(string)
			i++
		}
	}
	log.Println("[DEBUG]] build variablev2 data: ", variablev2)
	return variablev2
}

func buildLocationIdData(d []interface{}) []string {
	locationsList := make([]string, len(d))
	for i, locations := range d {
		locationsList[i] = locations.(string)
	}
	return locationsList
}

func buildTestIdData(d []interface{}) []int {
	testsList := make([]int, len(d))
	for i, tests := range d {
		testsList[i] = tests.(int)
	}
	return testsList
}

func buildRecurrenceData(recurrence *schema.Set) *sc2.Recurrence {
	var recurrenceData sc2.Recurrence

	as_list := recurrence.List()
	if len(as_list) > 0 {
		as_map := as_list[0].(map[string]interface{})

		if repeatsPtr := buildRepeatsData(as_map["repeats"].(*schema.Set)); repeatsPtr != nil {
			recurrenceData.Repeats = *repeatsPtr
		}
		if endPtr := buildEndData(as_map["end"].(*schema.Set)); endPtr != nil {
			recurrenceData.End = endPtr
		}

	}
	return &recurrenceData
}

func buildRepeatsData(repeats *schema.Set) *sc2.Repeats {
	repeatsList := repeats.List()

	if len(repeatsList) > 0 {
		repeatsMap := repeatsList[0].(map[string]interface{})
		repeatsData := &sc2.Repeats{
			Type: repeatsMap["type"].(string),
		}

		if val, ok := repeatsMap["custom_value"].(int); ok {
			repeatsData.Customvalue = &val
		}

		if val, ok := repeatsMap["custom_frequency"].(string); ok {
			repeatsData.Customfrequency = &val
		}

		return repeatsData
	}
	return nil
}

func buildEndData(end *schema.Set) *sc2.End {
	endList := end.List()

	if len(endList) > 0 {
		endMap := endList[0].(map[string]interface{})
		endData := &sc2.End{
			Type:  endMap["type"].(string),
			Value: endMap["value"].(string),
		}
		return endData
	}
	return nil
}

func buildRequestsData(requests []interface{}) []sc2.Requests {
	requestsList := make([]sc2.Requests, len(requests))
	for i, request := range requests {
		request := request.(map[string]interface{})
		req := sc2.Requests{
			Configuration: buildConfigurationData(request["configuration"].([]interface{})),
			Setup:         buildSetupData(request["setup"].([]interface{})),
			Validations:   buildValidationsData(request["validations"].([]interface{})),
		}
		requestsList[i] = req

	}
	return requestsList
}

func buildBusinessTransactionsData(businessTransactions []interface{}) ([]sc2.Transactions, error) {
	businessTransactionsList := make([]sc2.Transactions, len(businessTransactions))
	for i, bisTrans := range businessTransactions {
		bisTrans := bisTrans.(map[string]interface{})
		steps, err := buildStepV2Data(bisTrans["steps"].([]interface{}))
		if err != nil {
			return nil, err
		}
		transaction := sc2.Transactions{
			Name:    bisTrans["name"].(string),
			StepsV2: steps,
		}
		businessTransactionsList[i] = transaction
	}
	return businessTransactionsList, nil
}

func buildHttpHeadersData(httpHeaders *schema.Set) []sc2.HttpHeaders {
	httpHeadersList := make([]sc2.HttpHeaders, len(httpHeaders.List()))

	for i, httpHeads := range httpHeaders.List() {
		http := httpHeads.(map[string]interface{})
		if strings.Contains(http["name"].(string), " ") {
			log.Println("[ERROR] Header names cannot have spaces. Please check your header names")
		}
		headerValues := sc2.HttpHeaders{
			Name:  strings.TrimSpace(http["name"].(string)),
			Value: strings.TrimSpace(http["value"].(string)),
		}
		httpHeadersList[i] = headerValues

	}
	return httpHeadersList
}

func buildCustomPropertiesData(customProperties *schema.Set) []sc2.CustomProperties {
	customPropertiesList := make([]sc2.CustomProperties, len(customProperties.List()))

	for i, props := range customProperties.List() {
		prop := props.(map[string]interface{})
		propValues := sc2.CustomProperties{
			Key:   strings.TrimSpace(prop["key"].(string)),
			Value: strings.TrimSpace(prop["value"].(string)),
		}
		customPropertiesList[i] = propValues

	}
	return customPropertiesList
}

// dropStaleStepSelectors removes a single selectors block carried over from state
// when legacy selector fields were updated to different values in config.
func dropStaleStepSelectors(step map[string]interface{}) {
	legacyType := stepStringField(step, "selector_type")
	legacyVal := stepStringField(step, "selector")
	if legacyType == "" || legacyVal == "" {
		return
	}
	stale := parseSelectorsList(step["selectors"])
	if len(stale) != 1 {
		return
	}
	if stale[0].Type == legacyType && stale[0].Value == legacyVal {
		return
	}
	delete(step, "selectors")
}

func buildStepV2Data(steps []interface{}) ([]sc2.StepsV2, error) {
	stepsList := make([]sc2.StepsV2, len(steps))
	for i, step := range steps {
		step := step.(map[string]interface{})
		dropStaleStepSelectors(step)
		selectors, err := buildSelectorsFromStep(step)
		if err != nil {
			return nil, err
		}
		st := sc2.StepsV2{
			URL:                step["url"].(string),
			Name:               step["name"].(string),
			Type:               step["type"].(string),
			WaitForNav:         step["wait_for_nav"].(bool),
			WaitForNavTimeout:  step["wait_for_nav_timeout"].(int),
			MaxWaitTime:        step["max_wait_time"].(int),
			Selectors:          selectors,
			OptionSelectorType: step["option_selector_type"].(string),
			OptionSelector:     step["option_selector"].(string),
			VariableName:       step["variable_name"].(string),
			Value:              step["value"].(string),
			Duration:           step["duration"].(int),
		}
		stepsList[i] = st

	}
	return stepsList, nil
}

func buildSetupData(setups []interface{}) []sc2.Setup {
	setupsList := make([]sc2.Setup, len(setups))

	for i, setup := range setups {
		setup := setup.(map[string]interface{})
		set := sc2.Setup{
			Extractor: setup["extractor"].(string),
			Name:      setup["name"].(string),
			Source:    setup["source"].(string),
			Type:      setup["type"].(string),
			Variable:  setup["variable"].(string),
			Code:      setup["code"].(string),
			Value:     setup["value"].(string),
		}
		setupsList[i] = set

	}
	return setupsList
}

func buildValidationsData(validations []interface{}) []sc2.Validations {
	validationsList := make([]sc2.Validations, len(validations))

	for i, validation := range validations {
		validation := validation.(map[string]interface{})
		val := sc2.Validations{
			Actual:     stringMapValue(validation, "actual"),
			Comparator: stringMapValue(validation, "comparator"),
			Expected:   stringMapValue(validation, "expected"),
			Name:       stringMapValue(validation, "name"),
			Type:       stringMapValue(validation, "type"),
			Extractor:  stringMapValue(validation, "extractor"),
			Source:     stringMapValue(validation, "source"),
			Variable:   stringMapValue(validation, "variable"),
			Code:       stringMapValue(validation, "code"),
			Value:      stringMapValue(validation, "value"),
		}

		validationsList[i] = val

	}
	return validationsList
}

func stringMapValue(values map[string]interface{}, key string) string {
	if value, ok := values[key].(string); ok {
		return value
	}
	return ""
}

func buildConfigurationData(configuration []interface{}) sc2.Configuration {
	var configurationData sc2.Configuration

	config_list := configuration
	config_map := config_list[0].(map[string]interface{})

	configurationData.Body = config_map["body"].(string)
	configurationData.Headers = config_map["headers"].(map[string]interface{})
	configurationData.Name = config_map["name"].(string)
	configurationData.RequestMethod = config_map["request_method"].(string)
	configurationData.URL = config_map["url"].(string)

	return configurationData
}

func buildAdvancedSettingsData(advancedSettings *schema.Set) sc2.Advancedsettings {
	var advancedSettingsData sc2.Advancedsettings

	as_list := advancedSettings.List()
	if len(as_list) > 0 {
		as_map := as_list[0].(map[string]interface{})

		userAgentString := as_map["user_agent"].(string)
		advancedSettingsData.UserAgent = &userAgentString
		advancedSettingsData.Verifycertificates = as_map["verify_certificates"].(bool)
		advancedSettingsData.CollectInteractiveMetrics = as_map["collect_interactive_metrics"].(bool)
		advancedSettingsData.Authentication = buildAuthenticationData(as_map["authentication"].(*schema.Set))
		advancedSettingsData.BrowserHeaders = buildBrowserHeadersData(as_map["headers"].(*schema.Set))
		advancedSettingsData.Cookiesv2 = buildCookiesData(as_map["cookies"].(*schema.Set))
		advancedSettingsData.HostOverrides = buildHostOverridesData(as_map["host_overrides"].(*schema.Set))
		advancedSettingsData.ChromeFlags = buildChromeFlagsData(as_map["chrome_flags"].(*schema.Set))

	}
	return advancedSettingsData
}

func buildBrowserHeadersData(headers *schema.Set) []sc2.BrowserHeaders {
	headersList := make([]sc2.BrowserHeaders, len(headers.List()))

	for i, header := range headers.List() {
		header := header.(map[string]interface{})
		set := sc2.BrowserHeaders{
			Name:   header["name"].(string),
			Value:  header["value"].(string),
			Domain: header["domain"].(string),
		}
		headersList[i] = set

	}
	return headersList
}

func buildChromeFlagsData(d *schema.Set) []sc2.ChromeFlag {
	var flags []sc2.ChromeFlag
	for _, item := range d.List() {
		data := item.(map[string]interface{})
		flags = append(flags, sc2.ChromeFlag{
			Name:  data["name"].(string),
			Value: data["value"].(string),
		})
	}
	return flags
}

func buildCookiesData(cookies *schema.Set) []sc2.Cookiesv2 {
	cookiesList := make([]sc2.Cookiesv2, len(cookies.List()))

	for i, cookie := range cookies.List() {
		cookie := cookie.(map[string]interface{})
		if cookie != nil {
			set := sc2.Cookiesv2{
				Key:    cookie["key"].(string),
				Value:  cookie["value"].(string),
				Domain: cookie["domain"].(string),
				Path:   cookie["path"].(string),
			}
			cookiesList[i] = set
		}

	}
	return cookiesList
}

func buildHostOverridesData(hostOverrides *schema.Set) []sc2.HostOverrides {
	hostOverridesList := make([]sc2.HostOverrides, len(hostOverrides.List()))

	for i, hostOverride := range hostOverrides.List() {
		hostOverride := hostOverride.(map[string]interface{})
		set := sc2.HostOverrides{
			Source:         hostOverride["source"].(string),
			Target:         hostOverride["target"].(string),
			KeepHostHeader: hostOverride["keep_host_header"].(bool),
		}

		hostOverridesList[i] = set

	}
	return hostOverridesList
}

func buildAuthenticationData(authentication *schema.Set) *sc2.Authentication {
	authentication_list := authentication.List()

	if len(authentication_list) > 0 {
		authentication_map := authentication_list[0].(map[string]interface{})
		authenticationData := &sc2.Authentication{
			Username: authentication_map["username"].(string),
			Password: authentication_map["password"].(string),
		}
		return authenticationData
	}
	return nil
}

func flattenLinkData(checkLinks *sc.Links) []interface{} {
	links := make(map[string]interface{})

	if checkLinks.Self != "" {
		links["self"] = checkLinks.Self
	}
	if checkLinks.SelfHTML != "" {
		links["self_html"] = checkLinks.SelfHTML
	}
	if checkLinks.Metrics != "" {
		links["metrics"] = checkLinks.Metrics
	}
	if checkLinks.LastRun != "" {
		links["last_run"] = checkLinks.LastRun
	}

	return []interface{}{links}
}

func flattenStatusData(checkStatus *sc.Status) []interface{} {
	status := make(map[string]interface{})

	status["last_code"] = checkStatus.LastCode
	status["last_message"] = checkStatus.LastMessage
	status["last_response_time"] = checkStatus.LastResponseTime
	status["last_run_at"] = checkStatus.LastRunAt
	status["last_failure_at"] = checkStatus.LastFailureAt
	status["last_alert_at"] = checkStatus.LastAlertAt
	status["has_failure"] = checkStatus.HasFailure
	status["has_location_failure"] = checkStatus.HasLocationFailure

	return []interface{}{status}
}

func buildTagsData(d *schema.ResourceData) []string {
	tags := d.Get("tags").([]interface{})
	tagsList := make([]string, len(tags))
	for i, tag := range tags {
		tagsList[i] = tag.(string)
	}
	return tagsList
}

func flattenTagsData(checkTags *sc.Tags) []interface{} {
	if checkTags != nil {
		cls := make([]interface{}, len(*checkTags))

		for i, checkTags := range *checkTags {
			cl := make(map[string]interface{})

			cl["id"] = checkTags.ID
			cl["name"] = checkTags.Name

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)

}

func flattenBlackoutData(checkBlackout *sc.BlackoutPeriods) []interface{} {
	if checkBlackout != nil {
		cls := make([]interface{}, len(*checkBlackout))

		for i, checkBlackout := range *checkBlackout {
			cl := make(map[string]interface{})

			cl["start_date"] = checkBlackout.StartDate
			cl["end_date"] = checkBlackout.EndDate
			cl["timezone"] = checkBlackout.Timezone
			cl["start_time"] = checkBlackout.StartTime
			cl["end_time"] = checkBlackout.EndTime
			cl["repeat_type"] = checkBlackout.RepeatType
			cl["duration_in_minutes"] = checkBlackout.DurationInMinutes
			cl["is_repeat"] = checkBlackout.IsRepeat
			cl["monthly_repeat_type"] = checkBlackout.MonthlyRepeatType
			cl["created_at"] = checkBlackout.CreatedAt
			cl["updated_at"] = checkBlackout.UpdatedAt

			cls[i] = cl
		}
		return cls
	}

	return make([]interface{}, 0)
}

func buildNotificationsData(notifications sc.Notifications, d *schema.ResourceData) sc.Notifications {
	notificationData := d.Get("notifications").(*schema.Set).List()
	for _, notif := range notificationData {
		notif := notif.(map[string]interface{})
		notifications.Sms = notif["sms"].(bool)
		notifications.Call = notif["call"].(bool)
		notifications.Email = notif["email"].(bool)
		notifications.NotifyAfterFailureCount = notif["notify_after_failure_count"].(int)
		notifications.NotifyOnLocationFailure = notif["notify_on_location_failure"].(bool)
		notifications.NotifyWho = buildNotifyWhoData(notif["notify_who"].(*schema.Set).List())
		notifications.Escalations = buildEscalationsData(notif["escalations"].(*schema.Set).List())
	}
	return notifications
}

func flattenNotificationsData(checkNotifications *sc.Notifications) []interface{} {
	notifications := make(map[string]interface{})

	notifications["sms"] = checkNotifications.Sms
	notifications["call"] = checkNotifications.Call
	notifications["email"] = checkNotifications.Email
	notifications["notify_after_failure_count"] = checkNotifications.NotifyAfterFailureCount
	notifications["notify_on_location_failure"] = checkNotifications.NotifyOnLocationFailure
	notifications["muted"] = checkNotifications.Muted

	checkNotifyWho := flattenNotifyWhoData(checkNotifications.NotifyWho)
	notifications["notify_who"] = checkNotifyWho

	checkNotificationWindows := flattenNotificationWindowsData(&checkNotifications.NotificationWindows)
	notifications["notification_windows"] = checkNotificationWindows

	checkEscalations := flattenEscalationsData(checkNotifications.Escalations)
	notifications["escalations"] = checkEscalations

	return []interface{}{notifications}
}

func buildNotifyWhoData(notifyWhoCrit []interface{}) []sc.NotifyWho {
	notifyWhoList := make([]sc.NotifyWho, len(notifyWhoCrit))
	for i, notifyWho := range notifyWhoCrit {
		notifyWho := notifyWho.(map[string]interface{})
		notif := sc.NotifyWho{
			Sms:             notifyWho["sms"].(bool),
			Call:            notifyWho["call"].(bool),
			Email:           notifyWho["email"].(bool),
			CustomUserEmail: notifyWho["custom_user_email"].(string),
			Type:            notifyWho["type"].(string),
			ID:              notifyWho["id"].(int),
		}
		notifyWhoList[i] = notif

	}
	return notifyWhoList
}

func flattenNotifyWhoData(checkNotifyWho []sc.NotifyWho) []interface{} {
	if checkNotifyWho != nil {
		cls := make([]interface{}, len(checkNotifyWho))

		for i, checkNotifyWho := range checkNotifyWho {
			cl := make(map[string]interface{})

			if val := checkNotifyWho.Sms; val {
				cl["sms"] = checkNotifyWho.Sms
			}
			if val := checkNotifyWho.Call; val {
				cl["call"] = checkNotifyWho.Call
			}
			if val := checkNotifyWho.Email; val {
				cl["email"] = checkNotifyWho.Email
			}
			if checkNotifyWho.CustomUserEmail != "" {
				cl["custom_user_email"] = checkNotifyWho.CustomUserEmail
			}
			if checkNotifyWho.Type != "" {
				cl["type"] = checkNotifyWho.Type
			}
			if checkNotifyWho.ID != 0 {
				cl["id"] = checkNotifyWho.ID
			}

			checkNotifyWhoLinks := flattenLinkData(&checkNotifyWho.Links)
			cl["links"] = checkNotifyWhoLinks

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenNotificationWindowsData(checkNotificationWindows *sc.NotificationWindows) []interface{} {
	if checkNotificationWindows != nil {
		cls := make([]interface{}, len(*checkNotificationWindows))

		for i, checkNotificationWindows := range *checkNotificationWindows {
			cl := make(map[string]interface{})

			cl["start_time"] = checkNotificationWindows.StartTime
			cl["end_time"] = checkNotificationWindows.EndTime
			cl["duration_in_minutes"] = checkNotificationWindows.DurationInMinutes
			cl["time_zone"] = checkNotificationWindows.TimeZone

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func flattenNotificationWindowData(checkNotificationWindow *sc.NotificationWindow) []interface{} {
	notificationWindow := make(map[string]interface{})

	notificationWindow["start_time"] = checkNotificationWindow.StartTime
	notificationWindow["end_time"] = checkNotificationWindow.EndTime
	notificationWindow["duration_in_minutes"] = checkNotificationWindow.DurationInMinutes
	notificationWindow["time_zone"] = checkNotificationWindow.TimeZone

	return []interface{}{notificationWindow}
}

func buildConnectionData(d *schema.ResourceData) sc.Connection {
	connectionData := d.Get("check_connection").(*schema.Set).List()
	var connection sc.Connection
	for _, connect := range connectionData {
		connect := connect.(map[string]interface{})
		connection.DownloadBandwidth = connect["download_bandwidth"].(int)
		connection.UploadBandwidth = connect["upload_bandwidth"].(int)
		connection.Latency = connect["latency"].(int)
		connection.PacketLoss = connect["packet_loss"].(float64)
	}
	return connection
}

func flattenConnectionData(checkConnection *sc.Connection) []interface{} {
	connection := make(map[string]interface{})

	connection["download_bandwidth"] = checkConnection.DownloadBandwidth
	connection["upload_bandwidth"] = checkConnection.UploadBandwidth
	connection["latency"] = checkConnection.Latency
	connection["packet_loss"] = checkConnection.PacketLoss

	return []interface{}{connection}
}

func buildIntegrationsData(d *schema.ResourceData) []int {
	integrations := d.Get("integrations").([]interface{})
	integrationList := make([]int, len(integrations))
	for i, integration := range integrations {
		integrationList[i] = integration.(int)
	}
	return integrationList
}

func flattenIntegrationsData(checkIntegrations *sc.Integrations) []interface{} {
	if checkIntegrations != nil {
		cls := make([]interface{}, len(*checkIntegrations))

		for i, checkIntegrations := range *checkIntegrations {
			cl := make(map[string]interface{})

			cl["id"] = checkIntegrations.ID
			cl["name"] = checkIntegrations.Name

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)

}

func buildLocationsData(d *schema.ResourceData) []int {
	locations := d.Get("locations").([]interface{})
	locationList := make([]int, len(locations))
	for i, location := range locations {
		locationList[i] = location.(int)
	}
	return locationList
}

func flattenLocationsData(checkLocations *sc.Locations) []interface{} {
	if checkLocations != nil {
		cls := make([]interface{}, len(*checkLocations))

		for i, checkLocations := range *checkLocations {
			cl := make(map[string]interface{})

			cl["id"] = checkLocations.ID
			cl["name"] = checkLocations.Name
			cl["world_region"] = checkLocations.WorldRegion
			cl["region_code"] = checkLocations.RegionCode

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)
}

func buildSuccessCriteriaData(d *schema.ResourceData) []sc.SuccessCriteria {

	successCrit := d.Get("success_criteria").(*schema.Set).List()
	successList := make([]sc.SuccessCriteria, len(successCrit))
	for i, success := range successCrit {
		success := success.(map[string]interface{})
		suc := sc.SuccessCriteria{
			ActionType:       success["action_type"].(string),
			ComparisonString: success["comparison_string"].(string),
			CreatedAt:        success["created_at"].(string),
			UpdatedAt:        success["updated_at"].(string),
		}
		successList[i] = suc
	}
	return successList
}

func flattenSuccessCriteriaData(checkSuccessCriteria *[]sc.SuccessCriteria) []interface{} {
	if checkSuccessCriteria != nil {
		cls := make([]interface{}, len(*checkSuccessCriteria))

		for i, checkSuccessCriteria := range *checkSuccessCriteria {
			cl := make(map[string]interface{})

			cl["action_type"] = checkSuccessCriteria.ActionType
			cl["created_at"] = checkSuccessCriteria.CreatedAt
			cl["updated_at"] = checkSuccessCriteria.UpdatedAt
			cl["comparison_string"] = checkSuccessCriteria.ComparisonString

			cls[i] = cl
		}

		return cls
	}

	return make([]interface{}, 0)

}

func buildEscalationsData(escalations []interface{}) []sc.Escalations {
	escalationsList := make([]sc.Escalations, len(escalations))
	for i, escalation := range escalations {
		escalation := escalation.(map[string]interface{})
		esca := sc.Escalations{
			Sms:          escalation["sms"].(bool),
			Email:        escalation["email"].(bool),
			Call:         escalation["call"].(bool),
			AfterMinutes: escalation["after_minutes"].(int),
			NotifyWho:    buildNotifyWhoData(escalation["notify_who"].(*schema.Set).List()),
		}
		escalationsList[i] = esca

	}
	return escalationsList
}

func flattenEscalationsData(checkEscalations []sc.Escalations) []interface{} {
	if checkEscalations != nil {
		cls := make([]interface{}, len(checkEscalations))

		for i, checkEscalations := range checkEscalations {
			cl := make(map[string]interface{})

			cl["sms"] = checkEscalations.Sms
			cl["call"] = checkEscalations.Call
			cl["email"] = checkEscalations.Email
			cl["after_minutes"] = checkEscalations.AfterMinutes
			cl["is_repeat"] = checkEscalations.IsRepeat
			checkNotifyWho := flattenNotifyWhoData(checkEscalations.NotifyWho)
			cl["notify_who"] = checkNotifyWho
			checkNotificationWindow := flattenNotificationWindowData(&checkEscalations.NotificationWindow)
			cl["notification_window"] = checkNotificationWindow

			cls[i] = cl
		}
		return cls
	}

	return make([]interface{}, 0)
}

func buildViewportData(d *schema.ResourceData) sc.Viewport {
	viewportData := d.Get("viewport").(*schema.Set).List()
	var viewport sc.Viewport
	for _, view := range viewportData {
		view := view.(map[string]interface{})
		viewport.Height = view["height"].(int)
		viewport.Width = view["width"].(int)
	}
	return viewport
}

func buildStepData(d *schema.ResourceData) []sc.Steps {
	// This part of Rigor is not accessible from the public API and does not currently work.
	steps := d.Get("steps").(*schema.Set).List()
	stepsList := make([]sc.Steps, len(steps))
	for i, step := range steps {
		step := step.(map[string]interface{})
		ste := sc.Steps{
			ItemMethod:   step["item_method"].(string),
			Value:        step["value"].(string),
			How:          step["how"].(string),
			What:         step["what"].(string),
			VariableName: step["variable_name"].(string),
			Name:         step["name"].(string),
			Position:     step["position"].(int),
		}
		stepsList[i] = ste
	}
	return stepsList
}

func flattenStepData(checkSteps []sc.Steps) []interface{} {
	if checkSteps != nil {
		steps := make([]interface{}, len(checkSteps))

		for i, checkStep := range checkSteps {
			step := make(map[string]interface{})

			step["item_method"] = checkStep.ItemMethod
			step["value"] = checkStep.Value
			step["how"] = checkStep.How
			step["what"] = checkStep.What
			step["variable_name"] = checkStep.VariableName
			step["name"] = checkStep.Name
			step["position"] = checkStep.Position

			steps[i] = step
		}

		return steps
	}

	return make([]interface{}, 0)
}

func buildCookieData(d *schema.ResourceData) []sc.Cookies {

	cookies := d.Get("cookies").(*schema.Set).List()
	cookiesList := make([]sc.Cookies, len(cookies))
	for i, cookie := range cookies {
		cookie := cookie.(map[string]interface{})
		cke := sc.Cookies{
			Key:    cookie["key"].(string),
			Value:  cookie["value"].(string),
			Domain: cookie["domain"].(string),
			Path:   cookie["path"].(string),
		}
		cookiesList[i] = cke
	}
	return cookiesList
}

func buildDnsOverridesData(d *schema.ResourceData) sc.DNSOverrides {
	dnsOverridesData := d.Get("dns_overrides").(*schema.Set).List()
	var dnsOverrides sc.DNSOverrides
	for _, dns := range dnsOverridesData {
		dns := dns.(map[string]interface{})
		dnsOverrides.OriginalDomainCom = dns["original_domain"].(string)
		dnsOverrides.OriginalHostCom = dns["original_host"].(string)
	}
	return dnsOverrides
}

func buildThresholdMonitorsData(d *schema.ResourceData) []sc.ThresholdMonitors {

	thresholdMonitors := d.Get("threshold_monitors").(*schema.Set).List()
	thresholdMonitorsList := make([]sc.ThresholdMonitors, len(thresholdMonitors))
	for i, thresholdMonitor := range thresholdMonitors {
		thresholdMonitor := thresholdMonitor.(map[string]interface{})
		thm := sc.ThresholdMonitors{
			Matcher:        thresholdMonitor["matcher"].(string),
			MetricName:     thresholdMonitor["metric_name"].(string),
			ComparisonType: thresholdMonitor["comparison_type"].(string),
			Value:          thresholdMonitor["value"].(int),
		}
		thresholdMonitorsList[i] = thm
	}
	return thresholdMonitorsList
}

func buildJavascriptFilesData(d *schema.ResourceData) []sc.JavascriptFiles {
	// This part of Rigor is not accessible from the public API and does not currently work.
	javascriptFiles := d.Get("javascript_files").(*schema.Set).List()
	javascriptFilesList := make([]sc.JavascriptFiles, len(javascriptFiles))
	for i, javascriptFile := range javascriptFiles {
		javascriptFile := javascriptFile.(map[string]interface{})
		thm := sc.JavascriptFiles{
			ID:   javascriptFile["id"].(int),
			Name: javascriptFile["name"].(string),
		}
		javascriptFilesList[i] = thm
	}
	return javascriptFilesList
}

func buildExcludedFilesData(d *schema.ResourceData) []sc.ExcludedFiles {
	excludedFiles := d.Get("excluded_files").(*schema.Set).List()
	excludedFilesList := make([]sc.ExcludedFiles, len(excludedFiles))
	for i, excludedFile := range excludedFiles {
		excludedFile := excludedFile.(map[string]interface{})
		exf := sc.ExcludedFiles{
			ExclusionType: excludedFile["exclusion_type"].(string),
			PresetName:    excludedFile["preset_name"].(string),
			URL:           excludedFile["pattern"].(string),
		}
		excludedFilesList[i] = exf
	}
	return excludedFilesList
}
