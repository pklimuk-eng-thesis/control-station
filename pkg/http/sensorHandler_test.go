package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	"github.com/pklimuk-eng-thesis/control-station/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestGetInfo_Success(t *testing.T) {
	sensorService := new(service.MockSensorService)
	expectedSensorInfo := domain.SensorInfo{Enabled: true, Detected: false}
	sensorService.EXPECT().GetInfo().Return(expectedSensorInfo, nil)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sensorHandler.GetInfo(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"enabled": true, "detected": false}`, w.Body.String())
}

func TestGetInfo_ParsingFailure(t *testing.T) {
	sensorService := new(service.MockSensorService)
	sensorService.EXPECT().GetInfo().Return(domain.SensorInfo{Enabled: false, Detected: false}, service.ErrParsingFailed)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sensorHandler.GetInfo(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestToggleEnabled_Success(t *testing.T) {
	sensorService := new(service.MockSensorService)
	sensorService.EXPECT().ToggleEnabled().Return(domain.SensorInfo{Enabled: true, Detected: false}, nil)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sensorHandler.ToggleEnabled(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"enabled": true, "detected": false}`, w.Body.String())
}

func TestToggleEnabled_ParsingFailure(t *testing.T) {
	sensorService := new(service.MockSensorService)
	sensorService.EXPECT().ToggleEnabled().Return(domain.SensorInfo{Enabled: false, Detected: false}, service.ErrParsingFailed)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sensorHandler.ToggleEnabled(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestToggleDetected_Success(t *testing.T) {
	sensorService := new(service.MockSensorService)
	sensorService.EXPECT().ToggleDetected().Return(domain.SensorInfo{Enabled: true, Detected: false}, nil)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sensorHandler.ToggleDetected(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"enabled": true, "detected": false}`, w.Body.String())
}

func TestToggleDetected_ParsingFailure(t *testing.T) {
	sensorService := new(service.MockSensorService)
	sensorService.EXPECT().ToggleDetected().Return(domain.SensorInfo{Enabled: false, Detected: false}, service.ErrParsingFailed)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sensorHandler.ToggleDetected(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}
