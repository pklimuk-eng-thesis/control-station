package http

import "github.com/gin-gonic/gin"

var enabledEndpoint = "/enabled"
var detectedEndpoint = "/detected"
var infoEndpoint = "/info"

func SetupSensorRouter(r *gin.Engine, sH *SensorHandler, groupName string) {
	route := r.Group(groupName)
	route.GET(infoEndpoint, sH.GetInfo)
	route.POST(enabledEndpoint, sH.ToggleEnabled)
	route.POST(detectedEndpoint, sH.ToggleDetected)
}
