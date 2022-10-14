package tests

import (
	"encoding/json"
	"testing"

	"github.com/Fonzeca/Trackin/entry/manager"

	modeljson "github.com/Fonzeca/Trackin/entry/json"
)

func TestZonaProccesData(t *testing.T) {
	t.Log("Empezando test")
	manager := manager.NewGeofenceDetector()

	dummy1 := modeljson.SimplyData{}

	string1 := "{\"protocolHeadType\": 19, \"orignBytes\": [37, 37, 19, 0, 89, 85, 2, 8, 103, 115, 0, 84, 17, 33, 53, 0, 4, 0, 30, 10, 0, 200, 0, 76, 212, 192, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 0, 0, 0, 0, 121, 105, 0, 34, 16, 19, 32, 23, 66, 0, 0, 128, 65, 113, 134, 135, 194, 36, 56, 69, 194, 0, 0, 0, 10, 3, 144, 18, 32, 255, 255, 0, 0, 0, 41, 255, 255, 255, 255, 255], \"serialNo\": 21762, \"imei\": \"867730054112135\", \"samplingIntervalAccOn\": 4, \"samplingIntervalAccOff\": 30, \"angleCompensation\": 10, \"distanceCompensation\": 200, \"overspeedLimit\": 76, \"gpsWorking\": true, \"isHistoryData\": true, \"satelliteNumber\": 20, \"gSensorSensitivity\": 12, \"isManagerConfigured1\": false, \"isManagerConfigured2\": false, \"isManagerConfigured3\": false, \"isManagerConfigured4\": false, \"antitheftedStatus\": 0, \"heartbeatInterval\": 5, \"relayStatus\": 0, \"isRelayWaiting\": false, \"dragThreshold\": 0, \"IOP\": 0, \"iopIgnition\": false, \"iopPowerCutOff\": false, \"iopACOn\": false, \"analogInput1\": 0.0, \"analogInput2\": 0.0, \"originalAlarmCode\": 0, \"mileage\": 31081, \"externalPowerSupply\": false, \"input1\": false, \"input2\": false, \"input3\": false, \"input4\": false, \"input5\": false, \"analogInput4\": 0, \"analogInput5\": 0, \"output1\": false, \"isSmartUploadSupport\": false, \"supportChangeBattery\": false, \"batteryVoltage\": 3.9, \"deviceTemp\": 69, \"is_4g_lbs\": false, \"is_2g_lbs\": false, \"mcc_4g\": 0, \"mnc_4g\": 0, \"ci_4g\": 0, \"earfcn_4g_1\": 0, \"pcid_4g_1\": 0, \"earfcn_4g_2\": 0, \"pcid_4g_2\": 0, \"mcc_2g\": 0, \"mnc_2g\": 0, \"lac_2g_1\": 0, \"ci_2g_1\": 0, \"lac_2g_2\": 0, \"ci_2g_2\": 0, \"lac_2g_3\": 0, \"ci_2g_3\": 0, \"batteryCharge\": 100, \"date\": \"2022-10-13T20:17:42+00:00\", \"latlngValid\": true, \"altitude\": 16.0, \"latitude\": -49.30482482910156, \"longitude\": -67.76258087158203, \"speed\": 0.0, \"azimuth\": 10, \"externalPowerVoltage\": 12.2, \"networkSignal\": 76, \"output2\": false, \"output3\": false, \"output12V\": false, \"outputVout\": false, \"rpm\": 65535, \"analogInput3\": 0.0, \"rlyMode\": 0, \"smsLanguageType\": 0, \"accdetSettingStatus\": 0, \"isSendSmsAlarmWhenDigitalInput2Change\": false, \"isSendSmsAlarmToManagerPhone\": false, \"jammerDetectionStatus\": 0}"
	t.Log("Unmarshalling")
	json.Unmarshal([]byte(string1), &dummy1)

	manager.ProcessData(dummy1)
	manager.ProcessData(dummy1)
	manager.ProcessData(dummy1)
	dummy1.Latitude = -49.30398848548048
	dummy1.Longitude = -67.78385455435571
	manager.ProcessData(dummy1)
	manager.ProcessData(dummy1)
	manager.ProcessData(dummy1)
}
