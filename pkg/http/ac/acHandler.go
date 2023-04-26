package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	acService "github.com/pklimuk-eng-thesis/control-station/pkg/service/ac"
)

type ACHandler struct {
	service acService.ACService
}

func NewACHandler(service acService.ACService) *ACHandler {
	return &ACHandler{service: service}
}

func (h *ACHandler) GetInfo(c *gin.Context) {
	acInfo, err := h.service.GetInfo()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &acInfo)
}

func (h *ACHandler) ToggleEnabled(c *gin.Context) {
	acInfo, err := h.service.ToggleEnabled()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &acInfo)
}

func (h *ACHandler) UpdateACSettings(c *gin.Context) {
	var desiredSettings domain.ACInfo
	err := c.BindJSON(&desiredSettings)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	acInfo, err := h.service.UpdateACSettings(desiredSettings.Temperature, desiredSettings.Humidity)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &acInfo)
}

func (h *ACHandler) GetACLogsLimitN(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	acLogs, err := h.service.GetACLogsFromDataServiceLimitN(limit)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, &acLogs)
}
