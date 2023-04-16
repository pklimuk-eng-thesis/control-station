package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	domain "github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.SensorInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(domain.SensorInfo{Enabled: true, Detected: false})
			},
			)),
			want:    domain.SensorInfo{Enabled: true, Detected: false},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    domain.SensorInfo{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: test.ts.URL}}
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
		want    domain.SensorInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(domain.SensorInfo{Enabled: true, Detected: false})
			},
			)),
			want:    domain.SensorInfo{Enabled: true, Detected: false},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    domain.SensorInfo{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: test.ts.URL}}
			got, err := service.ToggleEnabled()

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestToggleDetected(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.SensorInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(domain.SensorInfo{Enabled: true, Detected: false})
			},
			)),
			want:    domain.SensorInfo{Enabled: true, Detected: false},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    domain.SensorInfo{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: test.ts.URL}}
			got, err := service.ToggleDetected()

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestGetSensorLogsFromDataServiceLimitN(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    []domain.SensorData
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]domain.SensorData{
					{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true, Detected: false},
					{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false, Detected: false},
				})
			},
			)),
			want: []domain.SensorData{
				{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true, Detected: false},
				{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false, Detected: false},
			},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    []domain.SensorData(nil),
			wantErr: true,
		},
		{
			name: "Unmarshal Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("invalid json"))
			},
			)),
			want:    []domain.SensorData(nil),
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: "address"}}
			os.Setenv("DATA_SERVICE_ADDRESS", test.ts.URL)
			got, err := service.GetSensorLogsFromDataServiceLimitN(2)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}
