package service

import (
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	controlStationUtils "github.com/pklimuk-eng-thesis/control-station/pkg/service/utils"
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
	address := s.sensor.Address + controlStationUtils.InfoEndpoint
	return controlStationUtils.MakeGetRequest(address, s.sensor.Name, domain.SensorInfo{Enabled: false, Detected: false})
}

func (s *sensorService) ToggleEnabled() (domain.SensorInfo, error) {
	address := s.sensor.Address + controlStationUtils.EnabledEndpoint
	return controlStationUtils.MakePatchRequest(address, s.sensor.Name, domain.SensorInfo{Enabled: false, Detected: false})
}

func (s *sensorService) ToggleDetected() (domain.SensorInfo, error) {
	address := s.sensor.Address + controlStationUtils.DetectedEndpoint
	return controlStationUtils.MakePatchRequest(address, s.sensor.Name, domain.SensorInfo{Enabled: false, Detected: false})
}

func (s *sensorService) GetSensorLogsFromDataServiceLimitN(limit int) ([]domain.SensorData, error) {
	return controlStationUtils.GetLogsFromDataServiceLimitN[domain.SensorData](s.sensor.Name, limit)
}
