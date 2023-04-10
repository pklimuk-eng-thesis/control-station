package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sService "github.com/pklimuk-eng-thesis/control-station/pkg/service"
)

type SensorHandler struct {
	service sService.SensorService
}

func NewSensorHandler(service sService.SensorService) *SensorHandler {
	return &SensorHandler{service: service}
}

func (h *SensorHandler) GetInfo(c *gin.Context) {
	sensorInfo, err := h.service.GetInfo()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &sensorInfo)
}

func (h *SensorHandler) ToggleEnabled(c *gin.Context) {
	sensorInfo, err := h.service.ToggleEnabled()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &sensorInfo)
}

func (h *SensorHandler) ToggleDetected(c *gin.Context) {
	sensorInfo, err := h.service.ToggleDetected()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &sensorInfo)
}

func (h *SensorHandler) GetSensorLogsLimitN(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	sensorLogs, err := h.service.GetSensorLogsFromDataServiceLimitN(limit)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &sensorLogs)
}
