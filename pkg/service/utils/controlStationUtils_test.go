package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestMakeGetRequestSensor(t *testing.T) {
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
			sensorInfo, err := MakeGetRequest(test.ts.URL, "test-sensor", domain.SensorInfo{Enabled: false, Detected: false})
			assert.Equal(t, test.want.Enabled, sensorInfo.Enabled)
			assert.Equal(t, test.want.Detected, sensorInfo.Detected)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestMakeGetRequestSensor_FailureConnection(t *testing.T) {
	sensorInfo, err := MakeGetRequest("http://localhost:1234",
		"test-sensor",
		domain.SensorInfo{Enabled: false, Detected: false})
	assert.Error(t, err)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakeGetRequestDevice(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.DeviceInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"enabled": true}`)
			},
			)),
			want:    domain.DeviceInfo{Enabled: true},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Device is disabled")
			},
			)),
			want:    domain.DeviceInfo{Enabled: false},
			wantErr: true,
		},
		{
			name: "FailureParsing",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "Invalid JSON")
			},
			)),
			want:    domain.DeviceInfo{Enabled: false},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deviceInfo, err := MakeGetRequest(test.ts.URL, "test-device", domain.DeviceInfo{Enabled: false})
			assert.Equal(t, test.want.Enabled, deviceInfo.Enabled)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestMakeGetRequestDevice_FailureConnection(t *testing.T) {
	deviceInfo, err := MakeGetRequest("http://localhost:1234",
		"test-device",
		domain.DeviceInfo{Enabled: false})
	assert.Error(t, err)
	assert.Equal(t, false, deviceInfo.Enabled)
}

func TestMakePatchRequestSensor(t *testing.T) {
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
			sensorInfo, err := MakePatchRequest(test.ts.URL, "test-sensor", domain.SensorInfo{Enabled: false, Detected: false})
			assert.Equal(t, test.want.Enabled, sensorInfo.Enabled)
			assert.Equal(t, test.want.Detected, sensorInfo.Detected)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestMakePatchRequestSensor_FailureConnection(t *testing.T) {
	sensorInfo, err := MakePatchRequest("http://localhost:1234",
		"test-sensor",
		domain.SensorInfo{Enabled: false, Detected: false})
	assert.Error(t, err)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakePatchRequestDevice(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    domain.DeviceInfo
		wantErr bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, `{"enabled": true}`)
			},
			)),
			want:    domain.DeviceInfo{Enabled: true},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Device is disabled")
			},
			)),
			want:    domain.DeviceInfo{Enabled: false},
			wantErr: true,
		},
		{
			name: "FailureParsing",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "Invalid JSON")
			},
			)),
			want:    domain.DeviceInfo{Enabled: false},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deviceInfo, err := MakePatchRequest(test.ts.URL, "test-device", domain.DeviceInfo{Enabled: false})
			assert.Equal(t, test.want.Enabled, deviceInfo.Enabled)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestMakePatchRequestDevice_FailureConnection(t *testing.T) {
	deviceInfo, err := MakePatchRequest("http://localhost:1234",
		"test-device",
		domain.DeviceInfo{Enabled: false})
	assert.Error(t, err)
	assert.Equal(t, false, deviceInfo.Enabled)
}

func TestSendLogsToDataService_Sensor(t *testing.T) {
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
			os.Setenv("DATA_SERVICE_ADDRESS", test.ts.URL)
			err := sendLogsToDataService("test-sensor", test.sensorInfo)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestSendLogsToDataService_Device(t *testing.T) {
	tests := []struct {
		name       string
		ts         *httptest.Server
		deviceInfo domain.DeviceInfo
		wantErr    bool
	}{
		{
			name: "Success",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			)),
			deviceInfo: domain.DeviceInfo{Enabled: true},
			wantErr:    false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			deviceInfo: domain.DeviceInfo{Enabled: true},
			wantErr:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("DATA_SERVICE_ADDRESS", test.ts.URL)
			err := sendLogsToDataService("test-device", test.deviceInfo)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestGetLogsFromDataServiceLimitN_Sensor(t *testing.T) {
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
			os.Setenv("DATA_SERVICE_ADDRESS", test.ts.URL)
			got, err := GetLogsFromDataServiceLimitN[domain.SensorData]("test-sensor", 2)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetLogsFromDataServiceLimitN_Device(t *testing.T) {
	tests := []struct {
		name    string
		ts      *httptest.Server
		want    []domain.DeviceData
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
			},
			)),
			want: []domain.DeviceData{
				{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true},
				{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false},
			},
			wantErr: false,
		},
		{
			name: "Failure",
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			)),
			want:    []domain.DeviceData(nil),
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
			want:    []domain.DeviceData(nil),
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv("DATA_SERVICE_ADDRESS", test.ts.URL)
			got, err := GetLogsFromDataServiceLimitN[domain.DeviceData]("test-device", 2)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}
