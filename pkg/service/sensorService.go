package service

import (
	"errors"
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
	// Detected() (bool, error)
	// ToggleEnabled() (bool, error)
	// ToggleDetected() (bool, error)
}

type sensorService struct {
	sensor *domain.Sensor
}

func NewSensorService(sensor *domain.Sensor) SensorService {
	return &sensorService{sensor: sensor}
}

func (s *sensorService) IsEnabled() (bool, error) {
	resp, err := http.Get(s.sensor.Address + sensorEnabledEndpoint)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, ErrParsingFailed
	}

	enabled, err := strconv.ParseBool(string(body))
	if err != nil {
		return false, ErrParsingFailed
	}

	return enabled, nil
}
