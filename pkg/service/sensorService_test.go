package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestMakeGetRequest_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"enabled": true, "detected": false}`)
	}))
	defer ts.Close()

	sensorInfo, err := makeGetRequest(ts.URL, "test-sensor")
	assert.NoError(t, err)
	assert.Equal(t, true, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakeGetRequest_Failure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sensor is disabled")
	}))
	defer ts.Close()

	sensorInfo, err := makeGetRequest(ts.URL, "test-sensor")
	assert.Error(t, err)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakeGetRequest_FailureParsing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "invalid json")
	}))
	defer ts.Close()

	sensorInfo, err := makeGetRequest(ts.URL, "test-sensor")
	assert.ErrorIs(t, err, ErrParsingFailed)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakeGetRequest_FailureConnection(t *testing.T) {
	sensorInfo, err := makeGetRequest("http://localhost:1234", "test-sensor")
	assert.Error(t, err)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakePostRequest_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"enabled": true, "detected": false}`)
	}))
	defer ts.Close()

	sensorInfo, err := makePostRequest(ts.URL, "test-sensor")
	assert.NoError(t, err)
	assert.Equal(t, true, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakePostRequest_Failure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sensor is disabled")
	}))
	defer ts.Close()

	sensorInfo, err := makePostRequest(ts.URL, "test-sensor")
	assert.Error(t, err)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
}

func TestMakePostRequest_FailureParsing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "invalid json")
	}))
	defer ts.Close()

	sensorInfo, err := makePostRequest(ts.URL, "test-sensor")
	assert.ErrorIs(t, err, ErrParsingFailed)
	assert.Equal(t, false, sensorInfo.Enabled)
	assert.Equal(t, false, sensorInfo.Detected)
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
