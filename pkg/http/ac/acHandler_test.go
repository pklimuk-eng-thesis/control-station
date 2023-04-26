package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	service "github.com/pklimuk-eng-thesis/control-station/pkg/service/ac"
	"github.com/pklimuk-eng-thesis/control-station/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetInfo_Success(t *testing.T) {
	acService := new(service.MockACService)
	expectedACInfo := domain.ACInfo{Enabled: true, Temperature: 20, Humidity: 50}
	acService.EXPECT().GetInfo().Return(expectedACInfo, nil)

	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	acHandler.GetInfo(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"enabled": true, "temperature": 20, "humidity": 50}`, w.Body.String())
}

func TestGetInfo_ParsingFailure(t *testing.T) {
	acService := new(service.MockACService)
	acService.EXPECT().GetInfo().Return(domain.ACInfo{Enabled: false, Temperature: 0, Humidity: 0}, utils.ErrParsingFailed)

	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	acHandler.GetInfo(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestToggleEnabled_Success(t *testing.T) {
	acService := new(service.MockACService)
	acService.EXPECT().ToggleEnabled().Return(domain.ACInfo{Enabled: true, Temperature: 20, Humidity: 50}, nil)

	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	acHandler.ToggleEnabled(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"enabled": true, "temperature": 20, "humidity": 50}`, w.Body.String())
}

func TestToggleEnabled_ParsingFailure(t *testing.T) {
	acService := new(service.MockACService)
	acService.EXPECT().ToggleEnabled().Return(domain.ACInfo{Enabled: false, Temperature: 0, Humidity: 0}, utils.ErrParsingFailed)

	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	acHandler.ToggleEnabled(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestGetACLogsLimitN_Success(t *testing.T) {
	acService := new(service.MockACService)
	expectedLogs := []domain.ACData{
		{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true, Temperature: 20, Humidity: 50},
		{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false, Temperature: 25, Humidity: 45},
	}
	acService.EXPECT().GetACLogsFromDataServiceLimitN(2).Return(expectedLogs, nil)

	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/ac/logs?limit=2", nil)
	acHandler.GetACLogsLimitN(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[
		{
			"id": 1,
			"created_at": "2023-01-01T00:00:00Z",
			"is_enabled": true,
			"temperature": 20,
			"humidity": 50
		},
		{
			"id": 2,
			"created_at": "2023-01-02T00:00:00Z",
			"is_enabled": false,
			"temperature": 25,
			"humidity": 45
		}
	]`, w.Body.String())
}

func TestGetACLogsLimitN_ParsingFailure(t *testing.T) {
	acService := new(service.MockACService)
	acService.EXPECT().GetACLogsFromDataServiceLimitN(2).Return([]domain.ACData{}, utils.ErrParsingFailed)

	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?limit=2", nil)
	acHandler.GetACLogsLimitN(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestGetACLogsLimitN_InvalidLimit(t *testing.T) {
	acService := new(service.MockACService)

	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?limit=invalid", nil)
	acHandler.GetACLogsLimitN(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid limit parameter", w.Body.String())
}

func TestUpdateACSettings_Success(t *testing.T) {
	desiredTemp := float32(20)
	desiredHum := float32(50)
	acService := new(service.MockACService)
	acService.EXPECT().UpdateACSettings(desiredTemp, desiredHum).Return(domain.ACInfo{Enabled: true, Temperature: desiredTemp, Humidity: desiredHum}, nil)

	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"enabled":true,"temperature":20,"humidity":50}`))

	acHandler.UpdateACSettings(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"enabled": true,
		"temperature": 20,
		"humidity": 50
	}`, w.Body.String())
}

func TestUpdateACSettings_ParsingFailure(t *testing.T) {
	acService := new(service.MockACService)
	acHandler := NewACHandler(acService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"enabled":true,"temperature":20,"humidity":50`))

	acHandler.UpdateACSettings(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
