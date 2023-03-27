package http

import "github.com/gin-gonic/gin"

var enabledEndpoint = "/enabled"
var detectedEndpoint = "/detected"

func SetupSensorRouter(r *gin.Engine, sH *SensorHandler, groupName string) {
	route := r.Group(groupName)
	route.GET(enabledEndpoint, sH.IsEnabled)
	route.GET(detectedEndpoint, sH.Detected)
	route.POST(enabledEndpoint, sH.ToggleEnabled)
	route.POST(detectedEndpoint, sH.ToggleDetected)
}
