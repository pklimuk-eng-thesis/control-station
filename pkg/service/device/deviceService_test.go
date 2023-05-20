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

func TestNewDeviceService(t *testing.T) {
	device := domain.Device{Name: "test", Address: "http://test"}
	service := NewDeviceService(&device)
	assert.NotNil(t, service)
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.DeviceInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(domain.DeviceInfo{Enabled: true})
			})),
			want:    domain.DeviceInfo{Enabled: true},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    domain.DeviceInfo{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			service := &deviceService{device: &domain.Device{Name: "test", Address: test.ts.URL}}
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
		want    domain.DeviceInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(domain.DeviceInfo{Enabled: true})
			})),
			want:    domain.DeviceInfo{Enabled: true},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    domain.DeviceInfo{},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			service := &deviceService{device: &domain.Device{Name: "test", Address: test.ts.URL}}
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
		want    []domain.DeviceData
		limit   int
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]domain.DeviceData{
					{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true},
					{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false},
				})
			})),
			want: []domain.DeviceData{
				{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true},
				{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false},
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
			want:    []domain.DeviceData(nil),
			limit:   2,
			wantErr: true,
		},
		{
			name: "Unmarshal failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("invalid json")
			},
			)),
			want:    []domain.DeviceData(nil),
			limit:   10,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.ts.Close()
			t.Setenv("DATA_SERVICE_ADDRESS", test.ts.URL)
			service := &deviceService{device: &domain.Device{Name: "test", Address: test.ts.URL}}
			got, err := service.GetDeviceLogsFromDataServiceLimitN(test.limit)

			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
