package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestGetSensorLogsLimitN(t *testing.T) {
	sensorService := new(service.MockSensorService)
	expectedSensorLogs := []domain.SensorData{
		{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true, Detected: false},
		{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false, Detected: false},
	}
	sensorService.EXPECT().GetSensorLogsFromDataServiceLimitN(2).Return(expectedSensorLogs, nil)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?limit=2", nil)
	sensorHandler.GetSensorLogsLimitN(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[
		{
			"id": 1,
			"created_at": "2023-01-01T00:00:00Z",
			"is_enabled": true,
			"detected": false
		},
		{
			"id": 2,
			"created_at": "2023-01-02T00:00:00Z",
			"is_enabled": false,
			"detected": false
		}]`, w.Body.String())
}

func TestGetSensorLogsLimitN_ParsingFailure(t *testing.T) {
	sensorService := new(service.MockSensorService)
	sensorService.EXPECT().GetSensorLogsFromDataServiceLimitN(2).Return([]domain.SensorData{}, service.ErrParsingFailed)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?limit=2", nil)
	sensorHandler.GetSensorLogsLimitN(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestGetSensorLogsLimitN_InvalidLimit(t *testing.T) {
	sensorService := new(service.MockSensorService)

	sensorHandler := NewSensorHandler(sensorService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?limit=invalid", nil)
	sensorHandler.GetSensorLogsLimitN(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid limit parameter", w.Body.String())
}
