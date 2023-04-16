package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	"github.com/pklimuk-eng-thesis/control-station/utils"
)

const DeviceInfoEndpoint = "/info"
const DeviceEnabledEndpoint = "/enabled"

type SmartHomeDeviceInfo interface {
	domain.SensorInfo | domain.DeviceInfo
}

type SmartHomeDeviceData interface {
	domain.SensorData | domain.DeviceData
}

func MakeGetRequest[V SmartHomeDeviceInfo](address string, deviceName string, defaultValueOnError V) (V, error) {
	resp, err := http.Get(address)
	if err != nil {
		return defaultValueOnError, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return defaultValueOnError, utils.ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return defaultValueOnError, fmt.Errorf("%s: %s", deviceName, string(body))
	}

	var deviceInfo V
	err = json.Unmarshal(body, &deviceInfo)
	if err != nil {
		return defaultValueOnError, utils.ErrParsingFailed
	}

	err = sendLogsToDataService(deviceName, deviceInfo)
	if err != nil {
		errStr := fmt.Sprintf("Failed to send '%s' logs to data service: %s", deviceName, err)
		log.Println(errStr)
	}

	return deviceInfo, nil
}

func MakePatchRequest[V SmartHomeDeviceInfo](address string, deviceName string, defaultValueOnError V) (V, error) {
	req, err := http.NewRequest(http.MethodPatch, address, nil)
	if err != nil {
		return defaultValueOnError, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return defaultValueOnError, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return defaultValueOnError, utils.ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return defaultValueOnError, fmt.Errorf("%s: %s", deviceName, string(body))
	}

	var deviceInfo V
	err = json.Unmarshal(body, &deviceInfo)
	if err != nil {
		return defaultValueOnError, utils.ErrParsingFailed
	}

	err = sendLogsToDataService(deviceName, deviceInfo)
	if err != nil {
		errStr := fmt.Sprintf("Failed to send '%s' logs to data service: %s", deviceName, err)
		log.Println(errStr)
	}

	return deviceInfo, nil
}

func GetLogsFromDataServiceLimitN[K SmartHomeDeviceData](deviceName string, limit int) ([]K, error) {
	dataServiceAddress := utils.GetEnvVariableOrDefault("DATA_SERVICE_ADDRESS", "http://localhost:8087")
	url := fmt.Sprintf("%s/%s/latest?limit=%d", dataServiceAddress, deviceName, limit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", deviceName, string(body))
	}

	var deviceData []K
	err = json.Unmarshal(body, &deviceData)
	if err != nil {
		return nil, utils.ErrParsingFailed
	}

	return deviceData, nil
}

func sendLogsToDataService[V SmartHomeDeviceInfo](deviceName string, deviceInfo V) error {
	dataServiceAddress := utils.GetEnvVariableOrDefault("DATA_SERVICE_ADDRESS", "http://localhost:8087")
	jsonValue, err := json.Marshal(deviceInfo)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/add", dataServiceAddress, deviceName)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return utils.ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", deviceName, string(body))
	}

	return nil
}
