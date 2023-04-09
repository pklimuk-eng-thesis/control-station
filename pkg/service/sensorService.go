package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
)

var ErrParsingFailed = errors.New("Parsing failed")

var sensorEnabledEndpoint = "/enabled"
var sensorDetectedEndpoint = "/detected"
var sensorInfoEndpoint = "/info"

//go:generate --name SensorService --output mock_sensorService.go
type SensorService interface {
	GetInfo() (domain.SensorInfo, error)
	ToggleEnabled() (domain.SensorInfo, error)
	ToggleDetected() (domain.SensorInfo, error)
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
	return makePostRequest(s.sensor.Address+sensorEnabledEndpoint, s.sensor.Name)
}

func (s *sensorService) ToggleDetected() (domain.SensorInfo, error) {
	return makePostRequest(s.sensor.Address+sensorDetectedEndpoint, s.sensor.Name)
}

func makeGetRequest(address string, sensorName string) (domain.SensorInfo, error) {
	resp, err := http.Get(address)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return domain.SensorInfo{Enabled: false, Detected: false}, fmt.Errorf("%s: %s", sensorName, string(body))
	}

	var sensorInfo domain.SensorInfo
	err = json.Unmarshal(body, &sensorInfo)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, ErrParsingFailed
	}

	return sensorInfo, nil
}

func makePostRequest(address string, sensorName string) (domain.SensorInfo, error) {
	resp, err := http.Post(address, "application/json", nil)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return domain.SensorInfo{Enabled: false, Detected: false}, fmt.Errorf("%s: %s", sensorName, string(body))
	}

	var sensorInfo domain.SensorInfo
	err = json.Unmarshal(body, &sensorInfo)
	if err != nil {
		return domain.SensorInfo{Enabled: false, Detected: false}, ErrParsingFailed
	}

	return sensorInfo, nil
}
