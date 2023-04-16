package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	deviceService "github.com/pklimuk-eng-thesis/control-station/pkg/service/device"
)

type DeviceHandler struct {
	service deviceService.DeviceService
}

func NewDeviceHandler(service deviceService.DeviceService) *DeviceHandler {
	return &DeviceHandler{service: service}
}

func (h *DeviceHandler) GetInfo(c *gin.Context) {
	deviceInfo, err := h.service.GetInfo()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &deviceInfo)
}

func (h *DeviceHandler) ToggleEnabled(c *gin.Context) {
	deviceInfo, err := h.service.ToggleEnabled()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &deviceInfo)
}

func (h *DeviceHandler) GetDeviceLogsLimitN(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	deviceLogs, err := h.service.GetDeviceLogsFromDataServiceLimitN(limit)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &deviceLogs)
}
