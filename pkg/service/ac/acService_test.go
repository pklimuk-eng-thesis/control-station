package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewACService(t *testing.T) {
	ac := domain.AC{Name: "test", Address: "http://test"}
	service := NewACService(&ac)
	assert.NotNil(t, service)
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.ACInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(domain.ACInfo{Enabled: true, Temperature: 20, Humidity: 50})
			})),
			want:    domain.ACInfo{Enabled: true, Temperature: 20, Humidity: 50},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    domain.ACInfo{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			service := &acService{ac: &domain.AC{Name: "test", Address: test.ts.URL}}
			got, err := service.GetInfo()

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestToggleEnabled(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.ACInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(domain.ACInfo{Enabled: true, Temperature: 20, Humidity: 50})
			})),
			want:    domain.ACInfo{Enabled: true, Temperature: 20, Humidity: 50},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    domain.ACInfo{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			service := &acService{ac: &domain.AC{Name: "test", Address: test.ts.URL}}
			got, err := service.ToggleEnabled()

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestGetDeviceLogsFromDataServiceLimitN(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    []domain.ACData
		limit   int
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]domain.ACData{
					{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true, Temperature: 20, Humidity: 50},
					{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false, Temperature: 25, Humidity: 40},
				})
			})),
			want: []domain.ACData{
				{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true, Temperature: 20, Humidity: 50},
				{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false, Temperature: 25, Humidity: 40},
			},
			limit:   2,
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    []domain.ACData(nil),
			limit:   2,
			wantErr: true,
		},
		{
			name: "Unmarshal failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			},
			)),
			want:    []domain.ACData(nil),
			limit:   2,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			t.Setenv("DATA_SERVICE_ADDRESS", test.ts.URL)
			service := &acService{ac: &domain.AC{Name: "test", Address: test.ts.URL}}
			got, err := service.GetACLogsFromDataServiceLimitN(test.limit)

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestUpdateACSettings(t *testing.T) {
	tests := []struct {
		name        string
		desiredTemp float32
		desiredHum  float32
		ts          *httptest.Server
		want        domain.ACInfo
		wantErr     bool
	}{
		{
			name:        "Success",
			desiredTemp: float32(20),
			desiredHum:  float32(50),
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(domain.ACInfo{Enabled: true, Temperature: 20, Humidity: 50})
			})),
			want:    domain.ACInfo{Enabled: true, Temperature: 20, Humidity: 50},
			wantErr: false,
		},
		{
			name:        "Failure",
			desiredTemp: float32(20),
			desiredHum:  float32(50),
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    domain.ACInfo{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			service := &acService{ac: &domain.AC{Name: "test", Address: test.ts.URL}}
			got, err := service.UpdateACSettings(test.desiredTemp, test.desiredHum)

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
