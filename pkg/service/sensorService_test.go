package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	domain "github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestMakeGetRequest(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.SensorInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"enabled": true, "detected": false}`)
			},
			)),
			want:    domain.SensorInfo{Enabled: true, Detected: false},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Sensor is disabled")
			},
			)),
			want:    domain.SensorInfo{Enabled: false, Detected: false},
			wantErr: true,
		},
		{
			name: "FailureParsing",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "Invalid JSON")
			},
			)),
			want:    domain.SensorInfo{Enabled: false, Detected: false},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sensorInfo, err := makeGetRequest(test.ts.URL, "test-sensor")
			assert.Equal(t, test.want.Enabled, sensorInfo.Enabled)
			assert.Equal(t, test.want.Detected, sensorInfo.Detected)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestMakeGetRequest_FailureConnection(t *testing.T) {
	sensorInfo, err := makeGetRequest("http://localhost:1234", "test-sensor")
	assert.Error(t, err)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakePostRequest(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.SensorInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"enabled": true, "detected": false}`)
			},
			)),
			want:    domain.SensorInfo{Enabled: true, Detected: false},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Sensor is disabled")
			},
			)),
			want:    domain.SensorInfo{Enabled: false, Detected: false},
			wantErr: true,
		},
		{
			name: "FailureParsing",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "Invalid JSON")
			},
			)),
			want:    domain.SensorInfo{Enabled: false, Detected: false},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sensorInfo, err := makePostRequest(test.ts.URL, "test-sensor")
			assert.Equal(t, test.want.Enabled, sensorInfo.Enabled)
			assert.Equal(t, test.want.Detected, sensorInfo.Detected)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestMakePostRequest_FailureConnection(t *testing.T) {
	sensorInfo, err := makePostRequest("http://localhost:1234", "test-sensor")
	assert.Error(t, err)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestGetInfo_Success(t *testing.T) {
	expected := domain.SensorInfo{Enabled: true, Detected: false}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: server.URL}}
	sensorInfo, err := service.GetInfo()

	assert.NoError(t, err)
	assert.Equal(t, expected, sensorInfo)
}

func TestGetInfo_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: server.URL}}
	sensorInfo, err := service.GetInfo()

	assert.Error(t, err)
	assert.Equal(t, domain.SensorInfo{}, sensorInfo)
}

func TestToggleEnabled_Success(t *testing.T) {
	expected := domain.SensorInfo{Enabled: true, Detected: false}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: server.URL}}
	sensorInfo, err := service.ToggleEnabled()

	assert.NoError(t, err)
	assert.Equal(t, expected, sensorInfo)
}

func TestToggleEnabled_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: server.URL}}
	sensorInfo, err := service.ToggleEnabled()

	assert.Error(t, err)
	assert.Equal(t, domain.SensorInfo{}, sensorInfo)
}

func TestToggleDetected_Success(t *testing.T) {
	expected := domain.SensorInfo{Enabled: true, Detected: false}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: server.URL}}
	sensorInfo, err := service.ToggleDetected()

	assert.NoError(t, err)
	assert.Equal(t, expected, sensorInfo)
}

func TestToggleDetected_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := &sensorService{sensor: &domain.Sensor{Name: "test", Address: server.URL}}
	sensorInfo, err := service.ToggleDetected()

	assert.Error(t, err)
	assert.Equal(t, domain.SensorInfo{}, sensorInfo)
}

func TestSendSensorLogsToDataService(t *testing.T) {
	tests := []struct {
		name       string
		ts         *httptest.Server
		sensorInfo domain.SensorInfo
		wantErr    bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			)),
			sensorInfo: domain.SensorInfo{Enabled: true, Detected: false},
			wantErr:    false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			sensorInfo: domain.SensorInfo{Enabled: true, Detected: false},
			wantErr:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := sendSensorLogsToDataService(test.ts.URL, "test-sensor", test.sensorInfo)
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
