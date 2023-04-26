package service

import (
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	controlStationUtils "github.com/pklimuk-eng-thesis/control-station/pkg/service/utils"
)

//go:generate --name ACService --output mock_acService.go
type ACService interface {
	GetInfo() (domain.ACInfo, error)
	ToggleEnabled() (domain.ACInfo, error)
	UpdateACSettings(desiredTemp float32, desiredHum float32) (domain.ACInfo, error)
	GetACLogsFromDataServiceLimitN(limit int) ([]domain.ACData, error)
}

type acService struct {
	ac *domain.AC
}

func NewACService(ac *domain.AC) ACService {
	return &acService{ac: ac}
}

func (s *acService) GetInfo() (domain.ACInfo, error) {
	address := s.ac.Address + controlStationUtils.InfoEndpoint
	return controlStationUtils.MakeGetRequest(address, s.ac.Name,
		domain.ACInfo{Enabled: false, Temperature: 0.0, Humidity: 0.0})
}

func (s *acService) ToggleEnabled() (domain.ACInfo, error) {
	address := s.ac.Address + controlStationUtils.EnabledEndpoint
	return controlStationUtils.MakePatchRequest(address, s.ac.Name, nil,
		domain.ACInfo{Enabled: false, Temperature: 0.0, Humidity: 0.0})
}

func (s *acService) UpdateACSettings(desiredTemp float32, desiredHum float32) (domain.ACInfo, error) {
	address := s.ac.Address + controlStationUtils.UpdateEndpoint
	return controlStationUtils.MakePatchRequest(address, s.ac.Name,
		&domain.ACInfo{Enabled: true, Temperature: desiredTemp, Humidity: desiredHum},
		domain.ACInfo{Enabled: false, Temperature: 0.0, Humidity: 0.0})
}

func (s *acService) GetACLogsFromDataServiceLimitN(limit int) ([]domain.ACData, error) {
	return controlStationUtils.GetLogsFromDataServiceLimitN[domain.ACData](s.ac.Name, limit)
}
