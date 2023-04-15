package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	"github.com/pklimuk-eng-thesis/control-station/pkg/http"
	sensorHttp "github.com/pklimuk-eng-thesis/control-station/pkg/http/sensor"
	sensorService "github.com/pklimuk-eng-thesis/control-station/pkg/service/sensor"
	"github.com/pklimuk-eng-thesis/control-station/utils"
)

func main() {
	serviceAddress := utils.GetEnvVariableOrDefault("ADDRESS", ":8080")
	presenceSensorAddress := utils.GetEnvVariableOrDefault("PRESENCE_SENSOR_ADDRESS", "http://localhost:8081")
	gasSensorAddress := utils.GetEnvVariableOrDefault("GAS_SENSOR_ADDRESS", "http://localhost:8082")
	doorsSensorAddress := utils.GetEnvVariableOrDefault("DOORS_SENSOR_ADDRESS", "http://localhost:8083")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	initializeSensor("presenceSensor", presenceSensorAddress, "/presenceSensor", r)
	initializeSensor("gasSensor", gasSensorAddress, "/gasSensor", r)
	initializeSensor("doorsSensor", doorsSensorAddress, "/doorsSensor", r)

	log.Printf("Starting service at %s\n", serviceAddress)
	log.Fatal(r.Run(serviceAddress))
}

func initializeSensor(name string, address string, groupName string, r *gin.Engine) {
	sensor := domain.Sensor{Name: name, Address: address}
	sensorService := sensorService.NewSensorService(&sensor)
	sensorHandler := sensorHttp.NewSensorHandler(sensorService)
	http.SetupSensorRouter(r, sensorHandler, groupName)
}
