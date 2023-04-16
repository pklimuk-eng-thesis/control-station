package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	service "github.com/pklimuk-eng-thesis/control-station/pkg/service/device"
	"github.com/pklimuk-eng-thesis/control-station/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetInfo_Success(t *testing.T) {
	deviceService := new(service.MockDeviceService)
	expectedDeviceInfo := domain.DeviceInfo{Enabled: true}
	deviceService.EXPECT().GetInfo().Return(expectedDeviceInfo, nil)

	deviceHandler := NewDeviceHandler(deviceService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	deviceHandler.GetInfo(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"enabled": true}`, w.Body.String())
}

func TestGetInfo_ParsingFailure(t *testing.T) {
	deviceService := new(service.MockDeviceService)
	deviceService.EXPECT().GetInfo().Return(domain.DeviceInfo{Enabled: false}, utils.ErrParsingFailed)

	deviceHandler := NewDeviceHandler(deviceService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	deviceHandler.GetInfo(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestToggleEnabled_Success(t *testing.T) {
	deviceService := new(service.MockDeviceService)
	deviceService.EXPECT().ToggleEnabled().Return(domain.DeviceInfo{Enabled: true}, nil)

	deviceHandler := NewDeviceHandler(deviceService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	deviceHandler.ToggleEnabled(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"enabled": true}`, w.Body.String())
}

func TestToggleEnabled_ParsingFailure(t *testing.T) {
	deviceService := new(service.MockDeviceService)
	deviceService.EXPECT().ToggleEnabled().Return(domain.DeviceInfo{Enabled: false}, utils.ErrParsingFailed)

	deviceHandler := NewDeviceHandler(deviceService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	deviceHandler.ToggleEnabled(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestGetDeviceLogsFromDataServiceLimitN_Success(t *testing.T) {
	deviceService := new(service.MockDeviceService)
	expectedLogs := []domain.DeviceData{
		{ID: 1, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), IsEnabled: true},
		{ID: 2, CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), IsEnabled: false},
	}
	deviceService.EXPECT().GetDeviceLogsFromDataServiceLimitN(2).Return(expectedLogs, nil)

	deviceHandler := NewDeviceHandler(deviceService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?limit=2", nil)
	deviceHandler.GetDeviceLogsLimitN(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `[
		{
			"id": 1,
			"created_at": "2023-01-01T00:00:00Z",
			"is_enabled": true
		},
		{
			"id": 2,
			"created_at": "2023-01-02T00:00:00Z",
			"is_enabled": false
		}]`, w.Body.String())
}

func TestGetDeviceLogsFromDataServiceLimitN_ParsingFailure(t *testing.T) {
	deviceService := new(service.MockDeviceService)
	deviceService.EXPECT().GetDeviceLogsFromDataServiceLimitN(2).Return([]domain.DeviceData{}, utils.ErrParsingFailed)

	deviceHandler := NewDeviceHandler(deviceService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?limit=2", nil)
	deviceHandler.GetDeviceLogsLimitN(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Parsing failed", w.Body.String())
}

func TestGetDeviceLogsFromDataServiceLimitN_InvalidLimit(t *testing.T) {
	deviceService := new(service.MockDeviceService)

	deviceHandler := NewDeviceHandler(deviceService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/?limit=invalid", nil)
	deviceHandler.GetDeviceLogsLimitN(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid limit parameter", w.Body.String())
}
