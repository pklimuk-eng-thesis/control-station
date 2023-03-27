package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
)

var ErrParsingFailed = errors.New("Parsing failed")

var sensorEnabledEndpoint = "/enabled"
var sensorDetectedEndpoint = "/detected"

type SensorService interface {
	IsEnabled() (bool, error)
	Detected() (bool, error)
	ToggleEnabled() (bool, error)
	ToggleDetected() (bool, error)
}

type sensorService struct {
	sensor *domain.Sensor
}

func NewSensorService(sensor *domain.Sensor) SensorService {
	return &sensorService{sensor: sensor}
}

func (s *sensorService) IsEnabled() (bool, error) {
	return makeGetRequest(s.sensor.Address+sensorEnabledEndpoint, s.sensor.Name)
}

func (s *sensorService) Detected() (bool, error) {
	return makeGetRequest(s.sensor.Address+sensorDetectedEndpoint, s.sensor.Name)
}

func (s *sensorService) ToggleEnabled() (bool, error) {
	return makePostRequest(s.sensor.Address+sensorEnabledEndpoint, s.sensor.Name)
}

func (s *sensorService) ToggleDetected() (bool, error) {
	return makePostRequest(s.sensor.Address+sensorDetectedEndpoint, s.sensor.Name)
}

func makeGetRequest(address string, sensorName string) (bool, error) {
	resp, err := http.Get(address)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%s: %s", sensorName, string(body))
	}

	result, err := strconv.ParseBool(string(body))
	if err != nil {
		return false, ErrParsingFailed
	}

	return result, nil
}

func makePostRequest(address string, sensorName string) (bool, error) {
	resp, err := http.Post(address, "application/json", nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, ErrParsingFailed
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("%s: %s", sensorName, string(body))
	}

	result, err := strconv.ParseBool(string(body))
	if err != nil {
		return false, ErrParsingFailed
	}

	return result, nil
}
