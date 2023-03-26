package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pklimuk-eng-thesis/control-station/pkg/domain"
	sHttp "github.com/pklimuk-eng-thesis/control-station/pkg/http"
	sService "github.com/pklimuk-eng-thesis/control-station/pkg/service"
)

func main() {
	r := gin.Default()

	initializeSensor("PresenceSensor", "http://localhost:8081", "/presenceSensor", r)
	initializeSensor("GasSensor", "http://localhost:8082", "/gasSensor", r)

	// Gets a service address from the environment variable or uses the default one
	// serviceAddress := viper.GetString("ADDRESS")
	serviceAddress := os.Getenv("ADDRESS")
	if serviceAddress == "" {
		serviceAddress = ":8080"
	}
	log.Printf("Starting service at %s\n", serviceAddress)
	log.Fatal(r.Run(serviceAddress))
}

func initializeSensor(name string, address string, groupName string, r *gin.Engine) {
	sensor := domain.Sensor{Name: name, Address: address}
	sensorService := sService.NewSensorService(&sensor)
	sensorHandler := sHttp.NewSensorHandler(sensorService)
	sHttp.SetupSensorRouter(r, sensorHandler, groupName)
}
