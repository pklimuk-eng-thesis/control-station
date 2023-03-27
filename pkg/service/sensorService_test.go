package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
)

func TestMakeGetRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}))
	defer ts.Close()

	s := &sensorService{
		sensor: &domain.Sensor{
			Address: ts.URL,
			Name:    "testSensor",
		},
	}

	t.Run("detected", func(t *testing.T) {
		detected, err := makeGetRequest(s.sensor.Address, s.sensor.Name)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !detected {
			t.Errorf("Expected detected to be true, got %v", detected)
		}
	})

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("false"))
	}))
	defer ts2.Close()

	s2 := &sensorService{
		sensor: &domain.Sensor{
			Address: ts2.URL,
			Name:    "testSensor2",
		},
	}

	t.Run("not detected", func(t *testing.T) {
		detected, err := makeGetRequest(s2.sensor.Address, s2.sensor.Name)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if detected {
			t.Errorf("Expected detected to be false, got %v", detected)
		}
	})

	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Sensor is disabled"))
	}))
	defer ts3.Close()

	s3 := &sensorService{
		sensor: &domain.Sensor{
			Address: ts3.URL,
			Name:    "testSensor3",
		},
	}

	t.Run("http error", func(t *testing.T) {
		_, err := makeGetRequest(s3.sensor.Address, s3.sensor.Name)
		if err == nil {
			t.Errorf("Expected an error, got none")
		}
		if err.Error() != "testSensor3: Sensor is disabled" {
			t.Errorf("Expected error to be 'testSensor3: Sensor is disabled', got %v", err)
		}
	})

	ts4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid"))
	}))
	defer ts4.Close()

	s4 := &sensorService{
		sensor: &domain.Sensor{
			Address: ts4.URL,
			Name:    "testSensor4",
		},
	}

	t.Run("parse error", func(t *testing.T) {
		_, err := makeGetRequest(s4.sensor.Address, s4.sensor.Name)
		if !errors.Is(err, ErrParsingFailed) {
			t.Errorf("Expected parsing error, got %v", err)
		}
	})
}

func TestMakePostRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	}))
	defer ts.Close()

	s := &sensorService{
		sensor: &domain.Sensor{
			Address: ts.URL,
			Name:    "testSensor",
		},
	}

	t.Run("detected", func(t *testing.T) {
		detected, err := makePostRequest(s.sensor.Address, s.sensor.Name)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !detected {
			t.Errorf("Expected detected to be true, got %v", detected)
		}
	})

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("false"))
	}))
	defer ts2.Close()

	s2 := &sensorService{
		sensor: &domain.Sensor{
			Address: ts2.URL,
			Name:    "testSensor2",
		},
	}

	t.Run("not detected", func(t *testing.T) {
		detected, err := makePostRequest(s2.sensor.Address, s2.sensor.Name)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if detected {
			t.Errorf("Expected detected to be false, got %v", detected)
		}
	})

	ts3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Sensor is disabled"))
	}))
	defer ts3.Close()

	s3 := &sensorService{
		sensor: &domain.Sensor{
			Address: ts3.URL,
			Name:    "testSensor3",
		},
	}

	t.Run("http error", func(t *testing.T) {
		_, err := makePostRequest(s3.sensor.Address, s3.sensor.Name)
		if err == nil {
			t.Errorf("Expected an error, got none")
		}
		if err.Error() != "testSensor3: Sensor is disabled" {
			t.Errorf("Expected error to be 'testSensor3: Sensor is disabled', got %v", err)
		}
	})

	ts4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid"))
	}))
	defer ts4.Close()

	s4 := &sensorService{
		sensor: &domain.Sensor{
			Address: ts4.URL,
			Name:    "testSensor4",
		},
	}

	t.Run("parse error", func(t *testing.T) {
		_, err := makePostRequest(s4.sensor.Address, s4.sensor.Name)
		if !errors.Is(err, ErrParsingFailed) {
			t.Errorf("Expected parsing error, got %v", err)
		}
	})
}
