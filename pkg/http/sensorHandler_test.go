package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	sService "github.com/pklimuk-eng-thesis/control-station/pkg/service"
)

type mockSensorService struct {
	mock.Mock
}

func (m *mockSensorService) IsEnabled() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *mockSensorService) Detected() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *mockSensorService) ToggleEnabled() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *mockSensorService) ToggleDetected() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func TestSensorHandler_IsEnabled_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockSensorService)
	handler := NewSensorHandler(mockSvc)

	r := gin.Default()
	r.GET(enabledEndpoint, handler.IsEnabled)

	t.Run("success", func(t *testing.T) {
		mockSvc.On("IsEnabled").Return(true, nil)
		req, _ := http.NewRequest(http.MethodGet, enabledEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "true", w.Body.String())
	})
}

func TestSensorHandler_IsEnabled_Error_ParsingFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockSensorService)
	handler := NewSensorHandler(mockSvc)

	r := gin.Default()
	r.GET(enabledEndpoint, handler.IsEnabled)

	t.Run("error", func(t *testing.T) {
		mockSvc.On("IsEnabled").Return(false, sService.ErrParsingFailed)
		req, _ := http.NewRequest(http.MethodGet, enabledEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, sService.ErrParsingFailed.Error(), w.Body.String())
	})
}

func TestSensorHandler_Detected_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockSensorService)
	handler := NewSensorHandler(mockSvc)

	r := gin.Default()
	r.GET(detectedEndpoint, handler.Detected)

	t.Run("success", func(t *testing.T) {
		mockSvc.On("Detected").Return(true, nil)
		req, _ := http.NewRequest(http.MethodGet, detectedEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "true", w.Body.String())
	})
}

func TestSensorHandler_Detected_Error_ParsingFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockSensorService)
	handler := NewSensorHandler(mockSvc)

	r := gin.Default()
	r.GET(detectedEndpoint, handler.Detected)

	t.Run("error", func(t *testing.T) {
		mockSvc.On("Detected").Return(false, sService.ErrParsingFailed)
		req, _ := http.NewRequest(http.MethodGet, detectedEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, sService.ErrParsingFailed.Error(), w.Body.String())
	})
}

func TestSensorHandler_ToggleEnabled_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockSensorService)
	handler := NewSensorHandler(mockSvc)

	r := gin.Default()
	r.POST(enabledEndpoint, handler.ToggleEnabled)

	t.Run("success", func(t *testing.T) {
		mockSvc.On("ToggleEnabled").Return(true, nil)
		req, _ := http.NewRequest(http.MethodPost, enabledEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "true", w.Body.String())
	})
}

func TestSensorHandler_ToggleEnabled_Error_ParsingFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockSensorService)
	handler := NewSensorHandler(mockSvc)

	r := gin.Default()
	r.POST(enabledEndpoint, handler.ToggleEnabled)

	t.Run("error", func(t *testing.T) {
		mockSvc.On("ToggleEnabled").Return(false, sService.ErrParsingFailed)
		req, _ := http.NewRequest(http.MethodPost, enabledEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, sService.ErrParsingFailed.Error(), w.Body.String())
	})
}

func TestSensorHandler_ToggleDetected_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockSensorService)
	handler := NewSensorHandler(mockSvc)

	r := gin.Default()
	r.POST(detectedEndpoint, handler.ToggleDetected)

	t.Run("success", func(t *testing.T) {
		mockSvc.On("ToggleDetected").Return(true, nil)
		req, _ := http.NewRequest(http.MethodPost, detectedEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "true", w.Body.String())
	})
}

func TestSensorHandler_ToggleDetected_Error_ParsingFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockSensorService)
	handler := NewSensorHandler(mockSvc)

	r := gin.Default()
	r.POST(detectedEndpoint, handler.ToggleDetected)

	t.Run("error", func(t *testing.T) {
		mockSvc.On("ToggleDetected").Return(false, sService.ErrParsingFailed)
		req, _ := http.NewRequest(http.MethodPost, detectedEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, sService.ErrParsingFailed.Error(), w.Body.String())
	})
}
