package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	sService "github.com/pklimuk-eng-thesis/control-station/pkg/service"
)

type SensorHandler struct {
	service sService.SensorService
}

func NewSensorHandler(service sService.SensorService) *SensorHandler {
	return &SensorHandler{service: service}
}

func (h *SensorHandler) IsEnabled(c *gin.Context) {
	enabled, err := h.service.IsEnabled()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &enabled)
}

func (h *SensorHandler) Detected(c *gin.Context) {
	detected, err := h.service.Detected()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &detected)
}

func (h *SensorHandler) ToggleEnabled(c *gin.Context) {
	enabled, err := h.service.ToggleEnabled()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &enabled)
}

func (h *SensorHandler) ToggleDetected(c *gin.Context) {
	detected, err := h.service.ToggleDetected()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &detected)
}
