package http

import (
	"github.com/gin-gonic/gin"
	device "github.com/pklimuk-eng-thesis/control-station/pkg/http/device"
	sensor "github.com/pklimuk-eng-thesis/control-station/pkg/http/sensor"
)

var enabledEndpoint = "/enabled"
var detectedEndpoint = "/detected"
var infoEndpoint = "/info"
var logsEndpoint = "/logs"

func SetupSensorRouter(r *gin.Engine, sH *sensor.SensorHandler, groupName string) {
	route := r.Group(groupName)
	route.GET(infoEndpoint, sH.GetInfo)
	route.PATCH(enabledEndpoint, sH.ToggleEnabled)
	route.PATCH(detectedEndpoint, sH.ToggleDetected)
	route.GET(logsEndpoint, sH.GetSensorLogsLimitN)
}

func SetupDeviceRouter(r *gin.Engine, dH *device.DeviceHandler, groupName string) {
	route := r.Group(groupName)
	route.GET(infoEndpoint, dH.GetInfo)
	route.PATCH(enabledEndpoint, dH.ToggleEnabled)
	route.GET(logsEndpoint, dH.GetDeviceLogsLimitN)
}
