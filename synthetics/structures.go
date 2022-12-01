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
	sc2 "syntheticsclientv2"

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

func flattenApiV2Data(checkApiV2 *sc2.ApiCheckV2Response) []interface{} {
	apiV2 := make(map[string]interface{})

	apiV2["active"] = checkApiV2.Test.Active
	
	if checkApiV2.Test.Createdat.IsZero() {
		}else{
			apiV2["created_at"] = checkApiV2.Test.Createdat.String()	
	}
	
	if checkApiV2.Test.Updatedat.IsZero() {
		}else{
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

	log.Printf("[WARN] XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	log.Println(checkApiV2)
	
	log.Printf("[WARN] XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	log.Println(apiV2)

	return []interface{}{apiV2}
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

func flattenConfigurationData(checkConfiguration *sc2.Configuration) []interface{} {
	configuration := make(map[string]interface{})

	if checkConfiguration.Body != "" {
		configuration["body"] = checkConfiguration.Body
	}
	if checkConfiguration.Name != "" {
		configuration["name"] = checkConfiguration.Name
	}

	if checkConfiguration.Requestmethod != "" {
		configuration["request_method"] = checkConfiguration.Requestmethod
	}
	if checkConfiguration.URL != "" {
		configuration["url"] = checkConfiguration.URL
	}

	headers := flattenHeaderData(&checkConfiguration.Headers)
	configuration["headers"] = headers

	return []interface{}{configuration}
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

func flattenNetworkConnectionData(checkNetworkConnection *sc2.Networkconnection) []interface{} {
	networkConnection := make(map[string]interface{})

	networkConnection["description"] = checkNetworkConnection.Description
	networkConnection["download_bandwidth"] = checkNetworkConnection.Downloadbandwidth
	networkConnection["latency"] = checkNetworkConnection.Latency
	networkConnection["packet_loss"] = checkNetworkConnection.Packetloss
	networkConnection["upload_bandwidth"] = checkNetworkConnection.Uploadbandwidth

	return []interface{}{networkConnection}
}

func buildApiV2Data(d *schema.ResourceData) sc2.ApiCheckV2Input {
	var apiv2 sc2.ApiCheckV2Input
	apiv2Data := d.Get("test").(*schema.Set).List()
	// log.Printf("[WARN] 11111***** apiv2Data in *****11111*****11111*****11111*****11111*****11111*****")
	// log.Println(apiv2Data)
	var i = 0
	for _, api := range apiv2Data {
		if i < 1 {
			api := api.(map[string]interface{})
			apiv2.Test.Active = api["active"].(bool)
			apiv2.Test.Deviceid = api["device_id"].(int)
			apiv2.Test.Frequency = api["frequency"].(int)
			apiv2.Test.Locationids = buildLocationIdData(api["location_ids"].([]interface{}))
			apiv2.Test.Name = api["name"].(string)
			apiv2.Test.Requests = buildRequestsData(api["requests"].(*schema.Set))
			apiv2.Test.Schedulingstrategy = api["scheduling_strategy"].(string)
			i++
			log.Println("[WARN] ------------------ I is now: ", i)
			log.Println(apiv2Data)
		}
	}
	log.Printf("[WARN] 11111***** apiv2 out *****11111*****11111*****11111*****11111*****11111*****")
	log.Println(apiv2)
	return apiv2
}

func buildLocationIdData(d []interface{}) []string {
	log.Printf("[WARN] *****&&*& d *schema.ResourceData ****************")
	log.Println(d)
	locationsList := make([]string, len(d))
	for i, locations := range d {
		log.Printf("[WARN] *****&&*& d *schema.ResourceData ****************")
		log.Println(i)
		log.Println(locations)
		locationsList[i] = locations.(string)
	}
	return locationsList
}

func buildRequestsData(requests *schema.Set) []sc2.Requests {
	requestsList := make([]sc2.Requests, len(requests.List()))
	log.Printf("[WARN] requests *************")
	log.Println(requests.List())

	for i, request := range requests.List() {
		request := request.(map[string]interface{})
		req := sc2.Requests{
			Configuration: buildConfigurationData(request["configuration"].(*schema.Set)),
			Setup: buildSetupData(request["setup"].(*schema.Set)),
			Validations: buildValidationsData(request["validations"].(*schema.Set)),
		}
		log.Printf("[WARN] iiiiiii")
		log.Println(i)
		requestsList[0] = req

	}
	return requestsList
}

func buildSetupData(setups *schema.Set) []sc2.Setup {
	setupsList := make([]sc2.Setup, len(setups.List()))
	log.Printf("[WARN] SETUP *************")
	log.Println(setups.List())

	for i, setup := range setups.List() {
		setup := setup.(map[string]interface{})
		set := sc2.Setup{
			Extractor:          setup["extractor"].(string),
			Name:          setup["name"].(string),
			Source:          setup["source"].(string),
			Type:          setup["type"].(string),
			Variable:          setup["variable"].(string),
		}
		log.Printf("[WARN] *********$$$$$$  SETUP!!! $$$$$**********")
		log.Println(i)
		setupsList[0] = set

	}
	return setupsList
}

func buildValidationsData(validations *schema.Set) []sc2.Validations {
	validationsList := make([]sc2.Validations, len(validations.List()))
	log.Printf("[WARN] Validation *************")
	log.Println(validations.List())

	for i, validation := range validations.List() {
		validation := validation.(map[string]interface{})
		val := sc2.Validations{
			Actual:          validation["actual"].(string),
			Comparator:          validation["comparator"].(string),
			Expected:          validation["expected"].(string),
			Name:          validation["name"].(string),
			Type:          validation["type"].(string),
		}
		log.Printf("[WARN] *********$$$$$$  Validation!!! $$$$$**********")
		log.Println(i)
		validationsList[0] = val

	}
	return validationsList
}

func buildConfigurationData(configuration *schema.Set) sc2.Configuration {
	var configurationData sc2.Configuration
	
	// for k, v := range configuration.List() {
	// log.Printf("[WARN] 11111***** configuration K, V *****11111*****11111*****11111*****11111*****11111*****")
	// configurationData := v
	// log.Println(k)
	// log.Println(v)
	// log.Println(configurationData)
	// }

	config_list := configuration.List()
	config_map := config_list[0].(map[string]interface{})

	configurationData.Body = config_map["body"].(string)
	configurationData.Headers = config_map["headers"].(map[string]interface{})
	configurationData.Name = config_map["name"].(string)
	configurationData.Requestmethod = config_map["request_method"].(string)
	configurationData.URL = config_map["url"].(string)
	
	return configurationData
}

// func buildHeadersData(headers map[string]interface{}) sc2.Headers {
// 	if headers != nil {
// 		cls := make(map[string]interface{}, len(headers))
// 		for k, v := range headers {
// 			log.Printf("[WARN] 11111***** headersData K , V *****11111*****11111*****11111*****11111*****11111*****")
// 			log.Println(k)
// 			log.Println(v)
// 			cls[k] = v
// 		}
// 		log.Printf("[WARN] 11111***** headers cls *****11111*****11111*****11111*****11111*****11111*****")
//  		log.Println(cls)
// 		return cls
// 	}	
// 	return headers
// }


// func buildRequestsData(requests interface{}) sc2.Requests {
// 	requestsList := make(sc2.Requests, 0) //len(requests))
// 	log.Printf("[WARN] 11111***** requests input *****11111*****11111*****11111*****11111*****11111*****")
// 	log.Println(requests)

// 	test_schema := requests.(*schema.Set)

// 	test_list := test_schema.List()
// 	test_map := test_list[0].(map[string]interface{})
// 	id := test_map["configuration"]
// 	log.Printf("[WARN] 11111***** requests test_list *****11111*****11111*****11111*****11111*****11111*****")
// 	log.Println(test_list)
// 	log.Println(id)

// 	// for i, request := range requests {
// 	// 	request := request.(map[string]interface{})
// 	// 	req := sc2.Requests{
			
// 	// 	}
// 	// 	// requestsList[i] = req

// 	// }
// 	return requestsList
// }



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
