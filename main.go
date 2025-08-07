package main

import (
	"fmt"

	db "github.com/Fonzeca/Trackin/internal/infrastructure/database"
	messagingInfra "github.com/Fonzeca/Trackin/internal/infrastructure/messaging"
	server "github.com/Fonzeca/Trackin/internal/interfaces/http"
	entry "github.com/Fonzeca/Trackin/internal/interfaces/messaging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	InitConfig()

	// go monitor.System()
	db.InitDB()
	defer db.CloseDB()

	startRabbit := false
	if startRabbit {
		_, closeFunc := messagingInfra.SetupRabbitMq()
		defer closeFunc()
		entry.NewRabbitMqDataEntry()
	}

	e := echo.New()

	api := server.NewApi()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/getLastLogByImei", api.GetLastLogByImei)
	e.POST("/getVehiclesStateByImeis", api.GetVehiclesStateByImeis)
	e.POST("/getRouteByImei", api.GetRouteByImei)
	e.GET("/getSummaryRoutesAndZones", api.GetSummaryRoutesAndZones)

	e.GET("/getZonesByEmpresaId", api.GetZonesByEmpresaId)
	e.POST("/createZone", api.CreateZone)
	e.PUT("/editZoneById", api.EditZoneById)
	e.DELETE("/deleteZoneById", api.DeleteZoneById)

	e.Logger.Fatal(e.Start(":4762"))
}

func InitConfig() {
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
