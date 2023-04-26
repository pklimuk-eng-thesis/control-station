package service

import (
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	controlStationUtils "github.com/pklimuk-eng-thesis/control-station/pkg/service/utils"
)

//go:generate --name DeviceService --output mock_deviceService.go
type DeviceService interface {
	GetInfo() (domain.DeviceInfo, error)
	ToggleEnabled() (domain.DeviceInfo, error)
	GetDeviceLogsFromDataServiceLimitN(limit int) ([]domain.DeviceData, error)
}

type deviceService struct {
	device *domain.Device
}

func NewDeviceService(device *domain.Device) DeviceService {
	return &deviceService{device: device}
}

func (s *deviceService) GetInfo() (domain.DeviceInfo, error) {
	address := s.device.Address + controlStationUtils.InfoEndpoint
	return controlStationUtils.MakeGetRequest(address, s.device.Name, domain.DeviceInfo{Enabled: false})
}

func (s *deviceService) ToggleEnabled() (domain.DeviceInfo, error) {
	address := s.device.Address + controlStationUtils.EnabledEndpoint
	return controlStationUtils.MakePatchRequest(address, s.device.Name, nil, domain.DeviceInfo{Enabled: false})
}

func (s *deviceService) GetDeviceLogsFromDataServiceLimitN(limit int) ([]domain.DeviceData, error) {
	return controlStationUtils.GetLogsFromDataServiceLimitN[domain.DeviceData](s.device.Name, limit)
}
