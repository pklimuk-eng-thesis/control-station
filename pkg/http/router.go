package http

import "github.com/gin-gonic/gin"

var enabledEndpoint = "/enabled"

func SetupSensorRouter(r *gin.Engine, sH *SensorHandler, groupName string) {
	route := r.Group(groupName)
	route.GET(enabledEndpoint, sH.IsEnabled)
}
