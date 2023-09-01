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
	"log"
	"strings"

	sc2 "github.com/splunk/syntheticsclient/v2/syntheticsclientv2"

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

	if checkApiV2.Test.Frequency != 0 {
		apiV2["frequency"] = checkApiV2.Test.Frequency
	}

	if checkApiV2.Test.Name != "" {
		apiV2["name"] = checkApiV2.Test.Name
	}

	if checkApiV2.Test.Schedulingstrategy != "" {
		apiV2["scheduling_strategy"] = checkApiV2.Test.Schedulingstrategy
	}

	apiV2["device_id"] = checkApiV2.Test.Device.ID
	

	locationIds := flattenLocationData(&checkApiV2.Test.Locationids)
	apiV2["location_ids"] = locationIds

	requests := flattenRequestData(&checkApiV2.Test.Requests)
	apiV2["requests"] = requests

	log.Println("[DEBUG] apiv2 data: ", apiV2)

	return []interface{}{apiV2}
}

func flattenApiV2Data(checkApiV2 *sc2.ApiCheckV2Response) []interface{} {
	apiV2 := make(map[string]interface{})

	apiV2["active"] = checkApiV2.Test.Active

	if checkApiV2.Test.Createdat.IsZero() {
	} else {
		apiV2["created_at"] = checkApiV2.Test.Createdat.String()
	}

	if checkApiV2.Test.Updatedat.IsZero() {
	} else {
		apiV2["updated_at"] = checkApiV2.Test.Updatedat.String()
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

	device := flattenDeviceData(&checkApiV2.Test.Device)
	apiV2["device"] = device

	locationIds := flattenLocationData(&checkApiV2.Test.Locationids)
	apiV2["location_ids"] = locationIds

	requests := flattenRequestData(&checkApiV2.Test.Requests)
	apiV2["requests"] = requests

	log.Println("[DEBUG] apiv2 data: ", apiV2)

	return []interface{}{apiV2}
}

func flattenVariableV2Read(checkVariableV2 *sc2.VariableV2Response) []interface{} {
	variableV2 := make(map[string]interface{})

	variableV2["name"] = checkVariableV2.Variable.Name

	variableV2["description"] = checkVariableV2.Variable.Description

	variableV2["value"] = checkVariableV2.Variable.Value

	variableV2["secret"] = checkVariableV2.Variable.Secret

	log.Println("[DEBUG] VARIABLE V2 data: ", variableV2)

	return []interface{}{variableV2}
}

func flattenVariableV2Data(checkVariableV2 *sc2.VariableV2Response) []interface{} {
	variableV2 := make(map[string]interface{})

	variableV2["name"] = checkVariableV2.Variable.Name

	variableV2["id"] = checkVariableV2.Variable.ID

	variableV2["description"] = checkVariableV2.Variable.Description

	variableV2["value"] = checkVariableV2.Variable.Value

	variableV2["secret"] = checkVariableV2.Variable.Secret

	if checkVariableV2.Variable.Createdat.IsZero() {
	} else {
		variableV2["created_at"] = checkVariableV2.Variable.Createdat.String()
	}

	if checkVariableV2.Variable.Updatedat.IsZero() {
	} else {
		variableV2["updated_at"] = checkVariableV2.Variable.Updatedat.String()
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

	browserV2["device_id"] = checkBrowserV2.Test.Device.ID

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

	log.Println("[DEBUG] read browserv2 data: ", browserV2)

	return []interface{}{browserV2}
}

func flattenBrowserV2Data(checkBrowserV2 *sc2.BrowserCheckV2Response) []interface{} {
	browserV2 := make(map[string]interface{})

	browserV2["active"] = checkBrowserV2.Test.Active

	if checkBrowserV2.Test.Createdat.IsZero() {
	} else {
		browserV2["created_at"] = checkBrowserV2.Test.Createdat.String()
	}

	if checkBrowserV2.Test.Updatedat.IsZero() {
	} else {
		browserV2["updated_at"] = checkBrowserV2.Test.Updatedat.String()
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

	device := flattenDeviceData(&checkBrowserV2.Test.Device)
	browserV2["device"] = device

	advancedSettings := flattenAdvancedSettingsData(&checkBrowserV2.Test.Advancedsettings)
	browserV2["advanced_settings"] = advancedSettings

	businessTranscations := flattenBusinessTransactionsData(&checkBrowserV2.Test.Transactions)
	browserV2["transactions"] = businessTranscations

	transcations := flattenTransactionsData(&checkBrowserV2.Test.Transactions)
	browserV2["transactions"] = transcations

	log.Println("[DEBUG] flatten browserv2 data: ", browserV2)

	return []interface{}{browserV2}
}

func flattenHttpV2Read(checkHttpV2 *sc2.HttpCheckV2Response) []interface{} {
	httpV2 := make(map[string]interface{})


	if checkHttpV2.Test.Name != "" {
		httpV2["name"] = checkHttpV2.Test.Name
	}

	httpV2["active"] = checkHttpV2.Test.Active

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

	log.Println("[DEBUG] httpV2 data: ", httpV2)

	return []interface{}{httpV2}
}

func flattenHttpV2Data(checkHttpV2 *sc2.HttpCheckV2Response) []interface{} {
	httpV2 := make(map[string]interface{})

	if checkHttpV2.Test.ID != 0 {
		httpV2["id"] = checkHttpV2.Test.ID
	}

	if checkHttpV2.Test.Name != "" {
		httpV2["name"] = checkHttpV2.Test.Name
	}

	httpV2["active"] = checkHttpV2.Test.Active

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

	if checkHttpV2.Test.SchedulingStrategy != "" {
		httpV2["scheduling_strategy"] = checkHttpV2.Test.SchedulingStrategy
	}

	if checkHttpV2.Test.Type != "" {
		httpV2["type"] = checkHttpV2.Test.Type
	}

	if checkHttpV2.Test.URL != "" {
		httpV2["url"] = checkHttpV2.Test.URL
	}

	if checkHttpV2.Test.RequestMethod != "" {
		httpV2["request_method"] = checkHttpV2.Test.RequestMethod
	}

	if checkHttpV2.Test.Body != "" {
		httpV2["body"] = checkHttpV2.Test.Body
	}

	if *checkHttpV2.Test.UserAgent != "" {
		httpV2["user_agent"] = checkHttpV2.Test.UserAgent
	}

	httpV2["verify_certificates"] = checkHttpV2.Test.Verifycertificates

	locationIds := flattenLocationData(&checkHttpV2.Test.LocationIds)
	httpV2["location_ids"] = locationIds

	httpHeaders := flattenHttpHeadersData(&checkHttpV2.Test.HttpHeaders)
	httpV2["headers"] = httpHeaders

	validations := flattenValidationsData(&checkHttpV2.Test.Validations)
	httpV2["validations"] = validations

	log.Println("[DEBUG] httpV2 data: ", httpV2)

	return []interface{}{httpV2}
}

func flattenPortCheckV2Read(checkPortV2 *sc2.PortCheckV2Response) []interface{} {
	portV2 := make(map[string]interface{})

	if checkPortV2.Test.Name != "" {
		portV2["name"] = checkPortV2.Test.Name
	}

	portV2["active"] = checkPortV2.Test.Active

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

	log.Println("[DEBUG] portv2 data: ", portV2)

	return []interface{}{portV2}
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

			if checkStep.Selector != "" {
				cl["selector"] = checkStep.Selector
			}

			if checkStep.SelectorType != "" {
				cl["selector_type"] = checkStep.SelectorType
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

			if checkStep.SelectorType != "" {
				cl["selector_type"] = checkStep.SelectorType
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

func flattenOptionsData(checkAuthentications *sc2.Options) []interface{} {
	authentication := make(map[string]interface{})

	authentication["url"] = checkAuthentications.URL

	return []interface{}{authentication}
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
			apiv2.Test.Locationids = buildLocationIdData(api["location_ids"].([]interface{}))
			apiv2.Test.Name = api["name"].(string)
			apiv2.Test.Requests = buildRequestsData(api["requests"].(*schema.Set))
			apiv2.Test.Schedulingstrategy = api["scheduling_strategy"].(string)
		}
	}
	log.Println("[DEBUG] build apiv2 data: ", apiv2)
	return apiv2
}

func buildBrowserV2Data(d *schema.ResourceData) sc2.BrowserCheckV2Input {
	var browserv2 sc2.BrowserCheckV2Input
	browserv2Data := d.Get("test").(*schema.Set).List()
	for _, browser := range browserv2Data {
		browser := browser.(map[string]interface{})
		if browser["name"].(string) != "" {
			browserv2.Test.Active = browser["active"].(bool)
			browserv2.Test.DeviceID = browser["device_id"].(int)
			browserv2.Test.Frequency = browser["frequency"].(int)
			browserv2.Test.LocationIds = buildLocationIdData(browser["location_ids"].([]interface{}))
			browserv2.Test.Name = browser["name"].(string)
			browserv2.Test.Transactions = buildBusinessTransactionsData(browser["transactions"].([]interface{}))
			browserv2.Test.Schedulingstrategy = browser["scheduling_strategy"].(string)
			browserv2.Test.Advancedsettings = buildAdvancedSettingsData(browser["advanced_settings"].(*schema.Set))
		}
	}
	log.Println("[DEBUG] build browserv2 data:")
	return browserv2
}

func buildHttpV2Data(d *schema.ResourceData) sc2.HttpCheckV2Input {
	var httpv2 sc2.HttpCheckV2Input
	httpv2Data := d.Get("test").(*schema.Set).List()
	var i = 0
	for _, http := range httpv2Data {
		if i < 1 {
			http := http.(map[string]interface{})
			httpv2.Test.Name = http["name"].(string)
			httpv2.Test.Type = http["type"].(string)
			httpv2.Test.URL = http["url"].(string)
			httpv2.Test.LocationIds = buildLocationIdData(http["location_ids"].([]interface{}))
			httpv2.Test.Frequency = http["frequency"].(int)
			httpv2.Test.SchedulingStrategy = http["scheduling_strategy"].(string)
			httpv2.Test.Active = http["active"].(bool)
			httpv2.Test.RequestMethod = http["request_method"].(string)
			httpv2.Test.Body = http["body"].(string)
			httpv2.Test.Verifycertificates = http["verify_certificates"].(bool)
			userAgentString := http["user_agent"].(string)
			httpv2.Test.UserAgent = &userAgentString
			httpv2.Test.HttpHeaders = buildHttpHeadersData(http["headers"].(*schema.Set))
			httpv2.Test.Validations = buildValidationsData(http["validations"].(*schema.Set))
			i++
		}
	}
	log.Println("[DEBUG] build httpv2 data: ", httpv2)
	return httpv2
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
			portv2.Test.SchedulingStrategy = port["scheduling_strategy"].(string)
			portv2.Test.Active = port["active"].(bool)
			i++

		}
	}
	log.Println("[DEBUG] build portv2 data: ", portv2)
	return portv2
}

func buildVariableV2Data(d *schema.ResourceData) sc2.VariableV2Input {
	var variablev2 sc2.VariableV2Input
	variablev2Data := d.Get("variable").(*schema.Set).List()
	var i = 0
	for _, variable := range variablev2Data {
		if i < 1 {
			variable := variable.(map[string]interface{})
			variablev2.Variable.Description = variable["description"].(string)
			variablev2.Variable.Name = variable["name"].(string)
			variablev2.Variable.Secret = variable["secret"].(bool)
			variablev2.Variable.Value = variable["value"].(string)
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

func buildRequestsData(requests *schema.Set) []sc2.Requests {
	requestsList := make([]sc2.Requests, len(requests.List()))

	for i, request := range requests.List() {
		request := request.(map[string]interface{})
		req := sc2.Requests{
			Configuration: buildConfigurationData(request["configuration"].(*schema.Set)),
			Setup:         buildSetupData(request["setup"].(*schema.Set)),
			Validations:   buildValidationsData(request["validations"].(*schema.Set)),
		}
		requestsList[i] = req

	}
	return requestsList
}

func buildBusinessTransactionsData(businessTransactions []interface{}) []sc2.Transactions {
	businessTransactionsList := make([]sc2.Transactions, len(businessTransactions))
	for i, bisTrans := range businessTransactions {
		bisTrans := bisTrans.(map[string]interface{})
		transaction := sc2.Transactions{
			Name:    bisTrans["name"].(string),
			StepsV2: buildStepV2Data(bisTrans["steps"].([]interface{})),
		}
		businessTransactionsList[i] = transaction
	}
	return businessTransactionsList
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


func buildStepV2Data(steps []interface{}) []sc2.StepsV2 {
	stepsList := make([]sc2.StepsV2, len(steps))
	for i, step := range steps {
		step := step.(map[string]interface{})
		st := sc2.StepsV2{
			URL:          step["url"].(string),
			Name:         step["name"].(string),
			Type:         step["type"].(string),
			WaitForNav:   step["wait_for_nav"].(bool),
			SelectorType: step["selector_type"].(string),
			Selector:     step["selector"].(string),
			OptionSelectorType: step["option_selector_type"].(string),
			OptionSelector:     step["option_selector"].(string),
			VariableName: step["variable_name"].(string),
			Value:     step["value"].(string),
			Duration:     step["duration"].(int),
		}
		stepsList[i] = st

	}
	return stepsList
}

func buildSetupData(setups *schema.Set) []sc2.Setup {
	setupsList := make([]sc2.Setup, len(setups.List()))

	for i, setup := range setups.List() {
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

func buildValidationsData(validations *schema.Set) []sc2.Validations {
	validationsList := make([]sc2.Validations, len(validations.List()))

	for i, validation := range validations.List() {
		validation := validation.(map[string]interface{})
		val := sc2.Validations{
			Actual:     validation["actual"].(string),
			Comparator: validation["comparator"].(string),
			Expected:   validation["expected"].(string),
			Name:       validation["name"].(string),
			Type:       validation["type"].(string),
			Extractor:  validation["extractor"].(string),
			Source:     validation["source"].(string),
			Variable:   validation["variable"].(string),
			Code:       validation["code"].(string),
			Value:      validation["value"].(string),
		}

		validationsList[i] = val

	}
	return validationsList
}

func buildConfigurationData(configuration *schema.Set) sc2.Configuration {
	var configurationData sc2.Configuration

	config_list := configuration.List()
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

func buildOptionsData(options *schema.Set) sc2.Options {
	var optionsData sc2.Options

	options_list := options.List()
	if len(options_list) > 0 {
		options_map := options_list[0].(map[string]interface{})

		optionsData.URL = options_map["url"].(string)
	}

	return optionsData
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
