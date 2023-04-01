package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	sHttp "github.com/pklimuk-eng-thesis/control-station/pkg/http"
	sService "github.com/pklimuk-eng-thesis/control-station/pkg/service"
)

func main() {
	serviceAddress := setupServiceAddress("ADDRESS", ":8080")
	presenceSensorAddress := setupServiceAddress("PRESENCE_SENSOR_ADDRESS", "http://localhost:8081")
	gasSensorAddress := setupServiceAddress("GAS_SENSOR_ADDRESS", "http://localhost:8082")
	doorsSensorAddress := setupServiceAddress("DOORS_SENSOR_ADDRESS", "http://localhost:8083")

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	initializeSensor("PresenceSensor", presenceSensorAddress, "/presenceSensor", r)
	initializeSensor("GasSensor", gasSensorAddress, "/gasSensor", r)
	initializeSensor("DoorsSensor", doorsSensorAddress, "/doorsSensor", r)

	log.Printf("Starting service at %s\n", serviceAddress)
	log.Fatal(r.Run(serviceAddress))
}

func initializeSensor(name string, address string, groupName string, r *gin.Engine) {
	sensor := domain.Sensor{Name: name, Address: address}
	sensorService := sService.NewSensorService(&sensor)
	sensorHandler := sHttp.NewSensorHandler(sensorService)
	sHttp.SetupSensorRouter(r, sensorHandler, groupName)
}

func setupServiceAddress(identifier string, defaultAddress string) string {
	address := os.Getenv(identifier)
	if address == "" {
		address = defaultAddress
	}
	return address
}
