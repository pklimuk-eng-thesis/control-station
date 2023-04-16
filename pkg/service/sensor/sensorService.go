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

var sensorEnabledEndpoint = "/enabled"
var sensorDetectedEndpoint = "/detected"
var sensorInfoEndpoint = "/info"
var dataServiceAddress = utils.GetEnvVariableOrDefault("DATA_SERVICE_ADDRESS", "http://localhost:8087")

//go:generate --name SensorService --output mock_sensorService.go
type SensorService interface {
	GetInfo() (domain.SensorInfo, error)
	ToggleEnabled() (domain.SensorInfo, error)
	ToggleDetected() (domain.SensorInfo, error)
	GetSensorLogsFromDataServiceLimitN(limit int) ([]domain.SensorData, error)
}

type sensorService struct {
	sensor *domain.Sensor
}

func NewSensorService(sensor *domain.Sensor) SensorService {
	return &sensorService{sensor: sensor}
}

func (s *sensorService) GetInfo() (domain.SensorInfo, error) {
	return makeGetRequest(s.sensor.Address+sensorInfoEndpoint, s.sensor.Name)
}

func (s *sensorService) ToggleEnabled() (domain.SensorInfo, error) {
	return makePatchRequest(s.sensor.Address+sensorEnabledEndpoint, s.sensor.Name)
}

func (s *sensorService) ToggleDetected() (domain.SensorInfo, error) {
	return makePatchRequest(s.sensor.Address+sensorDetectedEndpoint, s.sensor.Name)
}

func (s *sensorService) GetSensorLogsFromDataServiceLimitN(limit int) ([]domain.SensorData, error) {
	dataServiceAddress = utils.GetEnvVariableOrDefault("DATA_SERVICE_ADDRESS", "http://localhost:8087")
	url := fmt.Sprintf("%s/%s/latest?limit=%d", dataServiceAddress, s.sensor.Name, limit)
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
		return nil, fmt.Errorf("%s: %s", s.sensor.Name, string(body))
	}

	var sensorLogs []domain.SensorData
	err = json.Unmarshal(body, &sensorLogs)
	if err != nil {
		return nil, utils.ErrParsingFailed
	}

	return sensorLogs, nil
}

func makeGetRequest(address string, sensorName string) (domain.SensorInfo, error) {
	resp, err := http.Get(address)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, utils.ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return domain.SensorInfo{Enabled: false, Detected: false}, fmt.Errorf("%s: %s", sensorName, string(body))
	}

	var sensorInfo domain.SensorInfo
	err = json.Unmarshal(body, &sensorInfo)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, utils.ErrParsingFailed
	}

	err = sendSensorLogsToDataService(dataServiceAddress, sensorName, sensorInfo)
	if err != nil {
		log.Println("Failed to send sensor logs to data service: ", err)
	}

	return sensorInfo, nil
}

func makePatchRequest(address string, sensorName string) (domain.SensorInfo, error) {
	req, err := http.NewRequest(http.MethodPatch, address, nil)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, utils.ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return domain.SensorInfo{Enabled: false, Detected: false}, fmt.Errorf("%s: %s", sensorName, string(body))
	}

	var sensorInfo domain.SensorInfo
	err = json.Unmarshal(body, &sensorInfo)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, utils.ErrParsingFailed
	}

	err = sendSensorLogsToDataService(dataServiceAddress, sensorName, sensorInfo)
	if err != nil {
		log.Println("Failed to send sensor logs to data service: ", err)
	}

	return sensorInfo, nil
}

func sendSensorLogsToDataService(address string, sensorName string, sensorInfo domain.SensorInfo) error {
	jsonValue, err := json.Marshal(sensorInfo)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/add", address, sensorName)
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
		return fmt.Errorf("%s: %s", sensorName, string(body))
	}

	return nil
}
